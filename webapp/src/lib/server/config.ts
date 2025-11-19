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
	},
	createUploadSession: {
		method: 'POST',
		url: `${API_BASE}/api/uploadsessions`
	},
	closeUploadSession: {
		method: 'DELETE',
		url: `${API_BASE}/api/uploadsessions`
	}
};

export const getConfigForMethod = (method: string | null) => {
	if (method === null) {
		return null;
	}
	switch (method) {
		case 'createSite':
			return ProxyConfig.createSite;
		case 'getSites':
			return ProxyConfig.getSites;
		case 'deleteSite':
			return ProxyConfig.deleteSite;
		case 'createUploadSession':
			return ProxyConfig.createUploadSession;
		case 'closeUploadSession':
			return ProxyConfig.closeUploadSession;
		default:
			return null;
	}
};
