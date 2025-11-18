import type { RequestHandler } from './$types';
import { Config } from '$lib/config';

export const POST: RequestHandler = async ({ request }) => {
	try {
		const response = await fetch(Config.sites.createUrl, {
			method: 'POST',
			headers: {
				Cookie: request.headers.get('Cookie') || '',
				'Content-Type': 'application/json'
			},
			body: request.body
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
