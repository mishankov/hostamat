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
		color: var(--color-primary-90);
		text-decoration: underline;
		text-underline-offset: 3px;
		cursor: pointer;
	}

	summary::marker {
		content: none;
	}
</style>
