const STORAGE_KEY = "bfs";

class BrowserFS {
    constructor() {
        this._load();
        if (Object.keys(this.data.children).length === 0) {
            this._reset();
        }
    }

    _load() {
        try {
            const raw = localStorage.getItem(STORAGE_KEY);
            this.data = raw ? JSON.parse(raw) : { isDir: true, children: {} };
        } catch {
            this.data = { isDir: true, children: {} };
        }
    }

    _reset() {
        this.data = { isDir: true, children: {} };
        const d = {
            "main.bal": `
import ballerina/io;

public function main() {
    int a = "10";
    io:println("Hello, World!");
}
`,
        };

        for (const [name, content] of Object.entries(d)) {
            const { parent, name: fileName } = this._getParentAndName(name);
            if (!parent) continue;
            if (!parent.children) parent.children = {};
            parent.children[fileName] = {
                isDir: false,
                content,
                modTime: Date.now(),
            };
        }

        this._save();
    }

    _save() {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(this.data));
    }

    _getNode(path, autoCreateDirs = true) {
        if (!path || path === "." || path === "/") return this.data;
        const parts = path.split("/").filter(Boolean);

        let node = this.data;
        for (const part of parts) {
            if (!node.children || !node.children[part]) {
                if (autoCreateDirs) {
                    if (!node.children) node.children = {};
                    node.children[part] = { isDir: true, children: {} };
                } else {
                    return null;
                }
            }
            node = node.children[part];
        }
        return node;
    }

    _getParentNode(path) {
        const parts = path.split("/").filter(Boolean);
        if (parts.length === 0) return null;
        const parentPath = parts.slice(0, -1).join("/");
        const parentNode = this._getNode(parentPath);
        return parentNode;
    }

    _getParentAndName(path) {
        const parts = path.split("/").filter(Boolean);
        if (parts.length === 0) return { parent: null, name: "/" };
        const name = parts.slice(-1)[0];
        const parentPath = parts.slice(0, -1).join("/");
        const parentNode = this._getNode(parentPath);
        return { parent: parentNode, name };
    }

    open(path) {
        const node = this._getNode(path);
        if (!node || node.isDir) return null;
        return {
            content: node.content,
            size: node.content.length,
            modTime: node.modTime || 0,
            isDir: false,
        };
    }

    stat(path) {
        const node = this._getNode(path);
        if (!node) return null;
        return {
            name: path.split("/").filter(Boolean).slice(-1)[0] || "/",
            size: node.isDir ? 0 : node.content ? node.content.length : 0,
            modTime: node.modTime || 0,
            isDir: node.isDir,
        };
    }

    readDir(path) {
        const node = this._getNode(path);
        if (!node || !node.isDir) return null;
        return Object.entries(node.children).map(([name, child]) => ({
            name,
            isDir: child.isDir,
        }));
    }

    writeFile(path, content) {
        try {
            const { parent, name } = this._getParentAndName(path);
            if (!parent) return false;
            if (
                parent.children &&
                parent.children[name] &&
                parent.children[name].isDir
            )
                return false;

            if (!parent.children) parent.children = {};
            parent.children[name] = {
                isDir: false,
                content,
                modTime: Date.now(),
            };
            this._save();
            return true;
        } catch {
            return false;
        }
    }

    remove(path) {
        const parts = path.split("/").filter(Boolean);
        if (parts.length === 0) return false;

        const name = parts.slice(-1)[0];
        const parentPath = parts.slice(0, -1).join("/");
        const parentNode = this._getNode(parentPath);
        if (!parentNode || !parentNode.children || !parentNode.children[name])
            return false;
        delete parentNode.children[name];

        this._save();
        return true;
    }

    move(oldPath, newPath) {
        try {
            const node = this._getNode(oldPath);
            if (!node) return false;

            const { parent: newParent, name: newName } =
                this._getParentAndName(newPath);
            const { parent: oldParent, name: oldName } =
                this._getParentAndName(oldPath);
            if (!newParent || !oldParent) return false;
            if (
                newParent.children &&
                newParent.children[newName] &&
                newParent.children[newName].isDir
            )
                return false;

            if (!newParent.children) newParent.children = {};
            newParent.children[newName] = node;
            delete oldParent.children[oldName];

            this._save();
            return true;
        } catch {
            return false;
        }
    }
}

