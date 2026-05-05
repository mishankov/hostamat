<script lang="ts">
	import { Alert, Button, Card, Details, IconButton, Link, ProgressBar } from '$lib/components';
	import { Upload, X } from '$lib/icons';
	import { closeUploadSession, createUploadSession, enableFormCollection } from '$lib/api-client';
	import type { UploadedPage } from '$lib/api-client/upload-session';
	import type {
		DiscoveredForm,
		DiscoveredFormField,
		FormDiscovery
	} from '$lib/api-client/form-submissions';
	import { uploadFile } from './utils';

	interface Props {
		uploadUrl?: string;
		maxFileSize?: number;
		multiple?: boolean;
		disabled?: boolean;
		additionalHeaders?: Record<string, string>;
		siteId: string;
		onfilesselected?: (event: { files: File[] }) => void;
		onuploadcomplete?: (event: { fileCount: number }) => void | Promise<void>;
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
	let availablePages: UploadedPage[] = $state([]);
	let hasHomePage: boolean | null = $state(null);
	let isDragOver = $state(false);
	let closedDiscoveryId = $state('');
	let formDiscovery = $state.raw<FormDiscovery | null>(null);
	let selectedFormIds = $state<string[]>([]);
	let enabledFormIds = $state<string[]>([]);
	let enablingForms = $state(false);
	let enableError = $state('');
	let enableStatus = $state('');

	const discoveredForms = $derived(formDiscovery?.forms ?? []);
	const supportedForms = $derived(discoveredForms.filter((form) => isSupportedForm(form)));
	const unsupportedForms = $derived(discoveredForms.filter((form) => !isSupportedForm(form)));
	const selectedSupportedCount = $derived(
		supportedForms.filter((form) => selectedFormIds.includes(form.id)).length
	);

	export function resetUploadState() {
		uploadedFiles = [];
		skippedFiles = [];
		availablePages = [];
		hasHomePage = null;
		uploadProgress = 0;
		error = '';
		closedDiscoveryId = '';
		formDiscovery = null;
		selectedFormIds = [];
		enabledFormIds = [];
		enablingForms = false;
		enableError = '';
		enableStatus = '';
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
				// Close upload session on error
				const { error: closeError } = await closeUploadSession(uploadSession);
				if (closeError) {
					error = closeError;
				}
				return;
			}

			uploadProgress = ((i + 1) / selectedFiles.length) * 100;
		}

		const uploadedCount = selectedFiles.length;
		selectedFiles = [];
		fileInput.value = '';

		const closeUploadResp = await closeUploadSession(uploadSession);

		if (closeUploadResp.data) {
			uploadedFiles = closeUploadResp.data.uploadedFiles.map(
				(file) => `${file.name} (${formatFileSize(file.size)})`
			);
			skippedFiles = closeUploadResp.data.skippedFiles.map((file) => file.name);
			availablePages = closeUploadResp.data.availablePages;
			hasHomePage = closeUploadResp.data.hasHomePage;
			const normalizedDiscovery = normalizeFormDiscovery(closeUploadResp.data);
			formDiscovery = normalizedDiscovery;
			closedDiscoveryId = normalizedDiscovery?.discoveryId ?? '';
			selectedFormIds = (normalizedDiscovery?.forms ?? [])
				.filter((form) => isSupportedForm(form))
				.map((form) => form.id);
			await onuploadcomplete?.({ fileCount: uploadedCount });
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
		uploadedFiles = [];
		skippedFiles = [];
		availablePages = [];
		hasHomePage = null;
		closedDiscoveryId = '';
		formDiscovery = null;
		selectedFormIds = [];
		enabledFormIds = [];
		enableError = '';
		enableStatus = '';
	}

	function handleKeyDown(event: KeyboardEvent): void {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			fileInput.click();
		}
	}

	function normalizeFormDiscovery(summary: {
		formDiscovery?: FormDiscovery | null;
		forms?: DiscoveredForm[];
	}): FormDiscovery | null {
		const forms = summary.formDiscovery?.forms ?? summary.forms ?? [];
		if (forms.length === 0) {
			return null;
		}

		const supported = forms.filter((form) => isSupportedForm(form)).length;

		return {
			totalForms: summary.formDiscovery?.totalForms ?? forms.length,
			supportedForms: summary.formDiscovery?.supportedForms ?? supported,
			unsupportedForms: summary.formDiscovery?.unsupportedForms ?? forms.length - supported,
			discoveryId: summary.formDiscovery?.discoveryId,
			forms
		};
	}

	function isSupportedForm(form: DiscoveredForm): boolean {
		return Boolean(form.id) && form.supported !== false && !form.unsupportedReason && !form.reason;
	}

	function formTitle(form: DiscoveredForm): string {
		return form.label || form.displayName || form.name || 'Форма без названия';
	}

	function formPage(form: DiscoveredForm): string {
		return form.pagePath || form.path || '/';
	}

	function formAction(form: DiscoveredForm): string {
		return form.action || form.originalAction || 'не указан';
	}

	function unsupportedReason(form: DiscoveredForm): string {
		return form.unsupportedReason || form.reason || 'эта форма пока не поддерживается';
	}

	function fieldSummary(form: DiscoveredForm): string {
		if (Array.isArray(form.fields) && form.fields.length > 0) {
			return form.fields.map((field) => fieldLabel(field)).join(', ');
		}

		if (Array.isArray(form.fieldSummary)) {
			return form.fieldSummary.join(', ');
		}

		return form.fieldSummary || 'поля без названий';
	}

	function fieldLabel(field: DiscoveredFormField | string): string {
		if (typeof field === 'string') {
			return field;
		}

		const type = field.type ? ` (${field.type})` : '';
		const required = field.required ? ' *' : '';
		return `${field.name}${type}${required}`;
	}

	function toggleForm(formId: string, checked: boolean): void {
		enableError = '';
		enableStatus = '';

		if (checked) {
			selectedFormIds = [...new Set([...selectedFormIds, formId])];
			return;
		}

		selectedFormIds = selectedFormIds.filter((id) => id !== formId);
	}

	async function enableSelectedForms(): Promise<void> {
		if (!closedDiscoveryId || selectedFormIds.length === 0) {
			return;
		}

		enablingForms = true;
		enableError = '';
		enableStatus = '';

		const { data, error: enableRespError } = await enableFormCollection({
			siteId,
			discoveryId: closedDiscoveryId,
			formIds: selectedFormIds
		});

		enablingForms = false;

		if (enableRespError) {
			enableError = enableRespError;
			return;
		}

		if (data) {
			enabledFormIds = selectedFormIds;
			enableStatus = `Сбор ответов включён для форм: ${selectedSupportedCount}`;
		}
	}
