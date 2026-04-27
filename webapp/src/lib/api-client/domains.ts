import { type DomainRecord } from '$lib/types';
import { safeFetch, type AsyncResult } from '$lib/safe';

type DomainsResponse = {
	domains: DomainRecord[];
};

const readErrorMessage = async (resp: Response): Promise<string> => {
	const errorData: { errorCode: string; errorMessage: string } = await resp.json();
	return errorData.errorMessage;
};

export const getDomains = async (
	fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>,
	request: Request
): AsyncResult<DomainsResponse, string> => {
	const headers = new Headers(request.headers);
	headers.set('X-Method', 'getDomains');

	const { data: resp, error } = await safeFetch(
		'/api/proxy',
		{
			method: 'GET',
			headers
		},
		fetch
	);

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: await resp.json(), error: null };
};

export const createDomain = async (
	siteId: string,
	hostname: string
): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createDomain'
		},
		body: JSON.stringify({ siteId, hostname })
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: true, error: null };
};

export const verifyDomain = async (domainId: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'verifyDomain',
			'H-Domain-Id': domainId
		}
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: true, error: null };
};

export const updateDomain = async (
	domainId: string,
	siteId: string
): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'PUT',
		headers: {
			'X-Method': 'updateDomain',
			'H-Domain-Id': domainId
		},
		body: JSON.stringify({ siteId })
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: true, error: null };
};

export const deleteDomain = async (domainId: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'DELETE',
		headers: {
			'X-Method': 'deleteDomain',
			'H-Domain-Id': domainId
		}
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: true, error: null };
};
