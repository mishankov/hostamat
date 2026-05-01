import openapi from '$lib/api-docs/openapi.yaml?raw';

export const GET = () => {
	return new Response(openapi, {
		headers: {
			'Content-Type': 'application/yaml; charset=utf-8',
			'Cache-Control': 'public, max-age=300'
		}
	});
};
