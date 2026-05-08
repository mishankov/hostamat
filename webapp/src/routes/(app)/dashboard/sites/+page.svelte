<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';

	import { Alert } from '$lib/components';
	import { Spinner } from '$lib/icons';

	import { type Site } from '$lib/types';
	import AddNewSite from './AddNewSite.svelte';
	import SiteCard from './SiteCard.svelte';

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
	<meta
		name="description"
		content="Управляйте статическими сайтами: создание, загрузка файлов и удаление."
	/>
</svelte:head>

<div class="dashboard-page">
	<div class="headline">
		<h1>Статические сайты</h1>
		<div class="headline-meta">
			<p>Добавляйте и загружайте статические сайты</p>
			{#if !loading && maxAmount > 0}
				<span class="site-counter">Количество сайтов: {sites.length}/{maxAmount}</span>
			{/if}
		</div>
	</div>

	{#if errorMessage}
		<Alert variant="error" message={errorMessage} />
	{/if}

	<div class="site-grid">
		{#await data.sitesResp}
			<div class="loading-state" aria-live="polite">
				<Spinner size={50} />
			</div>
		{:then sitesResult}
			{#if sitesResult.data}
				{#each sites ?? [] as site (site.id)}
					<SiteCard
						{site}
						bind:loading
						bind:errorMessage
						onsiteupdate={async () => {
							await invalidateAll();
							updateState();
						}}
					/>
				{/each}

				{#if loading === false && sites && sites.length >= maxAmount}
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

	{#if !sites || sites.length < maxAmount}
		<div class="add-action">
			<AddNewSite
				bind:loading
				onsuccess={async () => {
					await invalidateAll();
					updateState();
				}}
			/>
		</div>
	{/if}
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
	}

	.headline h1 {
		font-weight: 800;
	}

	.headline p {
		max-width: 64ch;
		color: var(--color-text-secondary);
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
		font-weight: 500;
		color: var(--color-text-tertiary);
		white-space: nowrap;
	}

	.site-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
		gap: var(--space-lg);
		min-width: 0;
	}

	.loading-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 10.625rem;
		grid-column: 1 / -1;
	}

	.limit-info-wrap {
		min-height: 6.875rem;
		display: flex;
		align-items: center;
		padding: var(--space-lg) 0;
		grid-column: 1 / -1;
	}

	.add-action {
		padding-top: var(--space-lg);
	}

	@media (max-width: 820px) {
		.limit-info-wrap {
			min-height: auto;
		}
	}
</style>
