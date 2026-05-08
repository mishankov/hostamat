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
	<circle cx="12" cy="6" r="1" fill="currentColor" stroke="none" />
	<circle cx="12" cy="12" r="1" fill="currentColor" stroke="none" />
	<circle cx="12" cy="18" r="1" fill="currentColor" stroke="none" />
</svg>
