import { UserClient } from "~/client/users";
import { env } from "~/env";

export async function POST(req: Request) {
	const formData = await req.formData();
	const email = formData.get('email') as string;
	const password = formData.get('password') as string;
	if (!email || !password) {
		return Response.error();
	}
	console.log(`email: ${email}, password: ${password}`);

	try {
		const client = new UserClient(env.USER_SERVICE_URL);
		const user = await client.register(email, password);
		return Response.json(user);
	} catch (error) {
		console.error(error);
		return Response.error();
	}
}