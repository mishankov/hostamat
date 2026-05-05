import { safeFetch, type AsyncResult } from '$lib/safe';
import type { DiscoveredForm, FormDiscovery } from './form-submissions';

export const createUploadSession = async (siteId: string): AsyncResult<string, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createUploadSession',
			'H-Site-Id': siteId
		}
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		const errorData: { errorCode: string; errorMessage: string } = await resp.json();
		return { data: null, error: errorData.errorMessage };
	}

	return { data: await resp.text(), error: null };
};

export type UploadedFile = {
	id: string;
	name: string;
	uploaded: boolean;
	size: number;
};

export type UploadedPage = {
	fileName: string;
	path: string;
	url: string;
};

export type UploadSummary = {
	uploadedFiles: UploadedFile[];
	skippedFiles: UploadedFile[];
	availablePages: UploadedPage[];
	hasHomePage: boolean;
	formDiscovery?: FormDiscovery | null;
	forms?: DiscoveredForm[];
};

export const closeUploadSession = async (
	uploadSessionId: string
): AsyncResult<UploadSummary, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'DELETE',
		headers: {
			'X-Method': 'closeUploadSession',
			'H-UploadSession-Id': uploadSessionId
		}
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		const errorData: { errorCode: string; errorMessage: string } = await resp.json();
		return { data: null, error: errorData.errorMessage };
	}

	return { data: (await resp.json()) as UploadSummary, error: null };
};
