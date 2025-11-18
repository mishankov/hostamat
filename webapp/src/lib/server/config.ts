import { env } from '$env/dynamic/private';

const API_BASE = env.API_BASE || 'http://localhost:1001';

export const ProxyConfig = {
	createSite: {
		method: 'POST',
		url: `${API_BASE}/api/sites`
	},
	getSites: {
		method: 'GET',
		url: `${API_BASE}/api/sites`
	},
	deleteSite: {
		method: 'DELETE',
		url: `${API_BASE}/api/sites`
	}
} as const;
