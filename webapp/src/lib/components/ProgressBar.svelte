<script lang="ts">
	let { uploadProgress = $bindable(0) } = $props();
</script>

<div class="progress-bar">
	<p class={[uploadProgress > 50 && 'dark']}>{Math.round(uploadProgress)}%</p>
	<div class="progress-fill" style="width: {uploadProgress}%"></div>
</div>

<style>
	.progress-bar {
		position: relative;
		width: 100%;
		height: 34px;

		background-color: var(--color-neutral-surface);

		border: var(--border-regular);
		border-radius: 0;
	}

	.progress-fill {
		height: 100%;
		background: var(--color-accent-strong);
		border-radius: 0;
		transition: var(--transition);
	}

	p {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		margin: 0;
		z-index: 1;

		transition: var(--transition);
		user-select: none;
	}

	p.dark {
		color: var(--color-neutral-strong);
	}
</style>
