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
	<path class="pencil" d="M5.25 16.25 15.75 5.75a2.3 2.3 0 0 1 3.25 3.25L8.5 19.5l-5 1.25z" />
	<path class="pencil" d="M14.4 7.1 17.65 10.35" />
	<path class="pencil" d="M5.25 16.25 8.5 19.5" />
</svg>

<style>
	.pencil {
		transition: transform 180ms ease-out;
	}

	svg:hover .pencil {
		transform: translate(0.5px, -0.5px);
	}
</style>
