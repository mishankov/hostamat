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
	getDomains: {
		method: 'GET',
		url: `${API_BASE}/api/domains`
	},
	updateSite: {
		method: 'PUT',
		url: `${API_BASE}/api/sites`
	},
	deleteSite: {
		method: 'DELETE',
		url: `${API_BASE}/api/sites`
	},
	createDomain: {
		method: 'POST',
		url: `${API_BASE}/api/domains`
	},
	updateDomain: {
		method: 'PUT',
		url: `${API_BASE}/api/domains`
	},
	verifyDomain: {
		method: 'POST',
		url: `${API_BASE}/api/domains/verify`
	},
	deleteDomain: {
		method: 'DELETE',
		url: `${API_BASE}/api/domains`
	},
	createFeedback: {
		method: 'POST',
		url: `${API_BASE}/api/feedback`
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

	return ProxyConfig[method as keyof typeof ProxyConfig] || null;
};
