<script lang="ts">
	import type { PageProps } from './$types';
	import { Alert, Button, Select } from '$lib/components';
	import { Spinner, Trash } from '$lib/icons';
	import {
		deleteFormSubmission,
		type FormSubmission,
		type SubmissionFields
	} from '$lib/api-client/form-submissions';

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
		const entries = fieldEntries(fields).slice(0, 3);
		if (entries.length === 0) {
			return 'без заполненных полей';
		}

		return entries.map((entry) => `${entry.name}: ${formatFieldValue(entry.value)}`).join('; ');
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
		}
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
		<div class="filters" aria-label="Фильтры ответов">
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

		<div class="inbox-shell">
			{#if filteredSubmissions.length === 0}
				<div class="empty-state compact">
					<h2>По фильтрам ничего нет</h2>
					<p>Выберите другой сайт или форму.</p>
				</div>
			{:else}
				<table class="submission-table">
					<thead>
						<tr>
							<th>Получено</th>
							<th>Сайт</th>
							<th>Форма</th>
							<th>Страница</th>
							<th>Поля</th>
							<th class="action-column">Действие</th>
						</tr>
					</thead>
					<tbody>
						{#each filteredSubmissions as submission (submission.id)}
							<tr>
								<td data-label="Получено">
									<time datetime={submission.createdAt}>{formatDate(submission.createdAt)}</time>
								</td>
								<td data-label="Сайт">{submissionSiteLabel(submission)}</td>
								<td data-label="Форма">
									<div class="form-cell">
										<span>{submissionFormLabel(submission)}</span>
										{#if submissionHistorical(submission)}
											<span class="status-badge">историческая</span>
										{:else if !submissionActive(submission)}
											<span class="status-badge muted">выключена</span>
										{/if}
									</div>
								</td>
								<td data-label="Страница" class="path-cell">{submissionPage(submission)}</td>
								<td data-label="Поля" class="preview-cell">{fieldPreview(submission.fields)}</td>
								<td data-label="Действие" class="action-cell">
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
											<Trash size={16} />
											Удалить
										</Button>
									{/if}
								</td>
							</tr>
							<tr class="details-row">
								<td colspan="6">
									<details>
										<summary>Показать поля ответа</summary>
										<dl class="field-details">
											{#each fieldEntries(submission.fields) as field (field.key)}
												<div>
													<dt>{field.name}</dt>
													<dd>{formatFieldValue(field.value)}</dd>
												</div>
											{/each}
										</dl>
									</details>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
	{/if}
</div>

<style>
	.dashboard-page {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 0;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}

	.headline h1 {
		font-weight: 800;
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
		font-size: 0.82rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		white-space: nowrap;
	}

	.filters {
		display: grid;
		grid-template-columns: minmax(12rem, 18rem) minmax(12rem, 18rem);
		gap: var(--space-md);
		align-items: end;
		padding-bottom: var(--space-sm);
	}

	.inbox-shell {
		min-width: 0;
		border-top: 1px solid var(--color-border);
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

	.submission-table {
		width: 100%;
		border-collapse: collapse;
		table-layout: fixed;
	}

	th,
	td {
		padding: var(--space-md) var(--space-sm);
		border-bottom: 1px solid var(--color-border);
		text-align: left;
		vertical-align: top;
	}

	th {
		color: var(--color-text-secondary);
		font-size: 0.76rem;
		font-weight: 800;
		text-transform: uppercase;
	}

	th:nth-child(1) {
		width: 9.5rem;
	}

	th:nth-child(2),
	th:nth-child(3) {
		width: 15%;
	}

	th:nth-child(4) {
		width: 13%;
	}

	.action-column {
		width: 8.5rem;
	}

	td {
		color: var(--color-text-regular);
		line-height: 1.45;
		overflow-wrap: anywhere;
	}

	time,
	.path-cell {
		font-family: var(--font-mono);
		font-size: var(--font-size-caption);
		color: var(--color-text-secondary);
	}

	.preview-cell {
		color: var(--color-text-secondary);
	}

	.form-cell {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}

	.status-badge {
		width: fit-content;
		border: 1px solid var(--color-warning-border);
		border-radius: 999px;
		padding: 0.08rem 0.45rem;
		background: var(--color-warning-surface);
		color: var(--color-warning-text);
		font-size: 0.72rem;
		font-weight: 700;
	}

	.status-badge.muted {
		border-color: var(--color-border);
		background: var(--color-neutral-surface);
		color: var(--color-text-secondary);
	}

	.action-cell {
		min-width: 8rem;
	}

	.confirm-actions {
		display: flex;
		gap: var(--space-xs);
		flex-wrap: wrap;
	}

	.details-row td {
		padding-top: 0;
		background: color-mix(in hsl, var(--color-background-object) 72%, transparent);
	}

	details {
		padding: var(--space-xs) 0;
	}

	summary {
		width: fit-content;
		color: var(--color-accent-text);
		font-size: 0.86rem;
		font-weight: 700;
		cursor: pointer;
	}

	.field-details {
		display: grid;
		gap: var(--space-sm);
		margin: var(--space-md) 0 var(--space-sm);
	}

	.field-details div {
		display: grid;
		grid-template-columns: minmax(8rem, 14rem) minmax(0, 1fr);
		gap: var(--space-md);
		padding-top: var(--space-sm);
		border-top: 1px solid var(--color-border);
	}

	dt {
		color: var(--color-text-secondary);
		font-weight: 800;
		overflow-wrap: anywhere;
	}

	dd {
		margin: 0;
		white-space: pre-wrap;
		overflow-wrap: anywhere;
	}

	@media (max-width: 900px) {
		.filters {
			grid-template-columns: 1fr;
		}

		.submission-table,
		.submission-table thead,
		.submission-table tbody,
		.submission-table tr,
		.submission-table th,
		.submission-table td {
			display: block;
		}

		.submission-table thead {
			position: absolute;
			width: 1px;
			height: 1px;
			overflow: hidden;
			clip: rect(0 0 0 0);
		}

		.submission-table tr {
			border-bottom: 1px solid var(--color-border);
		}

		.submission-table tbody tr:not(.details-row) {
			padding: var(--space-md) 0;
		}

		.submission-table td {
			display: grid;
			grid-template-columns: minmax(6.5rem, 38%) minmax(0, 1fr);
			gap: var(--space-md);
			padding: var(--space-xs) 0;
			border-bottom: none;
		}

		.submission-table td::before {
			content: attr(data-label);
			color: var(--color-text-tertiary);
			font-size: 0.76rem;
			font-weight: 800;
			text-transform: uppercase;
		}

		.details-row {
			padding-bottom: var(--space-md);
		}

		.details-row td {
			display: block;
			padding: 0 0 var(--space-md);
		}

		.details-row td::before {
			content: none;
		}

		.field-details div {
			grid-template-columns: 1fr;
			gap: var(--space-xs);
		}
	}
</style>
