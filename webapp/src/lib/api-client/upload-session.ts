export const createUploadSession = async (siteId: string) => {
	const resp = await fetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createUploadSession',
			'H-Site-Id': siteId
		}
	});

	return resp.ok ? await resp.text() : null;
};

export const closeUploadSession = async (uploadSessionId: string) => {
	const resp = await fetch('/api/proxy', {
		method: 'DELETE',
		headers: {
			'X-Method': 'closeUploadSession',
			'H-UploadSession-Id': uploadSessionId
		}
	});

	return resp.ok;
};
