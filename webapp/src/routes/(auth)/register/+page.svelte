<script lang="ts">
	import { goto } from '$app/navigation';
	import { Alert, Button, Card, Input } from '$lib/components';

	let login = $state('');
	let password = $state('');
	let passwordAgain = $state('');
	let error = $state('');
	let registerButtonDisabled = $derived(
		!(
			login.length >= 5 &&
			login.length <= 20 &&
			password.length >= 8 &&
			password.length <= 100 &&
			password === passwordAgain
		)
	);

	const registerAction = async () => {
		error = '';
		const res = await fetch('/api/auth/register', {
			method: 'POST',
			body: JSON.stringify({
				login,
				password
			})
		});

		if (res.ok) {
			await goto('/login');
			return;
		}

		const data = await res.text();
		if (data.trim() == 'invalid password') {
			error = 'Пароль не соответствует требованиям';
		} else if (data.trim() == 'invalid username') {
			error = 'Логин не соответствует требованиям';
		} else {
			error = 'Ошибка сервера';
		}
	};

	const onSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		await registerAction();
	};
</script>

<svelte:head>
	<title>Регистрация | Хостамат</title>
</svelte:head>

<h1>Регистрация</h1>
<p class="lead">Создайте аккаунт и загрузите первый сайт за пару минут.</p>

{#if error}
	<Alert variant="error" message={error} />
{/if}

<form class="auth-form" onsubmit={onSubmit}>
	<Input id="login" type="text" label="Логин" bind:value={login} placeholder="ivan" />
	<Input
		id="password"
		type="password"
		label="Пароль"
		bind:value={password}
		placeholder="verysecure"
	/>
	<Input
		id="passwordAgain"
		type="password"
		label="Пароль ещё раз"
		bind:value={passwordAgain}
		placeholder="verysecure"
	/>

	<Card kind="primary">
		<div class="password-rules">
			<p>Логин должен быть от 5 до 20 символов</p>
			<p>Пароль должен быть от 8 до 100 символов</p>
		</div>
	</Card>

	<div class="actions">
		<Button type="submit" disabled={registerButtonDisabled} filled>Создать</Button>
		<Button href="/login">Войти</Button>
	</div>
</form>

<style>
	.auth-form {
		display: flex;
		flex-direction: column;
		gap: 0.9rem;
	}

	.actions {
		width: 100%;

		display: flex;
		justify-content: space-between;
		gap: 10px;

		padding-top: 10px;
	}

	.password-rules {
		padding: 0.9rem;

		display: flex;
		flex-direction: column;
		gap: 10px;

		color: var(--color-text-secondary);
		font-size: 0.92rem;
		line-height: 1.35;
		border-left: 2px solid var(--color-accent-40);
	}

	.lead {
		font-size: 0.95rem;
		margin-top: -0.4rem;
		margin-bottom: 0.4rem;
	}
</style>
