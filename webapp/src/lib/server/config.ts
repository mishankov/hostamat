import { env } from '$env/dynamic/private';

export const ProxyConfig = {
	createSite: {
		url: `${env.API_BASE || 'http://localhost:1001'}/api/sites`,
		method: 'POST'
	}
} as const;
