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

<div class="password-rules">
	<div class="rule" class:valid={login.length >= 5 && login.length <= 20}>
		<span class="rule-icon">{login.length >= 5 && login.length <= 20 ? '✓' : '○'}</span>
		<span>Логин от 5 до 20 символов</span>
	</div>
	<div class="rule" class:valid={password.length >= 8 && password.length <= 100}>
		<span class="rule-icon">{password.length >= 8 && password.length <= 100 ? '✓' : '○'}</span>
		<span>Пароль от 8 до 100 символов</span>
	</div>
	<div class="rule" class:valid={password.length > 0 && password === passwordAgain}>
		<span class="rule-icon">{password.length > 0 && password === passwordAgain ? '✓' : '○'}</span>
		<span>Пароли совпадают</span>
	</div>
</div>

<div class="actions">
	<Button kind="accent" onclick={registerAction} disabled={registerButtonDisabled} filled>Создать аккаунт</Button>
	<Button link="/login">Уже есть аккаунт</Button>
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

	.password-rules {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 16px;
		background: var(--color-background-elevated);
		border-radius: var(--border-radius-regular);
		border: 1px solid var(--glass-border);
	}

	.rule {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 13px;
		color: var(--color-text-muted);
		transition: var(--transition-fast);
	}

	.rule.valid {
		color: var(--color-confirm-60);
	}

	.rule-icon {
		font-family: 'JetBrains Mono', monospace;
		width: 16px;
		text-align: center;
	}
</style>
