<script lang="ts">
	let { uploadProgress = $bindable(0) } = $props();
	const isFilledBehindLabel = $derived(uploadProgress >= 50);
</script>

<div class="progress-bar" class:is-complete={uploadProgress >= 100}>
	<p class={[!isFilledBehindLabel && 'on-surface', isFilledBehindLabel && 'on-fill']}>
		{Math.round(uploadProgress)}%
	</p>
	<div class="progress-fill" style="width: {uploadProgress}%"></div>
</div>

<style>
	.progress-bar {
		position: relative;
		width: 100%;
		height: 2.125rem;

		background-color: var(--color-neutral-surface);

		border: var(--border-regular);
		border-radius: var(--radius-md);
		overflow: hidden;
	}

	.progress-fill {
		height: 100%;
		background: var(--color-accent-strong);
		border-radius: 0;
		transition: width 200ms var(--ease-out-quart);
	}

	.progress-bar.is-complete .progress-fill {
		animation: progress-glow 600ms var(--ease-out-quart) both;
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
		font-weight: 700;
		text-shadow: 0 1px 0 hsl(0 0% 100% / 0.14);
	}

	p.on-surface {
		color: var(--color-neutral-strong);
	}

	p.on-fill {
		color: var(--color-background-object);
		text-shadow: 0 1px 0 hsl(0 0% 0% / 0.18);
	}
</style>
