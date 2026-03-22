<script lang="ts">
	import { page } from '$app/state';
	import { dev } from '$app/environment';
	import { Button, Card } from '$lib/components';

	let { data } = $props();
</script>

<svelte:head>
	<title>{page.status} - Хостамат</title>
</svelte:head>

<main class="page-shell">
	<div class="page-content error-page">
		<Card border>
			<div class="error-content">
				<h1 class="error-code">{page.status}</h1>
				<h2 class="title">
					{#if page.status === 404}
						Страница не найдена
					{:else}
						Произошла ошибка
					{/if}
				</h2>
				<p class="message">
					{#if page.status === 404}
						Запрашиваемая страница не существует или была удалена.<br />
						Проверьте правильность адреса или вернитесь на главную страницу.
					{:else if dev}
						{page.error?.message || 'Произошла непредвиденная ошибка.'}
					{:else}
						Произошла непредвиденная ошибка. Попробуйте позже или обратитесь в поддержку.
					{/if}
				</p>
				<div class="action-buttons">
					{#if data.isAuthenticated}
						<Button href="/dashboard/sites" emphasis="solid">Вернуться на дашборд</Button>
					{/if}
					<Button href="/">Вернуться на главную</Button>
				</div>
			</div>
		</Card>
	</div>
</main>

<style>
	main {
		min-height: calc(100vh - 220px);
		align-items: center;
	}

	.error-page {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.error-content {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 20px;
		padding: 1.4rem;
	}

	.error-code {
		font-weight: bold;
		font-size: clamp(3rem, 8vw, 5.5rem);
		line-height: 1;
	}

	.title {
		font-weight: bold;
	}

	.message {
		font-size: 16px;
		line-height: 1.6;
		/*color: var(--color-text-secondary);*/
		margin: 0;
	}

	.action-buttons {
		display: flex;
		gap: 20px;
	}
</style>
