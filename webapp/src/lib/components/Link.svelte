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
		color: var(--color-accent-80);
		text-decoration: none;
		text-decoration-thickness: 1px;
		text-underline-offset: 3px;
		transition: color 260ms ease-out;
	}

	a.underline {
		text-decoration: underline;
	}

	a:hover {
		transition-duration: 90ms;
		text-decoration: underline;
		color: var(--color-accent-60);
	}

	a:visited {
		color: var(--color-accent-80);
	}

	a:focus-visible {
		text-decoration: underline;
	}
</style>
