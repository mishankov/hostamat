<script lang="ts">
	let { children, onclick, class: className = '', kind = 'primary', ...restProps } = $props();
</script>

<button class={[kind, className]} {onclick} {...restProps}>
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

		transition:
			color 260ms ease-out,
			background-color 260ms ease-out,
			border-color 260ms ease-out;
	}

	button:hover {
		transition-duration: 90ms;
		background: var(--color-accent-20);
		cursor: pointer;
		transform: none;
	}

	button.danger:hover {
		transition-duration: 90ms;
		background: var(--color-danger-20);
		cursor: pointer;
	}
</style>
