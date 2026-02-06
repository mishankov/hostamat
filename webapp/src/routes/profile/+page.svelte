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

<main class="page-shell">
	<div class="page-content profile">
		{#if errorMessage || sucessMessage}
			<div style="max-width: 500px;">
				{#if errorMessage}
					<Error message={errorMessage} />
				{/if}

				{#if sucessMessage}
					<Success message={sucessMessage} />
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
					<IconButton
						onclick={() => {
							navigator.clipboard.writeText(data.username);
						}}
					>
						<Copy size={15} />
					</IconButton>
				</div>
			</div>
		</div>

		<Card border>
			<div class="change-password">
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

			<Button kind="danger" onclick={changePasswordAction} disabled={changePasswordButtonDisabled}
				>Сменить пароль</Button
			>
			</div>
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

		<Modal bind:this={deleteAccountModal} kind="danger">
			<div class="delete-modal">
				<p style="color: var(--color-danger-90);">
					Вы уверены, что хотите удалить свой аккаунт? Это приведёт к безвозвратному удалению всех
					ваших сайтов и файлов.
				</p>
				<form method="dialog">
					<Button kind="danger" filled onclick={deleteAccountAction}>Удалить</Button>
					<Button kind="danger">Отмена</Button>
				</form>
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

	.headline {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	.profile-data {
		padding-top: 10px;
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

	.delete-modal > form {
		display: flex;
		gap: 10px;
	}
</style>
