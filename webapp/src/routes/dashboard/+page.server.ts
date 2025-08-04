import type { PageServerLoad } from './$types';

interface Site {
  id: string;
  subdomain: string;
  domain: string;
  createdAt: string;
  updatedAt: string;
}

export const load = (async ({ fetch }) => {
  const response = await fetch("/api/listSites");
  const sites: Promise<Site[]> = response.json();
  return { sites };
}) satisfies PageServerLoad;
