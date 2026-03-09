import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load = (async ({ url }) => {
	const authError = url.searchParams.get('error');
	if (authError) {
		return redirect(302, `/?authError=${encodeURIComponent(authError)}`);
	}

	return redirect(302, '/');
}) satisfies PageServerLoad;
