import type { RequestHandler } from './$types';
import { Config } from '$lib/config';

export const POST: RequestHandler = async ({ request }) => {
  try {
    const backendUrl = Config.files.uploadUrl;
    const response = await fetch(backendUrl, {
      method: 'POST',
      headers: request.headers,
      body: request.body,
      // @ts-ignore
      duplex: 'half', 
    });

    return new Response(response.body, { status: response.status, headers: response.headers });
  } catch (error) {
    console.error('Proxy error:', error);
    return new Response('Proxy error', { status: 500 });
  }
};