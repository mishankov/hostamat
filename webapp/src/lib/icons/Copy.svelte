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
	<path
		class="back"
		d="M8.5 8V6.5A2 2 0 0 1 10.5 4.5H17A2 2 0 0 1 19 6.5v8A2 2 0 0 1 17 16.5h-1.5"
	/>
	<rect
		class="front"
		x="5"
		y="8"
		width="10.5"
		height="12"
		rx="2"
		fill="var(--icon-surface, var(--color-background-object))"
	/>
</svg>

<style>
	.back,
	.front {
		transition:
			transform 180ms ease-out,
			fill 120ms ease-out;
	}

	svg:hover .back {
		transform: translate(0.5px, -0.5px);
	}

	svg:active .front {
		fill: currentColor;
	}
</style>
