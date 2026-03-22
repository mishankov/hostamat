<script lang="ts">
	type CardTone = 'neutral' | 'success' | 'danger';
	type CardSurface = 'default' | 'chrome';
	type LegacyCardKind = 'primary' | 'confirm' | 'danger' | 'action';
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
	class={[
		'card',
		className,
		floating && 'floating-compat',
		resolvedSurface,
		resolvedTone
	]}
	{...restProps}
>
	{@render children()}
</div>

<style>
	.card {
		min-width: fit-content;
		position: relative;
		--card-surface: var(--color-background-object);
		--card-border-color: var(--color-border);
		--card-text-color: inherit;
		--card-corner-size: 14px;
		border: var(--border-width) var(--border-style) var(--card-border-color);
		border-radius: 0;
		background: var(--card-surface);
		color: var(--card-text-color);
		backdrop-filter: blur(14px);
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

	.default {
		box-shadow: none;
	}

	.chrome {
		border-radius: 0;
		border-color: transparent;
		isolation: isolate;
	}

	.chrome::before {
		content: '';
		position: absolute;
		inset: 0;
		pointer-events: none;
		background:
			linear-gradient(var(--card-border-color) 0 0) top left / var(--card-corner-size)
				var(--border-width) no-repeat,
			linear-gradient(var(--card-border-color) 0 0) top left / var(--border-width)
				var(--card-corner-size) no-repeat,
			linear-gradient(var(--card-border-color) 0 0) top right / var(--card-corner-size)
				var(--border-width) no-repeat,
			linear-gradient(var(--card-border-color) 0 0) top right / var(--border-width)
				var(--card-corner-size) no-repeat,
			linear-gradient(var(--card-border-color) 0 0) bottom left / var(--card-corner-size)
				var(--border-width) no-repeat,
			linear-gradient(var(--card-border-color) 0 0) bottom left / var(--border-width)
				var(--card-corner-size) no-repeat,
			linear-gradient(var(--card-border-color) 0 0) bottom right / var(--card-corner-size)
				var(--border-width) no-repeat,
			linear-gradient(var(--card-border-color) 0 0) bottom right / var(--border-width)
				var(--card-corner-size) no-repeat;
	}

	.floating-compat {
		box-shadow: var(--shadow-sm);
		transition: transform 180ms ease-out, box-shadow 180ms ease-out;
	}

	.floating-compat:hover {
		box-shadow: var(--shadow-lg);
		transform: translateY(-2px);
	}
</style>
