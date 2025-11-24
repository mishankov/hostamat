<script lang="ts">
	import { browser } from '$app/environment';
	import { Button, Card, Link } from '$lib/components';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	const logout = async () => {
		await fetch('/api/auth/logout', { method: 'POST' });
		window.location.reload();
	};

	if (browser) {
		console.info(
			'%c📝 Обратная связь %c\nЕсли что-то пошло не так, вот форма обратной связи:\nhttps://forms.yandex.ru/u/6907bbfc49af47f3cae7650a/',
			'color: hsl(194, 70%, 30%); font-weight: bold; font-size: 16px;',
			'font-size: 14px;'
		);
	}

	let cookieConsent = $state('');
	if (browser) {
		cookieConsent = localStorage.getItem('cookieConsent') || '';
	}
</script>

{#if !cookieConsent}
	<div class="cookie">
		<Card kind="confirm">
			<div
				class="cookie-inner"
				style="display: flex; flex-direction: column; gap: 20px; padding: 10px;"
			>
				<p>
					Мы используем куки <Link href="/cookies" target="_blank" alwaysUnderline>вот так</Link>
				</p>
				<Button
					onclick={() => {
						cookieConsent = 'true';
						localStorage.setItem('cookieConsent', 'true');
					}}
					kind="confirm">Принять</Button
				>
			</div>
		</Card>
	</div>
{/if}

<header>
	<div class="header-card">
		<div class="nav-links">
			<Link href="/"><span class="logo">ХОСТАМАТ</span></Link>
			<Link href="/dashboard">Дашборд</Link>
			<Link href="/about">О сервисе</Link>
		</div>

		<div class="user-info">
			{#if data.isAuthenticated == true}
				<Link href="/profile">{data.username}</Link>
				<Button onclick={logout}>Выйти</Button>
			{:else}
				<Button link="/login">Войти</Button>
			{/if}
		</div>
	</div>
</header>

{@render children()}

<footer>
	<div class="header-card">
		<div class="nav-links">
			<Link href="/"
				><span class="logo" style="color: var(--color-inactive-50);">ХОСТАМАТ</span></Link
			>
			<Link href="/cookies" alwaysUnderline style="color: var(--color-inactive-50);">Куки</Link>
			<Link href="/privacy" alwaysUnderline style="color: var(--color-inactive-50);"
				>Политика конифиденциальности</Link
			>
			<Link href="/rights" alwaysUnderline style="color: var(--color-inactive-50);"
				>Права потребителя</Link
			>
			<Link href="/terms" alwaysUnderline style="color: var(--color-inactive-50);"
				>Пользовательское соглашение</Link
			>
		</div>
	</div>
</footer>

<style>
	.cookie {
		position: fixed;
		bottom: 10px;
		right: 10px;
		width: 300px;
		height: 100px;
	}

	header {
		width: 100%;
		height: fit-content;

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;

		background-color: var(--color-primary-5);
		border-bottom: var(--border-regular);
	}

	footer {
		width: 100%;
		height: fit-content;

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;

		background-color: var(--color-primary-5);
		border-top: var(--border-regular);
	}

	.header-card {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		gap: 10px;

		padding: 20px;

		height: 75px;
		width: clamp(500px, 75%, 1200px);
	}

	.nav-links {
		display: flex;
		flex-direction: row;
		align-items: end;
		gap: 20px;
	}

	.user-info {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 20px;
	}

	span.logo {
		font-size: larger;
		font-weight: bold;
	}
</style>
