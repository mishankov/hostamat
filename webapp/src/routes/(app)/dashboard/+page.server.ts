import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load = (async () => {
	return redirect(302, '/dashboard/sites');
}) satisfies PageServerLoad;
