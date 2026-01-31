<script lang="ts">
	import { Button, Error, IconButton, Modal, Success, Card } from '$lib/components';
	import Input from '$lib/components/Input.svelte';
	import { Copy } from '$lib/icons';
	import { goto } from '$app/navigation';
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
	let sucessMessage = $state('');
	let deleteAccountModal = $state<Modal>();

	const clearMessages = () => {
		errorMessage = '';
		sucessMessage = '';
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
			sucessMessage = 'Пароль успешно изменён!';
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

			goto('/login');
		} catch (error) {
			errorMessage = 'Ошибка при удалении аккаунта';
			deleteAccountModal?.close();
		}
	};
</script>

<svelte:head>
	<title>Профиль | Хостамат</title>
</svelte:head>

<main>
	<div class="profile">
		<h1>Профиль</h1>

		{#if errorMessage || sucessMessage}
			<div class="messages">
				{#if errorMessage}
					<Error message={errorMessage} />
				{/if}
				{#if sucessMessage}
					<Success message={sucessMessage} />
				{/if}
			</div>
		{/if}

		<section class="profile-section">
			<h2>Аккаунт</h2>
			<div class="profile-data">
				<div class="profile-field">
					<span class="field-label">Логин</span>
					<div class="field-value">
						<span class="username">{data.username}</span>
						<IconButton
							onclick={() => {
								navigator.clipboard.writeText(data.username);
							}}
							title="Скопировать логин"
						>
							<Copy size={16} />
						</IconButton>
					</div>
				</div>
			</div>
		</section>

		<section class="profile-section">
			<h2>Смена пароля</h2>
			<div class="password-form">
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
				<Button kind="primary" onclick={changePasswordAction} disabled={changePasswordButtonDisabled}>
					Сменить пароль
				</Button>
			</div>
		</section>

		<section class="profile-section danger-section">
			<Card kind="danger" border>
				<div class="danger-zone">
					<h2>Опасная зона</h2>
					<p>
						Удаление аккаунта нельзя отменить. Все ваши сайты и файлы будут удалены навсегда.
					</p>
					<Button kind="danger" filled onclick={() => deleteAccountModal?.open()}>
						Удалить аккаунт
					</Button>
				</div>
			</Card>
		</section>

		<Modal bind:this={deleteAccountModal} kind="danger">
			<div class="delete-modal">
				<h3>Удалить аккаунт?</h3>
				<p>
					Вы уверены, что хотите удалить свой аккаунт? Это приведёт к безвозвратному удалению всех
					ваших сайтов и файлов.
				</p>
				<form method="dialog" class="delete-actions">
					<Button kind="danger" filled onclick={deleteAccountAction}>Удалить</Button>
					<Button kind="danger">Отмена</Button>
				</form>
			</div>
		</Modal>
	</div>
</main>

<style>
	main {
		width: 100%;
		min-height: calc(100vh - 140px);

		display: flex;
		justify-content: center;
		padding: 40px 24px;
	}

	.profile {
		display: flex;
		flex-direction: column;
		gap: 32px;

		width: 100%;
		max-width: 600px;
	}

	.profile h1 {
		font-size: 32px;
		font-weight: 700;
	}

	.messages {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.profile-section {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.profile-section h2 {
		font-size: 18px;
		font-weight: 600;
		color: var(--color-text-regular);
	}

	.profile-data {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.profile-field {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.field-label {
		font-size: 14px;
		color: var(--color-text-muted);
		min-width: 60px;
	}

	.field-value {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.username {
		font-family: 'JetBrains Mono', monospace;
		font-size: 15px;
		color: var(--color-accent-60);
	}

	.password-form {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 320px;
	}

	.danger-section {
		margin-top: 16px;
	}

	.danger-zone {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 24px;
	}

	.danger-zone h2 {
		color: var(--color-danger-80);
		margin: 0;
	}

	.danger-zone p {
		color: var(--color-danger-70);
		font-size: 14px;
		line-height: 1.5;
		margin: 0;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 360px;
	}

	.delete-modal h3 {
		font-size: 20px;
		color: var(--color-danger-80);
		margin: 0;
	}

	.delete-modal p {
		color: var(--color-danger-70);
		font-size: 15px;
		line-height: 1.5;
	}

	.delete-actions {
		display: flex;
		gap: 12px;
		padding-top: 8px;
	}

	@media (max-width: 600px) {
		main {
			padding: 24px 16px;
		}

		.profile h1 {
			font-size: 24px;
		}

		.profile-field {
			flex-direction: column;
			align-items: flex-start;
			gap: 4px;
		}
	}
</style>
