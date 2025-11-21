<script lang="ts">
	type CardKind = 'primary' | 'confirm' | 'danger' | 'action';
	let { children, floating = false, kind = 'primary', ...restProps } = $props();
</script>

<div class={['card', floating && 'floating', kind]} {...restProps}>
	{@render children()}
</div>

<style>
	.card {
		min-width: fit-content;

		border: var(--border-regular);
		border-radius: var(--border-radius-regular);

		background-color: var(--color-background-object);

		transition: var(--transition);
	}

	.primary {
		border-color: var(--color-border);
	}

	.confirm {
		color: var(--color-confirm-90);
		border-color: var(--color-confirm-60);
		background: var(--gradient-confirm);
	}

	.danger {
		color: var(--color-danger-90);
		border-color: var(--color-danger-60);
		background: var(--gradient-danger);
	}

	.action {
		color: var(--color-action-90);
		border-color: var(--color-action-60);
		background: var(--gradient-action);
	}

	.floating:hover {
		box-shadow: var(--shadow);
		transform: translateY(-2px);
	}
</style>
