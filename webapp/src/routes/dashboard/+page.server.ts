import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { type Site } from './types';

export const load = (async ({ fetch, request, parent }) => {
  const { isAuthenticated } = await parent()
  
  if (!isAuthenticated) {
    return redirect(302, '/login');
  }

  const response = await fetch("/api/listSites", {
    method: 'GET',
    headers: {
      'Cookie': request.headers.get('Cookie') || ''
    }
  });
  
  if (!response.ok) {
    return { sites: [] };
  }

  const sites: Promise<Site[]> = response.json();

  return { sites };
}) satisfies PageServerLoad;
