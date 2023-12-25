import "~/styles/globals.css";

import { Inter } from "next/font/google";
import { cookies } from "next/headers";

import { TRPCReactProvider } from "~/trpc/react";
import { Header } from "@components/header";
import { Sidebar } from "@components/sidebar";

const inter = Inter({
	subsets: ["latin"],
	variable: "--font-sans",
});

export const metadata = {
	title: "Hevonen",
	description: "Manage your horse riding club",
	icons: [{ rel: "icon", url: "/favicon.ico" }],
};

export default function RootLayout({
	children,
}: {
	children: React.ReactNode;
}) {
	return (
		<html lang="en" className="dark h-full w-full">
			<body
				className={`h-full w-full flex m-0 overflow-hidden font-sans ${inter.variable} bg-gray-100 dark:bg-gray-900`}
			>
				<TRPCReactProvider cookies={cookies().toString()}>
					<div id="left-sidebar" className="h-100 flex flex-col w-48">
						<div id="brand" className="h-12 pl-2 flex items-center border-b-[1px]">
							<h1 className="uppercase font-extrabold">Hevonen</h1>
						</div>
						<div className="flex-grow-1 overflow-y-auto">
							<Sidebar visible={true} />
						</div>
					</div>
					<div className="h-100 flex flex-col flex-auto">
						<div className="flex-shrink-0">
							<Header />
						</div>
						<div className="flex-grow-1 overflow-y-auto">
							{children}
						</div>
					</div>
				</TRPCReactProvider>
			</body>
		</html>
	);
}
