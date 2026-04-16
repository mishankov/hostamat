<script lang="ts">
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { onDestroy } from 'svelte';
	import { Alert, Button, Card, IconButton, Modal } from '$lib/components';
	import { Copy } from '$lib/icons';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let errorMessage = $state('');
	let deleteAccountModal = $state<Modal>();
	let loginCopied = $state(false);
	let copyTimer: ReturnType<typeof setTimeout> | null = null;

	const deleteAccountAction = async () => {
		errorMessage = '';
		try {
			const res = await fetch('/api/auth/me', {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json'
				}
			});

			if (!res.ok) {
				errorMessage = 'Ошибка при удалении аккаунта';
				deleteAccountModal?.close();
				return;
			}

			await navigateAfterAuthChange('/login');
		} catch {
			errorMessage = 'Ошибка при удалении аккаунта';
			deleteAccountModal?.close();
		}
	};

	const copyLogin = async () => {
		try {
			await navigator.clipboard.writeText(data.username);
			loginCopied = true;

			if (copyTimer) {
				clearTimeout(copyTimer);
			}

			copyTimer = setTimeout(() => {
				loginCopied = false;
			}, 1400);
		} catch {
			loginCopied = false;
		}
	};

	onDestroy(() => {
		if (copyTimer) {
			clearTimeout(copyTimer);
		}
	});
</script>

<svelte:head>
		<title>Профиль | Хостамат</title>
		<meta name="description" content="Управление учётной записью и настройками профиля." />
	</svelte:head>

<div class="profile">
	{#if errorMessage}
		<div class="alerts">
			<Alert variant="error" message={errorMessage} />
		</div>
	{/if}

	<div class="headline">
		<h1>Профиль</h1>
		<p>Управление учётной записью Яндекса и вашими данными.</p>
	</div>

	<div class="profile-data">
		<div class="data-row">
			<span class="data-label">login:</span>
			<div class="profile-data-value">
				<span>{data.username}</span>
				<IconButton label="Скопировать логин Яндекса" onclick={copyLogin}>
					<Copy size={15} />
				</IconButton>
			</div>
		</div>
		{#if loginCopied}
			<span class="copy-feedback" aria-live="polite">[copied]</span>
		{/if}
	</div>

	<Card tone="danger">
		<div class="danger-zone">
			<h2>Опасная зона</h2>
			<p class="danger-zone-description">
				Удаление аккаунта нельзя отменить. Все ваши сайты и файлы будут удалены навсегда.
			</p>
			<Button tone="danger" emphasis="solid" onclick={() => deleteAccountModal?.open()}>
				Удалить аккаунт
			</Button>
		</div>
	</Card>

	<Modal
		bind:this={deleteAccountModal}
		tone="danger"
		title="Подтвердите удаление"
		titleId="delete-account-title"
		descriptionId="delete-account-description"
		initialFocusSelector="#delete-account-confirm"
	>
		<div class="delete-modal">
			<p id="delete-account-description" class="delete-modal-text">
				Вы уверены, что хотите удалить свой аккаунт? Это приведёт к безвозвратному удалению всех
				ваших сайтов и файлов.
			</p>
			<div class="delete-actions">
				<Button
					id="delete-account-confirm"
					tone="danger"
					emphasis="solid"
					onclick={deleteAccountAction}>Удалить</Button
				>
				<Button tone="danger" onclick={() => deleteAccountModal?.close()}>Отмена</Button>
			</div>
		</div>
	</Modal>
</div>

<style>
	.profile {
		padding-bottom: var(--space-xl);
		gap: var(--space-lg);
		display: flex;
		flex-direction: column;
	}

	.alerts {
		max-width: 31.25rem;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		padding-bottom: var(--space-md);
		border-bottom: 1px dashed color-mix(in hsl, var(--color-border) 55%, transparent);
	}

	.profile-data {
		padding-top: var(--space-sm);
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}

	.profile-data-value {
		display: flex;
		flex-direction: row;
		gap: var(--space-xs);
		align-items: center;
	}

	.copy-feedback {
		font-size: 0.78rem;
		color: var(--color-success-strong);
	}

.danger-zone {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		max-width: 31.25rem;
		padding: var(--space-lg);
	}

	.danger-zone h2 {
		color: var(--color-danger-90);
	}

	.danger-zone-description {
		color: var(--color-danger-90);
		font-size: 0.9rem;
		margin: 0;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.delete-modal-text {
		color: var(--color-danger-90);
	}

	.delete-actions {
		display: flex;
		gap: var(--space-sm);
	}
</style>
