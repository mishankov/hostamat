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

		padding: 6px;
		border-radius: var(--border-radius-regular);
		border: var(--border-regular);
		background: var(--color-background-object);

		transition: var(--transition);
	}

	button:hover {
		background: var(--color-accent-20);
		cursor: pointer;
		transform: translateY(-1px);
	}

	button.danger:hover {
		background: var(--color-danger-20);
		cursor: pointer;
	}
</style>
