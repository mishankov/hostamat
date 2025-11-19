export const createUploadSession = async (siteId: string) => {
	const resp = await fetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createUploadSession',
			'H-Site-Id': siteId
		}
	});

	return resp.ok ? await resp.text() : '';
};
