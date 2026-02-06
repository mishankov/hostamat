<script lang="ts">
	import { onDestroy } from 'svelte';
	import { Alert, Button, Card, IconButton, Link, Modal } from '$lib/components';
	import { Copy, Trash, Upload } from '$lib/icons';

	import FileUploader from './components/FileUploader.svelte';
	import UpdateSite from './UpdateSite.svelte';
	import { type Site } from '$lib/types';
	import { deleteSite } from '$lib/api-client';

	interface Props {
		site: Site;
		loading: boolean;
		errorMessage: string;
		onsiteupdate: () => Promise<void>;
	}

	let {
		site,
		loading = $bindable(false),
		errorMessage = $bindable(''),
		onsiteupdate = async () => {}
	}: Props = $props();

	const subdomain = $derived(site.domain.split('.')[0]);

	let uploadModal = $state<Modal>();
	let deleteModal = $state<Modal>();
	let fileUploader = $state<FileUploader>();
	let urlCopied = $state(false);
	let copyTimer: ReturnType<typeof setTimeout> | null = null;

	const uploadTitleId = $derived(`upload-site-title-${site.id}`);
	const uploadDescId = $derived(`upload-site-description-${site.id}`);
	const deleteTitleId = $derived(`delete-site-title-${site.id}`);
	const deleteDescId = $derived(`delete-site-description-${site.id}`);
	const deleteButtonId = $derived(`delete-site-confirm-${site.id}`);

	const resetState = () => {
		loading = false;
	};

	const copySiteUrl = async () => {
		try {
			await navigator.clipboard.writeText(`https://${site.domain}`);
			urlCopied = true;

			if (copyTimer) {
				clearTimeout(copyTimer);
			}

			copyTimer = setTimeout(() => {
				urlCopied = false;
			}, 1400);
		} catch {
			urlCopied = false;
		}
	};

	onDestroy(() => {
		if (copyTimer) {
			clearTimeout(copyTimer);
		}
	});
</script>

<div class="site-content">
	<div class="site-url-container">
		<Link href="https://{site.domain}" target="_blank" alwaysUnderline>{site.name}</Link>
		<span title="Поддомен" class="subdomain">{subdomain}</span>
		{#if urlCopied}
			<span class="copy-feedback" aria-live="polite">URL скопирован</span>
		{/if}
	</div>

	<Card border class="site-toolbar">
		<UpdateSite
			siteName={site.name}
			siteId={site.id}
			onsuccess={async () => {
				resetState();
				await onsiteupdate();
			}}
		/>

		<IconButton label="Загрузить сайт" onclick={() => uploadModal?.open()} title="Загрузить сайт">
			<Upload size={25} />
		</IconButton>
		<Modal
			bind:this={uploadModal}
			kind="primary"
			titleId={uploadTitleId}
			descriptionId={uploadDescId}
			initialFocusSelector=".drop-zone"
			closeCallback={() => {
				fileUploader?.resetUploadState();
			}}
		>
			<div class="upload-modal">
				<h2 id={uploadTitleId} class="upload-title">Загрузить файлы для "{site.name}"</h2>
				<Alert variant="info" title="Информация о загрузке сайтов">
					<p>При загрузке файлов все старые файлы будут удалены</p>
					<p id={uploadDescId}>
						Разрешённые типы файлов <Link href="/about" target="_blank" alwaysUnderline
							>на странице</Link
						>
					</p>
				</Alert>
				<FileUploader
					bind:this={fileUploader}
					onclose={() => uploadModal?.close()}
					multiple={true}
					uploadUrl="/api/uploadFile"
					siteId={site.id}
					additionalHeaders={{
						'H-Site-Id': site.id
					}}
				/>
			</div>
		</Modal>

		<IconButton
			label="Скопировать URL"
			title="Скопировать URL"
			onclick={copySiteUrl}
		>
			<Copy size={25} background="var(--color-primary-10)" />
		</IconButton>

		<IconButton
			label="Удалить сайт"
			title="Удалить сайт"
			kind="danger"
			class="toolbar-delete"
			onclick={() => deleteModal?.open()}
		>
			<Trash size={25} />
		</IconButton>
		<Modal
			bind:this={deleteModal}
			kind="danger"
			titleId={deleteTitleId}
			descriptionId={deleteDescId}
			initialFocusSelector={"#" + deleteButtonId}
		>
			<div class="delete-modal">
				<h2 id={deleteTitleId}>Подтверждение удаления</h2>
				<p id={deleteDescId} class="delete-modal-text">
					Вы уверены, что хотите удалить сайт "{site.name}"?
				</p>
				<div class="delete-actions">
					<Button
						id={deleteButtonId}
						kind="danger"
						filled
						onclick={async () => {
							const { data, error } = await deleteSite(site.id);

							if (error) {
								errorMessage = error;
								deleteModal?.close();
								return;
							}

							if (data) {
								resetState();
								deleteModal?.close();
								await onsiteupdate();
							}
						}}>Удалить</Button
					>
					<Button kind="danger" onclick={() => deleteModal?.close()}>Отмена</Button>
				</div>
			</div>
		</Modal>
	</Card>
</div>

<style>
	.site-content {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		min-width: 150px;
		height: 100%;

		padding: 1rem;
	}

	.site-url-container {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 5px;

		font-size: 18px;
	}

	.subdomain {
		font-size: 0.88rem;
		color: var(--color-text-secondary);
	}

	.copy-feedback {
		font-size: 0.82rem;
		color: var(--color-confirm-60);
	}

	:global(.site-toolbar) {
		width: fit-content;
		justify-self: end;

		display: flex;
		padding: 3px;
		margin-top: auto;
	}

	:global(.site-toolbar > button) {
		border: none;
		border-radius: 0;
	}

	:global(:root[data-theme='dark'] .site-toolbar button.toolbar-delete:hover) {
		background: var(--color-danger-10);
	}

	.upload-modal {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.upload-title {
		max-width: 500px;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.delete-modal-text {
		color: var(--color-danger-90);
	}

	.delete-actions {
		display: flex;
		gap: 10px;
	}

	@media (max-width: 620px) {
		:global(.site-toolbar) {
			width: 100%;
			justify-content: space-between;
			gap: 4px;
		}
	}
</style>
