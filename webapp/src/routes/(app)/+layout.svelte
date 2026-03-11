<script lang="ts">
	import { page } from '$app/state';
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { AccountMenu, Badge, Link } from '$lib/components';
	import type { LayoutProps } from './$types';
	import { resolve } from '$app/paths';

	type AppShellNavItem = {
		href: string;
		label: string;
		badge?: string;
	};

	const navigationItems: AppShellNavItem[] = [
		{
			href: resolve('/dashboard/sites'),
			label: 'Статические сайты'
		},
		{
			href: resolve('/dashboard/domains'),
			label: 'Домены',
			badge: 'Soon'
		}
	];

	let { children, data }: LayoutProps = $props();

	const currentPath = $derived(page.url.pathname);

	const logout = async () => {
		await fetch('/api/auth/logout', { method: 'POST' });
		await navigateAfterAuthChange('/login');
	};
</script>

<main class="app-shell">
	<div class="app-layout">
		<aside class="app-sidebar">
			<div class="brand-block">
				<Link href="/" class="brand-link">
					<span class="brand">ХОСТАМАТ</span>
				</Link>
			</div>

			<nav aria-label="Навигация по приложению" class="nav-list">
				{#each navigationItems as item (item.href)}
					<a
						href={item.href}
						class={['nav-link', currentPath === item.href && 'active']}
						aria-current={currentPath === item.href ? 'page' : undefined}
					>
						<span class="nav-link-text">{item.label}</span>
						{#if item.badge}
							<Badge text={item.badge} />
						{/if}
					</a>
				{/each}
			</nav>

			<div class="sidebar-footer">
				<AccountMenu username={data.username} onlogout={logout} context="sidebar" />
			</div>
		</aside>

		<section class="app-panel">
			{@render children?.()}
		</section>
	</div>
</main>

<style>
	.app-shell {
		min-height: 100vh;
	}

	.app-layout {
		display: grid;
		grid-template-columns: 300px minmax(0, 1fr);
		min-height: 100vh;
		align-items: stretch;
	}

	.app-sidebar {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 1.6rem 1rem 1rem;
		background: color-mix(in hsl, var(--color-background) 88%, var(--color-background-object) 12%);
		border-right: 1px solid color-mix(in hsl, var(--color-border) 88%, transparent);
		min-width: 0;
	}

	.brand-block {
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
		padding: 0 0.5rem;
	}

	:global(.brand-link) {
		text-decoration: none;
	}

	:global(.brand-link:hover) {
		text-decoration: none;
	}

	.brand {
		font-family: var(--font-display);
		font-size: 1.65rem;
		font-weight: 700;
		letter-spacing: 0.02em;
		color: var(--color-primary-95);
	}

	.nav-list {
		display: flex;
		flex-direction: column;
		gap: 0.3rem;
	}

	.nav-link {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.7rem;
		width: 100%;
		padding: 0.82rem 0.95rem;
		border: 1px solid transparent;
		border-radius: 12px;
		background: transparent;
		color: var(--color-text-regular);
		text-decoration: none;
		transition:
			background-color 180ms ease-out,
			border-color 180ms ease-out,
			box-shadow 180ms ease-out,
			color 180ms ease-out;
	}

	.nav-link:hover {
		text-decoration: none;
		color: var(--color-accent-80);
		background: color-mix(in hsl, var(--color-background-object) 40%, transparent);
		border-color: color-mix(in hsl, var(--color-border) 55%, transparent);
	}

	.nav-link.active {
		color: var(--color-accent-90);
		background: color-mix(in hsl, var(--color-accent-20) 72%, var(--color-background-object));
		border-color: color-mix(in hsl, var(--color-accent-40) 55%, transparent);
	}

	.nav-link-text {
		font-weight: 650;
		line-height: 1.3;
	}

	.sidebar-footer {
		margin-top: auto;

		display: flex;
		flex-direction: column;
		justify-content: start;
		gap: 0.7rem;
	}

	.app-panel {
		min-width: 0;
		padding: 1.6rem 1.8rem 2rem;
	}

	@media (max-width: 980px) {
		.app-layout {
			grid-template-columns: minmax(0, 1fr);
		}

		.app-sidebar {
			border-right: none;
			border-bottom: 1px solid color-mix(in hsl, var(--color-border) 88%, transparent);
			padding: 1.1rem 1rem;
		}

		.app-panel {
			padding: 1.2rem 1rem 1.5rem;
		}

		.sidebar-footer {
			margin-top: 0.5rem;
		}
	}
</style>
