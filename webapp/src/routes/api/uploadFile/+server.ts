import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request }) => {
  try {
    // TODO: use env variable
    const backendUrl = `http://localhost:1000/files/api/upload`;
    const response = await fetch(backendUrl, {
      method: 'POST',
      headers: request.headers,
      body: request.body,
      duplex: 'half',
    });

    return new Response(response.body, { status: response.status, headers: response.headers });
  } catch (error) {
    console.error('Proxy error:', error);
    return new Response('Proxy error', { status: 500 });
  }
};