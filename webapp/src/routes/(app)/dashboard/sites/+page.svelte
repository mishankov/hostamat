<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';

	import { Alert, Card } from '$lib/components';
	import { Spinner } from '$lib/icons';

	import { type Site } from '$lib/types';
	import AddNewSite from './AddNewSite.svelte';
	import SiteContent from './SiteContent.svelte';

	let { data }: PageProps = $props();
	let sites = $state<Site[]>([]);
	let maxAmount = $state(0);
	let loading = $state(true);
	let errorMessage = $state('');

	const updateState = async () => {
		const sitesResp = await data.sitesResp;

		if (sitesResp.data) {
			sites = sitesResp.data.sites ?? [];
			maxAmount = sitesResp.data.maxAmount;
			errorMessage = '';
		} else {
			sites = [];
			maxAmount = 0;
			errorMessage = sitesResp.error;
		}

		loading = false;
	};
	updateState();
</script>

<svelte:head>
		<title>Статические сайты | Хостамат</title>
		<meta name="description" content="Управляйте статическими сайтами: создание, загрузка файлов и удаление." />
	</svelte:head>

<div class="dashboard-page">
	<div class="headline">
		<h1>Статические сайты</h1>
		<div class="headline-meta">
			<p>Добавляйте и загружайте статические сайты</p>
			{#if !loading && maxAmount > 0}
				<span class="site-counter">sites: {sites.length}/{maxAmount}</span>
			{/if}
		</div>
	</div>

	{#if errorMessage}
		<Alert variant="error" message={errorMessage} />
	{/if}

	<div class="site-list">
		{#await data.sitesResp}
			<div class="loading-state" aria-live="polite">
				<Spinner size={50} />
			</div>
		{:then sitesResult}
			{#if sitesResult.data}
				{#each sites ?? [] as site (site.id)}
					<div class="site-card-shell">
						<Card border style="height: 100%;">
							<SiteContent
								{site}
								bind:loading
								bind:errorMessage
								onsiteupdate={async () => {
									await invalidateAll();
									updateState();
								}}
							/>
						</Card>
					</div>
				{/each}

				{#if !sites || sites.length < maxAmount}
					<AddNewSite
						bind:loading
						onsuccess={async () => {
							await invalidateAll();
							updateState();
						}}
					/>
				{:else if loading === false}
					<div class="limit-info-wrap">
						<Alert variant="info" title="Максимальное количество сайтов в альфа версии">
							Для добавления нового сайта, удалите один из существующих
						</Alert>
					</div>
				{/if}
			{:else}
				<Alert variant="error" message="Ошибка при загрузке сайтов: {sitesResult.error}" />
			{/if}
		{:catch error}
			<Alert variant="error" message="Ошибка при загрузке сайтов: {error.message}" />
		{/await}
	</div>
</div>

<style>
	.dashboard-page {
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

	.headline-meta {
		display: flex;
		align-items: baseline;
		gap: var(--space-lg);
		flex-wrap: wrap;
	}

	.site-counter {
		font-family: var(--font-body);
		font-size: 0.82rem;
		color: var(--color-text-secondary);
		white-space: nowrap;
	}

	.site-list {
		--site-tile-width: 20rem;

		display: flex;
		flex-wrap: wrap;
		gap: var(--space-md);
		align-items: stretch;
	}

	.loading-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 10.625rem;
		width: var(--site-tile-width);
	}

	.limit-info-wrap {
		min-height: 6.875rem;
		display: flex;
		align-items: center;
	}

	.site-card-shell {
		display: flex;
		min-height: 10.625rem;
		min-width: 0;
	}

	@media (max-width: 820px) {
		.site-list {
			--site-tile-width: 100%;
		}

		.limit-info-wrap {
			min-height: auto;
		}
	}
</style>
