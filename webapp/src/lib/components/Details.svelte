<script lang="ts">
	let { summary = '', children = undefined } = $props();
</script>

<details class="details-animated">
	<summary>{summary}</summary>
	<div class="details-content">
		{@render children()}
	</div>
</details>

<style>
	.details-animated {
		display: grid;
		grid-template-rows: auto 0fr;
		transition: grid-template-rows 250ms var(--ease-out-quart);
	}

	.details-animated[open] {
		grid-template-rows: auto 1fr;
	}

	summary {
		padding-bottom: var(--space-xs);
		font-weight: 650;
		color: var(--color-accent-text);
		text-decoration: underline;
		text-underline-offset: 3px;
		cursor: pointer;
		list-style: none;
	}

	summary::marker,
	summary::-webkit-details-marker {
		display: none;
	}

	.details-content {
		overflow: hidden;
		opacity: 0;
		transform: translateY(-4px);
		transition:
			opacity 200ms var(--ease-out-quart),
			transform 200ms var(--ease-out-quart);
	}

	.details-animated[open] .details-content {
		opacity: 1;
		transform: translateY(0);
	}
</style>
