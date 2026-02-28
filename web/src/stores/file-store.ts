import { create } from "zustand";

import balTree from "@/data/bal-tree.json";

export type FileNode =
	| { kind: "file"; name: string; language?: string; content: string }
	| { kind: "dir"; name: string; children: FileNode[] };

export type FilePath = string;

const EMPTY_BALLERINA_TOML = `[package]
org = "wso2"
name = "ballerina"
version = "0.1.0"`;

const EMPTY_MAIN_BAL = `import ballerina/io;

public function main() {
		io:println("Hello, World!");
}`;

function makeEmptyProjectDir(name: string): FileNode {
	return {
		kind: "dir",
		name,
		children: [
			{
				kind: "file",
				name: "Ballerina.toml",
				language: "toml",
				content: EMPTY_BALLERINA_TOML,
			},
			{
				kind: "file",
				name: "main.bal",
				language: "ballerina",
				content: EMPTY_MAIN_BAL,
			},
		],
	};
}

function uniqueProjectName(tree: FileNode[], base = "new-project"): string {
	const topNames = new Set(tree.map((n) => n.name));
	let name = base;
	let i = 1;
	while (topNames.has(name)) {
		name = `${base}-${++i}`;
	}
	return name;
}

export interface FileState {
	tree: FileNode[];
	selectedFilePath: FilePath | null;

	setTree: (tree: FileNode[]) => void;
	selectFile: (path: FilePath) => void;
	updateFileContent: (path: FilePath, content: string) => void;
	addNode: (parentPath: FilePath | null, node: FileNode) => void;
	createEmptyProject: () => void;
	deleteNode: (path: FilePath) => void;
	renameNode: (path: FilePath, newName: string) => void;
}

function segments(path: FilePath): string[] {
	return path.split("/").filter(Boolean);
}

export function getNodeAt(tree: FileNode[], path: FilePath): FileNode | null {
	const parts = segments(path);
	let current: FileNode[] = tree;

	for (let i = 0; i < parts.length; i++) {
		const found = current.find((n) => n.name === parts[i]);
		if (!found) return null;
		if (i === parts.length - 1) return found;
		if (found.kind !== "dir") return null;
		current = found.children;
	}

	return null;
}

function updateNodeAt(
	tree: FileNode[],
	path: FilePath,
	updater: (node: FileNode) => FileNode,
): FileNode[] {
	const parts = segments(path);

	function recurse(nodes: FileNode[], depth: number): FileNode[] {
		return nodes.map((node) => {
			if (node.name !== parts[depth]) return node;
			if (depth === parts.length - 1) return updater(node);
			if (node.kind !== "dir") return node;
			return { ...node, children: recurse(node.children, depth + 1) };
		});
	}

	return recurse(tree, 0);
}

function deleteNodeAt(tree: FileNode[], path: FilePath): FileNode[] {
	const parts = segments(path);

	function recurse(nodes: FileNode[], depth: number): FileNode[] {
		if (depth === parts.length - 1) {
			return nodes.filter((n) => n.name !== parts[depth]);
		}
		return nodes.map((node) => {
			if (node.name !== parts[depth] || node.kind !== "dir") return node;
			return { ...node, children: recurse(node.children, depth + 1) };
		});
	}

	return recurse(tree, 0);
}

function addNodeAt(
	tree: FileNode[],
	parentPath: FilePath | null,
	node: FileNode,
): FileNode[] {
	if (!parentPath) return [...tree, node];

	return updateNodeAt(tree, parentPath, (parent) => {
		if (parent.kind !== "dir") {
			console.warn("Cannot add a child to a file node.");
			return parent;
		}
		return { ...parent, children: [...parent.children, node] };
	});
}

export function allFilePaths(tree: FileNode[], prefix = ""): FilePath[] {
	const paths: FilePath[] = [];
	for (const node of tree) {
		const full = prefix ? `${prefix}/${node.name}` : node.name;
		if (node.kind === "file") {
			paths.push(full);
		} else {
			paths.push(...allFilePaths(node.children, full));
		}
	}
	return paths;
}
const DEFAULT_TREE: FileNode[] = balTree as FileNode[];

const DEFAULT_SELECTED: FilePath = "01-orders/main.bal";

export const useFileStore = create<FileState>((set) => ({
	tree: DEFAULT_TREE,
	selectedFilePath: DEFAULT_SELECTED,

	setTree: (tree) =>
		set((state) => {
			const paths = allFilePaths(tree);
			return {
				tree,
				selectedFilePath:
					state.selectedFilePath && paths.includes(state.selectedFilePath)
						? state.selectedFilePath
						: (paths[0] ?? null),
			};
		}),

	selectFile: (path) => set({ selectedFilePath: path }),

	updateFileContent: (path, content) =>
		set((state) => ({
			tree: updateNodeAt(state.tree, path, (node) =>
				node.kind === "file" ? { ...node, content } : node,
			),
		})),

	addNode: (parentPath, node) =>
		set((state) => ({
			tree: addNodeAt(state.tree, parentPath, node),
		})),

	createEmptyProject: () =>
		set((state) => {
			const name = uniqueProjectName(state.tree);
			const project = makeEmptyProjectDir(name);
			const newTree = [project, ...state.tree];
			return {
				tree: newTree,
				selectedFilePath: `${name}/main.bal`,
			};
		}),

	deleteNode: (path) =>
		set((state) => {
			const newTree = deleteNodeAt(state.tree, path);
			const paths = allFilePaths(newTree);
			return {
				tree: newTree,
				selectedFilePath:
					state.selectedFilePath === path
						? (paths[0] ?? null)
						: state.selectedFilePath,
			};
		}),

	renameNode: (path, newName) =>
		set((state) => {
			const parts = segments(path);
			const parentPath = parts.slice(0, -1).join("/") || null;
			const newPath = parentPath ? `${parentPath}/${newName}` : newName;

			const newTree = updateNodeAt(state.tree, path, (node) => ({
				...node,
				name: newName,
			}));

			return {
				tree: newTree,
				selectedFilePath:
					state.selectedFilePath === path ? newPath : state.selectedFilePath,
			};
		}),
}));

export function useSelectedFile(): Extract<FileNode, { kind: "file" }> | null {
	const tree = useFileStore((state) => state.tree);
	const path = useFileStore((state) => state.selectedFilePath);
	if (!path) return null;
	const node = getNodeAt(tree, path);
	return node?.kind === "file" ? node : null;
}

export function useNodeAt(path: FilePath): FileNode | null {
	const tree = useFileStore((state) => state.tree);
	return getNodeAt(tree, path);
}
