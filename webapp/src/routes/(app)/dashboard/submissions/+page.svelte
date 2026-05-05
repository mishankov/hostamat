<script lang="ts">
	import type { PageProps } from './$types';
	import { Alert, Button, Select } from '$lib/components';
	import { Spinner, Trash } from '$lib/icons';
	import {
		deleteFormSubmission,
		type FormSubmission,
		type SubmissionFields
	} from '$lib/api-client/form-submissions';
	import { resolve } from '$app/paths';

	type FieldEntry = {
		key: string;
		name: string;
		value: unknown;
	};

	type FilterOption = {
		value: string;
		label: string;
	};

	let { data }: PageProps = $props();

	let submissions = $state.raw<FormSubmission[]>([]);
	let loading = $state(true);
	let errorMessage = $state('');
	let siteFilter = $state('all');
	let formFilter = $state('all');
	let pendingDeleteId = $state('');
	let deletingId = $state('');
	let expandedIds = $state.raw<Set<string>>(new Set());

	const siteOptions = $derived.by<FilterOption[]>(() => [
		{ value: 'all', label: 'Все сайты' },
		...uniqueBy(
			submissions.map((submission) => ({
				value: submissionSiteId(submission),
				label: submissionSiteLabel(submission)
			}))
		)
	]);

	const formOptions = $derived.by<FilterOption[]>(() => [
		{ value: 'all', label: 'Все формы' },
		...uniqueBy(
			submissions
				.filter((submission) => siteFilter === 'all' || submissionSiteId(submission) === siteFilter)
				.map((submission) => ({
					value: submissionFormId(submission),
					label: submissionFormLabel(submission)
				}))
		)
	]);

	const filteredSubmissions = $derived.by(() =>
		submissions.filter((submission) => {
			const matchesSite = siteFilter === 'all' || submissionSiteId(submission) === siteFilter;
			const matchesForm = formFilter === 'all' || submissionFormId(submission) === formFilter;
			return matchesSite && matchesForm;
		})
	);

	const activeCount = $derived(
		submissions.filter(
			(submission) => submissionActive(submission) && !submissionHistorical(submission)
		).length
	);

	const updateState = async () => {
		const submissionsResp = await data.submissionsResp;

		if (submissionsResp.data) {
			submissions = submissionsResp.data.submissions ?? [];
			errorMessage = '';
		} else {
			submissions = [];
			errorMessage = submissionsResp.error;
		}

		loading = false;
	};

	updateState();

	function uniqueBy(options: FilterOption[]): FilterOption[] {
		let seen: string[] = [];
		return options.filter((option) => {
			if (!option.value || seen.includes(option.value)) {
				return false;
			}

			seen = [...seen, option.value];
			return true;
		});
	}

	function submissionSiteId(submission: FormSubmission): string {
		return submission.siteId || submission.site?.id || 'unknown-site';
	}

	function submissionSiteLabel(submission: FormSubmission): string {
		const siteName = submission.siteName || submission.site?.name || 'Сайт без названия';
		const domain = submission.siteDomain || submission.site?.domain;
		return domain ? `${siteName} (${domain})` : siteName;
	}

	function submissionFormId(submission: FormSubmission): string {
		return submission.formId || submission.form?.id || 'unknown-form';
	}

	function submissionFormLabel(submission: FormSubmission): string {
		return (
			submission.formLabel ||
			submission.formName ||
			submission.form?.label ||
			submission.form?.name ||
			'Форма без названия'
		);
	}

	function submissionPage(submission: FormSubmission): string {
		return submission.pagePath || submission.page || submission.form?.pagePath || '/';
	}

	function submissionActive(submission: FormSubmission): boolean {
		return submission.active ?? submission.form?.active ?? true;
	}

	function submissionHistorical(submission: FormSubmission): boolean {
		return submission.historical ?? submission.form?.historical ?? false;
	}

	function fieldEntries(fields: SubmissionFields): FieldEntry[] {
		if (Array.isArray(fields)) {
			return fields.map((field, index) => ({
				key: `${field.name}:${index}`,
				name: field.name,
				value: field.value
			}));
		}

		return Object.entries(fields ?? {}).map(([name, value]) => ({ key: name, name, value }));
	}

	function fieldPreview(fields: SubmissionFields): string {
		const entries = fieldEntries(fields).slice(0, 2);
		if (entries.length === 0) {
			return 'без заполненных полей';
		}

		return entries.map((entry) => `${entry.name}: ${formatFieldValue(entry.value)}`).join(', ');
	}

	function formatFieldValue(value: unknown): string {
		if (value === null || value === undefined || value === '') {
			return 'пусто';
		}

		if (typeof value === 'boolean') {
			return value ? 'да' : 'нет';
		}

		if (Array.isArray(value)) {
			return value.length > 0
				? value.map((item) => formatFieldValue(item)).join(', ')
				: 'не выбрано';
		}

		if (typeof value === 'object') {
			return JSON.stringify(value);
		}

		return String(value);
	}

	function formatDate(value: string): string {
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) {
			return value;
		}

		return new Intl.DateTimeFormat('ru-RU', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		}).format(date);
	}

	function handleSiteFilterChange(event: Event): void {
		siteFilter = (event.currentTarget as HTMLSelectElement).value;
		formFilter = 'all';
	}

	async function deleteSubmission(submissionId: string): Promise<void> {
		deletingId = submissionId;
		errorMessage = '';

		const { data: deleted, error } = await deleteFormSubmission(submissionId);
		deletingId = '';

		if (error) {
			errorMessage = error;
			return;
		}

		if (deleted) {
			submissions = submissions.filter((submission) => submission.id !== submissionId);
			pendingDeleteId = '';
			const nextExpanded = new Set(expandedIds);
			nextExpanded.delete(submissionId);
			expandedIds = nextExpanded;
		}
	}

	function toggleExpanded(id: string): void {
		const next = new Set(expandedIds);
		if (next.has(id)) {
			next.delete(id);
		} else {
			next.add(id);
		}
		expandedIds = next;
	}

	function isExpanded(id: string): boolean {
		return expandedIds.has(id);
	}
