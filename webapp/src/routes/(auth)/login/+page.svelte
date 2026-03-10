<script lang="ts">
	import { Alert, Button } from '$lib/components';
	import { startYandexAuth } from '$lib/auth/start';
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

	const openAuth = async () => {
		isLoading = true;
		startErrorMessage = '';
		await startYandexAuth();
	};
</script>

<svelte:head>
	<title>Вход | Хостамат</title>
</svelte:head>

<h1>Вход</h1>
<p class="lead">Выберите способ входа, чтобы продолжить работу с сервисом.</p>

{#if errorMessage}
	<Alert variant="error" message={errorMessage} />
{/if}

<div class="auth-form">
	<Button onclick={openAuth} filled disabled={isLoading}>
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
