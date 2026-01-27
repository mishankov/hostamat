<script lang="ts">
	import { page } from '$app/stores';
	import { dev } from '$app/environment';
	import { Button, Card } from '$lib/components';
	import { X } from '$lib/icons';
</script>

<svelte:head>
	<title>{$page.status} - Хостамат</title>
</svelte:head>

<main>
	<div class="error-page">
		<Card border kind="danger" style="max-width: 600px; width: 100%; padding: 40px;">
			<div class="error-content">
				<div class="icon-wrapper">
					<X size={100} color="var(--color-danger-50)" />
				</div>
				<h1 class="error-code">{$page.status}</h1>
				<h2 class="title">
					{#if $page.status === 404}
						Страница не найдена
					{:else}
						Произошла ошибка
					{/if}
				</h2>
				<p class="message">
					{#if $page.status === 404}
						Запрашиваемая страница не существует или была удалена.<br />
						Проверьте правильность адреса или вернитесь на главную страницу.
					{:else if dev}
						{$page.error?.message || 'Произошла непредвиденная ошибка.'}
					{:else}
						Произошла непредвиденная ошибка. Попробуйте позже или обратитесь в поддержку.
					{/if}
				</p>
				<Button link="/" kind="primary" filled>Вернуться на главную</Button>
			</div>
		</Card>
	</div>
</main>

<style>
	main {
		min-height: calc(100vh - 200px);
		width: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 40px 20px;
	}

	.error-page {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 20px;
	}

	.error-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 20px;
		text-align: center;
	}

	.icon-wrapper {
		width: 60px;
		height: 60px;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 20px;
		border-radius: var(--border-radius-regular);
		background: linear-gradient(-10deg, var(--color-danger-10), var(--color-danger-15));
		border: var(--border-width) var(--border-style) var(--color-danger-50);
	}

	.error-code {
		font-size: 80px;
		font-weight: bold;
		margin: 0;
		line-height: 1;
		color: var(--color-danger-90);
	}

	.title {
		font-size: 24px;
		font-weight: 600;
		margin: 0;
		line-height: 1.5;
	}

	.message {
		font-size: 16px;
		line-height: 1.6;
		color: var(--color-text-secondary);
		margin: 0;
	}
</style>
