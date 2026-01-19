import { tryCatch } from '$lib/safe/try-catch';

// export const safeFetch = async (input: string | URL | Request, init?: RequestInit) => {
// 	return await tryCatch(fetch(input, init));
// };

export const safeFetch = async (
	input: string | URL | Request,
	init?: RequestInit,
	customFetch?: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>
) => {
	const f = customFetch || fetch;
	return await tryCatch(f(input, init));
};

export type { Result, AsyncResult } from '$lib/safe/result';
