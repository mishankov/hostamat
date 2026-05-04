import { safeFetch, type AsyncResult } from '$lib/safe';

const readErrorMessage = async (resp: Response): Promise<string> => {
	const fallback = 'Не удалось отправить сообщение';
	const contentType = resp.headers.get('content-type') || '';

	if (contentType.includes('application/json')) {
		try {
			const errorData: { errorMessage?: string } = await resp.json();
			return errorData.errorMessage || fallback;
		} catch {
			return fallback;
		}
	}

	try {
		const text = await resp.text();
		return text || fallback;
	} catch {
		return fallback;
	}
};

export const createFeedback = async (
	message: string,
	url: string
): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: {
			'X-Method': 'createFeedback'
		},
		body: JSON.stringify({ message, url })
	});

	if (error) {
		return { data: null, error: error.message };
	}

	if (!resp.ok) {
		return { data: null, error: await readErrorMessage(resp) };
	}

	return { data: true, error: null };
};
