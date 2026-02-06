<script lang="ts">
	type ButtonKind = 'primary' | 'confirm' | 'danger' | 'action';
	type ButtonType = 'button' | 'submit' | 'reset';

	type ButtonProps = {
		kind?: ButtonKind;
		disabled?: boolean;
		filled?: boolean;
		href?: string;
		link?: string; // Deprecated alias kept for compatibility.
		type?: ButtonType;
		ref?: any;
		children?: any;
		onclick?: (e: Event) => void;
		[key: string]: any;
	};

	let {
		kind = 'primary',
		disabled = false,
		filled = false,
		href = '',
		link = '',
		type = 'button',
		children = undefined,
		ref = $bindable(null),
		onclick = undefined,
		...restProps
	}: ButtonProps = $props();

	const resolvedHref = $derived(href || link);

	const handleAnchorClick = (event: MouseEvent) => {
		if (disabled) {
			event.preventDefault();
			event.stopPropagation();
			return;
		}

		onclick?.(event);
	};
</script>

{#if resolvedHref}
	<a
		href={resolvedHref}
		bind:this={ref}
		class={['button', kind, filled && 'filled']}
		aria-disabled={disabled ? 'true' : undefined}
		tabindex={disabled ? -1 : undefined}
		onclick={handleAnchorClick}
		{...restProps}
	>
		{@render children?.()}
	</a>
{:else}
	<button
		{type}
		{disabled}
		bind:this={ref}
		class={['button', kind, filled && 'filled']}
		{onclick}
		{...restProps}
	>
		{@render children?.()}
	</button>
{/if}

<style>
	.button {
		appearance: none;
		width: fit-content;
		display: inline-flex;
		align-items: center;
		justify-content: center;

		border-width: 1px;
		border-style: solid;
		border-radius: calc(var(--border-radius-regular) + 1px);

		padding: 0.64rem 1.05rem;

		text-align: center;
		user-select: none;

		font-weight: 650;
		letter-spacing: 0.01em;
		text-decoration: none;
		cursor: pointer;

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

	button.button:hover:enabled,
	a.button:hover:not([aria-disabled='true']) {
		box-shadow: none;
		transform: none;
		transition-duration: 90ms;
		color: var(--color);
		background-color: color-mix(in hsl, var(--color-background) 95%, var(--color) 5%);
		border-color: color-mix(in hsl, var(--color-border) 82%, var(--color) 18%);
	}

	button.button.filled:hover:enabled,
	a.button.filled:hover:not([aria-disabled='true']) {
		transition-duration: 90ms;
		color: var(--color-background);
		background-color: color-mix(in hsl, var(--color-background-filled) 88%, var(--color) 12%);
		border-color: color-mix(in hsl, var(--color-background-filled) 88%, var(--color) 12%);
	}

	button.button:active:enabled,
	a.button:active:not([aria-disabled='true']) {
		box-shadow: none;
		transform: none;
	}

	button.button:disabled,
	a.button[aria-disabled='true'] {
		background-color: var(--color-inactive-90);
		color: var(--color-inactive-50);
		border-color: var(--color-inactive-70);
		cursor: not-allowed;
	}

	a.button[aria-disabled='true'] {
		pointer-events: none;
	}
</style>
