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
		border-radius: var(--border-radius-regular);

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
		--color: var(--color-primary-90);
		--color-dim: var(--color-primary-5);
	}

	.confirm {
		--color: var(--color-confirm-50);
		--color-dim: var(--color-confirm-10);
	}

	.danger {
		--color: var(--color-danger-50);
		--color-dim: var(--color-danger-10);
	}

	.action {
		--color: var(--color-action-50);
		--color-dim: var(--color-action-10);
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
		background-color: var(--color-inactive-30);
		color: var(--color-inactive-90);
		border-color: var(--color-inactive-90);
	}
</style>
