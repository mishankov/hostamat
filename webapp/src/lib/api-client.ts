export const createSite = async (name: string) => {
	const resp = await fetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createSite'
		},
		body: JSON.stringify({ name: name })
	});

	return resp.ok;
};