</script>

<div class="file-upload-container">
	{#if error}
		<Alert variant="error" message={error} />
	{/if}

	{#if uploadProgress == 100 && !error}
		<Alert variant="success" title="Загрузка успешно завершена">
			<div class="success-message">
				{#if availablePages.length > 0}
					<div class="available-pages">
						<Details summary="Доступные страницы ({availablePages.length})">
							<ul class="page-list">
								{#each availablePages as page (page.url)}
									<li>
										<Link href={page.url} target="_blank" alwaysUnderline>{page.url}</Link>
										<span>{page.fileName}</span>
									</li>
								{/each}
							</ul>
						</Details>
					</div>
				{/if}
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

	{#if uploadProgress == 100 && hasHomePage === false && !error}
		<Alert variant="warning" title="Главная страница не найдена">
			<p>
				Загрузи файл <strong>index.html</strong> в корень сайта или ZIP-архива, чтобы главный адрес сайта
				открывался.
			</p>
		</Alert>
	{/if}

	{#if uploadProgress == 100 && formDiscovery && !error}
		<section class="form-discovery" aria-labelledby="form-discovery-title">
			<div class="form-discovery-header">
				<div>
					<h3 id="form-discovery-title">Найдены формы на сайте</h3>
					<p>
						Хостамат может собирать ответы с выбранных форм. При включении действия форм будут
						заменены на endpoint Хостамата.
					</p>
				</div>
				<span class="form-counter">{formDiscovery.totalForms ?? discoveredForms.length}</span>
			</div>

			{#if enableError}
				<Alert variant="error" message={enableError} />
			{/if}

			{#if enableStatus}
				<Alert variant="success" message={enableStatus} />
			{/if}

			{#if supportedForms.length > 0}
				<div class="detected-form-list" aria-label="Поддерживаемые формы">
					{#each supportedForms as form (form.id)}
						<label class={['detected-form', enabledFormIds.includes(form.id) && 'is-enabled']}>
							<input
								type="checkbox"
								checked={selectedFormIds.includes(form.id)}
								disabled={enablingForms || enabledFormIds.includes(form.id)}
								onchange={(event) =>
									toggleForm(form.id, (event.currentTarget as HTMLInputElement).checked)}
							/>
							<span class="form-card-body">
								<span class="form-mainline">
									<strong>{formTitle(form)}</strong>
									<span>{formPage(form)}</span>
								</span>
								<span class="form-meta">
									<span>Метод: {(form.method || 'get').toUpperCase()}</span>
									<span>Action: {formAction(form)}</span>
								</span>
								<span class="field-line">Поля: {fieldSummary(form)}</span>
							</span>
						</label>
					{/each}
				</div>

				<div class="form-actions">
					<Button
						emphasis="solid"
						disabled={enablingForms || selectedFormIds.length === 0 || enabledFormIds.length > 0}
						onclick={enableSelectedForms}
					>
						{enablingForms ? 'Включаем...' : `Включить сбор: ${selectedSupportedCount}`}
					</Button>
					<Button
						disabled={enablingForms || enabledFormIds.length > 0}
						onclick={() => {
							selectedFormIds = [];
							enableError = '';
							enableStatus = 'Формы оставлены без изменений';
						}}>Не включать</Button
					>
				</div>
			{/if}

			{#if unsupportedForms.length > 0}
				<Details summary="Недоступные формы ({unsupportedForms.length})">
					<ul class="unsupported-list">
						{#each unsupportedForms as form (form.id || formPage(form) + formTitle(form))}
							<li>
								<strong>{formTitle(form)}</strong>
								<span>{formPage(form)}: {unsupportedReason(form)}</span>
							</li>
						{/each}
					</ul>
				</Details>
			{/if}
		</section>
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

	.page-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		list-style: none;
		padding: 0;
	}

	.page-list li {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		min-width: 0;
	}

	.page-list :global(a) {
		color: var(--color-accent-text);
		font-weight: 700;
		overflow-wrap: anywhere;
		text-decoration-thickness: 0.08em;
		text-underline-offset: 0.16em;
	}

	.page-list span {
		color: var(--color-text-secondary);
		font-family: var(--font-mono);
		font-size: var(--font-size-caption);
		overflow-wrap: anywhere;
	}

	.file-upload-container {
		width: 100%;
		max-width: 42rem;
		min-width: 0;

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

	.file-item:nth-child(1) {
		animation-delay: 0ms;
	}
	.file-item:nth-child(2) {
		animation-delay: 30ms;
	}
	.file-item:nth-child(3) {
		animation-delay: 60ms;
	}
	.file-item:nth-child(4) {
		animation-delay: 90ms;
	}
	.file-item:nth-child(5) {
		animation-delay: 120ms;
	}
	.file-item:nth-child(6) {
		animation-delay: 150ms;
	}
	.file-item:nth-child(7) {
		animation-delay: 180ms;
	}
	.file-item:nth-child(8) {
		animation-delay: 210ms;
	}
	.file-item:nth-child(9) {
		animation-delay: 240ms;
	}
	.file-item:nth-child(10) {
		animation-delay: 270ms;
	}

	.actions {
		display: flex;
		gap: var(--space-lg);
		align-items: center;
	}

	.form-discovery {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		min-width: 0;
		padding: var(--space-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-background-object);
		box-shadow: var(--shadow-sm);
	}

	.form-discovery-header {
		display: grid;
		grid-template-columns: minmax(0, 1fr) auto;
		gap: var(--space-sm);
		align-items: start;
	}

	.form-discovery h3 {
		margin: 0 0 var(--space-xs);
		font-size: 1rem;
		font-weight: 800;
	}

	.form-discovery p {
		margin: 0;
		color: var(--color-text-secondary);
		line-height: 1.5;
	}

	.form-counter {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: 2.25rem;
		height: 2.25rem;
		border: 1px solid var(--color-accent-border);
		border-radius: 999px;
		background: var(--color-accent-surface);
		color: var(--color-accent-text);
		font-weight: 800;
	}

	.detected-form-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}

	.detected-form {
		display: grid;
		grid-template-columns: auto minmax(0, 1fr);
		gap: var(--space-sm);
		align-items: start;
		min-width: 0;
		padding: var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-background);
		cursor: pointer;
		transition:
			border-color 160ms var(--ease-out-quart),
			background-color 160ms var(--ease-out-quart),
			box-shadow 160ms var(--ease-out-quart);
	}

	.detected-form:hover {
		border-color: var(--color-accent-border);
		background: var(--color-accent-surface);
	}

	.detected-form.is-enabled {
		border-color: var(--color-success-border);
		background: var(--color-success-surface);
	}

	.detected-form input {
		width: 1rem;
		height: 1rem;
		margin-top: 0.2rem;
		accent-color: var(--color-accent-strong);
	}

	.form-card-body {
		display: flex;
		min-width: 0;
		flex-direction: column;
		gap: var(--space-xs);
	}

	.form-mainline {
		display: grid;
		gap: var(--space-xs);
		min-width: 0;
	}

	.form-mainline strong {
		font-weight: 800;
	}

	.form-mainline span,
	.form-meta,
	.field-line {
		color: var(--color-text-secondary);
		overflow-wrap: anywhere;
	}

	.form-meta {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		font-family: var(--font-mono);
		font-size: var(--font-size-caption);
	}

	.field-line {
		font-size: 0.86rem;
		line-height: 1.45;
	}

	.form-actions {
		display: flex;
		gap: var(--space-sm);
		align-items: center;
		flex-wrap: wrap;
	}

	.unsupported-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		margin: var(--space-sm) 0 0;
		padding: 0;
		list-style: none;
	}

	.unsupported-list li {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		padding: var(--space-sm) 0;
		border-top: 1px solid var(--color-border);
	}

	.unsupported-list span {
		color: var(--color-text-secondary);
		overflow-wrap: anywhere;
	}

	@media (max-width: 560px) {
		.form-discovery-header {
			grid-template-columns: 1fr;
		}

		.detected-form {
			grid-template-columns: 1fr;
		}
	}
</style>
