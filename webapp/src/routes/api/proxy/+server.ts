import type { RequestHandler } from './$types';
import { ProxyConfig } from '$lib/server/config';

export const POST: RequestHandler = async ({ request, fetch }) => {
	try {
		const method = request.headers.get('X-Method');
		let config;
		switch (method) {
			case 'createSite':
				config = ProxyConfig.createSite;
				break;
			case 'getSites':
				config = ProxyConfig.getSites;
				break;
			default:
				return new Response('Method not allowed', { status: 405 });
		}

		const response = await fetch(config.url, {
			method: config.method,
			headers: request.headers,
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
