"use server";

import Link from "next/link";
import { getServerAuthSession } from "~/server/auth";
import { Avatar, AvatarFallback, AvatarImage } from "@components/ui/avatar";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuShortcut,
	DropdownMenuTrigger,
} from "@components/ui/dropdown-menu";
import { Button } from "@components/ui/button";

export async function Header() {
	return (
		<div id="header" className="top-0 h-12 flex items-center justify-between border-b-[1px]">
			<div className="pl-2">
				<h1 className="uppercase font-extrabold">hevonen</h1>
			</div>
			<div className="pr-2">
				<UserArea />
			</div>
		</div>
	);
}

export async function UserArea() {
	const session = await getServerAuthSession();
	const username = session?.user?.name ?? "JD";
	const email = session?.user?.email ?? "-";
	return (
		<div id="user-area" className="flex items-center">
			{session && 
				<UserAvatar 
					username={username}
					email={email}
				/>
			 }
			{!session && (
				<Link
					href={session ? "/api/auth/signout" : "/api/auth/signin"}
					className="rounded-full bg-white/10 px-10 py-3 font-semibold no-underline transition hover:bg-white/20"
				>
					Sign in
				</Link>
			)}
		</div>
	);
}

interface UserAvatarProps {
	username: string;
	email: string;
}
export async function UserAvatar({email, username}: UserAvatarProps) {
	const fallbackName = username.substring(0, 2).toUpperCase() ?? "JD";
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="ghost" className="relative h-8 w-8 rounded-full">
					<Avatar className="h-8 w-8">
						<AvatarImage src="/avatars/01.png" alt="@shadcn" />
						<AvatarFallback>{fallbackName}</AvatarFallback>
					</Avatar>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-56" align="end" forceMount>
				<DropdownMenuLabel className="font-normal">
					<div className="flex flex-col space-y-1">
						<p className="text-sm font-medium leading-none">{username}</p>
						<p className="text-xs leading-none text-muted-foreground">
							{email}
						</p>
					</div>
				</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuGroup>
					<DropdownMenuItem>
						Profile
						<DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
					</DropdownMenuItem>
					<DropdownMenuItem>
						Billing
						<DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
					</DropdownMenuItem>
					<DropdownMenuItem>
						Settings
						<DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
					</DropdownMenuItem>
					<DropdownMenuItem>New Team</DropdownMenuItem>
				</DropdownMenuGroup>
				<DropdownMenuSeparator />
				<Link href="/api/auth/signout">
					<DropdownMenuItem>
						
							Sign out
							
						<DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
					</DropdownMenuItem>
				</Link>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
