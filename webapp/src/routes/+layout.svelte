<script lang="ts">
	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { startYandexAuth } from '$lib/auth/start';
	import { onMount } from 'svelte';
	import { Button, Card, IconButton, Link } from '$lib/components';
	import { Moon, Sun } from '$lib/icons';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	let cookieConsent = $state('notset');
	let menuOpen = $state(false);
	let isAuthStarting = $state(false);
	const resolveInitialTheme = (): 'light' | 'dark' => {
		if (!browser) {
			return 'light';
		}

		return document.documentElement.getAttribute('data-theme') === 'dark' ? 'dark' : 'light';
	};
	let theme = $state<'light' | 'dark'>(resolveInitialTheme());

	const logout = async () => {
		await fetch('/api/auth/logout', { method: 'POST' });
		closeMenu();
		await navigateAfterAuthChange('/');
	};

	const openLogin = async () => {
		closeMenu();
		isAuthStarting = true;
		await startYandexAuth();
	};

	const syncTheme = () => {
		if (!browser) {
			return;
		}

		document.documentElement.setAttribute('data-theme', theme);
		try {
			localStorage.setItem('hostamat-theme', theme);
		} catch {
			// Ignore storage write failures.
		}
	};

	const toggleTheme = () => {
		theme = theme === 'light' ? 'dark' : 'light';
		syncTheme();
	};

	onMount(() => {
		try {
			cookieConsent = localStorage.getItem('cookieConsent') || '';

			const savedTheme = localStorage.getItem('hostamat-theme');
			if (savedTheme === 'dark' || savedTheme === 'light') {
				theme = savedTheme;
			} else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
				theme = 'dark';
			}
		} catch {
			cookieConsent = '';
			if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
				theme = 'dark';
			}
		} finally {
			syncTheme();
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

<div class="site-chrome ambient-grid">
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

	<header>
		<div class="header-shell">
			<Card border cornerBorder class="header-card">
				<div class="brand-group">
					<Link href={resolve('/')} onclick={closeMenu}><span class="logo">ХОСТАМАТ</span></Link>
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
					<Link href={resolve('/dashboard')} onclick={closeMenu}>Дашборд</Link>
					<Link href={resolve('/about')} onclick={closeMenu}>О сервисе</Link>
					<Link href={resolve('/faq')} onclick={closeMenu}>FAQ</Link>
				</nav>

				<div class={['user-info', menuOpen && 'open']}>
					<IconButton
						label={theme === 'light' ? 'Включить тёмную тему' : 'Включить светлую тему'}
						onclick={toggleTheme}
						title={theme === 'light' ? 'Включить тёмную тему' : 'Включить светлую тему'}
						class="theme-toggle"
					>
						<span class="theme-icon theme-icon-light" aria-hidden="true">
							<Sun size={18} color="var(--color-accent-70)" />
						</span>
						<span class="theme-icon theme-icon-dark" aria-hidden="true">
							<Moon size={18} color="var(--color-accent-70)" />
						</span>
					</IconButton>
					{#if data.isAuthenticated}
						<Link href="/profile">{data.username}</Link>
						<Button onclick={logout}>Выйти</Button>
					{:else}
						<Button onclick={openLogin} disabled={isAuthStarting}>
							{isAuthStarting ? 'Переходим...' : 'Войти'}
						</Button>
					{/if}
				</div>
			</Card>
		</div>
	</header>

	{@render children()}

	<footer>
		<div class="footer-shell">
			<Card border cornerBorder class="footer-card">
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
			</Card>
		</div>
	</footer>
</div>

<style>
	.site-chrome {
		position: relative;
		z-index: 1;
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

	:global(.header-card) {
		width: var(--layout-width);
		display: grid;
		grid-template-columns: auto 1fr auto;
		align-items: center;
		gap: 1.2rem;
		padding: 0.8rem 1rem;
		box-shadow: var(--shadow-sm);
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

	:global(.theme-toggle) {
		width: 34px;
		height: 34px;
	}

	:global(.theme-toggle .theme-icon) {
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}

	:global(:root[data-theme='light'] .theme-toggle .theme-icon-dark) {
		display: none;
	}

	:global(:root[data-theme='dark'] .theme-toggle .theme-icon-light) {
		display: none;
	}

	:global(button.theme-toggle:hover) {
		background: color-mix(in hsl, var(--color-background-object) 95%, var(--color-accent-80) 5%);
		border-color: color-mix(in hsl, var(--color-accent-60) 82%, var(--color-accent-80) 18%);
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

	:global(.footer-card) {
		width: var(--layout-width);
		padding: 1rem;
		box-shadow: var(--shadow-sm);
	}

	:global(:root[data-theme='dark'] .header-card),
	:global(:root[data-theme='dark'] .footer-card) {
		box-shadow: var(--shadow-md);
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
		:global(.header-card) {
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
