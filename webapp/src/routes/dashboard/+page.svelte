<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';

	import { Card, Error, Info } from '$lib/components';
	import { Spinner, Plus, X } from '$lib/icons';

	import { type Site } from './types';
	import SiteContent from './SiteContent.svelte';

	let { data }: PageProps = $props();
	let sites = $state<Site[]>([]);
	let maxAmount = $state(0);
	let loading = $state(true);

	const updateState = async () => {
		let sitesResp = await data.sitesResp;
		if (sitesResp) {
			sites = sitesResp.sites;
			maxAmount = sitesResp.maxAmount || 0;
		}

		loading = false;
	};

	updateState();
</script>

<main>
	<div class="dashboard">
		<h1>Твои сайты</h1>

		<div class="site-list">
			{#await data.sitesResp}
				<Spinner />
			{:then _}
				{#each sites as site (site.id)}
					<Card style="height: 110px;">
						<SiteContent {site} />
					</Card>
				{/each}

				{#if sites.length < maxAmount}
					<Card floating style="height: 110px;">
						<button
							class="add-site"
							onclick={async () => {
								await fetch('/api/createSite', {
									method: 'POST'
								});

								loading = true;
								await invalidateAll();
								updateState();
							}}
						>
							<p class="add-site-icon"><Plus size={40} /></p>
							<p class="add-site-text">Добавить сайт</p>
						</button>
					</Card>
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

	.add-site {
		all: unset;

		height: 100%;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;

		padding-left: 20px;
		padding-right: 20px;
	}

	.add-site-text {
		color: #bbb;
		user-select: none;
	}
</style>
