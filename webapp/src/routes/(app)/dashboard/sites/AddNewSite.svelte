<script lang="ts">
	import { Button, Modal, Input } from '$lib/components';
	import { Spinner } from '$lib/icons';
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

<Button
	emphasis="solid"
	disabled={loading}
	onclick={() => {
		resetState();
		modal?.open();
	}}
>
	{#if loading}
		<Spinner size={16} />
		<span>Добавляем сайт</span>
	{:else}
		<span>Добавить сайт</span>
	{/if}
</Button>

<Modal
	bind:this={modal}
	title="Новый сайт"
	titleId="new-site-modal-title"
	initialFocusSelector="#new-site-name"
	closeLabel="Закрыть форму создания сайта"
>
	<div class="modal-body">
		<Input
			id="new-site-name"
			label="Название"
			errorMessage={newSiteNameError}
			bind:value={newSiteName}
			maxlength="25"
		></Input>
		<div class="modal-actions">
			<Button
				emphasis="solid"
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
	.modal-body {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 16.875rem;
	}

	.modal-actions {
		display: flex;
		gap: var(--space-lg);
		align-items: center;
	}

	.action-spinner {
		display: inline-flex;
		align-items: center;
	}
</style>
