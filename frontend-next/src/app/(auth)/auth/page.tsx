"use client";

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@components/ui/form";
import { Input } from "@components/ui/input";
import { Button } from "@components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@components/ui/tabs";
import { z } from "zod";
import { Toggle } from "@components/ui/toggle";
import { useRef } from "react";
import { InputPassword } from "~/app/_components/ui/input-password";

const formSchemaLogin = z.object({
	email: z.string().email().toLowerCase(),
	password: z.string().min(8, "Password must be at least 8 characters"),
});
const formSchemaRegister = z.object({
	email: z.string().email().toLowerCase(),
	password: z.string().min(8, "Password must be at least 8 characters")
		.regex(/[a-z]/, "Password must contain at least one lowercase letter")
		.regex(/[A-Z]/, "Password must contain at least one uppercase letter")
		.regex(/[0-9]/, "Password must contain at least one number")
		.regex(/[^a-zA-Z0-9]/, "Password must contain at least one special character"),
	passwordRepeat: z.string().min(8, "Password must be at least 8 characters"),
}).superRefine(({ password, passwordRepeat }, ctx) => {
	if (password !== passwordRepeat) {
		ctx.addIssue({
			code: z.ZodIssueCode.custom,
			message: "Passwords must match",
			path: ["passwordRepeat"],
		});
	}
});

export default function SignIn() {
	const formRegister = useForm<z.infer<typeof formSchemaRegister>>({
		resolver: zodResolver(formSchemaRegister),
		defaultValues: {
			email: "",
			password: "",
			passwordRepeat: "",
		},
	});
	const formLogin = useForm<z.infer<typeof formSchemaLogin>>({
		resolver: zodResolver(formSchemaLogin),
		defaultValues: {
			email: "",
			password: "",
		},
	});

	const passwordLogin = useRef(null);
	const toogleVisibility = () => {
		if (passwordLogin.current) {
			(passwordLogin.current as HTMLInputElement).type === "password"
				? (passwordLogin.current as HTMLInputElement).type = "text"
				: (passwordLogin.current as HTMLInputElement).type = "password";
		}
	}

	return (
		<div className="flex items-center align-middle p-4 rounded-md m-auto w-[400px]">
			<Tabs defaultValue="login" className="w-full">
				<TabsList className="w-full">
					<TabsTrigger value="login" className="uppercase font-bold w-full">Login</TabsTrigger>
					<TabsTrigger value="register" className="uppercase font-bold w-full">Register</TabsTrigger>
				</TabsList>
				<TabsContent value="login">
					<Form {...formLogin}>
						<form className="w-full" onSubmit={formLogin.handleSubmit((data) => console.log(data))}>
							<FormField control={formLogin.control}
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
							<FormField control={formLogin.control}
								name="password"
								render={({ field }) => (
									<FormItem className="pt-2">
										<FormLabel>Password</FormLabel>
										<FormControl>
											<Input placeholder="Password" type="password" {...field} ref={passwordLogin}/>		
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
					<Form {...formRegister}>
						<form onSubmit={formRegister.handleSubmit((data) => console.log(data))}>
							<FormField control={formRegister.control}
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
							<FormField control={formRegister.control}
								name="password"
								render={({ field }) => (
									<FormItem className="pt-2">
										<FormLabel>Password</FormLabel>
										<FormControl>
											<InputPassword placeholder="Password" type="password" {...field} />
										</FormControl>
										<FormMessage />
									</FormItem>
								)
								} />
							<FormField control={formRegister.control}
								name="passwordRepeat"
								render={({ field }) => (
									<FormItem className="pt-2">
										<FormLabel>Repeat Password</FormLabel>
										<FormControl>
											<InputPassword placeholder="Reepat Password" type="password" {...field} ref={passwordLogin}/>		
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