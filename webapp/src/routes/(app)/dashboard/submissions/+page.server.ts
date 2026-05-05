import type { PageServerLoad } from './$types';
import { getFormSubmissions } from '$lib/api-client';

export const load = (async ({ fetch, request }) => {
	return { submissionsResp: getFormSubmissions(fetch, request) };
}) satisfies PageServerLoad;
