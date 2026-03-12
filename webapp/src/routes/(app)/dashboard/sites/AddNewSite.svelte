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

<div class="add-site-card-shell">
	<Card border floating kind="action" style="height: 100%;">
		<button
			type="button"
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
</div>

<Modal
	bind:this={modal}
	titleId="new-site-modal-title"
	initialFocusSelector="#new-site-name"
	closeLabel="Закрыть форму создания сайта"
>
	<div class="modal-body">
		<h2 id="new-site-modal-title">Новый сайт</h2>
		<Input
			id="new-site-name"
			label="Название"
			errorMessage={newSiteNameError}
			bind:value={newSiteName}
			maxlength="25"
		></Input>
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
				<span class="action-spinner"><Spinner size={30} /></span>
			{/if}
		</div>
	</div>
</Modal>

<style>
	.add-site-card-shell {
		display: flex;
		min-height: 170px;
	}

	.add-site {
		all: unset;

		width: 100%;
		height: 100%;
		min-height: 170px;
		box-sizing: border-box;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		text-align: center;

		padding: 1rem 1.3rem;
		cursor: pointer;
	}

	.add-site p {
		width: 100%;
		margin: 0;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.add-site-text {
		color: var(--color-action-90);
		font-weight: 700;
		user-select: none;
	}

	.modal-body {
		display: flex;
		flex-direction: column;
		gap: 15px;
		min-width: 270px;
	}

	.modal-actions {
		display: flex;
		gap: 15px;
		align-items: center;
	}

	.action-spinner {
		display: inline-flex;
		align-items: center;
	}
</style>
