<script lang="ts">
	let { size = 100, color = 'var(--color-primary-90)' } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100" aria-hidden="true">
	<path
		d="M 0 -38 L 42 34 H -42 Z"
		fill="transparent"
		stroke={color}
		stroke-width="7"
		stroke-linejoin="round"
	/>

	<line x1="0" y1="-16" x2="0" y2="8" stroke={color} stroke-width="8" stroke-linecap="round" />

	<line x1="0" y1="24" x2="0" y2="24" stroke={color} stroke-width="8" stroke-linecap="round" />
</svg>
