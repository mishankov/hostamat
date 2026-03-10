import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load = (async ({ parent, url }) => {
	const { isAuthenticated } = await parent();

	if (isAuthenticated) {
		return redirect(302, '/dashboard');
	}

	return {
		error: url.searchParams.get('error') || ''
	};
}) satisfies PageServerLoad;
