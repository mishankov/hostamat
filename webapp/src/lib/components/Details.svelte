<script lang="ts">
	import Link from './Link.svelte';

	let { summary = '', children = undefined } = $props();
</script>

<details>
	<summary><Link alwaysUnderline>{summary}</Link></summary>
	{@render children()}
</details>

<style>
	summary {
		padding-bottom: 5px;
	}

	summary::marker {
		content: none;
	}
</style>
