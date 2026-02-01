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
		<h1>Твои сайты</h1>

		{#if errorMessage}
			<Error message={errorMessage} />
		{/if}

		<div class="site-list">
			{#await data.sitesResp}
				<Spinner size={50} />
			{:then sitesResult}
				{#if sitesResult.data}
					{#each sites as site (site.id)}
						<Card border style="height: 150px;">
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
						<div style="height: 110px; display: flex; align-items: center;">
							<Info header="Максимальное количество сайтов в альфа версии">
								Для добавления нового сайта, удалите один из существующих
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
		min-height: calc(100vh - 85px);

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: start;
	}

	.dashboard {
		display: flex;
		flex-direction: column;
		gap: 24px;

		width: clamp(500px, 85%, 1400px);
		padding: 40px 20px;
	}

	.dashboard h1 {
		font-size: clamp(32px, 5vw, 48px);
		margin-bottom: 8px;
		animation: slideInLeft 0.5s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes slideInLeft {
		from {
			transform: translateX(-30px);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	.site-list {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
		gap: 20px;
		animation: fadeInUp 0.6s cubic-bezier(0.4, 0, 0.2, 1) 0.2s backwards;
	}

	@keyframes fadeInUp {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}
</style>
