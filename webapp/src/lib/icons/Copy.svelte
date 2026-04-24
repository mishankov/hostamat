<script lang="ts">
	let { size = 100 } = $props();

	let strokeWidth = 7;
</script>

<svg width={size} viewBox="0 0 100 100">
	<path class="back" d="M33 10 H69 Q83 10 83 24 V64" stroke-width={strokeWidth} />

	<rect
		class="front"
		x="5"
		y="25"
		rx="10"
		ry="10"
		width="60"
		height="60"
		stroke-width={strokeWidth}
	/>
</svg>

<style>
	path,
	rect {
		stroke: var(--color-primary-90);
		transition: all 260ms ease-out;
	}

	path {
		fill: none;
		stroke-linecap: round;
		stroke-linejoin: round;
	}

	rect {
		fill: transparent;
	}

	svg:hover path,
	svg:hover rect {
		transition-duration: 90ms;
		stroke: var(--color-primary-90);
	}

	svg:hover .back {
		transition-duration: 90ms;
		transform: translateX(1px) translateY(-1px);
	}

	svg:active rect.front {
		fill: var(--color-primary-90);
	}
</style>
