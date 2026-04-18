import type { PageServerLoad } from './$types';
import { getDomains, getSites } from '$lib/api-client';

export const load = (async ({ fetch, request }) => {
	return {
		sitesResp: getSites(fetch, request),
		domainsResp: getDomains(fetch, request)
	};
}) satisfies PageServerLoad;
