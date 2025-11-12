<script>
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import Error from '$lib/components/Error.svelte';
	import Input from '$lib/components/Input.svelte';

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
			window.location.href = '/login';
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
</script>

<svelte:head>
	<title>Регистрация | Хостамат</title>
</svelte:head>

<svelte:window
	on:keydown={(e) => {
		if (e.key === 'Enter' && !registerButtonDisabled) {
			registerAction();
		}
	}}
/>

<h1>Регистрация</h1>

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
	<Button onclick={registerAction} disabled={registerButtonDisabled} filled>Создать</Button>
	<Button link="/login">Войти</Button>
</div>

<style>
	.actions {
		width: 100%;

		display: flex;
		justify-content: space-between;
		gap: 10px;

		padding-top: 10px;
	}

	.password-rules {
		padding: 10px;

		display: flex;
		flex-direction: column;
		gap: 10px;

		color: var(--color-text-secondary);
	}
</style>
