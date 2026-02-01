<script lang="ts">
	type Props = {
		text: string;
	};

	let { text = '' }: Props = $props();
</script>

<span class="primary">{text}</span>

<style>
	span {
		width: fit-content;
		display: inline-block;

		border-radius: var(--border-radius-regular);
		border: 2px solid var(--color-accent-60);

		padding: 6px 12px;

		font-size: 11px;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;

		text-align: center;
		user-select: none;

		color: var(--color-accent-60);
		background-color: var(--color-accent-10);

		position: relative;
		box-shadow: 0 2px 0 0 var(--color-accent-70);
		animation: pulse 2s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.02);
		}
	}
</style>
