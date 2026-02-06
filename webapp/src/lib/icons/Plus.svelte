<script lang="ts">
	let { size = 100 } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100">
	<line
		x1="0"
		y1="-30"
		x2="0"
		y2="30"
		stroke="var(--color-primary-90)"
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>

	<line
		x1="-30"
		y1="0"
		x2="30"
		y2="0"
		stroke="var(--color-primary-90)"
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>
</svg>

<style>
	line {
		transition: all 260ms ease-out;
	}

	svg:hover line {
		transition-duration: 90ms;
		transform: rotate(90deg);
	}
</style>
