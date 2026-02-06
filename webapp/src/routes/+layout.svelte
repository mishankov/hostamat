<script lang="ts">
	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { Button, Card, IconButton, Link } from '$lib/components';
	import { Moon, Sun } from '$lib/icons';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	let cookieConsent = $state('notset');
	let menuOpen = $state(false);
	let theme = $state<'light' | 'dark'>('light');

	const logout = async () => {
		await fetch('/api/auth/logout', { method: 'POST' });
		window.location.reload();
	};

	const syncTheme = () => {
		if (!browser) {
			return;
		}

		document.documentElement.setAttribute('data-theme', theme);
		localStorage.setItem('hostamat-theme', theme);
	};

	const toggleTheme = () => {
		theme = theme === 'light' ? 'dark' : 'light';
		syncTheme();
	};

	onMount(() => {
		cookieConsent = localStorage.getItem('cookieConsent') || '';

		const savedTheme = localStorage.getItem('hostamat-theme');
		if (savedTheme === 'dark' || savedTheme === 'light') {
			theme = savedTheme;
		} else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
			theme = 'dark';
		}

		syncTheme();

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

{#if cookieConsent !== 'notset' && cookieConsent !== 'true'}
	<div class="cookie">
		<Card border>
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
				<button class="menu-toggle" onclick={() => (menuOpen = !menuOpen)} aria-label="Открыть меню">
					Меню
				</button>
			</div>

			<nav class={[menuOpen && 'open']}>
				<Link href={resolve('/dashboard')} onclick={closeMenu}>Дашборд</Link>
				<Link href={resolve('/about')} onclick={closeMenu}>О сервисе</Link>
				<Link href={resolve('/faq')} onclick={closeMenu}>FAQ</Link>
			</nav>

			<div class="user-info">
				<IconButton
					onclick={toggleTheme}
					title={theme === 'light' ? 'Включить тёмную тему' : 'Включить светлую тему'}
					aria-label={theme === 'light' ? 'Включить тёмную тему' : 'Включить светлую тему'}
					class="theme-toggle"
				>
					{#if theme === 'light'}
						<Sun size={18} color="var(--color-accent-70)" />
					{:else}
						<Moon size={18} color="var(--color-accent-70)" />
					{/if}
				</IconButton>
				{#if data.isAuthenticated}
					<Link href="/profile">{data.username}</Link>
					<Button onclick={logout}>Выйти</Button>
				{:else}
					<Button link="/login">Войти</Button>
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

<style>
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

	.menu-toggle {
		display: none;
		border: var(--border-regular);
		border-radius: calc(var(--border-radius-regular) + 2px);
		background: var(--color-background-object);
		padding: 0.4rem 0.8rem;
		color: var(--color-text-regular);
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
			justify-content: flex-start;
			flex-wrap: wrap;
		}

		.footer-inner {
			flex-direction: column;
		}

		.footer-links {
			justify-content: flex-start;
		}
	}
</style>
