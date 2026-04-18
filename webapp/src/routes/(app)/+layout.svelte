<script lang="ts">
	import { page } from '$app/state';
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { AccountMenu, Badge } from '$lib/components';
	import type { LayoutProps } from './$types';
	import { resolve } from '$app/paths';

	type AppShellNavItem = {
		route: '/dashboard/sites' | '/dashboard/domains';
		label: string;
		badge?: string;
	};

	const navigationItems: AppShellNavItem[] = [
		{
			route: '/dashboard/sites',
			label: 'Статические сайты'
		},
		{
			route: '/dashboard/domains',
			label: 'Домены'
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
				<a href={resolve('/')} class="brand-link">
					<span class="brand">ХОСТАМАТ</span>
				</a>
			</div>

			<div class="sidebar-separator" aria-hidden="true"></div>

			<nav aria-label="Навигация по приложению" class="nav-list">
				{#each navigationItems as item (item.route)}
					<a
						href={resolve(item.route)}
						class={['nav-link', currentPath === resolve(item.route) && 'active']}
						aria-current={currentPath === resolve(item.route) ? 'page' : undefined}
					>
						<span class="nav-link-text">{item.label}</span>
						{#if item.badge}
							<Badge text={item.badge} />
						{/if}
					</a>
				{/each}
			</nav>

			<div class="sidebar-separator" aria-hidden="true"></div>

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
		grid-template-columns: 18.75rem minmax(0, 1fr);
		min-height: 100vh;
		align-items: stretch;
	}

	.app-sidebar {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		padding: var(--space-xl) var(--space-lg) var(--space-lg);
		background: var(--color-neutral-surface);
		border-right: 1px solid var(--color-neutral-border);
		min-width: 0;
	}

	.brand-block {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		padding: 0 var(--space-sm);
	}

	.brand-link {
		color: inherit;
		text-decoration: none;
	}

	.brand-link:hover {
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
		gap: var(--space-xs);
	}

	.nav-link {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-md);
		width: 100%;
		padding: var(--space-md) 0.95rem;
		border: 1px solid transparent;
		border-radius: 0;
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
		color: var(--color-accent-text);
		background: var(--color-background-object);
		border-color: var(--color-neutral-border);
	}

	.nav-link.active {
		color: var(--color-accent-text);
		background: var(--color-accent-surface);
		border-color: var(--color-accent-border);
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
		gap: var(--space-md);
	}

	.sidebar-separator {
		border-top: 1px dashed color-mix(in hsl, var(--color-neutral-border) 50%, transparent);
		margin: 0 var(--space-sm);
	}

	.app-panel {
		min-width: 0;
		padding: var(--space-xl) 1.8rem var(--space-xl);
	}

	@media (max-width: 980px) {
		.app-layout {
			grid-template-columns: minmax(0, 1fr);
		}

		.app-sidebar {
			border-right: none;
			border-bottom: 1px solid var(--color-neutral-border);
			padding: var(--space-lg) var(--space-lg);
		}

		.app-panel {
			padding: var(--space-md) var(--space-lg) var(--space-lg);
		}

		.sidebar-footer {
			margin-top: var(--space-sm);
		}
	}
</style>
