<script lang="ts">
	let { size = 100, color = 'var(--color-text-primary)' } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100">
	<line
		x1="0"
		y1="-30"
		x2="0"
		y2="30"
		stroke={color}
		stroke-width="6"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>

	<line
		x1="-30"
		y1="0"
		x2="30"
		y2="0"
		stroke={color}
		stroke-width="6"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>
</svg>

<style>
	line {
		transition: all var(--transition);
	}

	svg:hover line {
		transform: rotate(90deg);
	}
</style>
