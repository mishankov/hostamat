<script lang="ts">
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { Alert, Button, Input } from '$lib/components';

	let login = $state('');
	let password = $state('');
	let error = $state('');
	let loginButtonDisabled = $derived(!(login.length > 0 && password.length > 0));

	const loginAction = async () => {
		error = '';
		const res = await fetch('/api/auth/login', {
			method: 'POST',
			body: JSON.stringify({
				login,
				password
			})
		});

		if (res.ok) {
			await navigateAfterAuthChange('/dashboard');
			return;
		}

		if (res.status === 401) {
			error = 'Неверный логин или пароль';
		} else {
			error = 'Ошибка сервера';
		}
	};

	const onSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		await loginAction();
	};
</script>

<svelte:head>
	<title>Вход | Хостамат</title>
</svelte:head>

<h1>Вход</h1>
<p class="lead">Вернитесь к вашим сайтам и публикациям.</p>

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

	<div class="actions">
		<Button type="submit" filled disabled={loginButtonDisabled}>Войти</Button>
		<Button href="/register">Создать аккаунт</Button>
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

	.lead {
		font-size: 0.95rem;
		margin-top: -0.4rem;
		margin-bottom: 0.4rem;
	}
</style>
