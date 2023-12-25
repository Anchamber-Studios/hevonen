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
import { Button } from "~/app/_components/ui/button";
import { z } from "zod";

const formSchema = z.object({
	email: z.string().email(),
	password: z.string().min(8, "Password must be at least 8 characters"),
});

export default function SignIn() {
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			email: "",
			password: "",
		},
	});

	return (
		<div className="flex items-center align-middle bg-gray-300 p-4 rounded-md m-auto">
			<Form {...form}>
				<form onSubmit={form.handleSubmit((data) => console.log(data))}>
					<FormField control={form.control} 
						name="email"
						render={({field}) => (
							<FormItem>
								<FormLabel>Email</FormLabel>
								<FormControl>
									<Input placeholder="Email" type="email" {...field} />
								</FormControl>
								{/* <FormDescription>
									This email is used for login. 
									A confirmation email will be sent to this address.
								</FormDescription> */}
								<FormMessage />
							</FormItem>
						)
					}/>
					<FormField control={form.control} 
						name="password"
						render={({field}) => (
							<FormItem>
								<FormLabel>Password</FormLabel>
								<FormControl>
									<Input placeholder="Password" type="password" {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)
					}/>
					<Button type="submit">Sign In</Button>
				</form>
			</Form>
		</div>
	)
}