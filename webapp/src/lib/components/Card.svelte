<script lang="ts">
	type CardKind = 'primary' | 'confirm' | 'danger' | 'action';
	let { children, border = false, floating = false, kind = 'primary', ...restProps } = $props();
</script>

<div class={['card', floating && 'floating', border && 'border', kind]} {...restProps}>
	{@render children()}
</div>

<style>
	.card {
		position: relative;
		min-width: fit-content;

		border-radius: var(--border-radius-lg);

		background: var(--gradient-primary);
		backdrop-filter: blur(8px);
		-webkit-backdrop-filter: blur(8px);

		transition: all var(--transition);
	}

	/* Subtle inner highlight */
	.card::before {
		content: '';
		position: absolute;
		inset: 0;
		border-radius: inherit;
		padding: 1px;
		background: linear-gradient(
			135deg,
			var(--color-primary-20) 0%,
			transparent 50%,
			transparent 100%
		);
		mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
		mask-composite: exclude;
		-webkit-mask-composite: xor;
		pointer-events: none;
		opacity: 0.5;
	}

	.primary {
		border-color: var(--color-border);
	}

	.confirm {
		--card-accent: var(--color-confirm-50);
		color: var(--color-confirm-90);
		border-color: var(--color-confirm-40);
		background: var(--gradient-confirm);
	}

	.confirm::before {
		background: linear-gradient(
			135deg,
			var(--color-confirm-30) 0%,
			transparent 50%,
			transparent 100%
		);
	}

	.danger {
		--card-accent: var(--color-danger-50);
		color: var(--color-danger-90);
		border-color: var(--color-danger-40);
		background: var(--gradient-danger);
	}

	.danger::before {
		background: linear-gradient(
			135deg,
			var(--color-danger-30) 0%,
			transparent 50%,
			transparent 100%
		);
	}

	.action {
		--card-accent: var(--color-action-50);
		color: var(--color-action-90);
		border-color: var(--color-action-40);
		background: var(--gradient-action);
	}

	.action::before {
		background: linear-gradient(
			135deg,
			var(--color-action-30) 0%,
			transparent 50%,
			transparent 100%
		);
	}

	.floating {
		cursor: pointer;
	}

	.floating:hover {
		transform: translateY(-4px);
		box-shadow: var(--shadow-lg);
	}

	.floating:active {
		transform: translateY(-2px);
		box-shadow: var(--shadow);
	}

	.border {
		border-width: 1px;
		border-style: solid;
	}
</style>
