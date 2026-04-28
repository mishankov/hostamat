<script lang="ts">
	type LinkProps = {
		href: string;
		target?: string;
		rel?: string;
		alwaysUnderline?: boolean;
		children?: any;
		onclick?: (event: Event) => void;
		[key: string]: any;
	};

	let {
		href,
		target = '_self',
		rel = '',
		children,
		alwaysUnderline = false,
		...restProps
	}: LinkProps = $props();

	const safeRel = $derived.by(() => {
		if (target !== '_blank') {
			return rel || undefined;
		}

		if (!rel) {
			return 'noopener noreferrer';
		}

		const parts = new Set(rel.split(/\s+/g).filter(Boolean));
		parts.add('noopener');
		parts.add('noreferrer');
		return [...parts].join(' ');
	});
</script>

<a {href} {target} rel={safeRel} class={alwaysUnderline ? 'underline' : ''} {...restProps}
	>{@render children?.()}</a
>

<style>
	a {
		color: var(--color-accent-text);
		text-decoration: none;
		text-decoration-thickness: 1px;
		text-underline-offset: 3px;
		transition: color 160ms var(--ease-out-quart);
		font-weight: 500;
	}

	a.underline {
		text-decoration: underline;
	}

	a:hover {
		transition-duration: 90ms;
		color: var(--color-accent-strong);
	}

	a.underline:hover {
		text-decoration: underline;
	}

	a:not(.underline) {
		background-image: linear-gradient(currentColor, currentColor);
		background-position: 0% 100%;
		background-repeat: no-repeat;
		background-size: 0% 1px;
		transition:
			color 160ms var(--ease-out-quart),
			background-size 220ms var(--ease-out-quart);
	}

	a:not(.underline):hover {
		background-size: 100% 1px;
	}

	a:visited {
		color: var(--color-accent-text);
	}

	a:focus-visible {
		outline: none;
		box-shadow: var(--focus-ring);
		border-radius: var(--radius-sm);
	}
</style>
