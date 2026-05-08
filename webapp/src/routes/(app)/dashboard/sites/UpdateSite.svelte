<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Button, IconButton, Modal, Input } from '$lib/components';
	import { Edit, Spinner } from '$lib/icons';
	import { updateSite } from '$lib/api-client';

	interface Props {
		siteName?: string;
		siteId?: string;
		loading?: boolean;
		iconButtonGrouped?: boolean;
		iconSize?: number;
		onsuccess?: () => Promise<void>;
		trigger?: Snippet<[{ open: () => void }]>;
		ontrigger?: (open: () => void) => void;
	}

	let {
		siteName = '',
		siteId = '',
		loading = false,
		iconButtonGrouped = false,
		iconSize = 25,
		onsuccess = async () => {},
		trigger,
		ontrigger
	}: Props = $props();

	let modal = $state<Modal>();
	let newName = $state('');
	let errorMessage = $state('');

	$effect(() => {
		if (ontrigger) {
			ontrigger(() => modal?.open());
		}
	});

	const modalTitleId = $derived(`update-site-title-${siteId}`);
	const inputId = $derived(`update-site-name-${siteId}`);
</script>

{#if trigger}
	{@render trigger({ open: () => modal?.open() })}
{:else if !ontrigger}
	<IconButton
		label="Изменить сайт"
		grouped={iconButtonGrouped}
		onclick={() => modal?.open()}
		title="Изменить сайт"
	>
		<Edit size={iconSize} />
	</IconButton>
{/if}
<Modal
	bind:this={modal}
	tone="neutral"
	title={`Изменить "${siteName}"`}
	titleId={modalTitleId}
	initialFocusSelector={'#' + inputId}
	openCallback={() => (newName = siteName)}
>
	<div class="modal-body">
		<Input id={inputId} label="Название" bind:value={newName} maxlength="25" {errorMessage}></Input>
		<div class="modal-actions">
			<Button
				emphasis="solid"
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
		gap: var(--space-lg);
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
