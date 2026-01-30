<script lang="ts">
	let { size = 100, color = 'var(--color-accent-50)' } = $props();
</script>

<svg width={size} viewBox="-50 -50 100 100" class="spinner">
	<circle
		cx="0"
		cy="0"
		r="40"
		stroke={color}
		stroke-width="6"
		stroke-dasharray="60 188"
		stroke-dashoffset="0"
		stroke-linecap="round"
		fill="none"
	/>
</svg>

<style>
	.spinner {
		animation: spinnerRotate 1s linear infinite;
	}

	circle {
		animation: spinnerDash 1.4s ease-in-out infinite;
	}

	@keyframes spinnerRotate {
		100% {
			transform: rotate(360deg);
		}
	}

	@keyframes spinnerDash {
		0% {
			stroke-dasharray: 10 188;
			stroke-dashoffset: 0;
		}
		50% {
			stroke-dasharray: 80 188;
			stroke-dashoffset: -40;
		}
		100% {
			stroke-dasharray: 10 188;
			stroke-dashoffset: -125;
		}
	}
</style>
