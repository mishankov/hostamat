import { getDashboardStats } from '$lib/db';

export async function load() {
	const stats = await getDashboardStats();
	return { stats };
}
