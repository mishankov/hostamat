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

<Card border floating kind="action">
	<button
		class="add-site"
		onclick={() => {
			resetState();
			modal?.open();
		}}
	>
		{#if loading}
			<div class="add-site-icon">
				<Spinner size={32} />
			</div>
			<span class="add-site-text">Добавляем...</span>
		{:else}
			<div class="add-site-icon">
				<Plus size={32} />
			</div>
			<span class="add-site-text">Добавить сайт</span>
		{/if}
	</button>
</Card>

<Modal bind:this={modal}>
	<div class="new-site-modal">
		<h2>Новый сайт</h2>
		<Input
			id="site-name"
			label="Название"
			errorMessage={newSiteNameError}
			bind:value={newSiteName}
			maxlength="25"
		/>
		<div class="modal-actions">
			<Button
				kind="accent"
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
				<Spinner size={24} />
			{/if}
		</div>
	</div>
</Modal>

<style>
	.add-site {
		all: unset;
		cursor: pointer;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 12px;

		min-height: 140px;
		padding: 24px;
		width: 100%;
	}

	.add-site-icon {
		width: 56px;
		height: 56px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-action-15);
		border-radius: var(--border-radius-regular);
		color: var(--color-action-60);
		transition: var(--transition);
	}

	.add-site:hover .add-site-icon {
		background: var(--color-action-20);
		transform: scale(1.05);
	}

	.add-site-text {
		font-size: 15px;
		font-weight: 500;
		color: var(--color-action-70);
		user-select: none;
	}

	.new-site-modal {
		display: flex;
		flex-direction: column;
		gap: 20px;
		min-width: 300px;
	}

	.new-site-modal h2 {
		font-size: 20px;
		margin: 0;
	}

	.modal-actions {
		display: flex;
		align-items: center;
		gap: 12px;
		padding-top: 8px;
	}
</style>
