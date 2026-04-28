<script lang="ts">
	import Alert from './Alert.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<Alert variant="success" title={header} {message}>
	{#if children}
		{@render children?.()}
	{/if}
</Alert>
