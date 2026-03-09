import type { PageServerLoad } from './$types';

export const load = (async ({ url }) => {
	return {
		authError: url.searchParams.get('authError') || ''
	};
}) satisfies PageServerLoad;
