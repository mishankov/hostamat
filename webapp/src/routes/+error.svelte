<script lang="ts">
	import { page } from '$app/state';
	import { dev } from '$app/environment';
	import { Button } from '$lib/components';

	let { data } = $props();
</script>

<svelte:head>
	<title>{page.status} - Хостамат</title>
</svelte:head>

<main>
	<div class="error-page">
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
					<Button link="/dashboard" kind="primary" filled>Вернуться на дашборд</Button>
				{/if}
				<Button link="/" kind="primary">Вернуться на главную</Button>
			</div>
		</div>
	</div>
</main>

<style>
	main {
		min-height: calc(100vh - 80px);
		width: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
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
	}

	.error-code {
		font-weight: bold;
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
