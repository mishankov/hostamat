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
	import { AccountMenu, Button, Card, Link } from '$lib/components';
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
	const chromeCardStyle = 'box-shadow: var(--chrome-card-shadow);';

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
		!isAppRoute && 'ambient-grid',
		isAppRoute && 'dashboard-chrome'
	]}
>
	{#if cookieConsent !== 'notset' && cookieConsent !== 'true'}
		<div class="cookie">
			<Card border>
				<div class="cookie-inner">
					<p>
						Мы используем cookie. Подробности в <Link
							href="/cookies"
							target="_blank"
							alwaysUnderline>политике</Link
						>.
					</p>
					<Button
						onclick={() => {
							cookieConsent = 'true';
							localStorage.setItem('cookieConsent', 'true');
						}}
						filled
					>
						Принять
					</Button>
				</div>
			</Card>
		</div>
	{/if}

	{#if !isAppRoute}
		<header>
			<div class="header-shell">
				<div class="header-card-shell">
					<Card border cornerBorder style={chromeCardStyle}>
						<div class="header-card">
							<div class="brand-group">
								<Link href={resolve('/')} onclick={closeMenu}
									><span class="logo">ХОСТАМАТ</span></Link
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
					</Card>
				</div>
			</div>
		</header>
	{/if}

	{@render children()}

	{#if !isAppRoute}
		<footer>
			<div class="footer-shell">
				<div class="footer-card-shell">
					<Card border cornerBorder style={chromeCardStyle}>
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
					</Card>
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
		--chrome-card-shadow: var(--shadow-sm);
	}

	.site-chrome.theme-dark {
		--chrome-card-shadow: var(--shadow-md);
	}

	.dashboard-chrome {
		background: linear-gradient(
			180deg,
			color-mix(in hsl, var(--color-background) 92%, var(--color-accent-20) 8%) 0%,
			var(--color-background) 100%
		);
	}

	header {
		position: sticky;
		top: 0;
		z-index: 50;
		padding: 0.8rem 0.8rem 0;
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
		gap: 1.2rem;
		padding: 0.8rem 1rem;
	}

	.brand-group {
		display: flex;
		align-items: center;
		gap: 0.8rem;
	}

	.logo {
		font-family: var(--font-display);
		font-size: 1.04rem;
		font-weight: 700;
		letter-spacing: 0.04em;
	}

	.logo.soft {
		opacity: 0.75;
	}

	nav {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.user-info {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 0.8rem;
	}

	.menu-toggle {
		display: none;
		border: var(--border-regular);
		border-radius: calc(var(--border-radius-regular) + 2px);
		background: var(--color-background-object);
		width: 2.25rem;
		height: 2.25rem;
		padding: 0;
		color: var(--color-text-regular);
		align-items: center;
		justify-content: center;
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
		right: 0.8rem;
		bottom: 0.8rem;
		z-index: 70;
		max-width: min(460px, calc(100vw - 1.2rem));
	}

	.cookie-inner {
		display: flex;
		flex-direction: column;
		gap: 0.9rem;
		padding: 1rem;
	}

	footer {
		padding: 2rem 1.2rem 1.8rem;
	}

	.footer-shell {
		display: flex;
		justify-content: center;
	}

	.footer-card-shell {
		width: var(--layout-width);
	}

	.footer-card {
		padding: 1rem;
	}

	.footer-inner {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
	}

	.footer-brand {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.footer-links {
		display: flex;
		flex-wrap: wrap;
		justify-content: flex-end;
		gap: 0.8rem 1.2rem;
	}

	@media (max-width: 900px) {
		.header-card {
			grid-template-columns: 1fr;
			gap: 0.8rem;
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
			padding-top: 0.4rem;
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
