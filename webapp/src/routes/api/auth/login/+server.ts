import { Config } from "$lib/config";

export const POST = async ({ request }) => {
  try {
    const response = await fetch(Config.auth.loginUrl, {
      method: 'POST',
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
}