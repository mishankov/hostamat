<script lang="ts">
	import { resolve } from '$app/paths';
	import { Card, Link } from '$lib/components';

	import type { Site } from '$lib/types';

	interface Props {
		site: Site;
	}

	let { site }: Props = $props();
</script>

<Card border>
	<div class="domain-card-content">
		<div class="data-row">
			<span class="data-label">domain:</span>
			<Link href="https://{site.domain}" target="_blank" alwaysUnderline>{site.domain}</Link>
		</div>
		<div class="data-row">
			<span class="data-label">site:</span>
			<Link href={resolve('/dashboard/sites')} alwaysUnderline>{site.name}</Link>
		</div>
	</div>
</Card>

<style>
	.domain-card-content {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: var(--space-sm);
		height: 100%;
		padding: var(--space-lg);
		min-width: 0;
	}
</style>
