<script lang="ts">
	import { Button, Info, Link, Modal } from '$lib/components';
	import { Copy, Trash, Upload } from '$lib/icons';

	import FileUploader from './components/FileUploader.svelte';
	import UpdateSite from './UpdateSite.svelte';
	import { IconButton } from '$lib/components';
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

	const subdomain = site.domain.split('.')[0];

	let uploadModal = $state<Modal>();
	let deleteModal = $state<Modal>();
	let fileUploader = $state<FileUploader>();
	let copied = $state(false);

	const resetState = () => {
		loading = false;
	};

	const copyUrl = () => {
		navigator.clipboard.writeText(`https://${site.domain}`);
		copied = true;
		setTimeout(() => copied = false, 2000);
	};
</script>

<div class="site-content">
	<div class="site-info">
		<div class="site-url-container">
			<Link href="https://{site.domain}" target="_blank" alwaysUnderline>
				<span class="site-name">{site.name}</span>
			</Link>
			<span class="subdomain">{subdomain}</span>
		</div>
	</div>

	<div class="site-toolbar">
		<UpdateSite
			siteName={site.name}
			siteId={site.id}
			onsuccess={async () => {
				resetState();
				await onsiteupdate();
			}}
		/>

		<IconButton onclick={() => uploadModal?.open()} title="Загрузить сайт">
			<Upload size={20} />
		</IconButton>
		<Modal
			bind:this={uploadModal}
			kind="primary"
			closeCallback={() => {
				fileUploader?.resetUploadState();
			}}
		>
			<div class="upload-modal">
				<h2>Загрузить файлы</h2>
				<p class="modal-subtitle">для «{site.name}»</p>
				<Info header="Информация о загрузке">
					<p>При загрузке файлов все старые файлы будут удалены</p>
					<p>
						Разрешённые типы файлов <Link href="/about" target="_blank" alwaysUnderline
							>на странице</Link
						>
					</p>
				</Info>
				<FileUploader
					bind:this={fileUploader}
					onclose={uploadModal?.close}
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
			title={copied ? "Скопировано!" : "Скопировать URL"}
			onclick={copyUrl}
		>
			<Copy size={20} background="var(--color-background-object)" />
		</IconButton>

		<div class="toolbar-divider"></div>

		<IconButton title="Удалить сайт" kind="danger" onclick={() => deleteModal?.open()}>
			<Trash size={20} />
		</IconButton>
		<Modal bind:this={deleteModal} kind="danger">
			<div class="delete-modal">
				<h3>Удалить сайт?</h3>
				<p>
					Вы уверены, что хотите удалить <strong>«{site.name}»</strong>?
					Это действие нельзя отменить.
				</p>
				<div class="modal-actions">
					<Button
						kind="danger"
						filled
						onclick={async () => {
							const { data, error } = await deleteSite(site.id);

							if (error) {
								errorMessage = error;
								deleteModal?.close();
							}

							if (data) {
								window.location.reload();
							}
						}}>Удалить</Button
					>
					<Button kind="danger" onclick={() => deleteModal?.close()}>Отмена</Button>
				</div>
			</div>
		</Modal>
	</div>
</div>

<style>
	.site-content {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		gap: 16px;
		min-width: 150px;
		height: 100%;
		min-height: 160px;

		padding: 20px;
	}

	.site-info {
		flex: 1;
	}

	.site-url-container {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 6px;
	}

	.site-name {
		font-family: var(--font-display);
		font-size: 1.125rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.subdomain {
		font-size: 0.8125rem;
		color: var(--color-text-muted);
		font-family: var(--font-mono);
	}

	.site-toolbar {
		display: flex;
		align-items: center;
		gap: 4px;

		padding: 4px;
		margin-top: auto;

		background: var(--color-background-object);
		border: 1px solid var(--color-border-subtle);
		border-radius: var(--border-radius-regular);

		width: fit-content;
	}

	.toolbar-divider {
		width: 1px;
		height: 20px;
		background: var(--color-border);
		margin: 0 4px;
	}

	.upload-modal {
		display: flex;
		flex-direction: column;
		gap: 16px;
		min-width: 320px;
	}

	.upload-modal h2 {
		margin: 0;
	}

	.modal-subtitle {
		color: var(--color-text-secondary);
		font-size: 0.9375rem;
		margin: -8px 0 8px 0;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 400px;
	}

	.delete-modal h3 {
		margin: 0;
		font-family: var(--font-display);
	}

	.delete-modal p {
		color: var(--color-danger-70);
		line-height: 1.6;
		margin: 0;
	}

	.delete-modal strong {
		color: var(--color-danger-80);
	}

	.modal-actions {
		display: flex;
		gap: 12px;
		padding-top: 8px;
	}
</style>
