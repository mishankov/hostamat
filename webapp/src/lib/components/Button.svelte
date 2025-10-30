<script lang="ts">
	type ButtonKind = 'primary' | 'confirm' | 'danger' | 'action';

	type ButtonProps = {
		kind?: ButtonKind;
		disabled?: boolean;
		filled?: boolean;
		link?: string;
		ref?: any;
		children?: any;
		onclick?: (e: Event) => void;
	};

	let {
		kind = 'primary',
		disabled = false,
		filled = false,
		link = '',
		children = undefined,
		ref = $bindable(null),
		onclick = undefined
	}: ButtonProps = $props();
</script>

{#if link}
	<a href={link}>
		<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
			{@render children?.()}
		</button>
	</a>
{:else}
	<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
		{@render children?.()}
	</button>
{/if}

<style>
	button {
		appearance: none;
		width: fit-content;

		border-width: 1px;
		border-style: solid;
		border-radius: 15px;

		padding: 12px 20px;

		text-align: center;
		user-select: none;

		transition: all 300ms cubic-bezier(0.23, 1, 0.32, 1);

		color: var(--color);
		background-color: var(--color-dim);
		border-color: var(--color);
	}

	.filled {
		color: var(--color-dim);
		background-color: var(--color);
		border-color: var(--color);
	}

	.primary {
		--color: var(--color-primary);
		--color-dim: var(--color-primary-dim);
	}

	.confirm {
		--color: var(--color-confirm);
		--color-dim: var(--color-confirm-dim);
	}

	.danger {
		--color: var(--color-danger);
		--color-dim: var(--color-danger-dim);
	}

	.action {
		--color: var(--color-action);
		--color-dim: var(--color-action-dim);
	}

	button:hover:enabled {
		box-shadow: var(--shadow);
		transform: translateY(-1px);
	}

	button:active:enabled {
		box-shadow: none;
		transform: translateY(+1px);
	}

	button:disabled {
		background-color: #555;
		color: #fff;
		border-color: #fff;
	}
</style>
