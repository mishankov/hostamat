<script lang="ts">
	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import {
		provideThemeController,
		resolveDocumentTheme,
		ThemeController
	} from '$lib/theme/controller.svelte';
	import { onMount } from 'svelte';
	import { AccountMenu, Button, Link } from '$lib/components';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	let cookieConsent = $state('notset');
	let menuOpen = $state(false);
	const isAppRoute = $derived(
		page.url.pathname.startsWith('/dashboard') || page.url.pathname.startsWith('/profile')
	);
	const themeController = provideThemeController(
		new ThemeController(browser ? resolveDocumentTheme() : 'light')
	);

	const logout = async () => {
		await fetch('/api/auth/logout', { method: 'POST' });
		closeMenu();
		await navigateAfterAuthChange('/login');
	};

	const openLogin = async () => {
		closeMenu();
		await navigateAfterAuthChange('/login');
	};

	onMount(() => {
		try {
			cookieConsent = localStorage.getItem('cookieConsent') || '';
		} catch {
			cookieConsent = '';
		} finally {
			themeController.initializeFromBrowser();
		}

		console.info(
			'%cHostamat feedback%c\nЕсли что-то пошло не так, форма обратной связи:\nhttps://forms.yandex.ru/u/6907bbfc49af47f3cae7650a/',
			'color: hsl(212 60% 42%); font-weight: 700; font-size: 14px;',
			'font-size: 13px;'
		);
	});

	const closeMenu = () => {
		menuOpen = false;
	};
</script>

<div
	class={[
		'site-chrome',
		themeController.current === 'dark' && 'theme-dark',
		'ambient-grid',
		isAppRoute && 'ambient-grid-dim'
	]}
