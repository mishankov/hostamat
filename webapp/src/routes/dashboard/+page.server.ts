import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { type Site } from './types';

export const load = (async ({ fetch, request }) => {
  const meResponse = await fetch("/api/auth/me", {
    method: 'GET',
    headers: {
      'Cookie': request.headers.get('Cookie') || ''
    }
  });

  if (meResponse.status !== 200) {
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
