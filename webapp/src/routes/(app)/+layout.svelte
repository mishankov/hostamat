<script lang="ts">
	import { page } from '$app/state';
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { AccountMenu, Badge } from '$lib/components';
	import type { LayoutProps } from './$types';
	import { resolve } from '$app/paths';
	import FeedbackWidget from './FeedbackWidget.svelte';

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

			<FeedbackWidget email={data.email} />

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
		grid-template-columns: 17rem minmax(0, 1fr);
		height: 100vh;
		overflow: hidden;
		align-items: stretch;
	}

	.app-sidebar {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		padding: var(--space-xl) var(--space-lg) var(--space-lg);
		background: var(--color-background-object);
		border-right: 1px solid var(--color-border);
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
		font-size: 1.4rem;
		font-weight: 800;
		letter-spacing: -0.03em;
		color: var(--color-text-regular);
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
		padding: 0.65rem 0.85rem;
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		background: transparent;
		color: var(--color-text-secondary);
		text-decoration: none;
		font-size: 0.92rem;
		font-weight: 500;
		transition:
			background-color 180ms var(--ease-out-quart),
			border-color 180ms var(--ease-out-quart),
			box-shadow 180ms var(--ease-out-quart),
			color 180ms var(--ease-out-quart),
			transform 120ms var(--ease-out-quart);
	}

	.nav-link:hover {
		transform: translateX(2px);
	}

	.nav-link:hover {
		text-decoration: none;
		color: var(--color-text-regular);
		background: var(--color-neutral-surface);
		border-color: var(--color-border);
	}

	.nav-link.active {
		color: var(--color-accent-text);
		background: var(--color-accent-surface);
		border-color: var(--color-accent-border);
		box-shadow: var(--shadow-sm);
	}

	.nav-link-text {
		font-weight: 600;
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
		border-top: 1px solid var(--color-border);
		margin: 0 var(--space-sm);
		opacity: 0.6;
	}

	.app-panel {
		min-width: 0;
		padding: var(--space-2xl) 2rem var(--space-2xl);
		background: var(--color-background);
		overflow-y: auto;
	}

	@media (max-width: 980px) {
		.app-layout {
			grid-template-columns: minmax(0, 1fr);
			height: auto;
			min-height: 100vh;
			overflow: visible;
		}

		.app-sidebar {
			border-right: none;
			border-bottom: 1px solid var(--color-border);
			padding: var(--space-lg) var(--space-lg);
			background: var(--color-background);
		}

		.app-panel {
			padding: var(--space-lg) var(--space-lg) var(--space-xl);
			overflow-y: visible;
		}

		.sidebar-footer {
			margin-top: var(--space-sm);
		}
	}
</style>
