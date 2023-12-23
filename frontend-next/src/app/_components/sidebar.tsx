"use client"

interface SidebarProps {
	visible: boolean;
}
export function Sidebar({visible}: SidebarProps) {
	return (
		<nav className="w-48 pl-2 pt-2 absolute bg-gray-100 dark:bg-gray-900 h-full">
			<SidebarNavGroup label="Horses">
				<SidebarNavItem label="Overview" href="/horses"/>
			</SidebarNavGroup>
			<SidebarNavGroup label="Club">
				<SidebarNavItem label="General" href="/general"/>
				<SidebarNavItem label="Members" href="/members"/>
			</SidebarNavGroup>
			<SidebarNavGroup label="Admin">
				<SidebarNavItem label="Settings" href="/settings"/>
				<SidebarNavItem label="Users" href="/users"/>
			</SidebarNavGroup>		
		</nav>
	);
}

interface SidebarNavGroupProps {
	label: string;
	children: React.ReactNode;
}
export async function SidebarNavGroup({label, children}: SidebarNavGroupProps) {
	return (
		<div className="pb-8">
			<h2 className="uppercase font-light text-xs">{label}</h2>
			{children}
		</div>
	);
}

interface SidebarNavItemProps {
	label: string;
	href: string;
}
export async function SidebarNavItem({label, href}: SidebarNavItemProps) {
	return (
		<div className="pt-2">
			<h2>{label}</h2>
		</div>
	);
}