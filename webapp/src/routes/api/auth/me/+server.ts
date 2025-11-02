import { Config } from '$lib/config';

export const GET = async ({ request }) => {
	try {
		const response = await fetch(Config.auth.meUrl, {
			method: 'GET',
			headers: {
				Cookie: request.headers.get('Cookie') || ''
			}
		});

		return new Response(response.body, {
			status: response.status,
			headers: response.headers
		});
	} catch (error) {
		console.error('Proxy error:', error);
		return new Response('Proxy error', { status: 500 });
	}
};
