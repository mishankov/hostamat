<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';

	import { Card, Error, Info, Modal } from '$lib/components';
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

		sites = sitesResp?.sites || [];
		maxAmount = sitesResp?.maxAmount || 0;

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
			{:then}
				{#each sites as site (site.id)}
					<Card style="height: 150px;">
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

				{#if sites.length < maxAmount}
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
							Для добавления нового сайта, удали один из существующих
						</Info>
					</div>
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

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
	}

	.dashboard {
		display: flex;
		flex-direction: column;
		gap: 10px;

		width: clamp(500px, 75%, 1200px);

		padding: 20px;
	}

	.site-list {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
	}
</style>
