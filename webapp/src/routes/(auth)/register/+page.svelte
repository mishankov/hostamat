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

<div class="auth-header">
	<h1>Регистрация</h1>
	<p>Создайте новый аккаунт</p>
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
	<Input
		id="passwordAgain"
		type="password"
		label="Пароль ещё раз"
		bind:value={passwordAgain}
		placeholder="verysecure"
	/>
</div>

<Card kind="primary" border>
	<div class="password-rules">
		<div class="rule" class:valid={login.length >= 5 && login.length <= 20}>
			<span class="rule-icon">{login.length >= 5 && login.length <= 20 ? '✓' : '○'}</span>
			<span>Логин от 5 до 20 символов</span>
		</div>
		<div class="rule" class:valid={password.length >= 8 && password.length <= 100}>
			<span class="rule-icon">{password.length >= 8 && password.length <= 100 ? '✓' : '○'}</span>
			<span>Пароль от 8 до 100 символов</span>
		</div>
		<div class="rule" class:valid={password === passwordAgain && password.length > 0}>
			<span class="rule-icon">{password === passwordAgain && password.length > 0 ? '✓' : '○'}</span>
			<span>Пароли совпадают</span>
		</div>
	</div>
</Card>

<div class="actions">
	<Button onclick={registerAction} disabled={registerButtonDisabled} filled>Создать</Button>
	<Button link="/login">Войти</Button>
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

	.password-rules {
		padding: 16px;
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.rule {
		display: flex;
		align-items: center;
		gap: 10px;
		font-size: 0.875rem;
		color: var(--color-text-muted);
		transition: color var(--transition-fast);
	}

	.rule.valid {
		color: var(--color-confirm-60);
	}

	.rule-icon {
		font-size: 0.75rem;
		width: 16px;
		text-align: center;
	}
</style>
