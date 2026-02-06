<script lang="ts">
	let { href, target = '_self', children, alwaysUnderline = false, ...restProps } = $props();
</script>

<a {href} {target} class={alwaysUnderline ? 'underline' : ''} {...restProps}
	>{@render children?.()}</a
>

<style>
	a {
		color: var(--color-accent-80);
		text-decoration: none;
		text-decoration-thickness: 1px;
		text-underline-offset: 3px;
		transition: var(--transition);
	}

	a.underline {
		text-decoration: underline;
	}

	a:hover {
		text-decoration: underline;
		color: var(--color-accent-60);
	}

	a:visited {
		color: var(--color-accent-80);
	}
</style>
