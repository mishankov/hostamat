import type { PageServerLoad } from './$types';
import { getSites } from '$lib/api-client';

export const load = (async ({ fetch, request }) => {
	return { sitesResp: getSites(fetch, request) };
}) satisfies PageServerLoad;
