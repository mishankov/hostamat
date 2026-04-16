<script lang="ts">
	type ButtonTone = 'accent' | 'success' | 'danger';
	type LegacyButtonKind = 'primary' | 'confirm' | 'danger' | 'action';
	type ButtonType = 'button' | 'submit' | 'reset';
	type ButtonEmphasis = 'quiet' | 'solid';

	type ButtonProps = {
		tone?: ButtonTone;
		emphasis?: ButtonEmphasis;
		kind?: LegacyButtonKind;
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
		tone = undefined,
		emphasis = undefined,
		kind = undefined,
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
	const resolvedTone = $derived.by<ButtonTone>(() => {
		if (tone) {
			return tone;
		}

		if (kind === 'confirm') {
			return 'success';
		}

		if (kind === 'danger') {
			return 'danger';
		}

		return 'accent';
	});
	const resolvedEmphasis = $derived<ButtonEmphasis>(emphasis || (filled ? 'solid' : 'quiet'));

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
		class={['button', resolvedTone, resolvedEmphasis === 'solid' && 'solid']}
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
		class={['button', resolvedTone, resolvedEmphasis === 'solid' && 'solid']}
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
		border-radius: 0;

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

	.solid {
		color: var(--tone-solid-text);
		background-color: var(--tone-strong);
		border-color: var(--tone-strong);
	}

	.accent {
		--color: var(--color-accent-text);
		--color-border: var(--color-accent-border);
		--color-background: var(--color-background-object);
		--tone-strong: var(--color-accent-strong);
		--tone-solid-text: var(--color-background-object);
	}

	.success {
		--color: var(--color-success-text);
		--color-border: var(--color-success-border);
		--color-background: var(--color-success-surface);
		--tone-strong: var(--color-success-strong);
		--tone-solid-text: var(--color-background-object);
	}

	.danger {
		--color: var(--color-danger-text);
		--color-border: var(--color-danger-border);
		--color-background: var(--color-danger-surface);
		--tone-strong: var(--color-danger-strong);
		--tone-solid-text: var(--color-background-object);
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

	button.button.solid:hover:enabled,
	a.button.solid:hover:not([aria-disabled='true']) {
		transition-duration: 90ms;
		color: var(--tone-solid-text);
		background-color: color-mix(in hsl, var(--tone-strong) 90%, black 10%);
		border-color: color-mix(in hsl, var(--tone-strong) 90%, black 10%);
	}

	button.button:active:enabled,
	a.button:active:not([aria-disabled='true']) {
		box-shadow: none;
		transform: none;
	}

	button.button:disabled,
	a.button[aria-disabled='true'] {
		background-color: var(--color-disabled-surface);
		color: var(--color-disabled-text);
		border-color: var(--color-disabled-border);
		cursor: not-allowed;
	}

	a.button[aria-disabled='true'] {
		pointer-events: none;
	}
</style>
