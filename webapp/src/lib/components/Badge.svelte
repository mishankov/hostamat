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

		border-radius: 999px;
		border: 1px solid var(--color-accent-border);

		padding: 0.2rem 0.65rem;

		font-size: 0.72rem;
		font-weight: 600;
		letter-spacing: 0.02em;

		text-align: center;
		user-select: none;

		color: var(--color-accent-text);
		background-color: var(--color-accent-surface);
		border-color: var(--color-accent-border);
	}
</style>
