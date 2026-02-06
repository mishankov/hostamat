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

<main class="page-shell">
	<div class="page-content dashboard">
		<div class="headline">
			<h1>Твои сайты</h1>
			<p>Управляйте публикациями, обновляйте файлы и следите за лимитом альфы.</p>
		</div>

		{#if errorMessage}
			<Error message={errorMessage} />
		{/if}

		<div class="site-list">
			{#await data.sitesResp}
				<Spinner size={50} />
			{:then sitesResult}
				{#if sitesResult.data}
					{#each sites as site (site.id)}
						<Card border style="min-height: 170px;">
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
		padding-bottom: 2rem;
	}

	.dashboard {
		gap: 1rem;
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
		display: flex;
		flex-wrap: wrap;
		gap: 0.8rem;
	}
</style>
