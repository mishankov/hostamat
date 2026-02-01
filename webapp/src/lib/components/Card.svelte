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
		background-color: var(--color-primary-10);
		transition: var(--transition);
		position: relative;
		overflow: hidden;
	}

	.card::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 3px;
		background: linear-gradient(90deg,
			transparent 0%,
			var(--card-accent, var(--color-accent-50)) 50%,
			transparent 100%
		);
		opacity: 0;
		transition: var(--transition);
	}

	.primary {
		border-color: var(--color-border);
		--card-accent: var(--color-accent-50);
	}

	.confirm {
		color: var(--color-confirm-60);
		border-color: var(--color-confirm-40);
		background: var(--gradient-confirm);
		--card-accent: var(--color-confirm-50);
	}

	.danger {
		color: var(--color-danger-60);
		border-color: var(--color-danger-40);
		background: var(--gradient-danger);
		--card-accent: var(--color-danger-50);
	}

	.action {
		color: var(--color-action-60);
		border-color: var(--color-action-40);
		background: var(--gradient-action);
		--card-accent: var(--color-action-50);
	}

	.floating:hover {
		transform: translateY(-4px);
	}

	.floating:hover::before {
		opacity: 1;
	}

	.border {
		border: var(--border-regular);
	}
</style>
