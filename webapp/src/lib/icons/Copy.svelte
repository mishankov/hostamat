<script lang="ts">
	let { size = 100, color = 'currentColor', background = 'var(--color-background-object)' } = $props();
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
		stroke-width="6"
		stroke={color}
		style="--bg: {background}"
	/>

	<rect
		class="front"
		x="5"
		y="25"
		rx="10"
		ry="10"
		width="60"
		height="60"
		stroke-width="6"
		stroke={color}
		style="--bg: {background}"
	/>
</svg>

<style>
	rect {
		fill: var(--bg);
		transition: all var(--transition);
	}

	svg:hover rect.back {
		transform: translateX(5px) translateY(-5px);
	}

	svg:active rect.front {
		fill: currentColor;
	}
</style>
