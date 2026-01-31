<script lang="ts">
	type ButtonKind = 'primary' | 'confirm' | 'danger' | 'action' | 'accent';

	type ButtonProps = {
		kind?: ButtonKind;
		disabled?: boolean;
		filled?: boolean;
		link?: string;
		ref?: any;
		children?: any;
		onclick?: (e: Event) => void;
	};

	let {
		kind = 'primary',
		disabled = false,
		filled = false,
		link = '',
		children = undefined,
		ref = $bindable(null),
		onclick = undefined
	}: ButtonProps = $props();
</script>

{#if link}
	<a href={link} class="button-link">
		<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
			<span class="button-content">{@render children?.()}</span>
		</button>
	</a>
{:else}
	<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
		<span class="button-content">{@render children?.()}</span>
	</button>
{/if}

<style>
	.button-link {
		text-decoration: none;
	}

	button {
		appearance: none;
		position: relative;
		width: fit-content;

		border-width: 1px;
		border-style: solid;
		border-radius: var(--border-radius-regular);

		padding: 10px 20px;

		font-family: 'Outfit', sans-serif;
		font-size: 14px;
		font-weight: 500;
		text-align: center;
		user-select: none;
		cursor: pointer;

		transition: var(--transition);

		color: var(--btn-color);
		background: var(--btn-bg);
		border-color: var(--btn-border);

		overflow: hidden;
	}

	.button-content {
		position: relative;
		z-index: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
	}

	/* Shimmer effect on hover */
	button::before {
		content: '';
		position: absolute;
		inset: 0;
		background: linear-gradient(
			105deg,
			transparent 20%,
			rgba(255, 255, 255, 0.1) 40%,
			rgba(255, 255, 255, 0.15) 50%,
			rgba(255, 255, 255, 0.1) 60%,
			transparent 80%
		);
		transform: translateX(-100%);
		transition: transform 0.6s ease;
	}

	button:hover:enabled::before {
		transform: translateX(100%);
	}

	/* ─── PRIMARY ─── */
	.primary {
		--btn-color: var(--color-primary-90);
		--btn-border: var(--color-primary-30);
		--btn-bg: var(--color-primary-10);
	}

	.primary:hover:enabled {
		--btn-border: var(--color-primary-40);
		--btn-bg: var(--color-primary-15);
		box-shadow: var(--shadow);
		transform: translateY(-1px);
	}

	.primary.filled {
		--btn-color: var(--color-primary-5);
		--btn-border: var(--color-primary-70);
		--btn-bg: linear-gradient(135deg, var(--color-primary-70) 0%, var(--color-primary-60) 100%);
	}

	.primary.filled:hover:enabled {
		--btn-bg: linear-gradient(135deg, var(--color-primary-80) 0%, var(--color-primary-70) 100%);
		box-shadow: var(--glow-primary), var(--shadow);
	}

	/* ─── ACCENT ─── */
	.accent {
		--btn-color: var(--color-accent-90);
		--btn-border: var(--color-accent-40);
		--btn-bg: var(--color-accent-10);
	}

	.accent:hover:enabled {
		--btn-border: var(--color-accent-50);
		--btn-bg: var(--color-accent-15);
		box-shadow: var(--shadow);
		transform: translateY(-1px);
	}

	.accent.filled {
		--btn-color: var(--color-primary-5);
		--btn-border: var(--color-accent-50);
		--btn-bg: linear-gradient(135deg, var(--color-accent-50) 0%, var(--color-accent-40) 100%);
	}

	.accent.filled:hover:enabled {
		--btn-bg: linear-gradient(135deg, var(--color-accent-60) 0%, var(--color-accent-50) 100%);
		box-shadow: var(--glow-accent), var(--shadow);
	}

	/* ─── CONFIRM ─── */
	.confirm {
		--btn-color: var(--color-confirm-80);
		--btn-border: var(--color-confirm-40);
		--btn-bg: var(--color-confirm-10);
	}

	.confirm:hover:enabled {
		--btn-border: var(--color-confirm-50);
		--btn-bg: var(--color-confirm-15);
		box-shadow: var(--shadow);
		transform: translateY(-1px);
	}

	.confirm.filled {
		--btn-color: var(--color-primary-5);
		--btn-border: var(--color-confirm-50);
		--btn-bg: linear-gradient(135deg, var(--color-confirm-50) 0%, var(--color-confirm-40) 100%);
	}

	.confirm.filled:hover:enabled {
		--btn-bg: linear-gradient(135deg, var(--color-confirm-60) 0%, var(--color-confirm-50) 100%);
		box-shadow: var(--glow-confirm), var(--shadow);
	}

	/* ─── DANGER ─── */
	.danger {
		--btn-color: var(--color-danger-80);
		--btn-border: var(--color-danger-40);
		--btn-bg: var(--color-danger-10);
	}

	.danger:hover:enabled {
		--btn-border: var(--color-danger-50);
		--btn-bg: var(--color-danger-15);
		box-shadow: var(--shadow);
		transform: translateY(-1px);
	}

	.danger.filled {
		--btn-color: white;
		--btn-border: var(--color-danger-50);
		--btn-bg: linear-gradient(135deg, var(--color-danger-50) 0%, var(--color-danger-40) 100%);
	}

	.danger.filled:hover:enabled {
		--btn-bg: linear-gradient(135deg, var(--color-danger-60) 0%, var(--color-danger-50) 100%);
		box-shadow: var(--glow-danger), var(--shadow);
	}

	/* ─── ACTION ─── */
	.action {
		--btn-color: var(--color-action-80);
		--btn-border: var(--color-action-40);
		--btn-bg: var(--color-action-10);
	}

	.action:hover:enabled {
		--btn-border: var(--color-action-50);
		--btn-bg: var(--color-action-15);
		box-shadow: var(--shadow);
		transform: translateY(-1px);
	}

	.action.filled {
		--btn-color: var(--color-primary-5);
		--btn-border: var(--color-action-50);
		--btn-bg: linear-gradient(135deg, var(--color-action-50) 0%, var(--color-action-40) 100%);
	}

	.action.filled:hover:enabled {
		--btn-bg: linear-gradient(135deg, var(--color-action-60) 0%, var(--color-action-50) 100%);
		box-shadow: 0 0 15px hsla(200, 85%, 60%, 0.3), var(--shadow);
	}

	/* ─── ACTIVE STATE ─── */
	button:active:enabled {
		transform: translateY(1px);
		box-shadow: none;
	}

	/* ─── DISABLED ─── */
	button:disabled {
		cursor: not-allowed;
		opacity: 0.5;
		background: var(--color-inactive-20);
		color: var(--color-inactive-50);
		border-color: var(--color-inactive-30);
	}
</style>
