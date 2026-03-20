<script lang="ts">
	import { Button, IconButton, Modal, Input } from '$lib/components';
	import { Edit, Spinner } from '$lib/icons';
	import { updateSite } from '$lib/api-client';

	let {
		siteName = '',
		siteId = '',
		loading = false,
		iconButtonGrouped = false,
		onsuccess = async () => {}
	} = $props();

	let modal = $state<Modal>();
	let newName = $state('');
	let errorMessage = $state('');

	const modalTitleId = $derived(`update-site-title-${siteId}`);
	const inputId = $derived(`update-site-name-${siteId}`);
</script>

<IconButton
	label="Изменить сайт"
	grouped={iconButtonGrouped}
	onclick={() => modal?.open()}
	title="Изменить сайт"
>
	<Edit size={25} />
</IconButton>
<Modal
	bind:this={modal}
	kind="primary"
	title={`Изменить "${siteName}"`}
	titleId={modalTitleId}
	initialFocusSelector={'#' + inputId}
	openCallback={() => (newName = siteName)}
>
	<div class="modal-body">
		<Input id={inputId} label="Название" bind:value={newName} maxlength="25" {errorMessage}></Input>
		<div class="modal-actions">
			<Button
				filled
				disabled={loading}
				onclick={async () => {
					if (newName.length === 0) {
						errorMessage = 'Название не должно быть пустым';
						return;
					}

					loading = true;
					const { data, error } = await updateSite(siteId, newName);

					if (error) {
						errorMessage = error;
						loading = false;
					}

					if (data) {
						await onsuccess();
						modal?.close();
					}
				}}>Изменить</Button
			>
			<Button disabled={loading} onclick={() => modal?.close()}>Отмена</Button>
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
		gap: 15px;
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
