import { safeFetch, type AsyncResult } from '$lib/safe';

export type DiscoveredFormField = {
	name: string;
	type?: string;
	required?: boolean;
	options?: string[];
};

export type DiscoveredForm = {
	id: string;
	pagePath?: string;
	path?: string;
	pageUrl?: string;
	label?: string;
	displayName?: string;
	name?: string;
	method?: string;
	action?: string;
	originalAction?: string;
	fields?: (DiscoveredFormField | string)[];
	fieldSummary?: string[] | string;
	supported?: boolean;
	unsupportedReason?: string;
	reason?: string;
	active?: boolean;
	historical?: boolean;
};

export type FormDiscovery = {
	discoveryId?: string;
	totalForms?: number;
	supportedForms?: number;
	unsupportedForms?: number;
	forms: DiscoveredForm[];
};

export type EnableFormCollectionRequest = {
	siteId: string;
	discoveryId: string;
	formIds: string[];
};

export type SubmissionFields = Record<string, unknown> | { name: string; value: unknown }[];

export type FormSubmission = {
	id: string;
	siteId?: string;
	siteName?: string;
	siteDomain?: string;
	site?: {
		id?: string;
		name?: string;
		domain?: string;
	};
	formId?: string;
	formLabel?: string;
	formName?: string;
	form?: {
		id?: string;
		label?: string;
		name?: string;
		pagePath?: string;
		active?: boolean;
		historical?: boolean;
	};
	pagePath?: string;
	page?: string;
	fields: SubmissionFields;
	createdAt: string;
	active?: boolean;
	historical?: boolean;
};

export type FormSubmissionsResponse = {
	submissions: FormSubmission[];
};

type BackendError = {
	errorCode?: string;
	errorMessage?: string;
};

const responseError = async (resp: Response): Promise<string> => {
	const fallback = 'Ошибка при обращении к серверу';
	const contentType = resp.headers.get('content-type') ?? '';

	if (contentType.includes('application/json')) {
		const errorData = (await resp.json()) as BackendError;
		return errorData.errorMessage || fallback;
	}

	return (await resp.text()) || fallback;
};

export const enableFormCollection = async ({
	siteId,
	discoveryId,
	formIds
}: EnableFormCollectionRequest): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'X-Method': 'enableFormCollection'
		},
		body: JSON.stringify({ siteId, discoveryId, formIds })
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await responseError(resp) };
	}

	return { data: true, error: null };
};

export const getFormSubmissions = async (
	customFetch?: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>,
	request?: Request
): AsyncResult<FormSubmissionsResponse, string> => {
	const headers = request ? new Headers(request.headers) : new Headers();
	headers.set('X-Method', 'getFormSubmissions');

	const { data: resp, error } = await safeFetch(
		'/api/proxy',
		{
			method: 'GET',
			headers
		},
		customFetch
	);

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await responseError(resp) };
	}

	const payload = (await resp.json()) as FormSubmissionsResponse | FormSubmission[];
	return {
		data: Array.isArray(payload) ? { submissions: payload } : payload,
		error: null
	};
};

export const deleteFormSubmission = async (submissionId: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch(
		`/api/proxy?id=${encodeURIComponent(submissionId)}`,
		{
			method: 'DELETE',
			headers: {
				'X-Method': 'deleteFormSubmission'
			}
		}
	);

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await responseError(resp) };
	}

	return { data: true, error: null };
};
