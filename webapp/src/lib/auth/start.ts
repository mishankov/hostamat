import { goto } from '$app/navigation';

export const startYandexAuth = async (): Promise<void> => {
	try {
		const response = await fetch('/api/auth/yandex/start');
		if (!response.ok) {
			await goto('/login?error=oauth_failed');
			return;
		}

		const payload = (await response.json()) as { authUrl?: string };
		if (!payload.authUrl) {
			await goto('/login?error=oauth_failed');
			return;
		}

		window.location.href = payload.authUrl;
	} catch {
		await goto('/login?error=oauth_failed');
	}
};
