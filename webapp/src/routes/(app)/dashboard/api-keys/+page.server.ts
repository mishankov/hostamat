import type { PageServerLoad } from './$types';
import { getApiKeys } from '$lib/api-client';

export const load = (async ({ fetch, request }) => {
	return { apiKeysResp: getApiKeys(fetch, request) };
}) satisfies PageServerLoad;
