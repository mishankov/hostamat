<script lang="ts">
	import { tick } from 'svelte';
	import ThemeToggle from './ThemeToggle.svelte';

	type AccountMenuProps = {
		username: string;
		profileHref?: string;
		onlogout: () => void | Promise<void>;
		context?: 'header' | 'sidebar';
		class?: string;
	};

	let {
		username,
		profileHref = '/profile',
		onlogout,
		context = 'header',
		class: className = ''
	}: AccountMenuProps = $props();

	let open = $state(false);
	let triggerRef = $state<HTMLButtonElement | null>(null);
	let firstActionRef = $state<HTMLAnchorElement | null>(null);
	let wrapperRef = $state<HTMLDivElement | null>(null);

	const menuId = $derived.by(() => {
		const normalized = username
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-+|-+$/g, '');

		return `account-menu-${context}-${normalized || 'user'}`;
	});

	const openMenu = async () => {
		open = true;
		await tick();
		firstActionRef?.focus();
	};

	const closeMenu = () => {
		open = false;
	};

	const toggleMenu = async () => {
		if (open) {
			closeMenu();
			return;
		}

		await openMenu();
	};

	const handleDocumentClick = (event: MouseEvent) => {
		if (!open || !wrapperRef) {
			return;
		}

		const clickPath = event.composedPath();
		const target = event.target;
		if (clickPath.includes(wrapperRef) || (target instanceof Node && wrapperRef.contains(target))) {
			return;
		}

		closeMenu();
	};

	const handleDocumentKeydown = (event: KeyboardEvent) => {
		if (!open) {
			return;
		}

		if (event.key === 'Escape') {
			event.preventDefault();
			closeMenu();
			triggerRef?.focus();
		}
	};

	const handleProfileClick = () => {
		closeMenu();
	};

	const handleLogoutClick = async () => {
		closeMenu();
		await onlogout();
	};
</script>

<svelte:document onclick={handleDocumentClick} onkeydown={handleDocumentKeydown} />

