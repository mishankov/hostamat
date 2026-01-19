import { safeFetch } from '$lib/safe';

export const uploadFile = async (url: string, headers: Headers, data: FormData) => {
	const { data: resp, error } = await safeFetch(url, {
		method: 'POST',
		headers,
		body: data
	});

	if (error) {
		return {
			error: `Ошибка при загрузке файла`,
			errorMessage: error.message
		};
	}

	if (!resp.ok) {
		const error: { errorCode: string; errorMessage: string } = await resp.json();

		switch (error.errorCode) {
			case 'OTHER':
				return {
					error: 'Ошибка при загрузке файла',
					errorMessage: 'Неизвестная ошибка при загрузке файла'
				};
			default:
				return {
					error: 'Ошибка при загрузке файла',
					errorMessage: error.errorMessage
				};
		}
	}

	return true;
};
