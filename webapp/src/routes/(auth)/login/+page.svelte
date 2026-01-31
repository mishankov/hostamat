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

<h1>Вход</h1>

{#if error}
	<Error message={error} />
{/if}

<Input id="login" type="text" label="Логин" bind:value={login} placeholder="ivan" />
<Input
	id="password"
	type="password"
	label="Пароль"
	bind:value={password}
	placeholder="verysecure"
/>

<div class="actions">
	<Button kind="accent" filled onclick={loginAction} disabled={loginButtonDisabled}>Войти</Button>
	<Button link="/register">Создать аккаунт</Button>
</div>

<style>
	h1 {
		text-align: center;
	}

	.actions {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding-top: 8px;
	}

	.actions :global(button),
	.actions :global(a) {
		width: 100%;
	}

	.actions :global(button) {
		justify-content: center;
	}
</style>
