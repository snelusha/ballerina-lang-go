declare global {
	export interface Window {
		Go: any;
		wasmFibonacciSum: (n: number) => number;
		updateFile(fileName: string, content: string): void;
		runProject(path: string): void;
	}
}
export {};
