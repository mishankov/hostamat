<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';

	import { Button, Card, Error, Info, Input, Modal } from '$lib/components';
	import { Spinner, Plus } from '$lib/icons';

	import { type Site } from '$lib/types';
	import SiteContent from './SiteContent.svelte';
	import { createSite } from '$lib/api-client';

	let { data }: PageProps = $props();
	let sites = $state<Site[]>([]);
	let maxAmount = $state(0);
	let loading = $state(true);
	let newSiteName = $state('Новый сайт');
	let errorMessage = $state('');

	const updateState = async () => {
		const sitesResp = await data.sitesResp;

		sites = sitesResp?.sites || [];
		maxAmount = sitesResp?.maxAmount || 0;

		newSiteName = 'Новый сайт';
		loading = false;
		newSiteModal?.close();
	};

	let newSiteModal = $state<Modal>();

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
					<Card floating style="height: 150px;" kind="action">
						<button
							class="add-site"
							onclick={() => {
								newSiteModal?.open();
							}}
						>
							{#if loading}
								<p><Spinner size={40} /></p>
								<p class="add-site-text">Добавляем сайт</p>
							{:else}
								<p><Plus size={40} /></p>
								<p class="add-site-text">Добавить сайт</p>
							{/if}
						</button>
					</Card>

					<Modal bind:this={newSiteModal}>
						<div style="display: flex; flex-direction: column; gap: 15px;">
							<h3>Новый сайт</h3>
							<Input id="site-name" label="Название" bind:value={newSiteName}></Input>
							<div style="display: flex; gap: 15px; align-items: center; ">
								<Button
									filled
									disabled={loading}
									onclick={async () => {
										loading = true;
										const res = await createSite(newSiteName);

										if (res) {
											await invalidateAll();
											updateState();
										} else {
											errorMessage = 'Ошибка при создании сайта';
											updateState();
										}
									}}>Создать</Button
								>
								<Button disabled={loading} onclick={() => newSiteModal?.close()}>Отмена</Button>
								{#if loading}
									<Spinner size={30} />
								{/if}
							</div>
						</div>
					</Modal>
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
		color: var(--color-primary-90);
		user-select: none;
	}
</style>
