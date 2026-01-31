<script lang="ts">
	let {
		orientation = 'horizontal'
	}: {
		orientation?: 'horizontal' | 'vertical';
	} = $props();
</script>

<div class="separator {orientation}"></div>

<style>
	.separator {
		border: 0;
		background: linear-gradient(
			90deg,
			transparent,
			var(--color-border),
			transparent
		);
	}

	.horizontal {
		width: 100%;
		height: 1px;
	}

	.vertical {
		width: 1px;
		height: 100%;
		background: linear-gradient(
			180deg,
			transparent,
			var(--color-border),
			transparent
		);
	}
</style>
