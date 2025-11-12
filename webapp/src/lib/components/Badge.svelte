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

		border-radius: var(--border-radius-regular);

		padding: 4px 8px;

		font-size: small;
		font-weight: 600;

		text-align: center;
		user-select: none;

		color: var(--color-action);
		background-color: var(--color-action-dim);
	}
</style>
