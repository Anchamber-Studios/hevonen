"use client";

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@components/ui/form";
import { Input } from "@components/ui/input";
import { Button } from "@components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@components/ui/tabs";
import { z } from "zod";

const formSchema = z.object({
	email: z.string().email(),
	password: z.string().min(8, "Password must be at least 8 characters"),
	passwordRepeat: z.string().min(8, "Password must be at least 8 characters"),
});

export default function SignIn() {
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			email: "",
			password: "",
			passwordRepeat: "",
		},
	});

	return (
		<div className="flex items-center align-middle p-4 rounded-md m-auto w-[400px]">
			<Tabs defaultValue="login" className="w-full">
				<TabsList className="w-full">
					<TabsTrigger value="login" className="uppercase font-bold w-full">Login</TabsTrigger>
					<TabsTrigger value="register" className="uppercase font-bold w-full">Register</TabsTrigger>
				</TabsList>
				<TabsContent value="login">
					<Form {...form}>
						<form className="w-full" onSubmit={form.handleSubmit((data) => console.log(data))}>
							<FormField control={form.control}
								name="email"
								render={({ field }) => (
									<FormItem className="pt-2">
										<FormLabel>Email</FormLabel>
										<FormControl>
											<Input placeholder="Email" type="email" {...field} />
										</FormControl>
										<FormMessage />
									</FormItem>
								)
								} />
							<FormField control={form.control}
								name="password"
								render={({ field }) => (
									<FormItem className="pt-2">
										<FormLabel>Password</FormLabel>
										<FormControl>
											<Input placeholder="Password" type="password" {...field} />
										</FormControl>
										<FormMessage />
									</FormItem>
								)
								} />
							<Button type="submit" className="w-full mt-4">Sign In</Button>
						</form>
					</Form>
				</TabsContent>
				<TabsContent value="register">
					<Form {...form}>
						<form onSubmit={form.handleSubmit((data) => console.log(data))}>
							<FormField control={form.control}
								name="email"
								render={({ field }) => (
									<FormItem className="pt-2">
										<FormLabel>Email</FormLabel>
										<FormControl>
											<Input placeholder="Email" type="email" {...field} />
										</FormControl>
										<FormMessage />
									</FormItem>
								)
								} />
							<FormField control={form.control}
								name="password"
								render={({ field }) => (
									<FormItem className="pt-2">
										<FormLabel>Password</FormLabel>
										<FormControl>
											<Input placeholder="Password" type="password" {...field} />
										</FormControl>
										<FormMessage />
									</FormItem>
								)
								} />
								<FormField control={form.control}
									name="passwordRepeat"
									render={({ field }) => (
										<FormItem className="pt-2">
											<FormLabel>Repeat Password</FormLabel>
											<FormControl>
												<Input placeholder="Reepat Password" type="password" {...field} />
											</FormControl>
											<FormMessage />
										</FormItem>
									)
									} />
							<Button type="submit" className="w-full mt-4">Create Account</Button>
						</form>
					</Form>
				</TabsContent>
			</Tabs>

		</div>
	)
}