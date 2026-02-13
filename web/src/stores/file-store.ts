import { create } from "zustand";

// A node is either a file or a directory
export type FileNode =
	| { kind: "file"; name: string; language?: string; content: string }
	| { kind: "dir"; name: string; children: FileNode[] };

export type FilePath = string; // e.g. "src/utils/helpers.bal"

// ---------------------------------------------------------------------------
// Store interface — mutations only.
//
// Derived values (selectedFile, getNodeAt) are NOT stored as functions on the
// store object. Storing functions-as-values causes Zustand to return a new
// function reference on every state update, which makes React see a changed
// value every render → infinite re-render loop.
//
// Use the selector hooks exported below instead.
// ---------------------------------------------------------------------------

export interface FileState {
	tree: FileNode[];
	selectedFilePath: FilePath | null;

	setTree: (tree: FileNode[]) => void;
	selectFile: (path: FilePath) => void;
	updateFileContent: (path: FilePath, content: string) => void;
	addNode: (parentPath: FilePath | null, node: FileNode) => void;
	deleteNode: (path: FilePath) => void;
	renameNode: (path: FilePath, newName: string) => void;
}

// ---------------------------------------------------------------------------
// Pure tree helpers (module-level, not part of store state)
// ---------------------------------------------------------------------------

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
const DEFAULT_TREE: FileNode[] = [
	{
		kind: "dir",
		name: "hello",
		children: [
			{
				kind: "file",
				name: "main.bal",
				language: "ballerina",
				content: `import ballerina/io;

public function main() {
	io:println("Hello, World!");
}`,
			},
			{
				kind: "file",
				name: "Ballerina.toml",
				language: "toml",
				content: `[package]
org = "wso2"
name = "ballerina"
version = "0.1.0"`,
			},
		],
	},
	{
		kind: "dir",
		name: "fibonacci",
		children: [
			{
				kind: "file",
				name: "main.bal",
				language: "ballerina",
				content: `import ballerina/io;

public function main() {
    int n = 10;
    int i = 0;
    while (i < n) {
        io:println("F(", i, ") = ", fibonacci(i));
        i += 1;
    }
}

function fibonacci(int n) returns int {
    if (n <= 1) {
        return n;
    }
    int prev = 0;
    int curr = 1;
    int i = 2;
    while (i <= n) {
        int next = prev + curr;
        prev = curr;
        curr = next;
        i += 1;
    }
    return curr;
}`,
			},
			{
				kind: "file",
				name: "Ballerina.toml",
				language: "toml",
				content: `[package]
org = "wso2"
name = "ballerina"
version = "0.1.0"`,
			},
		],
	},
];

const DEFAULT_SELECTED: FilePath = "hello/main.bal";

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

// ---------------------------------------------------------------------------
// Selector hooks
//
// Core rule: Zustand compares selector return values with ===.
// Returning a freshly-derived object from inside a selector always produces
// a new reference → Zustand thinks state changed → re-render → infinite loop.
//
// Solution: selectors return only PRIMITIVES (strings, null).
// The actual node object is looked up *outside* the selector using the
// stable `tree` reference and the primitive path, so React only re-renders
// when the path string or the tree reference actually changes.
// ---------------------------------------------------------------------------

/**
 * Returns the currently selected file node, or null.
 *
 * Splits the work in two:
 *  1. A selector that returns a primitive (the path string) — stable ===
 *  2. A plain getNodeAt call outside the selector using the current tree
 */
export function useSelectedFile(): Extract<FileNode, { kind: "file" }> | null {
	const tree = useFileStore((state) => state.tree);
	const path = useFileStore((state) => state.selectedFilePath);
	if (!path) return null;
	const node = getNodeAt(tree, path);
	return node?.kind === "file" ? node : null;
}

/**
 * Returns the node at an arbitrary path, or null.
 */
export function useNodeAt(path: FilePath): FileNode | null {
	const tree = useFileStore((state) => state.tree);
	return getNodeAt(tree, path);
}
