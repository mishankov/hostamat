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
			sites = sitesResp.data.sites;
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
</svelte:head>

<div class="dashboard-page">
	<div class="headline">
		<h1>Статические сайты</h1>
		<p>Добавляйте и загружайте статические сайты</p>
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
				{#each sites as site (site.id)}
					<Card border class="site-card">
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
		gap: 1rem;
		min-width: 0;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	.headline p {
		max-width: 64ch;
	}

	.site-list {
		--site-tile-width: 320px;

		display: flex;
		flex-wrap: wrap;
		gap: 0.8rem;
		align-items: stretch;
	}

	.loading-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 170px;
		width: var(--site-tile-width);
	}

	.limit-info-wrap {
		min-height: 110px;
		display: flex;
		align-items: center;
	}

	:global(.site-card) {
		min-height: 170px;
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
