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
		<header class="page-header">
			<h1>Твои сайты</h1>
			<p class="site-count">
				{#if sites.length > 0}
					{sites.length} из {maxAmount} сайтов
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
					<span>Загрузка сайтов...</span>
				</div>
			{:then sitesResult}
				{#if sitesResult.data}
					{#each sites as site, index (site.id)}
						<div class="site-card" style="--delay: {index * 50}ms">
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
						<div class="site-card add-new" style="--delay: {sites.length * 50}ms">
							<AddNewSite
								bind:loading
								onsuccess={async () => {
									await invalidateAll();
									updateState();
								}}
							/>
						</div>
					{:else if loading === false}
						<div class="limit-notice">
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

		padding: 0 24px;
	}

	.dashboard {
		width: 100%;
		max-width: 1200px;

		display: flex;
		flex-direction: column;
		gap: 24px;

		padding: 40px 0 60px;
	}

	.page-header {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: 16px;
		flex-wrap: wrap;
	}

	.page-header h1 {
		margin: 0;
	}

	.site-count {
		font-size: 0.875rem;
		color: var(--color-text-muted);
		margin: 0;
	}

	.site-list {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 16px;
	}

	.site-card {
		animation: fadeInUp 0.4s var(--ease-out-expo) forwards;
		animation-delay: var(--delay, 0ms);
		opacity: 0;
	}

	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateY(16px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.loading-state {
		grid-column: 1 / -1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 16px;
		padding: 60px 20px;

		color: var(--color-text-secondary);
		font-size: 0.9375rem;
	}

	.limit-notice {
		grid-column: 1 / -1;
	}

	.add-new {
		min-height: 180px;
	}

	@media (max-width: 600px) {
		main {
			padding: 0 16px;
		}

		.dashboard {
			padding: 24px 0 40px;
		}

		.site-list {
			grid-template-columns: 1fr;
		}
	}
</style>
