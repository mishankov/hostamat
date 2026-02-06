<script lang="ts">
	import { Card, Modal, Input, Button } from '$lib/components';
	import { Spinner, Plus } from '$lib/icons';
	import { createSite } from '$lib/api-client';

	let { loading = $bindable(false), onsuccess = async () => {} } = $props();

	let modal = $state<Modal>();
	let newSiteName = $state('');
	let newSiteNameError = $state('');

	const resetState = () => {
		newSiteName = 'Новый сайт';
		newSiteNameError = '';
	};

	resetState();
</script>

<Card border floating style="min-height: 170px;" kind="action">
	<button
		class="add-site"
		onclick={() => {
			resetState();
			modal?.open();
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

<Modal bind:this={modal}>
	<div style="display: flex; flex-direction: column; gap: 15px; min-width: 270px;">
		<h2>Новый сайт</h2>
		<Input
			id="site-name"
			label="Название"
			errorMessage={newSiteNameError}
			bind:value={newSiteName}
			maxlength="25"
		></Input>
		<div style="display: flex; gap: 15px; align-items: center; ">
			<Button
				filled
				disabled={loading}
				onclick={async () => {
					if (newSiteName.length === 0) {
						newSiteNameError = 'Название не должно быть пустым';
						return;
					}

					loading = true;
					const { data, error } = await createSite(newSiteName);

					if (error) {
						newSiteNameError = error;
						loading = false;
					}

					if (data) {
						await onsuccess();
						modal?.close();
					}
				}}>Создать</Button
			>
			<Button
				disabled={loading}
				onclick={() => {
					modal?.close();
				}}>Отмена</Button
			>
			{#if loading}
				<Spinner size={30} />
			{/if}
		</div>
	</div>
</Modal>

<style>
	.add-site {
		all: unset;

		height: 100%;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;

		padding: 1rem 1.3rem;
		cursor: pointer;
	}

	.add-site-text {
		color: var(--color-action-90);
		font-weight: 700;
		user-select: none;
	}
</style>
