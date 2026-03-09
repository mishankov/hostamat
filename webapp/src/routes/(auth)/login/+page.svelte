<script lang="ts">
	import { Alert, Button } from '$lib/components';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let isLoading = $state(false);
	let startErrorMessage = $state('');

	let errorMessage = $derived.by(() => {
		if (startErrorMessage) {
			return startErrorMessage;
		}

		switch (data.error) {
			case 'oauth_denied':
				return 'Вход через Яндекс был отменён.';
			case 'oauth_failed':
				return 'Не удалось завершить вход через Яндекс.';
			default:
				return '';
		}
	});

	const startYandexAuth = async () => {
		isLoading = true;
		startErrorMessage = '';

		try {
			const response = await fetch('/api/auth/yandex/start');
			if (!response.ok) {
				startErrorMessage = 'Не удалось начать вход через Яндекс.';
				return;
			}

			const payload = (await response.json()) as { authUrl?: string };
			if (!payload.authUrl) {
				startErrorMessage = 'Не удалось начать вход через Яндекс.';
				return;
			}

			window.location.href = payload.authUrl;
		} catch {
			startErrorMessage = 'Не удалось начать вход через Яндекс.';
		} finally {
			isLoading = false;
		}
	};
</script>

<svelte:head>
	<title>Вход | Хостамат</title>
</svelte:head>

<h1>Вход</h1>
<p class="lead">Войдите через Яндекс, чтобы вернуться к вашим сайтам и публикациям.</p>

{#if errorMessage}
	<Alert variant="error" message={errorMessage} />
{/if}

<div class="auth-form">
	<Button onclick={startYandexAuth} filled disabled={isLoading}>
		{isLoading ? 'Переходим в Яндекс...' : 'Продолжить с Яндексом'}
	</Button>
</div>

<style>
	.auth-form {
		display: flex;
		flex-direction: column;
		gap: 0.9rem;
	}

	.lead {
		font-size: 0.95rem;
		margin-top: -0.4rem;
		margin-bottom: 0.4rem;
	}
</style>
