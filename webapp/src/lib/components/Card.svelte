<script lang="ts">
	type CardTone = 'neutral' | 'success' | 'danger' | 'warning';
	type CardSurface = 'default' | 'chrome';
	let {
		children,
		class: className = '',
		tone = undefined,
		surface = undefined,
		border = false,
		cornerBorder = false,
		floating = false,
		kind = undefined,
		...restProps
	} = $props();

	const resolvedTone = $derived.by<CardTone>(() => {
		if (tone) {
			return tone;
		}

		if (kind === 'confirm') {
			return 'success';
		}

		if (kind === 'danger') {
			return 'danger';
		}

		return 'neutral';
	});
	const resolvedSurface = $derived<CardSurface>(
		surface || (border && cornerBorder ? 'chrome' : 'default')
	);
</script>

<div
	class={['card', className, floating && 'floating-compat', resolvedSurface, resolvedTone]}
	{...restProps}
>
	{@render children()}
</div>

<style>
	.card {
		position: relative;
		--card-surface: var(--color-background-object);
		--card-border-color: var(--color-border);
		--card-text-color: inherit;
		--card-corner-size: 0.875rem;
		border: var(--border-width) var(--border-style) var(--card-border-color);
		border-radius: var(--radius-lg);
		background: var(--card-surface);
		color: var(--card-text-color);
		box-shadow: var(--shadow-sm);
	}

	.neutral {
		--card-surface: var(--color-background-object);
		--card-border-color: var(--color-border);
		--card-text-color: inherit;
	}

	.success {
		--card-surface: var(--color-success-surface);
		--card-border-color: var(--color-success-border);
		--card-text-color: var(--color-success-text);
	}

	.danger {
		--card-surface: var(--color-danger-surface);
		--card-border-color: var(--color-danger-border);
		--card-text-color: var(--color-danger-text);
	}

	.warning {
		--card-surface: var(--color-warning-surface);
		--card-border-color: var(--color-warning-border);
		--card-text-color: var(--color-warning-text);
	}

	.default {
		box-shadow: none;
	}

	.chrome {
		border-radius: var(--radius-lg);
		border-color: var(--color-border);
		background: var(--card-surface);
		box-shadow: var(--shadow-md);
		isolation: isolate;
	}

	.floating-compat {
		box-shadow: var(--shadow-sm);
	}
</style>
