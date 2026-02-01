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
		height: 36px;

		background-color: var(--color-primary-10);
		border: var(--border-regular);
		border-radius: var(--border-radius-regular);
		overflow: hidden;
		box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.progress-fill {
		height: 100%;
		background: linear-gradient(90deg, var(--color-accent-50), var(--color-accent-60));
		border-radius: 0;
		transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		position: relative;
	}

	.progress-fill::after {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: repeating-linear-gradient(
			45deg,
			transparent,
			transparent 10px,
			rgba(255, 255, 255, 0.1) 10px,
			rgba(255, 255, 255, 0.1) 20px
		);
		animation: barberPole 1s linear infinite;
	}

	@keyframes barberPole {
		from {
			background-position: 0 0;
		}
		to {
			background-position: 28px 0;
		}
	}

	p {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		margin: 0;
		z-index: 1;
		font-weight: 700;
		font-size: 13px;
		letter-spacing: 0.05em;
		transition: var(--transition);
		user-select: none;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
	}

	p.dark {
		color: white;
	}
</style>
