import { type Site } from '$lib/types';

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

type SitesResponse = {
	sites: Site[];
	maxAmount: number;
};

export const getSites = async (
	fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>,
	request: Request
): Promise<SitesResponse | null> => {
	const headers = request.headers;
	headers.append('X-Method', 'getSites');
	const resp = await fetch('/api/proxy', {
		method: 'POST',
		headers: headers
	});

	return resp.ok ? await resp.json() : null;
};

export const deleteSite = async (siteId: string) => {
	const resp = await fetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'deleteSite',
			'H-Site-Id': siteId
		}
	});

	return resp.ok;
};
