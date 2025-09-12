import type { PageServerLoad } from './$types';
import { type Site } from './types';

export const load = (async ({ fetch }) => {
  const response = await fetch("/api/listSites");
  if (!response.ok) {
    return { sites: [] };
  }

  const sites: Promise<Site[]> = response.json();

  return { sites };
}) satisfies PageServerLoad;
