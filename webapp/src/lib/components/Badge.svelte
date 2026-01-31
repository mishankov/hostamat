<script lang="ts">
	type BadgeKind = 'accent' | 'primary' | 'danger' | 'confirm' | 'action';

	type Props = {
		text: string;
		kind?: BadgeKind;
	};

	let { text = '', kind = 'accent' }: Props = $props();
</script>

<span class={['badge', kind]}>{text}</span>

<style>
	.badge {
		display: inline-flex;
		align-items: center;
		width: fit-content;
		height: fit-content;

		padding: 6px 12px;
		border-radius: 20px;

		font-family: 'JetBrains Mono', monospace;
		font-size: 11px;
		font-weight: 600;
		letter-spacing: 0.03em;
		text-transform: uppercase;

		user-select: none;
		white-space: nowrap;
	}

	.accent {
		color: var(--color-accent-80);
		background: linear-gradient(135deg, var(--color-accent-20) 0%, var(--color-accent-10) 100%);
		border: 1px solid var(--color-accent-30);
	}

	.primary {
		color: var(--color-primary-80);
		background: linear-gradient(135deg, var(--color-primary-20) 0%, var(--color-primary-10) 100%);
		border: 1px solid var(--color-primary-30);
	}

	.danger {
		color: var(--color-danger-80);
		background: linear-gradient(135deg, var(--color-danger-20) 0%, var(--color-danger-10) 100%);
		border: 1px solid var(--color-danger-30);
	}

	.confirm {
		color: var(--color-confirm-80);
		background: linear-gradient(135deg, var(--color-confirm-20) 0%, var(--color-confirm-10) 100%);
		border: 1px solid var(--color-confirm-30);
	}

	.action {
		color: var(--color-action-80);
		background: linear-gradient(135deg, var(--color-action-20) 0%, var(--color-action-10) 100%);
		border: 1px solid var(--color-action-30);
	}
</style>
