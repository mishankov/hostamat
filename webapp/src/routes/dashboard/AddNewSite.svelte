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
		<div class="add-site-icon">
			{#if loading}
				<Spinner size={32} />
			{:else}
				<Plus size={32} color="currentColor" />
			{/if}
		</div>
		<span class="add-site-text">
			{#if loading}
				Добавляем сайт...
			{:else}
				Добавить сайт
			{/if}
		</span>
	</button>
</Card>

<Modal bind:this={modal}>
	<div class="new-site-modal">
		<h2>Новый сайт</h2>
		<p class="modal-subtitle">Введите название для вашего сайта</p>
		<Input
			id="site-name"
			label="Название"
			errorMessage={newSiteNameError}
			bind:value={newSiteName}
			maxlength="25"
		/>
		<div class="modal-actions">
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
				<Spinner size={24} />
			{/if}
		</div>
	</div>
</Modal>

<style>
	.add-site {
		all: unset;
		box-sizing: border-box;

		width: 100%;
		height: 100%;
		min-height: 160px;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 12px;

		padding: 24px;
		cursor: pointer;
	}

	.add-site-icon {
		display: flex;
		align-items: center;
		justify-content: center;

		width: 56px;
		height: 56px;

		background: var(--color-action-15);
		border: 1px dashed var(--color-action-40);
		border-radius: var(--border-radius-regular);

		color: var(--color-action-60);
		transition: all var(--transition-fast);
	}

	.add-site:hover .add-site-icon {
		background: var(--color-action-20);
		border-style: solid;
		color: var(--color-action-70);
	}

	.add-site-text {
		font-size: 0.9375rem;
		font-weight: 500;
		color: var(--color-action-70);
		user-select: none;
	}

	.new-site-modal {
		display: flex;
		flex-direction: column;
		gap: 16px;
		min-width: 320px;
	}

	.new-site-modal h2 {
		margin: 0;
	}

	.modal-subtitle {
		color: var(--color-text-secondary);
		font-size: 0.9375rem;
		margin: -8px 0 8px 0;
	}

	.modal-actions {
		display: flex;
		align-items: center;
		gap: 12px;
		padding-top: 8px;
	}
</style>
