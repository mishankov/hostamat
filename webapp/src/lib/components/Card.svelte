<script lang="ts">
	type CardKind = 'primary' | 'confirm' | 'danger' | 'action';
	let { children, border = false, floating = false, kind = 'primary', ...restProps } = $props();
</script>

<div class={['card', floating && 'floating', border && 'border', kind]} {...restProps}>
	{@render children()}
</div>

<style>
	.card {
		min-width: fit-content;

		border-radius: var(--border-radius-regular);

		background: var(--color-background-object);
		backdrop-filter: blur(14px);

		transition: transform 260ms ease-out, box-shadow 260ms ease-out;
	}

	.primary {
		border-color: var(--color-border);
	}

	.confirm {
		color: var(--color-confirm-90);
		border-color: var(--color-confirm-60);
		background: color-mix(in hsl, var(--color-confirm-10) 88%, transparent);
	}

	.danger {
		color: var(--color-danger-90);
		border-color: var(--color-danger-60);
		background: color-mix(in hsl, var(--color-danger-10) 88%, transparent);
	}

	.action {
		color: var(--color-action-90);
		border-color: var(--color-action-60);
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
	}
</style>
