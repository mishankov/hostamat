import { env } from '$env/dynamic/private';

export const ProxyConfig = {
	createSite: {
		method: 'POST',
		url: `${env.API_BASE || 'http://localhost:1001'}/api/sites`
	},
	getSites: {
		method: 'GET',
		url: `${env.API_BASE || 'http://localhost:1001'}/api/sites`
	}
} as const;
