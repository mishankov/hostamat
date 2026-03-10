import { Config } from '$lib/config';

export const GET = async ({ request, url }) => {
	try {
		const response = await fetch(Config.auth.yandexStartUrl, {
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
