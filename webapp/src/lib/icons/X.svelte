<script lang="ts">
	type IconProps = {
		size?: number | string;
		animated?: boolean;
		color?: string;
		class?: string;
	};

	let {
		size = 100,
		animated = false,
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
	<line class={[animated && 'animated']} x1="6" y1="6" x2="18" y2="18" />
	<line class={[animated && 'animated']} x1="18" y1="6" x2="6" y2="18" />
</svg>

<style>
	line {
		transform-origin: center;
		transition: transform 180ms ease-out;
	}

	svg:hover line.animated {
		transform: rotate(90deg);
	}
</style>
