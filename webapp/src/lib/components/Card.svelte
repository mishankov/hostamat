<script lang="ts">
	type CardKind = 'primary' | 'confirm' | 'danger' | 'action';
	let {
		children,
		class: className = '',
		border = false,
		cornerBorder = false,
		floating = false,
		kind = 'primary',
		...restProps
	} = $props();
</script>

<div
	class={[
		'card',
		className,
		floating && 'floating',
		border && 'border',
		border && cornerBorder && 'border-corners',
		kind
	]}
	{...restProps}
>
	{@render children()}
</div>

<style>
	.card {
		min-width: fit-content;
		position: relative;
		--card-border-color: var(--color-border);
		--card-corner-size: 14px;

		border-radius: 0;

		background: var(--color-background-object);
		backdrop-filter: blur(14px);

		transition: transform 260ms ease-out, box-shadow 260ms ease-out;
	}

	.primary {
		--card-border-color: var(--color-border);
	}

	.confirm {
		color: var(--color-confirm-90);
		--card-border-color: var(--color-confirm-60);
		background: color-mix(in hsl, var(--color-confirm-10) 88%, transparent);
	}

	.danger {
		color: var(--color-danger-90);
		--card-border-color: var(--color-danger-60);
		background: color-mix(in hsl, var(--color-danger-10) 88%, transparent);
	}

	.action {
		color: var(--color-action-90);
		--card-border-color: var(--color-action-60);
		background: color-mix(in hsl, var(--color-action-10) 88%, transparent);
	}

	.floating:hover {
		transition-duration: 90ms;
		box-shadow: var(--shadow-soft);
		transform: translateY(-2px);
	}

	.border {
		border-width: var(--border-width);
		border-style: var(--border-style);
		border-color: var(--card-border-color);
	}

	.border-corners {
		border-color: transparent;
		isolation: isolate;
	}

	.border-corners::before {
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
</style>
