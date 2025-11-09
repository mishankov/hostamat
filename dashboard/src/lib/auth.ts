import { env } from '$env/dynamic/private';

export const DASHBOARD_AUTH_TOKEN = env.DASHBOARD_AUTH_TOKEN || 'default-token';

export function validateAuthToken(token: string): boolean {
	return token === DASHBOARD_AUTH_TOKEN;
}
