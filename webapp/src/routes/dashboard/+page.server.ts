import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { getSites } from '$lib/api-client';

export const load = (async ({ fetch, request, parent }) => {
	const { isAuthenticated } = await parent();

	if (!isAuthenticated) {
		return redirect(302, '/login');
	}

	return { sitesResp: getSites(fetch, request) };
}) satisfies PageServerLoad;
