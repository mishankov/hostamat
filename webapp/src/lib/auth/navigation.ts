import { goto } from '$app/navigation';

export const navigateAfterAuthChange = async (path: string): Promise<void> => {
	await goto(path, { invalidateAll: true });
};
