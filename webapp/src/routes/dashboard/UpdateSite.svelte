<script lang="ts">
	import { Button, IconButton, Modal, Input } from '$lib/components';
	import { Edit, Spinner } from '$lib/icons';
	import { updateSite } from '$lib/api-client';

	let { siteName = '', siteId = '', loading = false, onsuccess = async () => {} } = $props();

	let modal = $state<Modal>();
	let newName = $state('');
	let errorMessage = $state('');
</script>

<IconButton onclick={() => modal?.open()} title="Изменить сайт">
	<Edit size={25} />
</IconButton>
<Modal bind:this={modal} kind="primary" openCallback={() => (newName = siteName)}>
	<div style="display: flex; flex-direction: column; gap: 15px;">
		<h2>Изменить "{siteName}"</h2>
		<Input id="site-name" label="Название" bind:value={newName} maxlength="25" {errorMessage}
		></Input>
		<div style="display: flex; gap: 15px; align-items: center;">
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
				<Spinner size={30} />
			{/if}
		</div>
	</div>
</Modal>
