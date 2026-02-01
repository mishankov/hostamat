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

		padding: 8px;
		border-radius: var(--border-radius-regular);
		border: 2px solid transparent;

		transition: var(--transition);
		cursor: pointer;
	}

	button:hover {
		background: var(--color-primary-15);
		border-color: var(--color-primary-30);
		transform: translateY(-1px);
	}

	button:active {
		transform: translateY(0);
	}

	button.danger:hover {
		background: var(--color-danger-10);
		border-color: var(--color-danger-40);
	}
</style>
