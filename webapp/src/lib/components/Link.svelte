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
		font-weight: 500;
		transition: all var(--transition-fast);
		position: relative;
	}

	a::after {
		content: '';
		position: absolute;
		bottom: -2px;
		left: 0;
		width: 0;
		height: 1px;
		background: var(--color-accent-50);
		transition: width var(--transition-fast);
	}

	a:hover {
		color: var(--color-accent-50);
	}

	a:hover::after {
		width: 100%;
	}

	a.underline {
		text-decoration: underline;
		text-underline-offset: 3px;
		text-decoration-color: var(--color-border);
	}

	a.underline:hover {
		text-decoration-color: var(--color-accent-50);
	}

	a.underline::after {
		display: none;
	}
</style>
