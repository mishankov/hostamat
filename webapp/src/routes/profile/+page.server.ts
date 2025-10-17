import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load = (async ({ parent }) => {
  const { isAuthenticated } = await parent()
  
  if (!isAuthenticated) {
    return redirect(302, '/login');
  }
}) satisfies PageServerLoad;
