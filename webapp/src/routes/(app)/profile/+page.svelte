<script lang="ts">
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { Alert, Button, Card, CopyButton, Modal, Toolbar } from '$lib/components';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let errorMessage = $state('');
	let deleteAccountModal = $state<Modal>();

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
				<CopyButton textToCopy={data.username} label="Скопировать логин Яндекса" size={15} />
			</div>
		</div>
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
		padding-bottom: var(--space-lg);
		border-bottom: 1px solid var(--color-border);
	}

	.headline h1 {
		font-weight: 800;
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

	.danger-zone {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		max-width: 31.25rem;
		padding: var(--space-lg);
	}

	.danger-zone h2 {
		color: var(--color-danger-90);
		font-size: 1.1rem;
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
