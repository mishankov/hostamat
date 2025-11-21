type ErrorResp = {
	errorText: string;
};

export type Resp<T> = T | ErrorResp;
