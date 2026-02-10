export const GET = ({ request }) => {
	console.log(request);
	return new Response('ok');
};
