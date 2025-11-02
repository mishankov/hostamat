<script lang="ts">
	let { size = 100, animated = false, color = 'var(--color-primary)' } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100">
	<line
		class={[animated && 'animated']}
		x1="-30"
		y1="-30"
		x2="30"
		y2="30"
		stroke={color}
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>

	<line
		class={[animated && 'animated']}
		x1="30"
		y1="-30"
		x2="-30"
		y2="30"
		stroke={color}
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>
</svg>

<style>
	line {
		transition: all 500ms;
	}

	svg:hover line.animated {
		transform: rotate(90deg);
	}
</style>
