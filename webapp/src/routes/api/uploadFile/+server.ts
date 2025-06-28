import type { RequestHandler } from './$types';
import { env } from '$env/dynamic/private';

export const POST: RequestHandler = async ({ request }) => {
  try {
    const backendUrl = `${env.API_BASE}/files/api/upload`;
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