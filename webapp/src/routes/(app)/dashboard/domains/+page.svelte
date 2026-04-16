<script lang="ts">
	import type { PageProps } from './$types';

	import { Alert, Info } from '$lib/components';
	import { Spinner } from '$lib/icons';

	import DomainCard from './DomainCard.svelte';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Домены | Хостамат</title>
</svelte:head>

<div class="domains-page">
	<div class="headline">
		<h1>Домены</h1>
		<p>Просматривайте и управляете доменами и поддоменами</p>
	</div>

	<div class="domains-list">
		{#await data.sitesResp}
			<div class="loading-state" aria-live="polite">
				<Spinner size={50} />
			</div>
		{:then sitesResult}
			{#if sitesResult.data}
				{#each sitesResult.data.sites as site (site.id)}
					<DomainCard {site} />
				{/each}
			{:else}
				<div class="domains-alert">
					<Alert variant="error" message="Ошибка при загрузке доменов: {sitesResult.error}" />
				</div>
			{/if}
		{:catch error}
			<div class="domains-alert">
				<Alert variant="error" message="Ошибка при загрузке доменов: {error.message}" />
			</div>
		{/await}
	</div>

	<div class="domains-alert">
		<Info header="Подключение собственных доменов">
			<p>Поддомены для ваших статических сайтов уже создаются автоматически и отображаются выше.</p>
			<p>
				Управление подключением собственных доменов к существующим сайтам появится в этом разделе
				позже.
			</p>
		</Info>
	</div>
</div>

<style>
	.domains-page {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 0;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		padding-bottom: var(--space-md);
		border-bottom: 1px dashed color-mix(in hsl, var(--color-border) 55%, transparent);
	}

	.headline p {
		max-width: 64ch;
	}

	.domains-list {
		--site-tile-width: 320px;

		display: flex;
		flex-wrap: wrap;
		gap: var(--space-md);
		align-items: stretch;
	}

	.loading-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 170px;
		width: var(--site-tile-width);
	}

	.domains-alert {
		width: 100%;
	}

	@media (max-width: 820px) {
		.domains-list {
			--site-tile-width: 100%;
		}
	}
</style>
