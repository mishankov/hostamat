<script lang="ts">
	import { browser } from '$app/environment';
	import { Button, Card, Link } from '$lib/components';
	import type { LayoutProps } from './$types';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';

	let { data, children }: LayoutProps = $props();

	const logout = async () => {
		await fetch('/api/auth/logout', { method: 'POST' });
		window.location.reload();
	};

	if (browser) {
		console.info(
			'%c📝 Обратная связь %c\nЕсли что-то пошло не так, вот форма обратной связи:\nhttps://forms.yandex.ru/u/6907bbfc49af47f3cae7650a/',
			'color: hsl(160, 80%, 50%); font-weight: bold; font-size: 16px;',
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
			<div class="cookie-inner">
				<p>
					Мы используем куки <Link href="/cookies" target="_blank" alwaysUnderline>вот так</Link>
				</p>
				<Button
					onclick={() => {
						cookieConsent = 'true';
						localStorage.setItem('cookieConsent', 'true');
					}}
					kind="confirm"
					filled
				>
					Принять
				</Button>
			</div>
		</Card>
	</div>
{/if}

<header>
	<div class="header-inner">
		<div class="nav-links">
			<Link href="/">
				<span class="logo">
					<span class="logo-bracket">[</span>
					<span class="logo-text">ХОСТАМАТ</span>
					<span class="logo-bracket">]</span>
				</span>
			</Link>
			<nav class="nav-items">
				<Link href="/dashboard">Дашборд</Link>
				<Link href="/about">О сервисе</Link>
			</nav>
		</div>

		<div class="user-info">
			<ThemeToggle />
			{#if data.isAuthenticated == true}
				<Link href="/profile">{data.username}</Link>
				<Button onclick={logout}>Выйти</Button>
			{:else}
				<Button link="/login" filled>Войти</Button>
			{/if}
		</div>
	</div>
</header>

<div class="page-content">
	{@render children()}
</div>

<footer>
	<div class="footer-inner">
		<div class="footer-brand">
			<Link href="/">
				<span class="logo logo-muted">
					<span class="logo-bracket">[</span>
					<span class="logo-text">ХОСТАМАТ</span>
					<span class="logo-bracket">]</span>
				</span>
			</Link>
			<p class="footer-tagline">Хостинг статических сайтов</p>
		</div>

		<div class="footer-links">
			<Link href="/cookies" alwaysUnderline>Куки</Link>
			<Link href="/privacy" alwaysUnderline>Политика конфиденциальности</Link>
			<Link href="/rights" alwaysUnderline>Права потребителя</Link>
			<Link href="/terms" alwaysUnderline>Пользовательское соглашение</Link>
		</div>
	</div>
</footer>

<style>
	.cookie {
		position: fixed;
		bottom: 20px;
		right: 20px;
		z-index: 1000;
		animation: slideUp 0.4s ease-out;
	}

	.cookie-inner {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 20px;
		max-width: 300px;
	}

	header {
		position: sticky;
		top: 0;
		width: 100%;
		z-index: 100;

		display: flex;
		justify-content: center;

		background: var(--glass-background);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border-bottom: 1px solid var(--glass-border);
	}

	.header-inner {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 20px;

		padding: 16px 24px;
		width: 100%;
		max-width: 1400px;
	}

	.nav-links {
		display: flex;
		align-items: center;
		gap: 32px;
	}

	.nav-items {
		display: flex;
		gap: 24px;
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.logo {
		font-family: 'JetBrains Mono', monospace;
		font-size: 18px;
		font-weight: 700;
		letter-spacing: 0.05em;
		display: inline-flex;
		align-items: center;
		gap: 2px;
	}

	.logo-bracket {
		color: var(--color-accent-50);
		font-weight: 400;
	}

	.logo-text {
		background: linear-gradient(135deg, var(--color-accent-50) 0%, var(--color-primary-70) 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.logo-muted .logo-bracket {
		color: var(--color-inactive-40);
	}

	.logo-muted .logo-text {
		background: linear-gradient(135deg, var(--color-inactive-50) 0%, var(--color-inactive-40) 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.page-content {
		min-height: calc(100vh - 140px);
	}

	footer {
		width: 100%;
		display: flex;
		justify-content: center;

		background: var(--color-background-object);
		border-top: 1px solid var(--glass-border);
	}

	.footer-inner {
		display: flex;
		justify-content: space-between;
		align-items: start;
		gap: 40px;

		padding: 40px 24px;
		width: 100%;
		max-width: 1400px;
	}

	.footer-brand {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.footer-tagline {
		color: var(--color-text-muted);
		font-size: 14px;
	}

	.footer-links {
		display: flex;
		flex-wrap: wrap;
		gap: 16px 32px;
		font-size: 14px;
		color: var(--color-text-secondary);
	}

	/* Responsive */
	@media (max-width: 768px) {
		.header-inner {
			padding: 12px 16px;
		}

		.nav-links {
			gap: 16px;
		}

		.nav-items {
			display: none;
		}

		.footer-inner {
			flex-direction: column;
			gap: 24px;
			padding: 24px 16px;
		}

		.footer-links {
			flex-direction: column;
			gap: 12px;
		}
	}
</style>
