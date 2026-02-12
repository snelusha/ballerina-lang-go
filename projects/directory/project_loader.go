/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package directory provides project loading functionality from filesystem directories.
package directory

import (
	"io/fs"
	"path/filepath"
	"strings"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/internal"
	"ballerina-lang-go/tools/diagnostics"
)

// ProjectLoadConfig holds configuration for project loading.
// All fields are optional - nil values use defaults.
type ProjectLoadConfig struct {
	// BuildOptions configures compilation behavior. If nil, defaults are used.
	BuildOptions *projects.BuildOptions

	// Future fields can be added here, e.g.:
	// EnvironmentBuilder *EnvironmentBuilder
}

// LoadProject loads a project from the given path.
// It detects the project type and delegates to the appropriate loader:
//   - Has Ballerina.toml -> loadBuildProject
//   - Is .bal file -> loadSingleFileProject
//   - Is .bala file -> error (not implemented)
//
// If no config is provided, default configuration is used.
func LoadProject(fsys fs.FS, path string, config ...ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	// Apply defaults
	var cfg ProjectLoadConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// We trust the path provided relative to fs.FS
	// os.Stat checks if path exists
	slashPath := filepath.ToSlash(path)
	if slashPath == "." {
		slashPath = "."
	} else if slashPath == "" {
		slashPath = "."
	}

	info, err := fs.Stat(fsys, slashPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// Detect project type
	if info.IsDir() {
		// Check for Ballerina.toml
		// For fs.FS, join with forward slash is safer, but filepath.Join then ToSlash works
		tomlPath := filepath.Join(path, projects.BallerinaTomlFile)
		slashTomlPath := filepath.ToSlash(tomlPath)
		if _, err := fs.Stat(fsys, slashTomlPath); err == nil {
			// Has Ballerina.toml - load as build project
			return loadBuildProject(fsys, path, cfg)
		}

		// Directory without Ballerina.toml - error
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "not a valid Ballerina project directory (missing Ballerina.toml): " + path,
		}
	}

	// Check file extension
	if strings.HasSuffix(path, projects.BalFileExtension) {
		// Single .bal file
		return loadSingleFileProject(fsys, path, cfg)
	}

	if strings.HasSuffix(path, projects.BalaFileExtension) {
		// .bala file - not implemented
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "loading from .bala files is not implemented: " + path,
		}
	}

	return projects.ProjectLoadResult{}, &projects.ProjectError{
		Message: "unsupported file type: " + path,
	}
}

// loadBuildProject loads a build project from the given path.
// It merges build options from Ballerina.toml (manifest defaults) with the caller's
// options using AcceptTheirs, so caller-provided options override manifest defaults.
func loadBuildProject(fsys fs.FS, path string, cfg ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	// Use internal.CreateBuildProjectConfig to scan and create package config
	packageConfig, err := internal.CreateBuildProjectConfig(fsys, path)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// Merge build options: manifest defaults are the base, caller's options override.
	manifestBuildOptions := packageConfig.PackageManifest().BuildOptions()
	var mergedOpts projects.BuildOptions
	if cfg.BuildOptions != nil {
		mergedOpts = manifestBuildOptions.AcceptTheirs(*cfg.BuildOptions)
	} else {
		mergedOpts = manifestBuildOptions
	}

	// Create the project
	// Note: We are using the path relative to FS as the source root.
	// This abstract path will be stored in the project.
	project := projects.NewBuildProject(path, mergedOpts)

	// Create package from config
	compilationOptions := mergedOpts.CompilationOptions()
	pkg := projects.NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	// Collect diagnostics from manifest
	var diags []diagnostics.Diagnostic
	manifestDiags := packageConfig.PackageManifest().Diagnostics()
	diags = append(diags, manifestDiags...)

	// Create diagnostic result
	diagResult := projects.NewDiagnosticResult(diags)

	return projects.NewProjectLoadResult(project, diagResult), nil
}

// loadSingleFileProject loads a single .bal file as a project.
func loadSingleFileProject(fsys fs.FS, path string, cfg ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	// Verify file exists and is a .bal file
	slashPath := filepath.ToSlash(path)
	info, err := fs.Stat(fsys, slashPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	if info.IsDir() {
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "expected a .bal file, got directory: " + path,
		}
	}

	if !strings.HasSuffix(path, projects.BalFileExtension) {
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "not a Ballerina source file: " + path,
		}
	}

	// Read file content
	content, err := fs.ReadFile(fsys, slashPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// Get build options or use defaults
	var buildOpts projects.BuildOptions
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	} else {
		buildOpts = projects.NewBuildOptions()
	}

	// Get directory and filename
	sourceDir := filepath.Dir(path)
	fileName := filepath.Base(path)

	// Derive package name from filename (without extension)
	packageName := strings.TrimSuffix(fileName, projects.BalFileExtension)

	// Create the project
	project := projects.NewSingleFileProject(sourceDir, buildOpts, path)

	// Create package descriptor with anonymous org and default version
	defaultVersion, _ := projects.NewPackageVersionFromString(projects.DefaultVersion)
	packageDesc := projects.NewPackageDescriptor(
		projects.NewPackageOrg(projects.DefaultOrg),
		projects.NewPackageName(packageName),
		defaultVersion,
	)

	// Create manifest with minimal info
	manifest := projects.NewPackageManifest(packageDesc)

	// Create package ID
	packageID := projects.NewPackageID(packageName)

	// Create module descriptor for default module
	moduleDesc := projects.NewModuleDescriptorForDefaultModule(packageDesc)

	// Create module ID
	moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)

	// Create document config
	docID := projects.NewDocumentID(fileName, moduleID)
	docConfig := projects.NewDocumentConfig(docID, fileName, string(content))

	// Create module config with single source file
	moduleConfig := projects.NewModuleConfig(
		moduleID,
		moduleDesc,
		[]projects.DocumentConfig{docConfig},
		nil, // no test docs
		nil, // no readme
		nil, // no dependencies
	)

	// Create package config
	packageConfig := projects.NewPackageConfig(projects.PackageConfigParams{
		PackageID:       packageID,
		PackageManifest: manifest,
		PackagePath:     sourceDir,
		DefaultModule:   moduleConfig,
		OtherModules:    nil,
		BallerinaToml:   nil,
		ReadmeMd:        nil,
	})

	// Create package from config
	compilationOptions := buildOpts.CompilationOptions()
	pkg := projects.NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	// Single file projects have no diagnostics
	diagResult := projects.NewDiagnosticResult([]diagnostics.Diagnostic{})

	return projects.NewProjectLoadResult(project, diagResult), nil
}
