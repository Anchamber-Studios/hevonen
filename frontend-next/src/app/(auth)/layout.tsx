import "~/styles/globals.css";

import { Inter } from "next/font/google";

const inter = Inter({
	subsets: ["latin"],
	variable: "--font-sans",
});

export const metadata = {
	title: 'Hevonen Authentication',
	description: '',
}

export default function RootLayout({
	children,
}: {
	children: React.ReactNode
}) {
	return (
		<html lang="en" className="dark h-full w-full">
			<body
				className={`h-full w-full flex items-center m-0 overflow-hidden font-sans ${inter.variable} bg-gray-100 dark:bg-gray-900`}>
				{children}
			</body>
		</html>
	)
}
