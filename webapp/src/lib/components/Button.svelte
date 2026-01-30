<script lang="ts">
	type ButtonKind = 'primary' | 'confirm' | 'danger' | 'action';

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
	<a href={link}>
		<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
			<span class="btn-content">{@render children?.()}</span>
		</button>
	</a>
{:else}
	<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
		<span class="btn-content">{@render children?.()}</span>
	</button>
{/if}

<style>
	button {
		position: relative;
		appearance: none;
		width: fit-content;

		border-width: 1px;
		border-style: solid;
		border-radius: var(--border-radius-regular);

		padding: 10px 20px;

		font-family: var(--font-mono);
		font-size: 0.875rem;
		font-weight: 500;
		letter-spacing: -0.01em;
		text-align: center;
		user-select: none;

		cursor: pointer;
		transition: all var(--transition-fast);

		color: var(--color);
		background-color: var(--color-background);
		border-color: var(--color-border);
	}

	.btn-content {
		position: relative;
		z-index: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
	}

	/* Filled state gradient overlay */
	button.filled::before {
		content: '';
		position: absolute;
		inset: 0;
		border-radius: inherit;
		background: var(--gradient-fill);
		opacity: 0;
		transition: opacity var(--transition-fast);
	}

	button.filled:hover::before {
		opacity: 1;
	}

	.filled {
		color: var(--color-filled);
		background-color: var(--color-background-filled);
		border-color: var(--color-background-filled);
	}

	/* Primary variant */
	.primary {
		--color: var(--color-primary-90);
		--color-border: var(--color-primary-30);
		--color-background: var(--color-primary-10);
		--color-background-filled: var(--color-accent-50);
		--color-filled: var(--color-primary-5);
		--gradient-fill: linear-gradient(135deg, var(--color-accent-40), var(--color-accent-60));
		--glow: var(--glow-accent);
	}

	/* Confirm variant */
	.confirm {
		--color: var(--color-confirm-70);
		--color-border: var(--color-confirm-40);
		--color-background: var(--color-confirm-10);
		--color-background-filled: var(--color-confirm-50);
		--color-filled: var(--color-confirm-5);
		--gradient-fill: linear-gradient(135deg, var(--color-confirm-40), var(--color-confirm-60));
		--glow: var(--glow-confirm);
	}

	/* Danger variant */
	.danger {
		--color: var(--color-danger-70);
		--color-border: var(--color-danger-40);
		--color-background: var(--color-danger-10);
		--color-background-filled: var(--color-danger-50);
		--color-filled: var(--color-danger-95);
		--gradient-fill: linear-gradient(135deg, var(--color-danger-40), var(--color-danger-60));
		--glow: var(--glow-danger);
	}

	/* Action variant */
	.action {
		--color: var(--color-action-70);
		--color-border: var(--color-action-40);
		--color-background: var(--color-action-10);
		--color-background-filled: var(--color-action-50);
		--color-filled: var(--color-action-5);
		--gradient-fill: linear-gradient(135deg, var(--color-action-40), var(--color-action-60));
		--glow: var(--glow-action);
	}

	/* Hover states */
	button:hover:enabled {
		transform: translateY(-2px);
		box-shadow: var(--shadow), var(--glow);
		border-color: var(--color);
	}

	button.filled:hover:enabled {
		box-shadow: var(--shadow-lg), var(--glow);
	}

	/* Active state */
	button:active:enabled {
		transform: translateY(0);
		box-shadow: var(--shadow-sm);
	}

	/* Disabled state */
	button:disabled {
		background-color: var(--color-inactive-20);
		color: var(--color-inactive-50);
		border-color: var(--color-inactive-30);
		cursor: not-allowed;
		opacity: 0.7;
	}

	/* Focus state */
	button:focus-visible {
		outline: 2px solid var(--color-accent-50);
		outline-offset: 2px;
	}
</style>
