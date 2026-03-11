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
		<Link href="https://{site.domain}" target="_blank" alwaysUnderline>{site.domain}</Link>

		<p class="connected">Подключен</p>

		<div style="display: inline;">
			<Link href={resolve('/dashboard/sites')} alwaysUnderline>Статический сайт</Link> -
			<span>{site.name}</span>
		</div>
	</div>
</Card>

<style>
	.domain-card-content {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 0.75rem;
		height: 100%;
		padding: 1rem;
		min-width: 0;
	}

	.connected {
		color: var(--color-confirm-50);
		font-weight: bold;
	}
</style>
