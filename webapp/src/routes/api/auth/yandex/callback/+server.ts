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

		if (!response.ok) {
			return Response.redirect(new URL('/?authError=oauth_failed', url), 302);
		}

		const payload = (await response.json()) as { redirectTo?: string };
		const location = payload.redirectTo || '/?authError=oauth_failed';
		const headers = new Headers();

		for (const [key, value] of response.headers.entries()) {
			if (key.toLowerCase() === 'set-cookie') {
				headers.append(key, value);
			}
		}

		headers.set('Location', new URL(location, url).toString());

		return new Response(null, {
			status: 302,
			headers
		});
	} catch (error) {
		console.error('Proxy error:', error);
		return Response.redirect(new URL('/?authError=oauth_failed', url), 302);
	}
};
