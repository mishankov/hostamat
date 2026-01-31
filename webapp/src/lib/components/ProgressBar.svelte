<script lang="ts">
	let { uploadProgress = $bindable(0) } = $props();
</script>

<div class="progress-bar">
	<div class="progress-track">
		<div class="progress-fill" style="width: {uploadProgress}%">
			<div class="progress-glow"></div>
		</div>
	</div>
	<span class="progress-text">{Math.round(uploadProgress)}%</span>
</div>

<style>
	.progress-bar {
		display: flex;
		align-items: center;
		gap: 12px;
		width: 100%;
	}

	.progress-track {
		position: relative;
		flex: 1;
		height: 8px;
		background: var(--color-primary-15);
		border-radius: 4px;
		overflow: hidden;
	}

	.progress-fill {
		position: relative;
		height: 100%;
		background: linear-gradient(90deg, var(--color-accent-50), var(--color-accent-60));
		border-radius: 4px;
		transition: width 0.3s ease;
	}

	.progress-glow {
		position: absolute;
		inset: 0;
		background: linear-gradient(
			90deg,
			transparent 0%,
			rgba(255, 255, 255, 0.4) 50%,
			transparent 100%
		);
		animation: shimmer 1.5s infinite;
	}

	.progress-text {
		font-family: 'JetBrains Mono', monospace;
		font-size: 13px;
		font-weight: 500;
		color: var(--color-accent-60);
		min-width: 42px;
		text-align: right;
	}
</style>
