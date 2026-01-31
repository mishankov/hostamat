<script lang="ts">
	type CardKind = 'primary' | 'confirm' | 'danger' | 'action';
	let { children, border = false, floating = false, glass = false, kind = 'primary', ...restProps } = $props();
</script>

<div class={['card', floating && 'floating', border && 'border', glass && 'glass', kind]} {...restProps}>
	{@render children()}
</div>

<style>
	.card {
		min-width: fit-content;
		border-radius: var(--border-radius-lg);
		background: var(--card-bg);
		border: 1px solid transparent;
		transition: var(--transition);
		position: relative;
		overflow: hidden;
	}

	/* Subtle inner glow */
	.card::before {
		content: '';
		position: absolute;
		inset: 0;
		border-radius: inherit;
		background: radial-gradient(
			ellipse at 50% 0%,
			rgba(255, 255, 255, 0.03) 0%,
			transparent 60%
		);
		pointer-events: none;
	}

	/* ─── KIND VARIANTS ─── */
	.primary {
		--card-bg: var(--color-background-object);
		--card-border: var(--color-border);
		--card-glow: var(--glow-primary);
	}

	.confirm {
		--card-bg: var(--gradient-confirm);
		--card-border: var(--color-confirm-30);
		--card-glow: var(--glow-confirm);
		color: var(--color-confirm-90);
	}

	.danger {
		--card-bg: var(--gradient-danger);
		--card-border: var(--color-danger-30);
		--card-glow: var(--glow-danger);
		color: var(--color-danger-90);
	}

	.action {
		--card-bg: var(--gradient-action);
		--card-border: var(--color-action-30);
		--card-glow: 0 0 15px hsla(200, 85%, 60%, 0.2);
		color: var(--color-action-90);
	}

	/* ─── MODIFIERS ─── */
	.border {
		border-color: var(--card-border);
	}

	.glass {
		background: var(--glass-background);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border-color: var(--glass-border);
	}

	.floating {
		cursor: pointer;
	}

	.floating:hover {
		transform: translateY(-4px);
		box-shadow: var(--shadow-lg), var(--card-glow);
	}

	.floating:active {
		transform: translateY(-2px);
	}
</style>
