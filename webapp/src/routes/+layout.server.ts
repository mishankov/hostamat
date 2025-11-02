import type { LayoutServerLoad } from './$types';

export const load = (async ({ request, fetch }) => {
	const meResponse = await fetch('/api/auth/me', {
		method: 'GET',
		headers: {
			Cookie: request.headers.get('Cookie') || ''
		}
	});

	if (meResponse.status !== 200) {
		return { isAuthenticated: false };
	}

	const data = await meResponse.json();

	return { isAuthenticated: true, username: data.username };
}) satisfies LayoutServerLoad;
