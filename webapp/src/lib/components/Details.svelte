<script lang="ts">
	let { summary = '', children = undefined } = $props();
</script>

<details>
	<summary>{summary}</summary>
	{@render children()}
</details>

<style>
	summary {
		padding-bottom: 5px;
		text-decoration: underline;
	}

	summary::marker {
		content: none;
	}
</style>
