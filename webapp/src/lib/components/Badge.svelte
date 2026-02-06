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

		border-radius: calc(var(--border-radius-regular) - 1px);
		border: var(--border-regular);

		padding: 0.24rem 0.64rem;

		font-size: 0.74rem;
		font-weight: 600;
		letter-spacing: 0.03em;

		text-align: center;
		user-select: none;

		color: var(--color-accent-80);
		background-color: var(--color-accent-20);
		border-color: var(--color-accent-40);
	}
</style>
