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

		padding: 4px;
		border-radius: var(--border-radius-regular);

		transition: var(--transition);
	}

	button:hover {
		background: var(--color-primary-15);
		cursor: pointer;
	}

	button.danger:hover {
		background: var(--color-danger-15);
		cursor: pointer;
	}
</style>
