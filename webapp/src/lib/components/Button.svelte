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
			{@render children?.()}
		</button>
	</a>
{:else}
	<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
		{@render children?.()}
	</button>
{/if}

<style>
	button {
		appearance: none;
		width: fit-content;

		border-width: 1px;
		border-style: solid;
		border-radius: calc(var(--border-radius-regular) + 1px);

		padding: 0.64rem 1.05rem;

		text-align: center;
		user-select: none;

		font-weight: 650;
		letter-spacing: 0.01em;

		transition:
			color 260ms ease-out,
			background-color 260ms ease-out,
			border-color 260ms ease-out;

		color: var(--color);
		background-color: var(--color-background);
		border-color: var(--color-border);
	}

	.filled {
		color: var(--color-background);
		background-color: var(--color-background-filled);
		border-color: var(--color-background-filled);
	}

	.primary {
		--color: var(--color-accent-80);
		--color-border: var(--color-accent-60);
		--color-background: var(--color-background-object);
		--color-background-filled: var(--color-accent-70);
	}

	.confirm {
		--color: var(--color-confirm-90);
		--color-border: var(--color-confirm-60);
		--color-background: var(--color-confirm-10);
		--color-background-filled: var(--color-confirm-60);
	}

	.danger {
		--color: var(--color-danger-90);
		--color-border: var(--color-danger-60);
		--color-background: var(--color-danger-10);
		--color-background-filled: var(--color-danger-60);
	}

	.action {
		--color: var(--color-accent-80);
		--color-border: var(--color-accent-60);
		--color-background: var(--color-background-object);
		--color-background-filled: var(--color-accent-70);
	}

	button:hover:enabled {
		box-shadow: none;
		transform: none;
		transition-duration: 90ms;
		color: var(--color);
		background-color: color-mix(in hsl, var(--color-background) 95%, var(--color) 5%);
		border-color: color-mix(in hsl, var(--color-border) 82%, var(--color) 18%);
	}

	button.filled:hover:enabled {
		transition-duration: 90ms;
		color: var(--color-background);
		background-color: color-mix(in hsl, var(--color-background-filled) 88%, var(--color) 12%);
		border-color: color-mix(in hsl, var(--color-background-filled) 88%, var(--color) 12%);
	}

	button:active:enabled {
		box-shadow: none;
		transform: none;
	}

	button:disabled {
		background-color: var(--color-inactive-90);
		color: var(--color-inactive-50);
		border-color: var(--color-inactive-70);
		cursor: not-allowed;
	}
</style>
