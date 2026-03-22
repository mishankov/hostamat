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

		border-radius: 0;
		border: var(--border-regular);

		padding: 0.24rem 0.64rem;

		font-size: 0.74rem;
		font-weight: 600;
		letter-spacing: 0.03em;

		text-align: center;
		user-select: none;

		color: var(--color-accent-text);
		background-color: var(--color-accent-surface);
		border-color: var(--color-accent-border);
	}
</style>
