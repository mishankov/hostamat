import { goto } from '$app/navigation';

export const startYandexAuth = async (): Promise<void> => {
	try {
		const response = await fetch('/api/auth/yandex/start');
		if (!response.ok) {
			await goto('/?authError=oauth_failed');
			return;
		}

		const payload = (await response.json()) as { authUrl?: string };
		if (!payload.authUrl) {
			await goto('/?authError=oauth_failed');
			return;
		}

		window.location.href = payload.authUrl;
	} catch {
		await goto('/?authError=oauth_failed');
	}
};
