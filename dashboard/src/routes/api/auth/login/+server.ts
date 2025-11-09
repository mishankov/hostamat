import { json } from '@sveltejs/kit';
import { validateAuthToken, DASHBOARD_AUTH_TOKEN } from '$lib/auth';

export async function POST({ request, cookies }) {
	try {
		const { token } = await request.json();

		console.log('Login attempt with token:', token ? '***' : 'empty');
		console.log('Expected token:', DASHBOARD_AUTH_TOKEN ? '***' : 'empty');

		if (validateAuthToken(token)) {
			cookies.set('dashboard-auth', token, {
				path: '/',
				httpOnly: true,
				secure: process.env.NODE_ENV === 'production',
				sameSite: 'strict',
				maxAge: 60 * 60 * 24 * 7 // 1 week
			});

			console.log('Login successful');
			return json({ success: true });
		}

		console.log('Login failed - invalid token');
		return json({ success: false }, { status: 401 });
	} catch (error) {
		console.error('Login endpoint error:', error);
		return json({ success: false, error: 'Server error' }, { status: 500 });
	}
}