</script>

<svelte:head>
	<title>Ответы форм | Хостамат</title>
	<meta
		name="description"
		content="Просматривайте ответы, которые посетители отправили через формы на статических сайтах."
	/>
</svelte:head>

<div class="dashboard-page">
	<div class="headline">
		<h1>Ответы форм</h1>
		<div class="headline-meta">
			<p>Все отправки с подключённых HTML-форм</p>
			{#if !loading}
				<span class="submission-counter">
					ответов: {submissions.length} · активных форм: {activeCount}
				</span>
			{/if}
		</div>
	</div>

	{#if errorMessage}
		<Alert variant="error" message={errorMessage} />
	{/if}

	{#if loading}
		<div class="inbox-shell">
			<div class="loading-state" aria-live="polite">
				<Spinner size={50} />
			</div>
		</div>
	{:else if submissions.length === 0}
		<div class="inbox-shell">
			<div class="empty-state">
				<h2>Ответов пока нет</h2>
				<p>
					После загрузки сайта включите сбор для найденных форм. Новые отправки появятся здесь сразу
					после обработки сервером.
				</p>
			</div>
		</div>
	{:else}
		<div class="filter-toolbar" aria-label="Фильтры ответов">
			<div class="filter-controls">
				<Select
					id="submission-site-filter"
					label="Сайт"
					options={siteOptions}
					bind:value={siteFilter}
					onchange={handleSiteFilterChange}
				/>
				<Select
					id="submission-form-filter"
					label="Форма"
					options={formOptions}
					bind:value={formFilter}
				/>
			</div>
			{#if filteredSubmissions.length !== submissions.length}
				<span class="filter-results">
					{filteredSubmissions.length} из {submissions.length}
				</span>
			{/if}
		</div>

		<div class="inbox-shell">
			{#if filteredSubmissions.length === 0}
				<div class="empty-state compact">
					<h2>По фильтрам ничего нет</h2>
					<p>Выберите другой сайт или форму.</p>
				</div>
			{:else}
				<div class="submissions-list" role="list">
					{#each filteredSubmissions as submission, index (submission.id)}
						<article
							class="submission-card"
							class:expanded={isExpanded(submission.id)}
							style="animation-delay: {Math.min(index * 45, 400)}ms"
							role="listitem"
						>
							<div class="card-header">
								<div class="header-col timestamp-col">
									<time datetime={submission.createdAt}>
										{formatDate(submission.createdAt)}
									</time>
								</div>

								<div class="header-col meta-col">
									<div class="meta-line">
										<span class="site-label">{submissionSiteLabel(submission)}</span>
										<span class="meta-divider" aria-hidden="true">/</span>
										<span class="form-label">{submissionFormLabel(submission)}</span>
										{#if submissionHistorical(submission)}
											<span class="status-pill historical">историческая</span>
										{:else if !submissionActive(submission)}
											<span class="status-pill inactive">выключена</span>
										{/if}
									</div>
									{#if fieldEntries(submission.fields).length > 0}
										<div class="preview-line">
											{fieldPreview(submission.fields)}
											{#if fieldEntries(submission.fields).length > 2}
												<span class="preview-more">
													+{fieldEntries(submission.fields).length - 2}
												</span>
											{/if}
										</div>
									{/if}
								</div>

								<div class="header-col path-col">
									<code class="page-path">{submissionPage(submission)}</code>
								</div>

								<div class="header-col action-col">
									<button
										type="button"
										class="expand-toggle"
										aria-expanded={isExpanded(submission.id)}
										aria-label={isExpanded(submission.id)
											? 'Скрыть поля ответа'
											: 'Показать поля ответа'}
										onclick={() => toggleExpanded(submission.id)}
									>
										<span class="expand-label">
											{isExpanded(submission.id) ? 'Свернуть' : 'Поля'}
										</span>
										<svg
											class="expand-icon"
											class:rotated={isExpanded(submission.id)}
											viewBox="0 0 24 24"
											width="14"
											height="14"
										>
											<polyline
												points="6 9 12 15 18 9"
												fill="none"
												stroke="currentColor"
												stroke-width="2.5"
												stroke-linecap="round"
												stroke-linejoin="round"
											/>
										</svg>
									</button>

									{#if pendingDeleteId === submission.id}
										<div class="confirm-actions">
											<Button
												size="xs"
												tone="danger"
												emphasis="solid"
												disabled={deletingId === submission.id}
												onclick={() => deleteSubmission(submission.id)}
											>
												Удалить
											</Button>
											<Button
												size="xs"
												disabled={deletingId === submission.id}
												onclick={() => (pendingDeleteId = '')}
											>
												Отмена
											</Button>
										</div>
									{:else}
										<Button
											size="xs"
											tone="danger"
											onclick={() => (pendingDeleteId = submission.id)}
										>
											<Trash size={14} />
											<span class="delete-label">Удалить</span>
										</Button>
									{/if}
								</div>
							</div>

							<div class="card-body" class:open={isExpanded(submission.id)}>
								<div class="card-body-inner">
									<div class="field-grid">
										{#each fieldEntries(submission.fields) as field (field.key)}
											<div class="field-row">
												<span class="data-label field-name">{field.name}</span>
												<span class="field-value">{formatFieldValue(field.value)}</span>
											</div>
										{/each}
									</div>
								</div>
							</div>
						</article>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	.dashboard-page {
		display: flex;
		flex-direction: column;
		gap: var(--space-xl);
		min-width: 0;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}

	.headline h1 {
		font-weight: 800;
		letter-spacing: -0.03em;
	}

	.headline p {
		max-width: 64ch;
		color: var(--color-text-secondary);
	}

	.headline-meta {
		display: flex;
		align-items: baseline;
		gap: var(--space-lg);
		flex-wrap: wrap;
	}

	.submission-counter {
		font-family: var(--font-mono);
		font-size: var(--font-size-caption);
		font-weight: 500;
		color: var(--color-text-tertiary);
		white-space: nowrap;
		letter-spacing: -0.01em;
	}

	.filter-toolbar {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: var(--space-lg);
		flex-wrap: wrap;
		padding: var(--space-md) var(--space-lg);
		background: var(--color-background-object);
		border: var(--border-regular);
		border-radius: var(--radius-md);
		box-shadow: var(--shadow-sm);
	}

	.filter-controls {
		display: grid;
		grid-template-columns: minmax(12rem, 18rem) minmax(12rem, 18rem);
		gap: var(--space-md);
		align-items: end;
	}

	.filter-results {
		font-family: var(--font-mono);
		font-size: var(--font-size-caption);
		font-weight: 500;
		color: var(--color-text-tertiary);
		white-space: nowrap;
		padding-bottom: 0.5rem;
	}

	.inbox-shell {
		min-width: 0;
	}

	.loading-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 10.625rem;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		padding: var(--space-xl) 0;
		max-width: 42rem;
	}

	.empty-state.compact {
		padding: var(--space-lg) 0;
	}

	.empty-state h2 {
		margin: 0;
		font-size: 1.15rem;
		font-weight: 800;
	}

	.empty-state p {
		margin: 0;
		color: var(--color-text-secondary);
		line-height: 1.55;
	}

	/* ===== SUBMISSION LIST ===== */

	.submissions-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}

	.submission-card {
		--card-border: var(--color-border);
		--card-bg: var(--color-background-object);
		--card-hover-bg: var(--color-background-object);

		position: relative;
		display: flex;
		flex-direction: column;
		background: var(--card-bg);
		border: 1px solid var(--card-border);
		border-radius: var(--radius-md);
		box-shadow: var(--shadow-sm);
		transition:
			border-color 160ms var(--ease-out-quart),
			box-shadow 160ms var(--ease-out-quart),
			transform 160ms var(--ease-out-quart);
		animation: card-enter 320ms var(--ease-out-quint) both;
	}

	@keyframes card-enter {
		from {
			opacity: 0;
			transform: translateY(8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.submission-card:hover {
		--card-border: var(--color-border-strong);
		box-shadow: var(--shadow-md);
		transform: translateY(-1px);
	}

	.submission-card.expanded {
		--card-border: var(--color-border-strong);
		box-shadow: var(--shadow-md);
	}

	/* ===== CARD HEADER ===== */

	.card-header {
		display: grid;
		grid-template-columns: 9.5rem 1fr auto auto;
		align-items: start;
		gap: var(--space-md);
		padding: var(--space-md) var(--space-lg);
		min-width: 0;
	}

	.header-col {
		min-width: 0;
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		padding-top: 0.15rem;
	}

	.timestamp-col {
		flex-shrink: 0;
	}

	.timestamp-col time {
		font-family: var(--font-mono);
		font-size: var(--font-size-caption);
		font-weight: 500;
		color: var(--color-text-secondary);
		letter-spacing: -0.01em;
		white-space: nowrap;
	}

	.meta-col {
		min-width: 0;
	}

	.meta-line {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		flex-wrap: wrap;
		min-width: 0;
	}

	.site-label {
		font-size: 0.82rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		white-space: nowrap;
	}

	.meta-divider {
		color: var(--color-border-strong);
		font-weight: 700;
	}

	.form-label {
		font-size: 0.92rem;
		font-weight: 700;
		color: var(--color-text-regular);
		letter-spacing: -0.01em;
	}

	.status-pill {
		width: fit-content;
		border-radius: 999px;
		padding: 0.1rem 0.5rem;
		font-size: 0.7rem;
		font-weight: 700;
		letter-spacing: 0.02em;
		white-space: nowrap;
		user-select: none;
	}

	.status-pill.historical {
		border: 1px solid var(--color-warning-border);
		background: var(--color-warning-surface);
		color: var(--color-warning-text);
	}

	.status-pill.inactive {
		border: 1px solid var(--color-border);
		background: var(--color-neutral-surface);
		color: var(--color-text-secondary);
	}

	.preview-line {
		font-size: 0.8rem;
		font-weight: 400;
		color: var(--color-text-tertiary);
		line-height: 1.4;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.preview-more {
		font-family: var(--font-mono);
		font-size: 0.72rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		margin-left: 0.15rem;
	}

	.path-col {
		flex-shrink: 0;
		align-items: flex-end;
		text-align: right;
	}

	.page-path {
		font-family: var(--font-mono);
		font-size: var(--font-size-caption);
		font-weight: 500;
		color: var(--color-text-tertiary);
		background: none;
		padding: 0;
		letter-spacing: -0.01em;
		white-space: nowrap;
	}

	.action-col {
		flex-shrink: 0;
		flex-direction: row;
		align-items: center;
		gap: var(--space-sm);
		padding-top: 0;
	}

	.expand-toggle {
		display: inline-flex;
		align-items: center;
		gap: var(--space-xs);
		padding: 0.35rem 0.6rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		background: transparent;
		color: var(--color-text-secondary);
		font-family: var(--font-body);
		font-size: 0.78rem;
		font-weight: 600;
		cursor: pointer;
		transition:
			background-color 140ms var(--ease-out-quart),
			border-color 140ms var(--ease-out-quart),
			color 140ms var(--ease-out-quart);
	}

	.expand-toggle:hover {
		background: var(--color-neutral-surface);
		border-color: var(--color-border-strong);
		color: var(--color-text-regular);
	}

	.expand-label {
		white-space: nowrap;
	}

	.expand-icon {
		flex-shrink: 0;
		transition: transform 200ms var(--ease-out-quart);
	}

	.expand-icon.rotated {
		transform: rotate(180deg);
	}

	.confirm-actions {
		display: flex;
		gap: var(--space-xs);
		flex-wrap: wrap;
	}

	.delete-label {
		white-space: nowrap;
	}

	/* ===== CARD BODY (expandable) ===== */

	.card-body {
		display: grid;
		grid-template-rows: 0fr;
		transition: grid-template-rows 280ms var(--ease-out-quart);
	}

	.card-body.open {
		grid-template-rows: 1fr;
	}

	.card-body-inner {
		overflow: hidden;
	}

	.submission-card.expanded .card-header {
		border-bottom: 1px solid var(--color-border);
	}

	.field-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(min(100%, 20rem), 1fr));
		gap: var(--space-xs) var(--space-xl);
		padding: var(--space-md) var(--space-lg);
		background: color-mix(in hsl, var(--color-background) 60%, transparent);
		border-radius: 0 0 var(--radius-md) var(--radius-md);
	}

	.field-row {
		display: grid;
		grid-template-columns: 7rem 1fr;
		gap: var(--space-md);
		align-items: baseline;
		padding: var(--space-sm) 0;
		border-bottom: 1px solid var(--color-border);
		min-width: 0;
	}

	.field-row:last-child {
		border-bottom: none;
	}

	.field-name {
		overflow-wrap: anywhere;
		min-width: 0;
	}

	.field-value {
		font-size: 0.88rem;
		font-weight: 500;
		color: var(--color-text-regular);
		white-space: pre-wrap;
		overflow-wrap: anywhere;
		min-width: 0;
		line-height: 1.45;
	}

	/* ===== RESPONSIVE ===== */

	@media (max-width: 980px) {
		.filter-controls {
			grid-template-columns: 1fr;
			width: 100%;
		}

		.card-header {
			grid-template-columns: 1fr auto;
			gap: var(--space-sm) var(--space-md);
		}

		.timestamp-col {
			order: 2;
			grid-column: 1;
		}

		.meta-col {
			order: 1;
			grid-column: 1 / -1;
		}

		.path-col {
			order: 3;
			grid-column: 1;
			align-items: flex-start;
			text-align: left;
		}

		.action-col {
			order: 4;
			grid-column: 2;
			grid-row: 2 / span 2;
			align-self: center;
		}

		.field-grid {
			grid-template-columns: 1fr;
		}
	}

	@media (max-width: 520px) {
		.card-header {
			grid-template-columns: 1fr;
		}

		.action-col {
			grid-column: 1;
			grid-row: auto;
			justify-content: space-between;
			width: 100%;
		}

		.path-col {
			grid-column: 1;
		}

		.timestamp-col {
			grid-column: 1;
		}

		.delete-label {
			display: none;
		}
	}
</style>
