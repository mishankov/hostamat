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
	<line x1="12" y1="5" x2="12" y2="19" />
	<line x1="5" y1="12" x2="19" y2="12" />
</svg>

<style>
	line {
		transform-origin: center;
		transition: transform 180ms ease-out;
	}

	svg:hover line {
		transform: rotate(90deg);
	}
</style>
