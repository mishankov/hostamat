import { env } from '$env/dynamic/private';

export const Config = {
	files: {
		uploadUrl: `${env.API_BASE || 'http://localhost:1001'}/api/files`
	},
	sites: {
		deleteUrl: `${env.API_BASE || 'http://localhost:1001'}/api/sites`
	},
	uploadSessionsUrl: `${env.API_BASE || 'http://localhost:1001'}/api/uploadsessions`,
	auth: {
		logoutUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/logout`,
		meUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/me`,
		yandexStartUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/yandex/start`,
		yandexCallbackUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/yandex/callback`
	}
} as const;
