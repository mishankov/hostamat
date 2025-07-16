import type { RequestHandler } from './$types';
import { Config } from '$lib/config';

export const POST: RequestHandler = async () => {
  try {
    const response = await fetch(Config.sites.createUrl, {
      method: 'POST'
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
