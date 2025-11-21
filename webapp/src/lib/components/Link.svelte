<script lang="ts">
	let { href, target = '_self', children, alwaysUnderline = false, ...restProps } = $props();
</script>

<a {href} {target} class={alwaysUnderline ? 'underline' : ''} {...restProps}
	>{@render children?.()}</a
>

<style>
	a {
		color: var(--color-primary-90);
		text-decoration: none;
		transition: all 500ms;
	}

	a.underline {
		text-decoration: underline;
	}

	a:hover {
		text-decoration: underline;
	}
</style>
