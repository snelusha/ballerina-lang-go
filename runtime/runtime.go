// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package runtime

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/runtime/internal/exec"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

// Runtime represents a Ballerina runtime instance that owns a module registry
// and is used as the execution context for interpreting BIR packages.
type Runtime struct {
	registry *modules.Registry
	host     Host
}

// ModuleInitializer is a function that can install modules (e.g. stdlibs) into
// a runtime instance during its construction.
type ModuleInitializer func(*Runtime)

var moduleInitializers []ModuleInitializer

// NewRuntime constructs a new runtime with DefaultHost, an empty registry,
// and runs all registered module initializers.
func NewRuntime() *Runtime {
	return NewRuntimeWithHost(DefaultHost())
}

// NewRuntimeWithHost constructs a new runtime with the given host. Nil FS or
// Stdout are replaced with the same defaults as DefaultHost.
func NewRuntimeWithHost(h Host) *Runtime {
	h = normalizeHost(h)
	rt := &Runtime{
		registry: modules.NewRegistry(),
		host:     h,
	}
	for _, init := range moduleInitializers {
		init(rt)
	}
	return rt
}

// Host returns the execution host (filesystem, stdout, and future capabilities).
func (rt *Runtime) Host() Host {
	return rt.host
}

// Interpret interprets a BIR package using this runtime instance.
func (rt *Runtime) Interpret(pkg bir.BIRPackage) (err error) {
	rt.registry.SetTypeEnv(pkg.TypeEnv)
	return exec.Interpret(pkg, rt.registry)
}

// RegisterModuleInitializer registers a module initializer that will be invoked
// for every newly created runtime.
func RegisterModuleInitializer(init ModuleInitializer) {
	moduleInitializers = append(moduleInitializers, init)
}

// GetTypeEnv returns the semantic type environment from the runtime's registry.
func (rt *Runtime) GetTypeEnv() semtypes.Env {
	return rt.registry.GetTypeEnv()
}

// RegisterExternFunction registers a native (extern) function implementation in
// the given runtime instance so it can be called from interpreted BIR code.
func RegisterExternFunction(rt *Runtime, orgName string, moduleName string, funcName string, impl func(args []values.BalValue) (values.BalValue, error)) {
	rt.registry.RegisterExternFunction(orgName, moduleName, funcName, impl)
}
