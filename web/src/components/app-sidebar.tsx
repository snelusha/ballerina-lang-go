import * as React from "react";

import { HugeiconsIcon } from "@hugeicons/react";
import {
	FolderIcon,
	File01Icon,
	ChevronDown,
} from "@hugeicons/core-free-icons";

import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
	Sidebar,
	SidebarContent,
	SidebarGroup,
	SidebarGroupContent,
	SidebarGroupLabel,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	useSidebar,
} from "@/components/ui/sidebar";

import { useFileStore } from "@/stores/file-store";

import type { FileNode, FilePath } from "@/stores/file-store";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
	const tree = useFileStore((s) => s.tree);

	return (
		<Sidebar {...props}>
			<SidebarContent>
				<SidebarGroup>
					<SidebarGroupLabel>Files</SidebarGroupLabel>
					<SidebarGroupContent className="mt-2">
						<SidebarMenu>
							{tree.map((node) => (
								<TreeNode key={node.name} node={node} path={node.name} />
							))}
						</SidebarMenu>
					</SidebarGroupContent>
				</SidebarGroup>
			</SidebarContent>
		</Sidebar>
	);
}

function TreeNode({ node, path }: { node: FileNode; path: FilePath }) {
	const selectFile = useFileStore((s) => s.selectFile);
	const selectedFilePath = useFileStore((s) => s.selectedFilePath);
	const { toggleSidebar } = useSidebar();

	if (node.kind === "file") {
		return (
			<SidebarMenuButton
				isActive={selectedFilePath === path}
				onClick={() => {
					selectFile(path);
					toggleSidebar();
				}}
			>
				<HugeiconsIcon icon={File01Icon} strokeWidth={1.5} />
				{node.name}
			</SidebarMenuButton>
		);
	}

	return (
		<SidebarMenuItem>
			<Collapsible
				defaultOpen
				className="group/collapsible [&[data-state=open]>button>svg:first-child]:rotate-90"
			>
				<CollapsibleTrigger className="w-full">
					<SidebarMenuButton>
						<HugeiconsIcon icon={ChevronDown} strokeWidth={1.5} />
						<HugeiconsIcon icon={FolderIcon} strokeWidth={1.5} />
						{node.name}
					</SidebarMenuButton>
				</CollapsibleTrigger>
				<CollapsibleContent>
					<SidebarMenuSub>
						{node.children.map((child) => (
							<TreeNode
								key={child.name}
								node={child}
								path={`${path}/${child.name}`}
							/>
						))}
					</SidebarMenuSub>
				</CollapsibleContent>
			</Collapsible>
		</SidebarMenuItem>
	);
}
