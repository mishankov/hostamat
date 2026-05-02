<script lang="ts">
	type IconProps = {
		size?: number | string;
		color?: string;
		class?: string;
	};

	let {
		size = 100,
		color = 'var(--color-primary-90)',
		class: className = ''
	}: IconProps = $props();
</script>

<svg
	class={['spinner', className]}
	width={size}
	height={size}
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="1.75"
	stroke-linecap="round"
	stroke-linejoin="round"
	style:color
	aria-hidden="true"
>
	<circle cx="12" cy="12" r="8.25" stroke-dasharray="18 34" />
</svg>

<style>
	circle {
		transform-origin: center;
		animation: rotate 900ms linear infinite;
	}

	@keyframes rotate {
		to {
			transform: rotate(360deg);
		}
	}
</style>
