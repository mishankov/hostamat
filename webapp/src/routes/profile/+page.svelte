<script lang="ts">
	import { Button, Error, IconButton, Success } from '$lib/components';
	import Input from '$lib/components/Input.svelte';
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
	let sucessMessage = $state('');

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
</script>

<svelte:head>
	<title>Профиль | Хостамат</title>
</svelte:head>

<main>
	<div class="profile">
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

		<h1>Профиль</h1>

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
	</div>
</main>

<style>
	main {
		width: 100%;

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
	}

	.profile {
		display: flex;
		flex-direction: column;
		gap: 20px;

		width: clamp(500px, 75%, 1200px);

		padding: 20px;
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
		gap: 15px;

		max-width: 300px;
	}
</style>
