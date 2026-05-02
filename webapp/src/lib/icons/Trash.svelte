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
	<g class="lid">
		<path d="M8.25 7V5.5a1.75 1.75 0 0 1 1.75-1.75h4A1.75 1.75 0 0 1 15.75 5.5V7" />
		<path d="M4 7h16" />
	</g>
	<path d="M6.25 9.25 7.25 20.25h9.5l1-11" />
	<path d="M10 11.5v5.75" />
	<path d="M14 11.5v5.75" />
</svg>

<style>
	.lid {
		transform-box: fill-box;
		transform-origin: center;
		transition: transform 180ms ease-out;
	}

	svg:hover .lid {
		transform: translate(0.25px, -1px) rotate(6deg);
	}

	svg:active .lid {
		transform: translate(0.5px, -1.5px) rotate(10deg);
	}
</style>
