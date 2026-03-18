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
		<div class="profile-data-row">
			<span class="label">Яндекс логин</span>
			<div class="profile-data-value">
				<span>{data.username}</span>
				<IconButton label="Скопировать логин Яндекса" onclick={copyLogin}>
					<Copy size={15} />
				</IconButton>
			</div>
		</div>
		{#if loginCopied}
			<p class="copy-feedback" aria-live="polite">Логин Яндекса скопирован</p>
		{/if}
	</div>

	<Card kind="danger" border={true}>
		<div class="danger-zone">
			<h2>Опасная зона</h2>
			<p class="danger-zone-description">
				Удаление аккаунта нельзя отменить. Все ваши сайты и файлы будут удалены навсегда.
			</p>
			<Button kind="danger" filled onclick={() => deleteAccountModal?.open()}>
				Удалить аккаунт
			</Button>
		</div>
	</Card>

	<Modal
		bind:this={deleteAccountModal}
		kind="danger"
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
				<Button id="delete-account-confirm" kind="danger" filled onclick={deleteAccountAction}
					>Удалить</Button
				>
				<Button kind="danger" onclick={() => deleteAccountModal?.close()}>Отмена</Button>
			</div>
		</div>
	</Modal>
</div>

<style>
	.profile {
		padding-bottom: 2rem;
		gap: 1rem;
		display: flex;
		flex-direction: column;
	}

	.alerts {
		max-width: 500px;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	.profile-data {
		padding-top: 10px;
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
	}

	.profile-data-row {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 20px;
	}

	span.label {
		color: var(--color-text-secondary);
	}

	.profile-data-value {
		display: flex;
		flex-direction: row;
		gap: 4px;
		align-items: center;
	}

	.copy-feedback {
		font-size: 0.86rem;
		color: var(--color-confirm-60);
	}

	.danger-zone {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		max-width: 500px;
		padding: 1rem;
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
		gap: 20px;
	}

	.delete-modal-text {
		color: var(--color-danger-90);
	}

	.delete-actions {
		display: flex;
		gap: 10px;
	}
</style>
