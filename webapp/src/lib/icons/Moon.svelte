<script lang="ts">
	let { size = 100, color = 'var(--color-primary-90)' } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100" aria-hidden="true">
	<path
		d="M 18 -30 A 32 32 0 1 0 18 30 A 23 23 0 1 1 18 -30 Z"
		fill="none"
		stroke={color}
		stroke-width="8"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>
</svg>
