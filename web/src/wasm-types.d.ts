declare global {
	export interface Window {
		Go: any;
		updateFile(fileName: string, content: string): void;
		runProject(path: string): string;
	}
}
export {};
