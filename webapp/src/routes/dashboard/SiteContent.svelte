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

	const resetState = () => {
		loading = false;
	};
</script>

<div class="site-content">
	<div class="site-info">
		<div class="site-name">
			<Link href="https://{site.domain}" target="_blank" alwaysUnderline>{site.name}</Link>
		</div>
		<span class="site-subdomain" title="Поддомен">{subdomain}</span>
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
				<p class="modal-subtitle">для "{site.name}"</p>
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
			title="Скопировать URL"
			onclick={() => {
				navigator.clipboard.writeText(`https://${site.domain}`);
			}}
		>
			<Copy size={20} />
		</IconButton>

		<div class="toolbar-divider"></div>

		<IconButton title="Удалить сайт" kind="danger" onclick={() => deleteModal?.open()}>
			<Trash size={20} />
		</IconButton>
		<Modal bind:this={deleteModal} kind="danger">
			<div class="delete-modal">
				<h3>Удалить сайт?</h3>
				<p>
					Вы уверены, что хотите удалить сайт "{site.name}"? Это действие нельзя отменить.
				</p>
				<form method="dialog" class="delete-actions">
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
					<Button kind="danger">Отмена</Button>
				</form>
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
		min-width: 200px;
		min-height: 140px;
		padding: 20px;
	}

	.site-info {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.site-name {
		font-size: 17px;
		font-weight: 500;
	}

	.site-subdomain {
		font-family: 'JetBrains Mono', monospace;
		font-size: 13px;
		color: var(--color-text-muted);
	}

	.site-toolbar {
		display: flex;
		align-items: center;
		gap: 4px;
		padding: 4px;
		background: var(--color-background-elevated);
		border: 1px solid var(--glass-border);
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
		min-width: 300px;
	}

	.upload-modal h2 {
		font-size: 20px;
		margin: 0;
	}

	.modal-subtitle {
		color: var(--color-text-secondary);
		font-size: 14px;
		margin-top: -12px;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 360px;
	}

	.delete-modal h3 {
		font-size: 20px;
		color: var(--color-danger-80);
		margin: 0;
	}

	.delete-modal p {
		color: var(--color-danger-70);
		font-size: 15px;
		line-height: 1.5;
	}

	.delete-actions {
		display: flex;
		gap: 12px;
		padding-top: 8px;
	}
</style>
