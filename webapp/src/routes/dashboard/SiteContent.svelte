<script lang="ts">
	import { Button, Info, Input, Link, Modal } from '$lib/components';
	import { Copy, Trash, Spinner, Upload, Edit } from '$lib/icons';

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
		<span title="Поддомен" style="font-size: 16px; color: var(--color-text-secondary);"
			>{subdomain}</span
		>
	</div>

	<div class="site-toolbar">
		<IconButton onclick={() => updateModal?.open()} title="Изменить сайт">
			<Edit size={25} />
		</IconButton>
		<Modal
			bind:this={updateModal}
			kind="primary"
			openCallback={() => (newSiteName = site.name)}
			closeCallback={() => {
				fileUploader?.resetUploadState();
			}}
		>
			<div style="display: flex; flex-direction: column; gap: 15px;">
				<h3>Изменить "{site.name}"</h3>
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
					<Button disabled={loading} onclick={() => updateModal?.close()}>Отмена</Button>
					{#if loading}
						<Spinner size={30} />
					{/if}
				</div>
			</div>
		</Modal>

		<IconButton onclick={() => uploadModal?.open()} title="Загрузить сайт">
			<Upload size={25} />
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
			<Copy size={25} background="var(--color-primary-7)" />
		</IconButton>

		<IconButton title="Удалить сайт" kind="danger" onclick={() => deleteModal?.open()}>
			<Trash size={25} />
		</IconButton>
		<Modal bind:this={deleteModal} kind="danger">
			<div class="delete-modal">
				<p>Вы уверены, что хотите удалить сайт "{site.name}"?</p>
				<form method="dialog">
					<Button
						kind="danger"
						filled
						onclick={async () => {
							await deleteSite(site.id);

							window.location.reload();
						}}>Удалить</Button
					>
					<Button kind="primary">Отмена</Button>
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
		min-width: 150px;
		height: 100%;

		padding-top: 20px;
		padding-left: 20px;
		padding-right: 20px;
		padding-bottom: 20px;
	}

	.site-url-container {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 5px;

		font-size: 18px;
	}

	.site-toolbar {
		width: fit-content;
		justify-self: end;

		display: flex;
		padding: 2px;
		margin-top: auto;
		background-color: var(--color-primary-7);

		border: var(--border-regular);
		border-radius: 7px;
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