<div bind:this={wrapperRef} class={['account-menu', context, className]}>
	<button
		bind:this={triggerRef}
		type="button"
		class={['menu-trigger', open && 'open']}
		aria-controls={menuId}
		aria-expanded={open}
		aria-haspopup="dialog"
		onclick={toggleMenu}
	>
		<span class="menu-trigger-label">{username}</span>
		<span class="menu-trigger-caret" aria-hidden="true"></span>
	</button>

	{#if open}
		<div
			id={menuId}
			class="menu-panel"
			aria-label={`Меню аккаунта ${username}`}
		>
			<a
				bind:this={firstActionRef}
				href={profileHref}
				class="menu-row"
				onclick={handleProfileClick}
			>
				<span>Аккаунт</span>
			</a>

			<button type="button" class="menu-row logout-row" onclick={handleLogoutClick}>
				<span>Выйти</span>
			</button>

			<div class="menu-row theme-row">
				<span class="theme-label">Тема</span>
				<ThemeToggle />
			</div>
		</div>
	{/if}
</div>

<style>
	.account-menu {
		position: relative;
		display: inline-flex;
		width: min(100%, 240px);
	}

	.account-menu.header {
		justify-content: flex-end;
	}

	.account-menu.sidebar {
		width: 100%;
	}

	.menu-trigger {
		appearance: none;
		display: inline-flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.7rem;
		width: 100%;
		min-width: 0;
		padding: 0.78rem 0.95rem;
		border: 1px solid color-mix(in hsl, var(--color-border) 90%, transparent);
		border-radius: calc(var(--border-radius-regular) + 1px);
		background: color-mix(in hsl, var(--color-background-object) 92%, var(--color-background) 8%);
		color: var(--color-accent-90);
		font-weight: 650;
		letter-spacing: 0.01em;
		cursor: pointer;
		transition:
			border-color 160ms ease-out,
			background-color 160ms ease-out,
			box-shadow 160ms ease-out,
			color 160ms ease-out;
	}

	.menu-trigger:hover,
	.menu-trigger.open {
		background: color-mix(in hsl, var(--color-accent-20) 65%, var(--color-background-object));
		border-color: color-mix(in hsl, var(--color-accent-50) 55%, transparent);
		color: var(--color-accent-95);
		box-shadow: var(--shadow-sm);
	}

	:global(:root[data-theme='dark']) .menu-trigger {
		color: var(--color-text-regular);
		background: color-mix(in hsl, var(--color-background-object) 94%, var(--color-accent-20) 6%);
	}

	:global(:root[data-theme='dark']) .menu-trigger:hover,
	:global(:root[data-theme='dark']) .menu-trigger.open {
		background: color-mix(in hsl, var(--color-background-object) 74%, var(--color-accent-30) 26%);
		border-color: color-mix(in hsl, var(--color-accent-60) 36%, var(--color-border));
		color: var(--color-primary-95);
		box-shadow: var(--shadow-md);
	}

	.menu-trigger:focus-visible {
		outline: 2px solid color-mix(in hsl, var(--color-accent-60) 80%, transparent);
		outline-offset: 2px;
	}

	.menu-trigger-label {
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.menu-trigger-caret {
		flex: none;
		width: 0.6rem;
		height: 0.6rem;
		border-right: 2px solid currentColor;
		border-bottom: 2px solid currentColor;
		transform: translateY(-0.08rem) rotate(45deg);
		transition: transform 160ms ease-out;
	}

	.menu-trigger.open .menu-trigger-caret {
		transform: translateY(0.08rem) rotate(-135deg);
	}

	.menu-panel {
		position: absolute;
		top: calc(100% + 0.6rem);
		right: 0;
		z-index: 80;
		width: min(240px, min(100vw - 1.5rem, 100%));
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
		padding: 0.45rem;
		border: 1px solid color-mix(in hsl, var(--color-border) 88%, transparent);
		border-radius: calc(var(--border-radius-regular) + 1px);
		background: color-mix(in hsl, var(--color-background-object) 96%, var(--color-background) 4%);
		box-shadow: var(--shadow-md);
	}

	.sidebar .menu-panel {
		top: auto;
		bottom: calc(100% + 0.75rem);
		left: 0;
		right: auto;
	}

	.menu-row {
		appearance: none;
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		min-height: 3rem;
		padding: 0.85rem 0.95rem;
		border: 0;
		border-radius: var(--border-radius-regular);
		background: transparent;
		color: var(--color-text-regular);
		font-size: 1rem;
		font-weight: 650;
		letter-spacing: 0.01em;
		text-align: left;
		text-decoration: none;
		cursor: pointer;
		transition:
			background-color 140ms ease-out,
			color 140ms ease-out;
	}

	a.menu-row:visited {
		color: var(--color-text-regular);
	}

	.menu-row:hover,
	.menu-row:focus-visible {
		background: color-mix(in hsl, var(--color-accent-20) 70%, var(--color-background-object));
		color: var(--color-accent-90);
		text-decoration: none;
		outline: none;
	}

	.logout-row {
		color: var(--color-accent-85);
	}

	.theme-row {
		cursor: default;
	}

	.theme-row:hover,
	.theme-row:focus-visible {
		background: color-mix(in hsl, var(--color-background-object) 84%, var(--color-background) 16%);
		color: var(--color-text-regular);
	}

	.theme-label {
		color: var(--color-text-secondary);
	}

	.theme-row :global(button) {
		flex: none;
	}

	@media (max-width: 980px) {
		.sidebar .menu-panel {
			top: calc(100% + 0.6rem);
			bottom: auto;
		}
	}

	@media (max-width: 640px) {
		.account-menu,
		.account-menu.header {
			width: 100%;
		}

		.menu-panel {
			left: 0;
			right: auto;
			width: min(100%, calc(100vw - 1.5rem));
		}
	}
</style>
