<script>
	import Button from '$lib/components/Button.svelte';
	import Error from '$lib/components/Error.svelte';
	import Input from '$lib/components/Input.svelte';

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
			window.location.href = '/dashboard';
			return;
		}

		if (res.status === 401) {
			error = 'Неверный логин или пароль';
		} else {
			error = 'Ошибка сервера';
		}
	};
</script>

<svelte:head>
	<title>Вход | Хостамат</title>
</svelte:head>

<svelte:window
	on:keydown={(e) => {
		if (e.key === 'Enter' && !loginButtonDisabled) {
			loginAction();
		}
	}}
/>

<div class="auth-header">
	<h1>Вход</h1>
	<p>Войдите в свой аккаунт</p>
</div>

{#if error}
	<Error message={error} />
{/if}

<div class="form-fields">
	<Input id="login" type="text" label="Логин" bind:value={login} placeholder="ivan" />
	<Input
		id="password"
		type="password"
		label="Пароль"
		bind:value={password}
		placeholder="verysecure"
	/>
</div>

<div class="actions">
	<Button filled onclick={loginAction} disabled={loginButtonDisabled}>Войти</Button>
	<Button link="/register">Создать аккаунт</Button>
</div>

<style>
	.auth-header {
		text-align: center;
		margin-bottom: 8px;
	}

	.auth-header h1 {
		font-family: var(--font-display);
		font-size: 1.75rem;
		font-weight: 700;
		margin: 0 0 8px 0;
	}

	.auth-header p {
		color: var(--color-text-secondary);
		font-size: 0.9375rem;
		margin: 0;
	}

	.form-fields {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 12px;
		padding-top: 8px;
	}
</style>
