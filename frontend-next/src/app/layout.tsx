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
		<html lang="en" className="dark">
			<body
				className={`font-sans ${inter.variable} bg-gray-100 dark:bg-gray-900`}
			>
				<TRPCReactProvider cookies={cookies().toString()}>
					<Header />
					<Sidebar visible={true} />
					{children}
				</TRPCReactProvider>
			</body>
		</html>
	);
}
