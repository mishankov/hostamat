<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<div class="toolbar-shell">
	<div class="toolbar">
		{@render children?.()}
	</div>
</div>

<style>
	.toolbar-shell {
		width: fit-content;
		min-width: 0;
	}

	.toolbar {
		display: flex;
		padding: var(--space-xs);
		gap: var(--space-xs);
		border: var(--border-width) var(--border-style) var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-background-object);
		box-shadow: var(--shadow-sm);
	}

	@media (max-width: 820px) {
		.toolbar-shell {
			width: 100%;
		}

		.toolbar {
			justify-content: space-between;
			gap: var(--space-xs);
		}
	}
</style>
