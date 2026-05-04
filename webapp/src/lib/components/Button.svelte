<script lang="ts">
	type ButtonTone = 'accent' | 'success' | 'danger';
	type LegacyButtonKind = 'primary' | 'confirm' | 'danger' | 'action';
	type ButtonType = 'button' | 'submit' | 'reset';
	type ButtonEmphasis = 'quiet' | 'solid';

	type ButtonSize = 'default' | 'xs' | 'sm' | 'lg' | 'xl';

	type ButtonProps = {
		tone?: ButtonTone;
		emphasis?: ButtonEmphasis;
		kind?: LegacyButtonKind;
		size?: ButtonSize;
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
		size = 'default',
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
	const resolvedVariant = $derived.by(() => {
		if (resolvedTone === 'danger') {
			return resolvedEmphasis === 'solid' ? 'destructive' : 'destructive-outline';
		}

		if (resolvedTone === 'success') {
			return resolvedEmphasis === 'solid' ? 'success' : 'success-outline';
		}

		return resolvedEmphasis === 'solid' ? 'default' : 'outline';
	});

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
		class={['button', `variant-${resolvedVariant}`, `size-${size}`, disabled && 'is-disabled']}
		data-slot="button"
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
		class={['button', `variant-${resolvedVariant}`, `size-${size}`]}
		data-slot="button"
		{onclick}
		{...restProps}
	>
		{@render children?.()}
	</button>
{/if}

<style>
	.button {
		--button-height: 2.25rem;
		--button-padding-x: 0.75rem;
		--button-gap: 0.5rem;
		--button-radius: 0.7rem;
		--button-font-size: 0.92rem;
		--button-icon-size: 1.125rem;
		--button-background: var(--color-background-object);
		--button-border: var(--color-border);
		--button-text: var(--color-text-regular);
		--button-hover-background: var(--button-background);
		--button-hover-border: var(--button-border);
		--button-hover-text: var(--button-text);
		--button-active-background: var(--button-hover-background);
		--button-active-border: var(--button-hover-border);
		--button-active-text: var(--button-hover-text);
		--button-shadow: none;
		--button-inner-shadow: none;
		--button-active-shadow: none;
		--button-active-inner-shadow: none;

		position: relative;
		display: inline-flex;
		width: fit-content;
		flex-shrink: 0;
		align-items: center;
		justify-content: center;
		gap: var(--button-gap);
		min-height: var(--button-height);
		padding: 0 calc(var(--button-padding-x) - 1px);
		border: 1px solid var(--button-border);
		border-radius: var(--button-radius);
		background: var(--button-background);
		box-shadow: var(--button-shadow);
		color: var(--button-text);
		font-size: var(--button-font-size);
		font-weight: 600;
		line-height: 1;
		text-align: center;
		text-decoration: none;
		white-space: nowrap;
		user-select: none;
		cursor: pointer;
		outline: none;
		transition:
			background-color 140ms var(--ease-out-quart),
			border-color 140ms var(--ease-out-quart),
			box-shadow 140ms var(--ease-out-quart),
			color 140ms var(--ease-out-quart),
			transform 120ms var(--ease-out-quart);
	}

	.button::before {
		content: '';
		position: absolute;
		inset: 0;
		border-radius: calc(var(--button-radius) - 1px);
		box-shadow: var(--button-inner-shadow);
		pointer-events: none;
		transition: box-shadow 140ms ease;
	}

	.button::after {
		content: '';
		position: absolute;
		inset: 0;
		min-width: 2.75rem;
		min-height: 2.75rem;
	}

	.button :global(svg) {
		width: var(--button-icon-size);
		height: var(--button-icon-size);
		flex-shrink: 0;
		pointer-events: none;
		opacity: 0.82;
	}

	.button:hover:not(:disabled):not([aria-disabled='true']) {
		background: var(--button-hover-background);
		border-color: var(--button-hover-border);
		box-shadow: var(--button-shadow);
		color: var(--button-hover-text);
		transform: translateY(-1px);
	}

	.button:active:not(:disabled):not([aria-disabled='true']) {
		background: var(--button-active-background);
		border-color: var(--button-active-border);
		box-shadow: var(--button-active-shadow);
		color: var(--button-active-text);
		transform: translateY(1px) scale(0.98);
		transition-duration: 80ms;
	}

	.button:active:not(:disabled):not([aria-disabled='true'])::before {
		box-shadow: var(--button-active-inner-shadow);
	}

	.button:focus-visible {
		box-shadow: var(--button-shadow), var(--focus-ring);
	}

	.button:focus-visible::before {
		box-shadow: var(--button-inner-shadow);
	}

	.button:disabled,
	.button.is-disabled {
		background: var(--color-disabled-surface);
		border-color: var(--color-disabled-border);
		box-shadow: none;
		color: var(--color-disabled-text);
		cursor: not-allowed;
		opacity: 0.64;
		pointer-events: none;
		transform: none;
	}

	.button:disabled::before,
	.button.is-disabled::before {
		box-shadow: none;
	}

	/* Solid variants: convex by default, concave when pressed */
	.variant-default {
		--button-background: var(--color-accent-strong);
		--button-border: color-mix(in hsl, var(--color-accent-strong) 82%, black 18%);
		--button-text: var(--color-background-object);
		--button-hover-background: color-mix(in hsl, var(--color-accent-strong) 90%, black 10%);
		--button-hover-border: color-mix(in hsl, var(--color-accent-strong) 78%, black 22%);
		--button-hover-text: var(--color-background-object);
		--button-active-background: color-mix(in hsl, var(--color-accent-strong) 84%, black 16%);
		--button-active-border: color-mix(in hsl, var(--color-accent-strong) 72%, black 28%);
		--button-active-text: var(--color-background-object);
		--button-shadow:
			0 1px 2px color-mix(in hsl, var(--color-accent-strong) 24%, transparent),
			0 6px 14px color-mix(in hsl, var(--color-accent-strong) 22%, transparent);
		--button-inner-shadow: inset 0 1px 0 hsl(0 0% 100% / 0.16);
		--button-active-inner-shadow: inset 0 1px 0 hsl(0 0% 0% / 0.1);
	}

	.variant-destructive {
		--button-background: var(--color-danger-strong);
		--button-border: color-mix(in hsl, var(--color-danger-strong) 82%, black 18%);
		--button-text: var(--color-background-object);
		--button-hover-background: color-mix(in hsl, var(--color-danger-strong) 90%, black 10%);
		--button-hover-border: color-mix(in hsl, var(--color-danger-strong) 78%, black 22%);
		--button-hover-text: var(--color-background-object);
		--button-active-background: color-mix(in hsl, var(--color-danger-strong) 84%, black 16%);
		--button-active-border: color-mix(in hsl, var(--color-danger-strong) 72%, black 28%);
		--button-active-text: var(--color-background-object);
		--button-shadow:
			0 1px 2px color-mix(in hsl, var(--color-danger-strong) 24%, transparent),
			0 6px 14px color-mix(in hsl, var(--color-danger-strong) 22%, transparent);
		--button-inner-shadow: inset 0 1px 0 hsl(0 0% 100% / 0.16);
		--button-active-inner-shadow: inset 0 1px 0 hsl(0 0% 0% / 0.1);
	}

	.variant-success {
		--button-background: var(--color-success-strong);
		--button-border: color-mix(in hsl, var(--color-success-strong) 82%, black 18%);
		--button-text: var(--color-background-object);
		--button-hover-background: color-mix(in hsl, var(--color-success-strong) 90%, black 10%);
		--button-hover-border: color-mix(in hsl, var(--color-success-strong) 78%, black 22%);
		--button-hover-text: var(--color-background-object);
		--button-active-background: color-mix(in hsl, var(--color-success-strong) 84%, black 16%);
		--button-active-border: color-mix(in hsl, var(--color-success-strong) 72%, black 28%);
		--button-active-text: var(--color-background-object);
		--button-shadow:
			0 1px 2px color-mix(in hsl, var(--color-success-strong) 24%, transparent),
			0 6px 14px color-mix(in hsl, var(--color-success-strong) 22%, transparent);
		--button-inner-shadow: inset 0 1px 0 hsl(0 0% 100% / 0.16);
		--button-active-inner-shadow: inset 0 1px 0 hsl(0 0% 0% / 0.1);
	}

	/* Outline variants: subtle convex lift, flatten on press */
	.variant-outline {
		--button-background: color-mix(in hsl, var(--color-background-object) 96%, transparent);
		--button-border: color-mix(in hsl, var(--color-border) 92%, transparent);
		--button-text: var(--color-text-regular);
		--button-hover-background: color-mix(
			in hsl,
			var(--color-accent-surface) 46%,
			var(--color-background-object)
		);
		--button-hover-border: color-mix(in hsl, var(--color-border) 86%, var(--color-text-secondary));
		--button-hover-text: var(--color-text-regular);
		--button-active-background: color-mix(
			in hsl,
			var(--color-accent-surface) 62%,
			var(--color-background-object)
		);
		--button-active-border: color-mix(in hsl, var(--color-border) 78%, var(--color-text-secondary));
		--button-active-text: var(--color-text-regular);
		--button-shadow: 0 1px 2px hsl(214 30% 12% / 0.06);
		--button-inner-shadow: inset 0 1px 0 hsl(0 0% 100% / 0.2);
		--button-active-shadow: none;
		--button-active-inner-shadow: inset 0 1px 0 hsl(0 0% 0% / 0.08);
	}

	.variant-destructive-outline {
		--button-background: color-mix(in hsl, var(--color-background-object) 96%, transparent);
		--button-border: color-mix(in hsl, var(--color-danger-border) 78%, var(--color-danger-surface));
		--button-text: var(--color-danger-text);
		--button-hover-background: color-mix(
			in hsl,
			var(--color-danger-surface) 62%,
			var(--color-background-object)
		);
		--button-hover-border: color-mix(
			in hsl,
			var(--color-danger-border) 86%,
			var(--color-danger-strong)
		);
		--button-hover-text: var(--color-danger-text);
		--button-active-background: color-mix(
			in hsl,
			var(--color-danger-surface) 76%,
			var(--color-background-object)
		);
		--button-active-border: color-mix(
			in hsl,
			var(--color-danger-border) 70%,
			var(--color-danger-strong)
		);
		--button-active-text: var(--color-danger-text);
		--button-shadow: 0 1px 2px hsl(214 30% 12% / 0.06);
		--button-inner-shadow: inset 0 1px 0 hsl(0 0% 100% / 0.2);
		--button-active-shadow: none;
		--button-active-inner-shadow: inset 0 1px 0 hsl(0 0% 0% / 0.08);
	}

	.variant-success-outline {
		--button-background: color-mix(in hsl, var(--color-background-object) 96%, transparent);
		--button-border: color-mix(in hsl, var(--color-success-border) 56%, var(--color-border));
		--button-text: var(--color-success-text);
		--button-hover-background: color-mix(
			in hsl,
			var(--color-success-surface) 62%,
			var(--color-background-object)
		);
		--button-hover-border: color-mix(in hsl, var(--color-success-border) 78%, transparent);
		--button-hover-text: var(--color-success-text);
		--button-active-background: color-mix(
			in hsl,
			var(--color-success-surface) 76%,
			var(--color-background-object)
		);
		--button-active-border: color-mix(in hsl, var(--color-success-border) 90%, transparent);
		--button-active-text: var(--color-success-text);
		--button-shadow: 0 1px 2px hsl(214 30% 12% / 0.06);
		--button-inner-shadow: inset 0 1px 0 hsl(0 0% 100% / 0.2);
		--button-active-shadow: none;
		--button-active-inner-shadow: inset 0 1px 0 hsl(0 0% 0% / 0.08);
	}

	/* Subtle secondary with tactile press */
	.variant-secondary {
		--button-background: var(--color-neutral-surface);
		--button-border: color-mix(in hsl, var(--color-neutral-border) 86%, transparent);
		--button-text: var(--color-text-regular);
		--button-hover-background: color-mix(in hsl, var(--color-neutral-surface) 84%, black 16%);
		--button-hover-border: color-mix(in hsl, var(--color-neutral-border) 90%, black 10%);
		--button-hover-text: var(--color-text-regular);
		--button-active-background: color-mix(in hsl, var(--color-neutral-surface) 76%, black 24%);
		--button-active-border: color-mix(in hsl, var(--color-neutral-border) 94%, black 6%);
		--button-active-text: var(--color-text-regular);
	}

	.variant-secondary::before {
		box-shadow: none;
	}

	/* Ghost flattens on interaction */
	.variant-ghost {
		--button-background: transparent;
		--button-border: transparent;
		--button-text: var(--color-text-regular);
		--button-hover-background: color-mix(
			in hsl,
			var(--color-accent-surface) 44%,
			var(--color-background-object)
		);
		--button-hover-border: transparent;
		--button-hover-text: var(--color-text-regular);
		--button-active-background: color-mix(
			in hsl,
			var(--color-accent-surface) 60%,
			var(--color-background-object)
		);
		--button-active-border: transparent;
		--button-active-text: var(--color-text-regular);
	}

	.variant-ghost::before {
		box-shadow: none;
	}

	/* Link is flat and underlined on hover */
	.variant-link {
		--button-background: transparent;
		--button-border: transparent;
		--button-text: var(--color-text-regular);
		--button-hover-background: transparent;
		--button-hover-border: transparent;
		--button-hover-text: var(--color-text-regular);
		--button-active-background: transparent;
		--button-active-border: transparent;
		--button-active-text: var(--color-text-regular);
		padding-inline: 0;
		box-shadow: none;
		text-underline-offset: 0.26em;
	}

	.variant-link:hover:not(:disabled):not([aria-disabled='true']) {
		text-decoration: underline;
	}

	.variant-link::before,
	.variant-ghost::before,
	.variant-secondary::before {
		box-shadow: none;
	}

	/* Sizes */
	.size-default {
		--button-height: 2.25rem;
		--button-padding-x: 0.75rem;
	}

	.size-xs {
		--button-height: 1.75rem;
		--button-padding-x: 0.5rem;
		--button-gap: 0.25rem;
		--button-radius: 0.5rem;
		--button-font-size: 0.78rem;
		--button-icon-size: 1rem;
	}

	.size-sm {
		--button-height: 2rem;
		--button-padding-x: 0.625rem;
		--button-gap: 0.375rem;
		--button-font-size: 0.84rem;
	}

	.size-lg {
		--button-height: 2.5rem;
		--button-padding-x: 0.875rem;
	}

	.size-xl {
		--button-height: 2.75rem;
		--button-padding-x: 1rem;
		--button-font-size: 1rem;
		--button-icon-size: 1.25rem;
	}

	@media (max-width: 640px) {
		.button {
			--button-font-size: 0.875rem;
		}

		.size-xs {
			--button-font-size: 0.75rem;
		}

		.size-xl {
			--button-font-size: 0.95rem;
		}
	}
</style>
