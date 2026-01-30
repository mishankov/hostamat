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
			'color: hsl(175, 70%, 45%); font-weight: bold; font-size: 16px;',
			'font-size: 14px;'
		);
	}

	let cookieConsent = $state('notset');
	if (browser) {
		cookieConsent = localStorage.getItem('cookieConsent') || '';
	}

	// Theme management
	let theme = $state<'dark' | 'light'>('dark');
	if (browser) {
		const savedTheme = localStorage.getItem('theme') as 'dark' | 'light' | null;
		if (savedTheme) {
			theme = savedTheme;
		} else if (window.matchMedia('(prefers-color-scheme: light)').matches) {
			theme = 'light';
		}
		document.documentElement.setAttribute('data-theme', theme);
	}

	const toggleTheme = () => {
		theme = theme === 'dark' ? 'light' : 'dark';
		if (browser) {
			localStorage.setItem('theme', theme);
			document.documentElement.setAttribute('data-theme', theme);
		}
	};
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
					<span class="logo-icon">◈</span>
					<span class="logo-text">ХОСТАМАТ</span>
				</span>
			</Link>
			<nav class="nav-menu">
				<Link href="/dashboard">Дашборд</Link>
				<Link href="/about">О сервисе</Link>
			</nav>
		</div>

		<div class="user-info">
			<button class="theme-toggle" onclick={toggleTheme} title={theme === 'dark' ? 'Светлая тема' : 'Тёмная тема'}>
				{#if theme === 'dark'}
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<circle cx="12" cy="12" r="5"/>
						<path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
					</svg>
				{:else}
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
					</svg>
				{/if}
			</button>
			{#if data.isAuthenticated == true}
				<Link href="/profile">{data.username}</Link>
				<Button onclick={logout}>Выйти</Button>
			{:else}
				<Button link="/login">Войти</Button>
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
				<span class="logo footer-logo">
					<span class="logo-icon">◈</span>
					<span class="logo-text">ХОСТАМАТ</span>
				</span>
			</Link>
			<p class="footer-tagline">Хостинг статических сайтов</p>
		</div>

		<nav class="footer-links">
			<Link href="/cookies" alwaysUnderline>Куки</Link>
			<Link href="/privacy" alwaysUnderline>Политика конфиденциальности</Link>
			<Link href="/rights" alwaysUnderline>Права потребителя</Link>
			<Link href="/terms" alwaysUnderline>Пользовательское соглашение</Link>
		</nav>
	</div>
</footer>

<style>
	.cookie {
		position: fixed;
		bottom: 20px;
		right: 20px;
		z-index: 1000;
		animation: slideUp 0.4s var(--ease-out-expo) forwards;
	}

	.cookie-inner {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 20px;
	}

	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	header {
		position: sticky;
		top: 0;
		z-index: 100;
		width: 100%;

		display: flex;
		justify-content: center;

		background: var(--color-background-elevated);
		border-bottom: 1px solid var(--color-border-subtle);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
	}

	.header-inner {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 24px;

		width: 100%;
		max-width: 1200px;
		padding: 16px 24px;
	}

	.nav-links {
		display: flex;
		align-items: center;
		gap: 32px;
	}

	.nav-menu {
		display: flex;
		align-items: center;
		gap: 24px;
	}

	.logo {
		display: flex;
		align-items: center;
		gap: 8px;
		font-family: var(--font-display);
		font-weight: 800;
		font-size: 1.1rem;
		letter-spacing: -0.02em;
		color: var(--color-text-primary);
		transition: all var(--transition-fast);
	}

	.logo:hover {
		color: var(--color-accent-50);
	}

	.logo-icon {
		font-size: 1.4rem;
		color: var(--color-accent-50);
		animation: pulse 3s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.6; }
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.theme-toggle {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;

		background: var(--color-background-object);
		border: 1px solid var(--color-border);
		border-radius: var(--border-radius-regular);

		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.theme-toggle:hover {
		background: var(--color-background-hover);
		color: var(--color-text-primary);
		border-color: var(--color-border-strong);
	}

	.page-content {
		min-height: calc(100vh - 180px);
	}

	footer {
		width: 100%;
		display: flex;
		justify-content: center;

		background: var(--color-background-object);
		border-top: 1px solid var(--color-border-subtle);
		margin-top: auto;
	}

	.footer-inner {
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
		gap: 24px;

		width: 100%;
		max-width: 1200px;
		padding: 32px 24px;
	}

	.footer-brand {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.footer-logo {
		opacity: 0.6;
	}

	.footer-logo:hover {
		opacity: 1;
	}

	.footer-tagline {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	.footer-links {
		display: flex;
		flex-wrap: wrap;
		gap: 16px 24px;
	}

	.footer-links :global(a) {
		font-size: 0.875rem;
		color: var(--color-text-secondary);
	}

	.footer-links :global(a:hover) {
		color: var(--color-text-primary);
	}

	/* Responsive */
	@media (max-width: 768px) {
		.header-inner {
			padding: 12px 16px;
		}

		.nav-links {
			gap: 16px;
		}

		.nav-menu {
			display: none;
		}

		.footer-inner {
			flex-direction: column;
			align-items: flex-start;
			gap: 24px;
			padding: 24px 16px;
		}

		.footer-links {
			flex-direction: column;
			gap: 12px;
		}
	}
</style>
