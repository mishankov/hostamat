<script lang="ts">
	import { browser } from '$app/environment';
	import { Button, Card, Link, ThemeToggle } from '$lib/components';
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

	let cookieConsent = $state('notset');
	if (browser) {
		cookieConsent = localStorage.getItem('cookieConsent') || '';
	}
</script>

{#if cookieConsent !== 'notset' && cookieConsent !== 'true'}
	<div class="cookie">
		<Card border kind="confirm">
			<div
				class="cookie-inner"
				style="display: flex; flex-direction: column; gap: 20px; padding: 20px;"
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
			<ThemeToggle />
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
		bottom: 20px;
		right: 20px;
		z-index: 1000;
		animation: slideInUp 0.4s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes slideInUp {
		from {
			transform: translateY(100px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	header {
		width: 100%;
		height: fit-content;

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;

		background-color: var(--color-primary-10);
		border-bottom: var(--border-regular);
		position: relative;
	}

	header::after {
		content: '';
		position: absolute;
		bottom: -4px;
		left: 0;
		right: 0;
		height: 4px;
		background: linear-gradient(90deg,
			var(--color-accent-50) 0%,
			var(--color-accent-60) 25%,
			var(--color-accent-50) 50%,
			var(--color-accent-60) 75%,
			var(--color-accent-50) 100%
		);
		opacity: 0.6;
	}

	footer {
		width: 100%;
		height: fit-content;

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;

		background-color: var(--color-primary-10);
		border-top: var(--border-regular);
		margin-top: auto;
	}

	.header-card {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		gap: 15px;

		padding: 24px 20px;

		min-height: 85px;
		width: clamp(500px, 85%, 1400px);
	}

	.nav-links {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 24px;
	}

	.user-info {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 16px;
	}

	span.logo {
		font-size: 24px;
		font-weight: 900;
		letter-spacing: -0.03em;
		position: relative;
		display: inline-block;
		padding: 4px 12px;
		background: linear-gradient(135deg, var(--color-accent-50), var(--color-accent-60));
		color: white;
		border-radius: var(--border-radius-regular);
		box-shadow:
			0 2px 0 0 var(--color-accent-70),
			0 4px 8px rgba(0, 0, 0, 0.15);
		transition: var(--transition);
	}

	span.logo:hover {
		transform: translateY(-2px);
		box-shadow:
			0 4px 0 0 var(--color-accent-70),
			0 6px 12px rgba(0, 0, 0, 0.2);
	}

	span.logo:active {
		transform: translateY(0);
		box-shadow:
			0 1px 0 0 var(--color-accent-70),
			0 2px 4px rgba(0, 0, 0, 0.15);
	}
</style>
