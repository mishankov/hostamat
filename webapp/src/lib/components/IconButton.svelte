<script lang="ts">
	let { children, onclick, ...restProps } = $props();
</script>

<button {onclick} {...restProps}>
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
		background-color: #222;
		cursor: pointer;
	}
</style>
