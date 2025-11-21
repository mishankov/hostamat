import { type Site } from '$lib/types';
import { type Resp } from './models';

export const createSite = async (name: string): Promise<Resp<true>> => {
	const resp = await fetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createSite'
		},
		body: JSON.stringify({ name: name })
	});

	if (resp.ok) return true;

	return resp.status >= 500
		? { errorText: 'Ошибка сервера при создании сайта' }
		: { errorText: await resp.text() };
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
		method: 'GET',
		headers: headers
	});

	return resp.ok ? await resp.json() : null;
};

export const updateSite = async (siteId: string, name: string) => {
	const resp = await fetch('/api/proxy', {
		method: 'PUT',
		headers: {
			'X-Method': 'updateSite',
			'H-Site-Id': siteId
		},
		body: JSON.stringify({ name: name })
	});

	return resp.ok;
};

export const deleteSite = async (siteId: string) => {
	await fetch('/api/proxy', {
		method: 'DELETE',
		headers: {
			'X-Method': 'deleteSite',
			'H-Site-Id': siteId
		}
	});
};
