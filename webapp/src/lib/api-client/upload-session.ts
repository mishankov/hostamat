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

type File = {
	id: string;
	name: string;
	uploaded: boolean;
	size: number;
};

type FileRespOk = {
	ok: true;
	files: File[];
};

type FileRespError = {
	ok: false;
	errorMessage: string;
};

type FileResp = FileRespOk | FileRespError;

export const closeUploadSession = async (uploadSessionId: string): Promise<FileResp> => {
	const resp = await fetch('/api/proxy', {
		method: 'DELETE',
		headers: {
			'X-Method': 'closeUploadSession',
			'H-UploadSession-Id': uploadSessionId
		}
	});

	if (resp.ok) {
		return {
			ok: true,
			files: (await resp.json()) as File[]
		};
	}

	return {
		ok: false,
		errorMessage: await resp.text()
	};
};
