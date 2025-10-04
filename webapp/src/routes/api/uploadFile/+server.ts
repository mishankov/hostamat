import type { RequestHandler } from './$types';
import { Config } from '$lib/config';

export const POST: RequestHandler = async ({ request }) => {
  try {
    const headers = new Headers(request.headers);

    const respDeleteFiles = await fetch(Config.files.uploadUrl, {
      method: 'DELETE',
      headers,
    });

    if (!respDeleteFiles.ok) {
      return new Response('Delete error', { status: 400 });
    }
    
    // Forward the request to backend
    const response = await fetch(Config.files.uploadUrl, {
      method: 'POST',
      headers,
      body: request.body,
      // @ts-ignore
      duplex: 'half',
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
