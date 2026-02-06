<script lang="ts">
	let { size = 100 } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100">
	<polyline
		points="
      -20 0
      -40 0
      -40 30
      40 30
      40 0
      20 0
    "
		stroke="var(--color-primary-90)"
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
		fill="none"
	/>

	<polyline
		class="arrow"
		points="
      0 15
      0 -30
      10 -20
      0 -30
      -10 -20
    "
		stroke="var(--color-primary-90)"
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
		fill="none"
	/>
</svg>

<style>
	.arrow {
		transition: all 140ms ease-out;
	}
	svg:hover .arrow {
		transform: translateY(-5px);
	}

	svg:active .arrow {
		transform: translateY(-8px);
	}
</style>
