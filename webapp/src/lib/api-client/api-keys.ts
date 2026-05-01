import { safeFetch, type AsyncResult } from '$lib/safe';
import type { ApiKeyRecord, ApiKeyScope } from '$lib/types';

type ApiKeysResponse = {
	apiKeys: ApiKeyRecord[];
	maxActiveAmount?: number;
};

type CreateApiKeyPayload = {
	name: string;
	scopes: ApiKeyScope[];
	expiresAt: string | null;
};

type CreateApiKeyResponse = {
	apiKey: ApiKeyRecord;
	secret: string;
};

const readErrorMessage = async (resp: Response): Promise<string> => {
	try {
		const errorData: {
			errorCode?: string;
			errorMessage?: string;
			error?: { code?: string; message?: string };
			message?: string;
		} = await resp.json();

		return errorData.errorMessage || errorData.error?.message || errorData.message || resp.statusText;
	} catch {
		return resp.statusText || 'Ошибка запроса';
	}
};

const dashboardHeaders = (method: string, request?: Request): Headers => {
	const headers = new Headers(request?.headers);
	headers.set('X-Method', method);
	headers.set('Content-Type', 'application/json');
	return headers;
};

const normalizeListResponse = (payload: unknown): ApiKeysResponse => {
	const data = payload as ApiKeysResponse & { keys?: ApiKeyRecord[] };
	return {
		apiKeys: data.apiKeys ?? data.keys ?? [],
		maxActiveAmount: data.maxActiveAmount
	};
};

const normalizeCreateResponse = (payload: unknown): CreateApiKeyResponse => {
	const data = payload as CreateApiKeyResponse & {
		key?: ApiKeyRecord;
		apiKeySecret?: string;
		token?: string;
	};

	return {
		apiKey: data.apiKey ?? data.key!,
		secret: data.secret ?? data.apiKeySecret ?? data.token ?? ''
	};
};

export const getApiKeys = async (
	fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>,
	request: Request
): AsyncResult<ApiKeysResponse, string> => {
	const { data: resp, error } = await safeFetch(
		'/api/proxy',
		{
			method: 'GET',
			headers: dashboardHeaders('getApiKeys', request)
		},
		fetch
	);

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: normalizeListResponse(await resp.json()), error: null };
};

export const createApiKey = async (
	payload: CreateApiKeyPayload
): AsyncResult<CreateApiKeyResponse, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: dashboardHeaders('createApiKey'),
		body: JSON.stringify(payload)
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: normalizeCreateResponse(await resp.json()), error: null };
};

export const updateApiKeyName = async (
	apiKeyId: string,
	name: string
): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'PATCH',
		headers: new Headers([
			['X-Method', 'updateApiKey'],
			['H-API-Key-Id', apiKeyId],
			['Content-Type', 'application/json']
		]),
		body: JSON.stringify({ name })
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: true, error: null };
};

export const revokeApiKey = async (apiKeyId: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'DELETE',
		headers: new Headers([
			['X-Method', 'revokeApiKey'],
			['H-API-Key-Id', apiKeyId]
		])
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: true, error: null };
};
