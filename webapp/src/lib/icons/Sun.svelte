<script lang="ts">
	let { size = 100, color = 'var(--color-primary-90)' } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100" aria-hidden="true">
	<circle cx="0" cy="0" r="18" fill="none" stroke={color} stroke-width="7" />
	<line x1="0" y1="-38" x2="0" y2="-28" stroke={color} stroke-width="7" stroke-linecap="round" />
	<line x1="0" y1="38" x2="0" y2="28" stroke={color} stroke-width="7" stroke-linecap="round" />
	<line x1="-38" y1="0" x2="-28" y2="0" stroke={color} stroke-width="7" stroke-linecap="round" />
	<line x1="38" y1="0" x2="28" y2="0" stroke={color} stroke-width="7" stroke-linecap="round" />
	<line x1="-27" y1="-27" x2="-20" y2="-20" stroke={color} stroke-width="7" stroke-linecap="round" />
	<line x1="27" y1="27" x2="20" y2="20" stroke={color} stroke-width="7" stroke-linecap="round" />
	<line x1="27" y1="-27" x2="20" y2="-20" stroke={color} stroke-width="7" stroke-linecap="round" />
	<line x1="-27" y1="27" x2="-20" y2="20" stroke={color} stroke-width="7" stroke-linecap="round" />
</svg>
