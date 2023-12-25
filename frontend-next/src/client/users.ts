interface User {
	id: string;
	email: string;
}

export class UserClient {
	private url: string;

	constructor(url: string) {
		this.url = url;
	}

	async login(email: string, password: string): Promise<User> {
		try {
			let response = await fetch(`${this.url}/login`, {
				method: 'POST',
				body: JSON.stringify({ email, password }),
			});

			if (!response.ok) {
				throw new Error('Login failed');
			}
			return response.json();
		} catch (e) {
			console.log(`login rqeuest failed: ${e}`)
			throw new Error('Login failed');
		}
	}

	async create() {
		// ...
	}
}