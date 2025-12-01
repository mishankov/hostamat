<script lang="ts">
	import { Button, Card, Details, Error, ProgressBar, Success } from '$lib/components';
	import { Upload, X } from '$lib/icons';
	import { IconButton } from '$lib/components';
	import { closeUploadSession, createUploadSession } from '$lib/api-client';

	interface Props {
		uploadUrl?: string;
		maxFileSize?: number;
		multiple?: boolean;
		disabled?: boolean;
		additionalHeaders?: Record<string, string>;
		siteId: string;
		onfilesselected?: (event: { files: File[] }) => void;
		onuploadcomplete?: (event: { fileCount: number }) => void;
		onuploaderror?: (event: { error: string }) => void;
		onclose?: () => void;
	}

	let {
		uploadUrl = '/api/upload',
		maxFileSize = 20 * 1024 * 1024, // 20MB default
		multiple = false,
		disabled = false,
		additionalHeaders = {},
		siteId,
		onfilesselected,
		onuploadcomplete,
		onuploaderror,
		onclose = () => {}
	}: Props = $props();

	let fileInput: HTMLInputElement;
	let selectedFiles: File[] = $state([]);
	let uploading: boolean = $state(false);
	let uploadProgress: number = $state(0);
	let error: string = $state('');
	let uploadedFiles: string[] = $state([]);
	let skippedFiles: string[] = $state([]);

	export function resetUploadState() {
		uploadedFiles = [];
		skippedFiles = [];
		uploadProgress = 0;
		error = '';
	}

	function handleFileSelect(event: Event): void {
		const target = event.target as HTMLInputElement;
		if (!target.files) return;

		resetUploadState();

		const files = Array.from(target.files);
		selectedFiles = files;

		// Validate file sizes
		const filesOverallSize = files.reduce((acc, file) => acc + file.size, 0);
		if (filesOverallSize > maxFileSize) {
			error = `Файлы слишком большие. Максимальный размер: ${formatFileSize(maxFileSize)}`;
			return;
		}

		onfilesselected?.({ files });
	}

	function handleDragOver(event: DragEvent): void {
		event.preventDefault();
		resetUploadState();
		if (event.dataTransfer) {
			event.dataTransfer.dropEffect = 'copy';
		}
	}

	function handleDrop(event: DragEvent): void {
		event.preventDefault();
		if (!event.dataTransfer) return;

		resetUploadState();

		const files = Array.from(event.dataTransfer.files);

		// Update file input
		const dt = new DataTransfer();
		files.forEach((file) => dt.items.add(file));
		fileInput.files = dt.files;

		selectedFiles = files;
		error = '';

		onfilesselected?.({ files });
	}

	async function uploadFiles(): Promise<void> {
		if (selectedFiles.length === 0) {
			error = 'Выбери файлы для загрузки';
			return;
		}

		uploading = true;
		uploadProgress = 0;
		error = '';

		const uploadSession = await createUploadSession(siteId);

		if (uploadSession === null) {
			error = 'Ошибка при создании сессии загрузки';
			uploading = false;
			return;
		}

		// Validate file sizes
		const filesOverallSize = selectedFiles.reduce((acc, file) => acc + file.size, 0);
		if (filesOverallSize > maxFileSize) {
			error = `Файлы слишком большие. Максимальный размер: ${formatFileSize(maxFileSize)}`;
			uploading = false;
			return;
		}

		try {
			for (let i = 0; i < selectedFiles.length; i++) {
				const file = selectedFiles[i];
				const formData = new FormData();
				formData.append('file', file);

				const headers = new Headers();
				Object.entries(additionalHeaders).forEach(([key, value]) => {
					headers.set(key, value);
				});
				headers.set('H-UploadSession-Id', uploadSession);

				const response = await fetch(uploadUrl, {
					method: 'POST',
					headers,
					body: formData
				});

				if (!response.ok) {
					error = 'Ошибка при загрузке файла';
					return;
				}

				uploadProgress = ((i + 1) / selectedFiles.length) * 100;
			}

			const uploadedCount = selectedFiles.length;
			selectedFiles = [];
			fileInput.value = '';

			onuploadcomplete?.({ fileCount: uploadedCount });
		} catch {
			error = 'Ошибка при загрузке файла';
			onuploaderror?.({ error });
		}

		const closeUploadResp = await closeUploadSession(uploadSession);

		if (closeUploadResp.ok) {
			uploadedFiles = closeUploadResp.files
				.filter((file) => file.uploaded)
				.map((file) => `${file.name} (${formatFileSize(file.size)})`);
			skippedFiles = closeUploadResp.files
				.filter((file) => !file.uploaded)
				.map((file) => file.name);
		} else {
			error = 'Ошибка при завершении загрузки';
		}

		uploading = false;
	}

	function removeFile(index: number): void {
		selectedFiles = selectedFiles.filter((_, i) => i !== index);

		// Update file input
		const dt = new DataTransfer();
		selectedFiles.forEach((file) => dt.items.add(file));
		fileInput.files = dt.files;

		if (selectedFiles.length === 0) {
			error = '';
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 байт';
		const k = 1024;
		const sizes = ['байт', 'Кбайт', 'Мбайт', 'Гбайт'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function clearAll(): void {
		selectedFiles = [];
		fileInput.value = '';
		error = '';
		uploadProgress = 0;
	}

	function handleKeyDown(event: KeyboardEvent): void {
		if (event.key === 'Enter') {
			fileInput.click();
		}
	}
</script>

<div class="file-upload-container">
	{#if error}
		<Error message={error} />
	{/if}

	{#if uploadProgress == 100}
		<Success header="Загрузка успешно завершена">
			<div class="success-message">
				{#if uploadedFiles.length > 0}
					<div class="uploaded-files">
						<Details summary="Загруженные файлы ({uploadedFiles.length})">
							<ul>
								{#each uploadedFiles as file (file)}
									<li>{file}</li>
								{/each}
							</ul>
						</Details>
					</div>
				{/if}
				{#if skippedFiles.length > 0}
					<div class="skipped-files">
						<Details summary="Пропущенные файлы ({skippedFiles.length})">
							<ul>
								{#each skippedFiles as file (file)}
									<li>{file}</li>
								{/each}
							</ul>
						</Details>
					</div>
				{/if}
			</div>
		</Success>
	{/if}

	<div
		class="drop-zone"
		class:drag-over={false}
		ondragover={handleDragOver}
		ondrop={handleDrop}
		role="button"
		tabindex="0"
		onclick={() => fileInput.click()}
		onkeydown={handleKeyDown}
	>
		<div class="drop-zone-content">
			<Upload size={80} />
			<p class="drop-text">
				<strong>Нажми тут для выбора файлов</strong> или перетащи их
			</p>
			<p class="drop-subtext">
				Максимальный размер: {formatFileSize(maxFileSize)}
			</p>
		</div>
	</div>

	<input
		bind:this={fileInput}
		type="file"
		{multiple}
		accept="text/html, text/css, application/javascript, image/png, image/jpeg, image/jpeg, image/svg+xml, application/x-zip-compressed"
		onchange={handleFileSelect}
		style="display: none;"
		{disabled}
	/>

	{#if selectedFiles.length > 0}
		<Card style="width: 100%;">
			<div class="selected-files">
				<div class="files-header">
					<h3>Выбранные файлы ({selectedFiles.length})</h3>
					<Button kind="danger" onclick={clearAll}>Очистить</Button>
				</div>

				<div class="files-list">
					{#each selectedFiles as file, index (index)}
						<div class="file-item">
							<div class="file-info">
								<span class="file-name">{file.name}</span>
								<span class="file-size">{formatFileSize(file.size)}</span>
							</div>

							<IconButton onclick={() => removeFile(index)} title="Удалить">
								<X size={20} />
							</IconButton>
						</div>
					{/each}
				</div>
			</div>
		</Card>
	{/if}

	<ProgressBar bind:uploadProgress />

	<div class="actions">
		<Button
			filled
			onclick={uploadFiles}
			disabled={disabled || uploading || selectedFiles.length === 0}
		>
			Загрузить
		</Button>
		<Button onclick={onclose} disabled={disabled || uploading}>Закрыть</Button>
	</div>
</div>

<style>
	.success-message {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.file-upload-container {
		max-width: 500px;

		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.drop-zone {
		border: 1px dashed var(--color-action-50);
		border-radius: var(--border-radius-regular);
		padding: 2rem;
		text-align: center;
		cursor: pointer;
		transition: var(--transition);
	}

	.drop-zone:focus {
		outline: var(--border-width) var(--border-style) var(--color-action);
		outline-offset: 2px;
	}

	.drop-zone-content {
		pointer-events: none;
	}

	.files-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		border-bottom: var(--border-regular);
	}

	.files-header h3 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
	}

	.files-list {
		max-height: 200px;
		overflow-y: auto;
	}

	.file-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1rem;
		border-bottom: var(--border-regular);
	}

	.file-item:last-child {
		border-bottom: none;
	}

	.file-info {
		flex: 1;
		min-width: 0;
	}

	.file-name {
		display: block;
		font-weight: 800;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.file-size {
		color: #aaa;
	}

	.actions {
		display: flex;
		gap: 15px;
		align-items: center;
	}
</style>
