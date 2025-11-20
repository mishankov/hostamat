<script lang="ts">
	let { size = 100 } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100">
	<polyline
		points="
		0 -30
		-30 -30
		-30 30
		30 30
		30 0
    "
		stroke="var(--color-primary)"
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
		fill="none"
	/>

	<polyline
		class="stick"
		points="
		35 -35
		5 -5
    "
		stroke="var(--color-primary)"
		stroke-width="7"
		stroke-linecap="round"
		stroke-linejoin="round"
		fill="none"
	/>
</svg>

<style>
	.stick {
		transition: all 500ms;
	}
	svg:hover .stick {
		transform: translateY(-5px);
	}

	svg:active .stick {
		transform: translateY(3px);
	}
</style>
