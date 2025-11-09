import postgres from 'postgres';
import { env } from '$env/dynamic/private';

export const sql = postgres(env.DATABASE_URL || '');

export interface DashboardStats {
	userCount: number;
	siteCount: number;
	uploadSessionCount: number;
}

export async function getDashboardStats(): Promise<DashboardStats> {
	// Count distinct users from sites and upload sessions
	const userCountResult = await sql`
		SELECT COUNT(DISTINCT "user") as count FROM (
			SELECT "user" FROM sites
			UNION
			SELECT "user" FROM uploadsessions
		) as combined_users
	`;

	const siteCountResult = await sql`SELECT COUNT(*) as count FROM sites`;
	const uploadSessionCountResult = await sql`SELECT COUNT(*) as count FROM uploadsessions`;

	return {
		userCount: Number(userCountResult[0]?.count || 0),
		siteCount: Number(siteCountResult[0]?.count || 0),
		uploadSessionCount: Number(uploadSessionCountResult[0]?.count || 0)
	};
}
