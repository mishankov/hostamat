import { Config } from '$lib/config';

export const GET = async ({ request, url }) => {
	try {
		const backendURL = new URL(Config.auth.yandexCallbackUrl);
		backendURL.search = url.search;

		const response = await fetch(backendURL, {
			method: 'GET',
			headers: {
				Cookie: request.headers.get('Cookie') || '',
				'X-Forwarded-Proto': url.protocol.replace(':', '')
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
