<script lang="ts">
	type Props = {
		text: string;
		kind?: 'accent' | 'confirm' | 'danger' | 'action' | 'muted';
	};

	let { text = '', kind = 'accent' }: Props = $props();
</script>

<span class={kind}>{text}</span>

<style>
	span {
		display: inline-flex;
		align-items: center;
		width: fit-content;

		border-radius: 6px;

		padding: 5px 10px;

		font-family: var(--font-mono);
		font-size: 0.75rem;
		font-weight: 600;
		letter-spacing: 0.03em;
		text-transform: uppercase;

		text-align: center;
		user-select: none;

		transition: all var(--transition-fast);
	}

	.accent {
		color: var(--color-accent-80);
		background: linear-gradient(135deg, var(--color-accent-10), var(--color-accent-15));
		border: 1px solid var(--color-accent-30);
	}

	.confirm {
		color: var(--color-confirm-70);
		background: linear-gradient(135deg, var(--color-confirm-10), var(--color-confirm-15));
		border: 1px solid var(--color-confirm-30);
	}

	.danger {
		color: var(--color-danger-70);
		background: linear-gradient(135deg, var(--color-danger-10), var(--color-danger-15));
		border: 1px solid var(--color-danger-30);
	}

	.action {
		color: var(--color-action-70);
		background: linear-gradient(135deg, var(--color-action-10), var(--color-action-15));
		border: 1px solid var(--color-action-30);
	}

	.muted {
		color: var(--color-text-secondary);
		background: var(--color-background-object);
		border: 1px solid var(--color-border);
	}
</style>
