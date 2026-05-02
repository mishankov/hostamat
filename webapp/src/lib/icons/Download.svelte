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
	<path d="M3.75 13.75v3.75A2.75 2.75 0 0 0 6.5 20.25h11a2.75 2.75 0 0 0 2.75-2.75v-3.75" />
	<path class="arrow" d="M12 4.25v10.25" />
	<polyline class="arrow" points="6.75 9.5 12 14.75 17.25 9.5" />
</svg>

<style>
	.arrow {
		transition: transform 180ms ease-out;
	}

	svg:hover .arrow {
		transform: translateY(1px);
	}

	svg:active .arrow {
		transform: translateY(1.75px);
	}
</style>
