<script lang="ts">
	import { Alert, Button, CopyButton, IconButton, Link, Modal, Toolbar } from '$lib/components';
	import { useThemeController } from '$lib/theme/controller.svelte';
	import { Trash, Upload } from '$lib/icons';

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
	const theme = useThemeController();

	let uploadModal = $state<Modal>();
	let deleteModal = $state<Modal>();
	let fileUploader = $state<FileUploader>();

	const uploadTitleId = $derived(`upload-site-title-${site.id}`);
	const uploadDescId = $derived(`upload-site-description-${site.id}`);
	const deleteTitleId = $derived(`delete-site-title-${site.id}`);
	const deleteDescId = $derived(`delete-site-description-${site.id}`);
	const deleteButtonId = $derived(`delete-site-confirm-${site.id}`);

	const resetState = () => {
		loading = false;
	};
</script>

<div class={['site-content', theme.current === 'dark' && 'theme-dark']}>
	<div class="site-data">
		<div class="data-row">
			<span class="data-label">name:</span>
			<Link href="https://{site.domain}" target="_blank" alwaysUnderline>{site.name}</Link>
		</div>
		<div class="data-row">
			<span class="data-label">domain:</span>
			<span class="subdomain-value" title="Поддомен">{subdomain}</span>
		</div>
	</div>

	<Toolbar>
		<UpdateSite
			siteName={site.name}
			siteId={site.id}
			iconButtonGrouped={true}
			onsuccess={async () => {
				resetState();
				await onsiteupdate();
			}}
		/>

		<IconButton
			grouped
			label="Загрузить сайт"
			onclick={() => uploadModal?.open()}
			title="Загрузить сайт"
		>
			<Upload size={25} />
		</IconButton>
		<CopyButton
			grouped
			textToCopy={`https://${site.domain}`}
			label="Скопировать URL"
			title="Скопировать URL"
			size={25}
		/>

		<IconButton
			grouped
			label="Удалить сайт"
			title="Удалить сайт"
			tone="danger"
			onclick={() => deleteModal?.open()}
		>
			<Trash size={25} />
		</IconButton>
	</Toolbar>

	<Modal
		bind:this={uploadModal}
		tone="neutral"
		title={`Загрузить файлы для "${site.name}"`}
		panelWidth="33.75rem"
		titleId={uploadTitleId}
		descriptionId={uploadDescId}
		initialFocusSelector=".drop-zone"
		closeCallback={() => {
			fileUploader?.resetUploadState();
		}}
	>
		<div class="upload-modal">
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

	<Modal
		bind:this={deleteModal}
		tone="danger"
		title="Подтверждение удаления"
		titleId={deleteTitleId}
		descriptionId={deleteDescId}
		initialFocusSelector={'#' + deleteButtonId}
	>
		<div class="delete-modal">
			<p id={deleteDescId} class="delete-modal-text">
				Вы уверены, что хотите удалить сайт "{site.name}"?
			</p>
			<div class="delete-actions">
				<Button
					id={deleteButtonId}
					tone="danger"
					emphasis="solid"
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
				<Button tone="danger" onclick={() => deleteModal?.close()}>Отмена</Button>
			</div>
		</div>
	</Modal>
</div>

<style>
	.site-content {
		display: grid;
		grid-template-columns: 1fr auto;
		align-items: center;
		gap: var(--space-lg);
		min-width: 0;
		padding: var(--space-lg) 0;
	}

	.site-data {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		min-width: 0;
	}

	.subdomain-value {
		font-size: 0.88rem;
		color: var(--color-text-secondary);
		overflow-wrap: anywhere;
	}

	.upload-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.delete-modal-text {
		color: var(--color-danger-90);
	}

	.delete-actions {
		display: flex;
		gap: var(--space-sm);
	}

	@media (max-width: 820px) {
		.site-content {
			grid-template-columns: 1fr;
			gap: var(--space-md);
		}
	}
</style>
