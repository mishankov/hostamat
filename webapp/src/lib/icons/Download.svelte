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

	<!-- 0 15
      0 -30
      15 -20
      0 -30
      -15 -20 -->

	<polyline
		class="arrow"
		points="
      0 -30
      0 15
      10 5
      0 15
      -10 5
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
		transition: all 260ms ease-out;
	}
	svg:hover .arrow {
		transition-duration: 90ms;
		transform: translateY(5px);
	}

	svg:active .arrow {
		transform: translateY(8px);
	}
</style>
