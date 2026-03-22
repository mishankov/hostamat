<script lang="ts">
	let { summary = '', children = undefined } = $props();
</script>

<details>
	<summary>{summary}</summary>
	{@render children()}
</details>

<style>
	summary {
		padding-bottom: 6px;
		font-weight: 650;
		color: var(--color-accent-text);
		text-decoration: underline;
		text-underline-offset: 3px;
		cursor: pointer;
	}

	summary::marker {
		content: none;
	}
</style>
