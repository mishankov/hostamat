import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { type Site } from './types';

type SitesResponse = {
	sites: Site[];
	maxAmount: number;
};

const loadSites = async (
	fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>,
	request: Request
): Promise<SitesResponse> => {
	const resp = await fetch('/api/listSites', {
		method: 'GET',
		headers: {
			Cookie: request.headers.get('Cookie') || ''
		}
	});

	return resp.json();
};

export const load = (async ({ fetch, request, parent }) => {
	const { isAuthenticated } = await parent();

	if (!isAuthenticated) {
		return redirect(302, '/login');
	}

	return { sitesResp: loadSites(fetch, request) };
}) satisfies PageServerLoad;