>
	{#if cookieConsent !== 'notset' && cookieConsent !== 'true'}
		<div class="cookie">
			<div class="cookie-inner">
				<p>
					Мы используем cookie. Подробности в <Link href="/cookies" target="_blank" alwaysUnderline
						>политике</Link
					>.
				</p>
				<Button
					onclick={() => {
						cookieConsent = 'true';
						localStorage.setItem('cookieConsent', 'true');
					}}
					emphasis="solid"
				>
					Принять
				</Button>
			</div>
		</div>
	{/if}

	{#if !isAppRoute}
		<header>
			<div class="header-shell">
				<div class="header-card-shell">
					<div class="header-card">
						<div class="brand-group">
							<Link href={resolve('/')} onclick={closeMenu}><span class="logo">ХОСТАМАТ</span></Link
							>
							<button
								type="button"
								class="menu-toggle"
								onclick={() => (menuOpen = !menuOpen)}
								aria-controls="primary-navigation"
								aria-expanded={menuOpen}
								aria-label={menuOpen ? 'Закрыть меню' : 'Открыть меню'}
							>
								<span class="burger" aria-hidden="true"></span>
							</button>
						</div>

						<nav id="primary-navigation" class={[menuOpen && 'open']}>
							<Link href={resolve('/dashboard/sites')} onclick={closeMenu}>Дашборд</Link>
							<Link href={resolve('/about')} onclick={closeMenu}>О сервисе</Link>
							<Link href={resolve('/faq')} onclick={closeMenu}>FAQ</Link>
						</nav>

						<div class={['user-info', menuOpen && 'open']}>
							{#if data.isAuthenticated}
								<AccountMenu username={data.username} onlogout={logout} context="header" />
							{:else}
								<Button onclick={openLogin}>Войти</Button>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</header>
	{/if}

	{@render children()}

	{#if !isAppRoute}
		<footer>
			<div class="footer-shell">
				<div class="footer-card-shell">
					<div class="footer-card">
						<div class="footer-inner">
							<div class="footer-brand">
								<Link href="/"><span class="logo soft">ХОСТАМАТ</span></Link>
								<p>Статический хостинг без рутины</p>
							</div>
							<div class="footer-links">
								<Link href="/cookies" alwaysUnderline>Куки</Link>
								<Link href="/privacy" alwaysUnderline>Политика конфиденциальности</Link>
								<Link href="/rights" alwaysUnderline>Права потребителя</Link>
								<Link href="/terms" alwaysUnderline>Пользовательское соглашение</Link>
							</div>
						</div>
					</div>
				</div>
			</div>
		</footer>
	{/if}
</div>

<style>
	.site-chrome {
		position: relative;
		z-index: 1;
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.site-chrome > :global(.page-shell),
	.site-chrome > :global(.app-shell) {
		flex: 1;
	}

	header {
		position: sticky;
		top: 0;
		z-index: 50;
		background: var(--color-background-object);
		border-bottom: var(--border-regular);
		box-shadow: var(--shadow-sm);
	}

	.site-chrome.theme-dark header {
		box-shadow: var(--shadow-md);
	}

	.header-shell {
		display: flex;
		justify-content: center;
	}

	.header-card-shell {
		width: var(--layout-width);
	}

	.header-card {
		width: var(--layout-width);
		display: grid;
		grid-template-columns: auto 1fr auto;
		align-items: center;
		gap: var(--space-lg);
		padding: var(--space-md) var(--space-lg);
	}

	.brand-group {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
	}

	.logo {
		font-family: var(--font-display);
		font-size: 1.1rem;
		font-weight: 800;
		letter-spacing: -0.02em;
	}

	.logo.soft {
		color: var(--color-text-secondary);
	}

	nav {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
	}

	:global(nav a) {
		padding: 0.4rem 0.7rem;
		border-radius: var(--radius-sm);
		font-weight: 500;
		font-size: 0.88rem;
		color: var(--color-text-secondary);
		transition: background 140ms ease-out, color 140ms ease-out;
	}

	:global(nav a:hover) {
		background: var(--color-neutral-surface);
		color: var(--color-text-regular);
		text-decoration: none;
	}

	.user-info {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: var(--space-md);
	}

	.menu-toggle {
		display: none;
		border: var(--border-regular);
		border-radius: var(--radius-md);
		background: var(--color-background-object);
		width: 2.25rem;
		height: 2.25rem;
		padding: 0;
		color: var(--color-text-regular);
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: background 140ms ease-out;
	}

	.menu-toggle:hover {
		background: var(--color-neutral-surface);
	}

	.menu-toggle .burger {
		position: relative;
		display: block;
		width: 1rem;
		height: 2px;
		background: currentColor;
		border-radius: 999px;
	}

	.menu-toggle .burger::before,
	.menu-toggle .burger::after {
		content: '';
		position: absolute;
		left: 0;
		width: 1rem;
		height: 2px;
		background: currentColor;
		border-radius: 999px;
	}

	.menu-toggle .burger::before {
		top: -0.32rem;
	}

	.menu-toggle .burger::after {
		top: 0.32rem;
	}

	.cookie {
		position: fixed;
		right: var(--space-md);
		bottom: var(--space-md);
		z-index: 70;
		max-width: min(28.75rem, calc(100vw - var(--space-xl)));
		background: var(--color-background-object);
		border: var(--border-regular);
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-lg);
	}

	.cookie-inner {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		padding: var(--space-lg);
	}

	footer {
		padding: var(--space-2xl) var(--space-xl) var(--space-xl);
		background: var(--color-background-object);
		border-top: var(--border-regular);
	}

	.footer-shell {
		display: flex;
		justify-content: center;
	}

	.footer-card-shell {
		width: var(--layout-width);
	}

	.footer-card {
		padding: var(--space-lg);
	}

	.footer-inner {
		display: flex;
		justify-content: space-between;
		gap: var(--space-lg);
	}

	.footer-brand {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}

	.footer-brand p {
		font-size: 0.88rem;
		color: var(--color-text-tertiary);
	}

	.footer-links {
		display: flex;
		flex-wrap: wrap;
		justify-content: flex-end;
		gap: var(--space-md) var(--space-xl);
	}

	.footer-links :global(a) {
		font-size: 0.85rem;
		color: var(--color-text-tertiary);
	}

	.footer-links :global(a:hover) {
		color: var(--color-text-regular);
	}

	@media (max-width: 900px) {
		.header-card {
			grid-template-columns: 1fr;
			gap: var(--space-md);
		}

		.brand-group {
			justify-content: space-between;
		}

		.menu-toggle {
			display: inline-flex;
			align-items: center;
		}

		nav {
			display: none;
			flex-direction: column;
			align-items: flex-start;
			padding-top: var(--space-sm);
		}

		nav.open {
			display: flex;
		}

		.user-info {
			display: none;
			justify-content: flex-start;
			flex-wrap: wrap;
		}

		.user-info.open {
			display: flex;
		}

		.footer-inner {
			flex-direction: column;
		}

		.footer-links {
			justify-content: flex-start;
		}
	}
</style>
