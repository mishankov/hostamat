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
		padding: 0;
		border-radius: var(--border-radius-sm);

		color: var(--color-text-secondary);
		transition: var(--transition-fast);
		cursor: pointer;
	}

	button:hover {
		background: var(--color-primary-15);
		color: var(--color-text-regular);
	}

	button:active {
		transform: scale(0.95);
	}

	button.danger:hover {
		background: var(--color-danger-15);
		color: var(--color-danger-60);
	}

	button:focus-visible {
		outline: 2px solid var(--color-accent-50);
		outline-offset: 2px;
	}
</style>
