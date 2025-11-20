<script lang="ts">
	import { Button, Info, Input, Link, Modal } from '$lib/components';
	import { Copy, Trash, Spinner, Upload } from '$lib/icons';

	import FileUploader from './components/FileUploader.svelte';
	import { IconButton } from '$lib/components';
	import { type Site } from '$lib/types';
	import { deleteSite, updateSite } from '$lib/api-client';

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

	let updateModal = $state<Modal>();
	let uploadModal = $state<Modal>();
	let deleteModal = $state<Modal>();
	let fileUploader = $state<FileUploader>();
	let newSiteName = $state('');

	const resetState = () => {
		loading = false;
		newSiteName = '';
		updateModal?.close();
	};
</script>

<div class="site-content">
	<div class="site-url-container">
		<Link href="https://{site.domain}" target="_blank" alwaysUnderline>{site.name}</Link>

		<IconButton
			title="Скопировать URL"
			onclick={() => {
				navigator.clipboard.writeText(`https://${site.domain}`);
			}}
		>
			<Copy size={20} />
		</IconButton>
	</div>

	<div class="site-toolbar">
		<IconButton onclick={() => updateModal?.open()} title="Изменить сайт">
			<Upload size={20} />
		</IconButton>
		<Modal
			bind:this={updateModal}
			kind="primary"
			closeCallback={() => {
				fileUploader?.resetUploadState();
			}}
		>
			<div style="display: flex; flex-direction: column; gap: 15px;">
				<h2>Сайт "{site.name}"</h2>
				<Input id="site-name" label="Название" bind:value={newSiteName}></Input>
				<div style="display: flex; gap: 15px; align-items: center;">
					<Button
						filled
						disabled={loading}
						onclick={async () => {
							loading = true;
							const res = await updateSite(site.id, newSiteName);

							if (res) {
								await onsiteupdate();
								resetState();
							} else {
								errorMessage = 'Ошибка при изменении сайта';
								resetState();
							}
						}}>Изменить</Button
					>
					{#if loading}
						<Spinner size={30} />
					{/if}
				</div>
			</div>
		</Modal>

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
				<h3>Загрузить файлы для "{site.name}"</h3>
				<Info header="Информация о загрузке сайтов">
					<p>При загрузке файлов все старые файлы будут удалены</p>
					<p>
						Разрешённые типы файлов <Link href="/about" target="_blank" alwaysUnderline
							>на странице</Link
						>
					</p>
				</Info>
				<FileUploader
					bind:this={fileUploader}
					multiple={true}
					uploadUrl="/api/uploadFile"
					siteId={site.id}
					additionalHeaders={{
						'H-Site-Id': site.id
					}}
				/>
			</div>
		</Modal>

		<IconButton title="Удалить сайт" onclick={() => deleteModal?.open()}>
			<Trash size={20} />
		</IconButton>
		<Modal bind:this={deleteModal} kind="danger">
			<div class="delete-modal">
				<p>Вы уверены, что хотите удалить сайт {subdomain}?</p>
				<form method="dialog">
					<Button kind="primary">Отмена</Button>
					<Button
						kind="danger"
						onclick={async () => {
							await deleteSite(site.id);

							window.location.reload();
						}}>Удалить</Button
					>
				</form>
			</div>
		</Modal>
	</div>
</div>

<style>
	.site-content {
		display: flex;
		flex-direction: column;
		gap: 5px;

		padding-top: 20px;
		padding-left: 20px;
		padding-right: 20px;
		padding-bottom: 10px;
	}

	.site-url-container {
		display: flex;
		align-items: center;
		gap: 3px;

		padding-bottom: 10px;
	}

	.site-toolbar {
		width: fit-content;
		align-self: flex-end;
		padding-top: 10px;

		display: flex;
		gap: 5px;
	}

	.upload-modal {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.delete-modal > form {
		display: flex;
		gap: 10px;
	}
</style>
