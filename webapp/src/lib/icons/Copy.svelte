<script lang="ts">
	let { size = 100, background = 'var(--color-background-object)' } = $props();

	let strokeWidth = 7;
</script>

<svg width={size} viewBox="0 0 100 100">
	<rect
		class="back"
		x="30"
		y="10"
		rx="10"
		ry="10"
		width="60"
		height="60"
		stroke-width={strokeWidth}
		style="--color-background: {background}"
	/>

	<rect
		class="front"
		x="5"
		y="25"
		rx="10"
		ry="10"
		width="60"
		height="60"
		stroke-width={strokeWidth}
		style="--color-background: {background}"
	/>
</svg>

<style>
	rect {
		stroke: var(--color-primary-90);
		transition: all 500ms;
	}

	rect {
		fill: var(--color-background);
	}

	svg:hover rect {
		stroke: var(--color-primary-90);
	}

	svg:hover rect.back {
		transform: translateX(5px) translateY(-5px);
	}

	svg:active rect.front {
		fill: var(--color-primary-90);
	}
</style>
