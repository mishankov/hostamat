<script lang="ts">
	import { Button, Error, IconButton, Modal, Success } from '$lib/components';
	import Input from '$lib/components/Input.svelte';
	import { Copy } from '$lib/icons';
	import { goto } from '$app/navigation';
	import type { PageProps } from './$types';
	import Card from '$lib/components/Card.svelte';

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
	let copied = $state(false);

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

	const copyUsername = () => {
		navigator.clipboard.writeText(data.username);
		copied = true;
		setTimeout(() => copied = false, 2000);
	};
</script>

<svelte:head>
	<title>Профиль | Хостамат</title>
</svelte:head>

<main>
	<div class="profile">
		<header class="page-header">
			<h1>Профиль</h1>
		</header>

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

		<section class="section">
			<h2>Данные аккаунта</h2>
			<div class="profile-data">
				<div class="data-item">
					<span class="data-label">Логин</span>
					<div class="data-value">
						<span class="username">{data.username}</span>
						<IconButton onclick={copyUsername} title={copied ? 'Скопировано!' : 'Копировать'}>
							<Copy size={16} />
						</IconButton>
					</div>
				</div>
			</div>
		</section>

		<section class="section">
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

				<Button kind="action" onclick={changePasswordAction} disabled={changePasswordButtonDisabled}>
					Сменить пароль
				</Button>
			</div>
		</section>

		<section class="section danger-section">
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
				<div class="modal-actions">
					<Button kind="danger" filled onclick={deleteAccountAction}>Удалить</Button>
					<Button kind="danger" onclick={() => deleteAccountModal?.close()}>Отмена</Button>
				</div>
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

		padding: 0 24px;
	}

	.profile {
		width: 100%;
		max-width: 800px;

		display: flex;
		flex-direction: column;
		gap: 32px;

		padding: 40px 0 60px;
	}

	.page-header h1 {
		margin: 0;
	}

	.messages {
		max-width: 500px;
	}

	.section {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.section h2 {
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 700;
		margin: 0;
	}

	.profile-data {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.data-item {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.data-label {
		font-size: 0.875rem;
		color: var(--color-text-secondary);
		min-width: 60px;
	}

	.data-value {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.username {
		font-family: var(--font-mono);
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.password-form {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 350px;
	}

	.danger-section {
		max-width: 500px;
	}

	.danger-zone {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 24px;
	}

	.danger-zone h2 {
		color: var(--color-danger-70);
	}

	.danger-zone p {
		color: var(--color-danger-60);
		font-size: 0.9375rem;
		line-height: 1.6;
		margin: 0;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 400px;
	}

	.delete-modal h3 {
		font-family: var(--font-display);
		margin: 0;
	}

	.delete-modal p {
		color: var(--color-danger-70);
		line-height: 1.6;
		margin: 0;
	}

	.modal-actions {
		display: flex;
		gap: 12px;
		padding-top: 8px;
	}

	@media (max-width: 600px) {
		main {
			padding: 0 16px;
		}

		.profile {
			padding: 24px 0 40px;
			gap: 24px;
		}

		.password-form {
			max-width: 100%;
		}
	}
</style>
