import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load = (async ({ parent }) => {
	const { isAuthenticated } = await parent();

	if (!isAuthenticated) {
		return redirect(302, '/login');
	}

	return {};
}) satisfies LayoutServerLoad;
