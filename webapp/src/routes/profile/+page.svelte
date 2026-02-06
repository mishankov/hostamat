<script lang="ts">
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { onDestroy } from 'svelte';
	import { Alert, Button, Card, IconButton, Input, Modal } from '$lib/components';
	import { Copy } from '$lib/icons';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let currentPassword = $state('');
	let newPassword = $state('');
	let newPasswordAgain = $state('');
	let changePasswordButtonDisabled = $derived(
		!(
			currentPassword.length > 0 &&
			newPassword.length >= 8 &&
			newPassword.length <= 100 &&
			newPasswordAgain.length >= 8 &&
			newPasswordAgain.length <= 100 &&
			newPassword === newPasswordAgain
		)
	);

	let errorMessage = $state('');
	let successMessage = $state('');
	let deleteAccountModal = $state<Modal>();
	let loginCopied = $state(false);
	let copyTimer: ReturnType<typeof setTimeout> | null = null;

	const clearMessages = () => {
		errorMessage = '';
		successMessage = '';
	};

	const clearPasswords = () => {
		currentPassword = '';
		newPassword = '';
		newPasswordAgain = '';
	};

	const changePasswordAction = async () => {
		clearMessages();
		const res = await fetch('/api/auth/change-password', {
			method: 'POST',
			body: JSON.stringify({
				currentPassword,
				newPassword
			})
		});

		if (res.ok) {
			successMessage = 'Пароль успешно изменён!';
			clearPasswords();
			return;
		}

		const data = await res.text();
		if (res.status === 401) {
			errorMessage = 'Неверный текущий пароль';
		} else if (data.trim() == 'invalid password') {
			errorMessage = 'Пароль не соответствует требованиям';
		} else {
			errorMessage = 'Ошибка сервера';
		}
	};

	const deleteAccountAction = async () => {
		clearMessages();
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
		} catch (error) {
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

	const submitChangePassword = async (event: SubmitEvent) => {
		event.preventDefault();
		await changePasswordAction();
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

<main class="page-shell">
	<div class="page-content profile">
		{#if errorMessage || successMessage}
			<div class="alerts">
				{#if errorMessage}
					<Alert variant="error" message={errorMessage} />
				{/if}

				{#if successMessage}
					<Alert variant="success" message={successMessage} />
				{/if}
			</div>
		{/if}

		<div class="headline">
			<h1>Профиль</h1>
			<p>Управление учётной записью и безопасностью.</p>
		</div>

		<div class="profile-data">
			<div class="profile-data-row">
				<span class="label">Логин</span>
				<div class="profile-data-value">
					<span>{data.username}</span>
					<IconButton label="Скопировать логин" onclick={copyLogin}>
						<Copy size={15} />
					</IconButton>
				</div>
			</div>
			{#if loginCopied}
				<p class="copy-feedback" aria-live="polite">Логин скопирован</p>
			{/if}
		</div>

		<Card border>
			<form class="change-password" onsubmit={submitChangePassword}>
				<h2>Смена пароля</h2>

				<Input
					id="oldPassword"
					type="password"
					label="Текущий пароль"
					bind:value={currentPassword}
					placeholder="verysecure"
				/>
				<Input
					id="newPassword"
					type="password"
					label="Новый пароль"
					bind:value={newPassword}
					placeholder="verysecure2"
				/>
				<Input
					id="newPasswordAgain"
					type="password"
					label="Новый пароль ещё раз"
					bind:value={newPasswordAgain}
					placeholder="verysecure2"
				/>

				<Button type="submit" kind="danger" disabled={changePasswordButtonDisabled}
					>Сменить пароль</Button
				>
			</form>
		</Card>

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
			titleId="delete-account-title"
			descriptionId="delete-account-description"
			initialFocusSelector="#delete-account-confirm"
		>
			<div class="delete-modal">
				<h2 id="delete-account-title">Подтвердите удаление</h2>
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
</main>

<style>
	main {
		padding-bottom: 2rem;
	}

	.profile {
		gap: 1rem;
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

	.change-password {
		display: flex;
		flex-direction: column;
		gap: 0.9rem;

		max-width: 420px;
		padding: 1rem;
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
