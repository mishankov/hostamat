<script lang="ts">
	type IconProps = {
		size?: number | string;
		color?: string;
		class?: string;
	};

	let {
		size = 100,
		color = 'var(--color-primary-90)',
		class: className = ''
	}: IconProps = $props();
</script>

<svg
	class={className}
	width={size}
	height={size}
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="1.75"
	stroke-linecap="round"
	stroke-linejoin="round"
	style:color
	aria-hidden="true"
>
	<circle cx="12" cy="12" r="4" />
	<path d="M12 2.75v2" />
	<path d="M12 19.25v2" />
	<path d="M4.25 4.25 5.75 5.75" />
	<path d="M18.25 18.25 19.75 19.75" />
	<path d="M2.75 12h2" />
	<path d="M19.25 12h2" />
	<path d="M4.25 19.75 5.75 18.25" />
	<path d="M18.25 5.75 19.75 4.25" />
</svg>
