import type { RequestHandler } from './$types';
import { Config } from '$lib/config';

export const GET: RequestHandler = async ({ request }) => {
  try {
    const response = await fetch(Config.sites.listUrl, {
      method: 'GET',
      headers: {
        'Cookie': request.headers.get('Cookie') || ''
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
