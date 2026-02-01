<script lang="ts">
	let { href, target = '_self', children, alwaysUnderline = false, ...restProps } = $props();
</script>

<a {href} {target} class={alwaysUnderline ? 'underline' : ''} {...restProps}
	>{@render children?.()}</a
>

<style>
	a {
		color: var(--color-text-regular);
		text-decoration: none;
		transition: var(--transition);
		position: relative;
		font-weight: 500;
		display: inline-block;
	}

	a::after {
		content: '';
		position: absolute;
		bottom: -2px;
		left: 0;
		width: 0;
		height: 2px;
		background: var(--color-accent-50);
		transition: var(--transition);
	}

	a.underline {
		text-decoration: underline;
		text-decoration-color: var(--color-accent-50);
		text-decoration-thickness: 2px;
		text-underline-offset: 4px;
	}

	a:hover {
		color: var(--color-accent-50);
	}

	a:hover::after {
		width: 100%;
	}

	a.underline:hover {
		text-decoration-color: var(--color-accent-60);
	}
</style>
