<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';

	import { Card, Error, Info } from '$lib/components';
	import { Spinner } from '$lib/icons';

	import { type Site } from '$lib/types';
	import SiteContent from './SiteContent.svelte';
	import AddNewSite from './AddNewSite.svelte';

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
	<title>Дашборд | Хостамат</title>
</svelte:head>

<main>
	<div class="dashboard">
		<header class="dashboard-header">
			<h1>Твои сайты</h1>
			<p class="sites-count">
				{#if !loading}
					<span class="count">{sites.length}</span> из <span class="count">{maxAmount}</span>
				{/if}
			</p>
		</header>

		{#if errorMessage}
			<Error message={errorMessage} />
		{/if}

		<div class="site-list">
			{#await data.sitesResp}
				<div class="loading-state">
					<Spinner size={40} />
					<p>Загружаем сайты...</p>
				</div>
			{:then sitesResult}
				{#if sitesResult.data}
					{#each sites as site, i (site.id)}
						<div class="site-card" style="animation-delay: {i * 0.05}s">
							<Card border>
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
						<div class="site-card" style="animation-delay: {sites.length * 0.05}s">
							<AddNewSite
								bind:loading
								onsuccess={async () => {
									await invalidateAll();
									updateState();
								}}
							/>
						</div>
					{:else if loading === false}
						<div class="max-reached">
							<Info header="Достигнут лимит сайтов">
								Для добавления нового сайта удалите один из существующих
							</Info>
						</div>
					{/if}
				{:else}
					<Error message="Ошибка при загрузке сайтов: {sitesResult.error}" />
				{/if}
			{:catch error}
				<Error message="Ошибка при загрузке сайтов: {error.message}" />
			{/await}
		</div>
	</div>
</main>

<style>
	main {
		width: 100%;
		min-height: calc(100vh - 140px);

		display: flex;
		justify-content: center;
		padding: 40px 24px;
	}

	.dashboard {
		display: flex;
		flex-direction: column;
		gap: 24px;

		width: 100%;
		max-width: 1200px;
	}

	.dashboard-header {
		display: flex;
		align-items: baseline;
		gap: 16px;
	}

	.dashboard-header h1 {
		font-size: 32px;
		font-weight: 700;
	}

	.sites-count {
		font-size: 14px;
		color: var(--color-text-muted);
	}

	.sites-count .count {
		font-family: 'JetBrains Mono', monospace;
		color: var(--color-accent-60);
		font-weight: 600;
	}

	.site-list {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 16px;
	}

	.site-card {
		animation: slideUp 0.4s ease-out both;
	}

	.loading-state {
		grid-column: 1 / -1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 16px;
		padding: 60px 20px;
		color: var(--color-text-muted);
	}

	.max-reached {
		grid-column: 1 / -1;
		max-width: 500px;
	}

	@media (max-width: 600px) {
		main {
			padding: 24px 16px;
		}

		.dashboard-header {
			flex-direction: column;
			gap: 4px;
		}

		.dashboard-header h1 {
			font-size: 24px;
		}
	}
</style>
