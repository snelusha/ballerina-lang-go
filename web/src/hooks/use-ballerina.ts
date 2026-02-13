import "@/wasm_exec";
// import "@/wasm-types";

import * as React from "react";

export function useBallerina() {
	const [loaded, setLoaded] = React.useState(false);

	React.useEffect(() => {
		async function loadWasm() {
			const go = new window.Go();
			const result = await WebAssembly.instantiateStreaming(
				fetch("ballerina.wasm"),
				go.importObject,
			);
			go.run(result.instance);
			setLoaded(true);
		}

		loadWasm();
	}, []);

	function updateFile(fileName: string, content: string) {
		window.updateFile(fileName, content);
	}

	function runProject(path: string) {
		window.runProject(path);
	}

	return { ready: loaded, updateFile, runProject };
}
