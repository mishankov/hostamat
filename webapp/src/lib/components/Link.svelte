<script lang="ts">
	let { href, target = '_self', children, alwaysUnderline = false, ...restProps } = $props();
</script>

<a {href} {target} class={alwaysUnderline ? 'underline' : ''} {...restProps}>
	{@render children?.()}
</a>

<style>
	a {
		color: var(--color-text-regular);
		text-decoration: none;
		transition: var(--transition-fast);
		position: relative;
	}

	a::after {
		content: '';
		position: absolute;
		bottom: -2px;
		left: 0;
		width: 0;
		height: 1px;
		background: linear-gradient(90deg, var(--color-accent-50), var(--color-primary-60));
		transition: width 0.3s ease;
	}

	a:hover::after {
		width: 100%;
	}

	a:hover {
		color: var(--color-accent-60);
	}

	a.underline::after {
		width: 100%;
		background: var(--color-text-secondary);
	}

	a.underline:hover::after {
		background: linear-gradient(90deg, var(--color-accent-50), var(--color-primary-60));
	}
</style>
