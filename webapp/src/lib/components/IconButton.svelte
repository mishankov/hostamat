<script lang="ts">
	let { children, onclick, kind = 'primary', ...restProps } = $props();
</script>

<button class={kind} {onclick} {...restProps}>
	{@render children()}
</button>

<style>
	button {
		all: unset;
		display: flex;
		align-items: center;
		justify-content: center;

		width: 36px;
		height: 36px;
		padding: 6px;
		border-radius: var(--border-radius-sm);

		color: var(--color-text-secondary);
		background: transparent;
		cursor: pointer;

		transition: all var(--transition-fast);
	}

	button:hover {
		color: var(--color-text-primary);
		background: var(--color-background-hover);
	}

	button:active {
		transform: scale(0.95);
	}

	button.danger:hover {
		color: var(--color-danger-60);
		background: var(--color-danger-10);
	}

	button:focus-visible {
		outline: 2px solid var(--color-accent-50);
		outline-offset: 2px;
	}
</style>
