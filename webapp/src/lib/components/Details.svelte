<script lang="ts">
	let { summary = '', children = undefined } = $props();
</script>

<details>
	<summary>{summary}</summary>
	{@render children()}
</details>

<style>
	details {
		border: var(--border-regular);
		border-radius: var(--border-radius-regular);
		padding: 12px 16px;
		background: var(--color-primary-10);
		transition: var(--transition);
	}

	details[open] {
		background: var(--color-primary-15);
	}

	summary {
		padding-bottom: 8px;
		font-weight: 600;
		cursor: pointer;
		user-select: none;
		list-style: none;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	summary::before {
		content: '▶';
		display: inline-block;
		transition: transform 0.2s;
		color: var(--color-accent-50);
		font-size: 12px;
	}

	details[open] summary::before {
		transform: rotate(90deg);
	}

	summary::marker {
		content: none;
	}

	summary::-webkit-details-marker {
		display: none;
	}
</style>
