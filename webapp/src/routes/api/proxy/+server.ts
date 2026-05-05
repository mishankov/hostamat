import type { RequestHandler } from './$types';
import { getConfigForMethod } from '$lib/server/config';

const proxyHandler = (): RequestHandler => {
	return async ({ request, fetch, url }) => {
		try {
			const method = request.headers.get('X-Method');
			const config = getConfigForMethod(method);

			if (config === null) {
				return new Response('Method not allowed', { status: 405 });
			}

			const targetUrl = `${config.url}${url.search}`;
			const response = await fetch(targetUrl, {
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
};

export const GET: RequestHandler = proxyHandler();
export const PUT: RequestHandler = proxyHandler();
export const POST: RequestHandler = proxyHandler();
export const DELETE: RequestHandler = proxyHandler();
