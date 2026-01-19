import { type Site } from '$lib/types';
import { safeFetch, type AsyncResult } from '$lib/safe';

type SitesResponse = {
	sites: Site[];
	maxAmount: number;
};

export const getSites = async (
	fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>,
	request: Request
): AsyncResult<SitesResponse, string> => {
	const headers = request.headers;
	headers.append('X-Method', 'getSites');
	const { data: resp, error } = await safeFetch(
		'/api/proxy',
		{
			method: 'GET',
			headers: headers
		},
		fetch
	);

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		const errorData: { errorCode: string; errorMessage: string } = await resp.json();
		return { data: null, error: errorData.errorMessage };
	}

	return { data: await resp.json(), error: null };
};

export const createSite = async (name: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createSite'
		},
		body: JSON.stringify({ name: name })
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		const errorData: { errorCode: string; errorMessage: string } = await resp.json();
		return { data: null, error: errorData.errorMessage };
	}

	return { data: true, error: null };
};

export const updateSite = async (siteId: string, name: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'PUT',
		headers: {
			'X-Method': 'updateSite',
			'H-Site-Id': siteId
		},
		body: JSON.stringify({ name: name })
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		const errorData: { errorCode: string; errorMessage: string } = await resp.json();
		return { data: null, error: errorData.errorMessage };
	}

	return { data: true, error: null };
};

export const deleteSite = async (siteId: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'DELETE',
		headers: {
			'X-Method': 'deleteSite',
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

	return { data: true, error: null };
};
