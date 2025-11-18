export const createSite = async (name: string) => {
	const resp = await fetch('/api/createSite', {
		method: 'POST',
		body: JSON.stringify({ name: name })
	});

	return resp.ok;
};
