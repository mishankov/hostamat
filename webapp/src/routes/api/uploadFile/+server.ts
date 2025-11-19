import type { RequestHandler } from './$types';
import { Config } from '$lib/config';

export const POST: RequestHandler = async ({ request }) => {
	try {
		const headers = new Headers(request.headers);

		// Forward the request to backend
		const response = await fetch(Config.files.uploadUrl, {
			method: 'POST',
			headers,
			body: request.body,
			// @ts-expect-error duplex needed to send files
			duplex: 'half'
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
