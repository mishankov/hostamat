<script lang="ts">
	import { Alert, Button, Card, Details, IconButton, ProgressBar } from '$lib/components';
	import { Upload, X } from '$lib/icons';
	import { closeUploadSession, createUploadSession } from '$lib/api-client';
	import { uploadFile } from './utils';

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
	let isDragOver = $state(false);

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
		selectedFiles = selectedFiles.concat(files);

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
		isDragOver = true;
		if (event.dataTransfer) {
			event.dataTransfer.dropEffect = 'copy';
		}
	}

	function handleDragEnter(event: DragEvent): void {
		event.preventDefault();
		isDragOver = true;
	}

	function handleDragLeave(event: DragEvent): void {
		const target = event.currentTarget as HTMLElement | null;
		const relatedTarget = event.relatedTarget as Node | null;
		if (target && relatedTarget && target.contains(relatedTarget)) {
			return;
		}

		isDragOver = false;
	}

	function handleDrop(event: DragEvent): void {
		event.preventDefault();
		isDragOver = false;
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

		// Validate file sizes
		const filesOverallSize = selectedFiles.reduce((acc, file) => acc + file.size, 0);
		if (filesOverallSize > maxFileSize) {
			error = `Файлы слишком большие. Максимальный размер: ${formatFileSize(maxFileSize)}`;
			uploading = false;
			return;
		}

		const { data: uploadSession, error: errorText } = await createUploadSession(siteId);

		if (errorText !== null) {
			error = 'Ошибка при создании сессии загрузки';
			uploading = false;
			return;
		}

		for (let i = 0; i < selectedFiles.length; i++) {
			const file = selectedFiles[i];
			const formData = new FormData();
			formData.append('file', file);

			const headers = new Headers();
			Object.entries(additionalHeaders).forEach(([key, value]) => {
				headers.set(key, value);
			});
			headers.set('H-UploadSession-Id', uploadSession);

			const result = await uploadFile(uploadUrl, headers, formData);

			if (result !== true) {
				error = result.errorMessage;
				uploading = false;
				onuploaderror?.({ error });
				return;
			}

			uploadProgress = ((i + 1) / selectedFiles.length) * 100;
		}

		const uploadedCount = selectedFiles.length;
		selectedFiles = [];
		fileInput.value = '';

		onuploadcomplete?.({ fileCount: uploadedCount });

		const closeUploadResp = await closeUploadSession(uploadSession);

		if (closeUploadResp.data) {
			uploadedFiles = closeUploadResp.data
				.filter((file) => file.uploaded)
				.map((file) => `${file.name} (${formatFileSize(file.size)})`);
			skippedFiles = closeUploadResp.data.filter((file) => !file.uploaded).map((file) => file.name);
		} else {
			error = closeUploadResp.error;
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
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			fileInput.click();
		}
	}
</script>

<div class="file-upload-container">
	{#if error}
		<Alert variant="error" message={error} />
	{/if}

	{#if uploadProgress == 100}
		<Alert variant="success" title="Загрузка успешно завершена">
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
		</Alert>
	{/if}

	<div
		class="drop-zone"
		class:drag-over={isDragOver}
		ondragenter={handleDragEnter}
		ondragover={handleDragOver}
		ondragleave={handleDragLeave}
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
		class="file-input-hidden"
		{disabled}
	/>

	{#if selectedFiles.length > 0}
		<div class="selected-files-shell">
			<Card>
				<div class="selected-files">
					<div class="files-header">
						<h3>Выбранные файлы ({selectedFiles.length})</h3>
						<Button tone="danger" onclick={clearAll}>Очистить</Button>
					</div>

					<div class="files-list">
						{#each selectedFiles as file, index (index)}
							<div class="file-item">
								<div class="file-info">
									<span class="file-name">{file.name}</span>
									<span class="file-size">{formatFileSize(file.size)}</span>
								</div>

								<IconButton label="Удалить файл" onclick={() => removeFile(index)} title="Удалить">
									<X size={20} />
								</IconButton>
							</div>
						{/each}
					</div>
				</div>
			</Card>
		</div>
	{/if}

	<ProgressBar bind:uploadProgress />

	<div class="actions">
		<Button
			emphasis="solid"
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
		gap: var(--space-sm);
	}

	.file-upload-container {
		max-width: 31.25rem;

		display: flex;
		flex-direction: column;
		gap: var(--space-xl);
	}

	.drop-zone {
		border: 1.5px dashed var(--color-accent-border);
		border-radius: var(--radius-lg);
		padding: var(--space-xl);
		text-align: center;
		cursor: pointer;
		transition:
			background-color 180ms var(--ease-out-quart),
			border-color 180ms var(--ease-out-quart),
			transform 180ms var(--ease-out-quart),
			box-shadow 180ms var(--ease-out-quart);
		background: var(--color-accent-surface);
	}

	.drop-zone:hover {
		border-color: var(--color-accent-strong);
		background: color-mix(in hsl, var(--color-accent-surface) 80%, var(--color-background-object));
	}

	.drop-zone.drag-over {
		border-color: var(--color-accent-strong);
		background: var(--color-background-object);
		border-style: solid;
		transform: scale(1.01);
		box-shadow: var(--shadow-md);
		animation: pulse-subtle 1.4s var(--ease-out-quart) infinite;
	}

	.drop-zone:focus-visible {
		outline: var(--border-width) var(--border-style) var(--color-accent-strong);
		outline-offset: 2px;
	}

	.drop-zone-content {
		pointer-events: none;
	}

	.files-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--space-lg) var(--space-lg);
		border-bottom: 1px solid var(--color-border);
	}

	.files-header h3 {
		margin: 0;
		font-size: 0.95rem;
		font-weight: 700;
	}

	.files-list {
		max-height: 12.5rem;
		overflow-y: auto;
		scrollbar-width: thin;
	}

	.file-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--space-md) var(--space-lg);
		border-bottom: 1px solid var(--color-border);
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
		color: var(--color-text-secondary);
	}

	.file-input-hidden {
		display: none;
	}

	.selected-files-shell {
		width: 100%;
		animation: fade-in-scale 220ms var(--ease-out-quart) both;
	}

	.file-item {
		animation: file-item-in 200ms var(--ease-out-quart) both;
	}

	.file-item:nth-child(1) { animation-delay: 0ms; }
	.file-item:nth-child(2) { animation-delay: 30ms; }
	.file-item:nth-child(3) { animation-delay: 60ms; }
	.file-item:nth-child(4) { animation-delay: 90ms; }
	.file-item:nth-child(5) { animation-delay: 120ms; }
	.file-item:nth-child(6) { animation-delay: 150ms; }
	.file-item:nth-child(7) { animation-delay: 180ms; }
	.file-item:nth-child(8) { animation-delay: 210ms; }
	.file-item:nth-child(9) { animation-delay: 240ms; }
	.file-item:nth-child(10) { animation-delay: 270ms; }

	.actions {
		display: flex;
		gap: var(--space-lg);
		align-items: center;
	}
</style>
