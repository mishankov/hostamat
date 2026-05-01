<script lang="ts">
	import { resolve } from '$app/paths';
	import { invalidateAll } from '$app/navigation';
	import type { PageProps } from './$types';

	import {
		Alert,
		Button,
		CopyButton,
		IconButton,
		Input,
		Link,
		Modal,
		Toolbar
	} from '$lib/components';
	import { Edit, Plus, Spinner, Trash } from '$lib/icons';
	import { createApiKey, revokeApiKey, updateApiKeyName } from '$lib/api-client';
	import type { ApiKeyRecord, ApiKeyScope, ApiKeyState } from '$lib/types';

	type ScopeOption = {
		value: ApiKeyScope;
		label: string;
		description: string;
	};

	type ScopePreset = {
		label: string;
		description: string;
		scopes: ApiKeyScope[];
	};

	type ExpirationMode = 'never' | '30d' | '90d' | '1y' | 'custom';

	const scopeOptions: ScopeOption[] = [
		{ value: 'sites:read', label: 'sites:read', description: 'Просмотр сайтов' },
		{
			value: 'sites:write',
			label: 'sites:write',
			description: 'Создание, переименование, удаление'
		},
		{ value: 'sites:upload', label: 'sites:upload', description: 'Публикация ZIP-архива' },
		{ value: 'domains:read', label: 'domains:read', description: 'Просмотр доменов' },
		{ value: 'domains:write', label: 'domains:write', description: 'Управление доменами' }
	];

	const scopePresets: ScopePreset[] = [
		{
			label: 'Только чтение',
			description: 'Инвентаризация сайтов и доменов',
			scopes: ['sites:read', 'domains:read']
		},
		{
			label: 'Публикация',
			description: 'CI/CD без удаления сайтов',
			scopes: ['sites:read', 'sites:upload']
		},
		{
			label: 'Сайты',
			description: 'Полное управление сайтами',
			scopes: ['sites:read', 'sites:write', 'sites:upload']
		},
		{
			label: 'Домены',
			description: 'Настройка и проверка DNS',
			scopes: ['domains:read', 'domains:write']
		},
		{
			label: 'Полный доступ',
			description: 'Все возможности Public API',
			scopes: ['sites:read', 'sites:write', 'sites:upload', 'domains:read', 'domains:write']
		}
	];

	const expirationOptions: { value: ExpirationMode; label: string; description: string }[] = [
		{ value: 'never', label: 'Никогда', description: 'never expires' },
		{ value: '30d', label: '30 дней', description: 'Короткий ключ для CI' },
		{ value: '90d', label: '90 дней', description: 'Обычная ротация' },
		{ value: '1y', label: '1 год', description: 'Долгий проект' },
		{ value: 'custom', label: 'Дата', description: 'Выбрать вручную' }
	];

	const apiDocsHref = resolve('/api-docs');
	const maxNameLength = 80;
	const defaultMaxActiveAmount = 20;
	const todayInputValue = new Date().toISOString().slice(0, 10);

	let { data }: PageProps = $props();

	let apiKeys = $state.raw<ApiKeyRecord[]>([]);
	let maxActiveAmount = $state(defaultMaxActiveAmount);
	let loading = $state(true);
	let pageError = $state('');
	let actionLoading = $state(false);

	let createModal = $state<Modal>();
	let secretModal = $state<Modal>();
	let editModal = $state<Modal>();
	let revokeModal = $state<Modal>();

	let createName = $state('');
	let createNameError = $state('');
	let selectedScopes = $state<ApiKeyScope[]>([]);
	let scopesError = $state('');
	let expirationMode = $state<ExpirationMode>('90d');
	let customExpirationDate = $state('');
	let expirationError = $state('');
	let createError = $state('');

	let createdSecret = $state('');
	let createdSecretName = $state('');
	let secretAcknowledged = $state(false);

	let editingKey = $state<ApiKeyRecord | null>(null);
	let editName = $state('');
	let editError = $state('');

	let revokingKey = $state<ApiKeyRecord | null>(null);
	let revokeError = $state('');

	const visibleKeys = $derived(apiKeys.filter((key) => resolveKeyState(key) !== 'revoked'));
	const activeKeys = $derived(visibleKeys.filter((key) => resolveKeyState(key) === 'active'));
	const expiredKeys = $derived(visibleKeys.filter((key) => resolveKeyState(key) === 'expired'));
	const canCreateKey = $derived(activeKeys.length < maxActiveAmount);

	const updateState = async () => {
		const apiKeysResult = await data.apiKeysResp;

		if (apiKeysResult.data) {
			apiKeys = apiKeysResult.data.apiKeys ?? [];
			maxActiveAmount = apiKeysResult.data.maxActiveAmount ?? defaultMaxActiveAmount;
			pageError = '';
		} else {
			apiKeys = [];
			maxActiveAmount = defaultMaxActiveAmount;
			pageError = apiKeysResult.error;
		}

		loading = false;
	};

	const refresh = async (showLoading = true) => {
		if (showLoading) {
			loading = true;
		}

		await invalidateAll();
		await updateState();
	};

	const resetCreateForm = () => {
		createName = '';
		createNameError = '';
		selectedScopes = ['sites:read'];
		scopesError = '';
		expirationMode = '90d';
		customExpirationDate = '';
		expirationError = '';
		createError = '';
	};

	const resolveKeyState = (key: ApiKeyRecord): ApiKeyState => {
		if (key.state === 'revoked' || key.revokedAt) {
			return 'revoked';
		}

		if (key.state === 'expired') {
			return 'expired';
		}

		if (key.expiresAt && new Date(key.expiresAt).getTime() <= Date.now()) {
			return 'expired';
		}

		return 'active';
	};

	const keyIdentifier = (key: ApiKeyRecord) => {
		return key.prefix || key.keyPrefix || key.identifier || key.id.slice(0, 8);
	};

	const formatDate = (value: string | null | undefined, emptyText = 'нет') => {
		if (!value) {
			return emptyText;
		}

		return new Intl.DateTimeFormat('ru-RU', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		}).format(new Date(value));
	};

	const statusLabel = (state: ApiKeyState) => {
		if (state === 'expired') {
			return 'истек';
		}

		if (state === 'revoked') {
			return 'отозван';
		}

		return 'активен';
	};

	const applyScopePreset = (preset: ScopePreset) => {
		selectedScopes = [...preset.scopes];
		scopesError = '';
	};

	const toggleScope = (scope: ApiKeyScope, checked: boolean) => {
		if (checked) {
			selectedScopes = [...new Set([...selectedScopes, scope])];
		} else {
			selectedScopes = selectedScopes.filter((selected) => selected !== scope);
		}

		if (selectedScopes.length > 0) {
			scopesError = '';
		}
	};

	const toEndOfDayIso = (date: Date) => {
		const endOfDay = new Date(date);
		endOfDay.setHours(23, 59, 59, 999);
		return endOfDay.toISOString();
	};

	const resolveExpiration = (): string | null => {
		if (expirationMode === 'never') {
			return null;
		}

		if (expirationMode === 'custom') {
			if (!customExpirationDate) {
				expirationError = 'Выберите дату окончания';
				return '';
			}

			const date = new Date(`${customExpirationDate}T00:00:00`);
			if (Number.isNaN(date.getTime()) || customExpirationDate < todayInputValue) {
				expirationError = 'Дата окончания не может быть в прошлом';
				return '';
			}

			return toEndOfDayIso(date);
		}

		const days = expirationMode === '30d' ? 30 : expirationMode === '90d' ? 90 : 365;
		const date = new Date();
		date.setDate(date.getDate() + days);
		return toEndOfDayIso(date);
	};

	const validateCreateForm = () => {
		const trimmedName = createName.trim();
		createNameError = '';
		scopesError = '';
		expirationError = '';
		createError = '';

		if (!trimmedName) {
			createNameError = 'Название обязательно';
		} else if (trimmedName.length > maxNameLength) {
			createNameError = `Максимум ${maxNameLength} символов`;
		}

		if (selectedScopes.length === 0) {
			scopesError = 'Выберите хотя бы один scope';
		}

		const expiresAt = resolveExpiration();

		if (createNameError || scopesError || expirationError || expiresAt === '') {
			return null;
		}

		return {
			name: trimmedName,
			scopes: selectedScopes,
			expiresAt
		};
	};

	const submitCreate = async () => {
		const payload = validateCreateForm();
		if (!payload) {
			return;
		}

		actionLoading = true;
		const { data: created, error } = await createApiKey(payload);

		if (error) {
			createError = error;
			actionLoading = false;
			return;
		}

		if (created) {
			createdSecret = created.secret;
			createdSecretName = created.apiKey.name;
			secretAcknowledged = false;
			createModal?.close();
			await refresh(false);
			queueMicrotask(() => secretModal?.open());
		}

		actionLoading = false;
	};

	const openEdit = (key: ApiKeyRecord) => {
		editingKey = key;
		editName = key.name;
		editError = '';
		editModal?.open();
	};

	const submitEdit = async () => {
		if (!editingKey) {
			return;
		}

		const trimmedName = editName.trim();
		if (!trimmedName) {
			editError = 'Название обязательно';
			return;
		}

		if (trimmedName.length > maxNameLength) {
			editError = `Максимум ${maxNameLength} символов`;
			return;
		}

		actionLoading = true;
		editError = '';

		const { data: updated, error } = await updateApiKeyName(editingKey.id, trimmedName);
		if (error) {
			editError = error;
			actionLoading = false;
			return;
		}

		if (updated) {
			editModal?.close();
			await refresh();
		}

		actionLoading = false;
	};

	const openRevoke = (key: ApiKeyRecord) => {
		revokingKey = key;
		revokeError = '';
		revokeModal?.open();
	};

	const submitRevoke = async () => {
		if (!revokingKey) {
			return;
		}

		actionLoading = true;
		revokeError = '';

		const { data: revoked, error } = await revokeApiKey(revokingKey.id);
		if (error) {
			revokeError = error;
			actionLoading = false;
			return;
		}

		if (revoked) {
			revokeModal?.close();
			await refresh();
		}

		actionLoading = false;
	};

	updateState();
</script>

<svelte:head>
	<title>API-ключи | Хостамат</title>
	<meta
		name="description"
		content="Создавайте и отзывайте API-ключи для автоматизации сайтов, загрузок и доменов."
	/>
</svelte:head>

<div class="dashboard-page">
	<div class="headline">
		<div>
			<h1>API-ключи</h1>
			<div class="headline-meta">
				<p>Управляйте доступом к Public API для скриптов, CI и интеграций.</p>
				{#if !loading}
					<span class="key-counter">активных: {activeKeys.length}/{maxActiveAmount}</span>
					<span class="key-counter">истекших: {expiredKeys.length}</span>
				{/if}
			</div>
		</div>
		<div class="headline-actions">
			<Button href={apiDocsHref}>Документация API</Button>
			<Button
				emphasis="solid"
				disabled={loading || !canCreateKey}
				onclick={() => {
					resetCreateForm();
					createModal?.open();
				}}
			>
				<Plus size={16} />
				<span>Создать ключ</span>
			</Button>
		</div>
	</div>

	{#if pageError}
		<Alert variant="error" message={pageError} />
	{/if}

	{#if !loading && !canCreateKey}
		<Alert variant="info" title="Достигнут лимит активных ключей">
			Отзовите один из активных ключей или дождитесь окончания срока действия. Истекшие и отозванные
			ключи не учитываются в лимите.
		</Alert>
	{/if}

	<div class="key-list">
		{#if loading}
			<div class="loading-state" aria-live="polite">
				<Spinner size={50} />
			</div>
		{:else if visibleKeys.length === 0}
			<div class="empty-state">
				<h2>Ключей пока нет</h2>
				<p>Создайте ключ, выберите scopes и используйте его как Bearer-токен в Public API.</p>
			</div>
		{:else}
			<div class="key-table" role="table" aria-label="API-ключи">
				<div class="key-table-head" role="row">
					<span role="columnheader">Ключ</span>
					<span role="columnheader">Scopes</span>
					<span role="columnheader">Даты</span>
					<span role="columnheader">Состояние</span>
					<span role="columnheader" class="actions-heading">Действия</span>
				</div>

				{#each visibleKeys as key (key.id)}
					{@const state = resolveKeyState(key)}
					<div class="key-row" role="row">
						<div class="key-main" role="cell">
							<div class="key-name-line">
								<strong>{key.name}</strong>
								<code>{keyIdentifier(key)}</code>
							</div>
							<span class="key-id">id: {key.id}</span>
						</div>

						<div class="scope-list" role="cell">
							{#each key.scopes as scope (scope)}
								<span class="scope-chip">{scope}</span>
							{/each}
						</div>

						<div class="date-stack" role="cell">
							<span><span class="date-label">создан</span>{formatDate(key.createdAt)}</span>
							<span
								><span class="date-label">истекает</span>{formatDate(
									key.expiresAt,
									'никогда'
								)}</span
							>
							<span
								><span class="date-label">использован</span>{formatDate(
									key.lastUsedAt,
									'нет'
								)}</span
							>
						</div>

						<div class="state-cell" role="cell">
							<span class={['status-chip', state]}>{statusLabel(state)}</span>
						</div>

						<div class="key-actions" role="cell">
							<Toolbar>
								<IconButton
									grouped
									label="Переименовать ключ"
									title="Переименовать"
									disabled={actionLoading}
									onclick={() => openEdit(key)}
								>
									<Edit size={20} />
								</IconButton>
								<IconButton
									grouped
									tone="danger"
									label="Отозвать ключ"
									title="Отозвать"
									disabled={actionLoading}
									onclick={() => openRevoke(key)}
								>
									<Trash size={20} />
								</IconButton>
							</Toolbar>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<Modal
	bind:this={createModal}
	title="Новый API-ключ"
	titleId="create-api-key-title"
	panelWidth="min(50rem, calc(100vw - 2rem))"
	initialFocusSelector="#api-key-name"
	closeLabel="Закрыть форму создания API-ключа"
>
	<div class="modal-form">
		<Input
			id="api-key-name"
			label="Название"
			placeholder="Мой любимый API ключ"
			maxlength={maxNameLength}
			errorMessage={createNameError}
			bind:value={createName}
		/>

		<div class="field-group">
			<div class="field-heading">
				<h3>Scopes</h3>
				<span>Нужен минимум один scope</span>
			</div>

			<div class="preset-grid">
				{#each scopePresets as preset (preset.label)}
					<button
						type="button"
						class="preset-button"
						onclick={() => applyScopePreset(preset)}
						aria-pressed={preset.scopes.every((scope) => selectedScopes.includes(scope))}
					>
						<strong>{preset.label}</strong>
						<span>{preset.description}</span>
					</button>
				{/each}
			</div>

			<div class="scope-checkboxes" class:invalid={Boolean(scopesError)}>
				{#each scopeOptions as scope (scope.value)}
					<label class="scope-option">
						<input
							type="checkbox"
							checked={selectedScopes.includes(scope.value)}
							onchange={(event) =>
								toggleScope(scope.value, (event.currentTarget as HTMLInputElement).checked)}
						/>
						<span>
							<strong>{scope.label}</strong>
							<small>{scope.description}</small>
						</span>
					</label>
				{/each}
			</div>
			{#if scopesError}
				<p class="field-error">{scopesError}</p>
			{/if}
		</div>

		<div class="field-group">
			<div class="field-heading">
				<h3>Срок действия</h3>
				<span>Выберите дату ротации или вариант без срока</span>
			</div>

			<div class="expiration-grid">
				{#each expirationOptions as option (option.value)}
					<button
						type="button"
						class={['expiration-button', expirationMode === option.value && 'selected']}
						onclick={() => {
							expirationMode = option.value;
							expirationError = '';
						}}
					>
						<strong>{option.label}</strong>
						<span>{option.description}</span>
					</button>
				{/each}
			</div>

			{#if expirationMode === 'custom'}
				<Input
					id="api-key-custom-expiration"
					type="date"
					label="Дата окончания"
					min={todayInputValue}
					errorMessage={expirationError}
					bind:value={customExpirationDate}
				/>
			{:else if expirationError}
				<p class="field-error">{expirationError}</p>
			{/if}
		</div>

		{#if createError}
			<Alert variant="error" message={createError} />
		{/if}

		<div class="modal-actions">
			<Button emphasis="solid" disabled={actionLoading || !canCreateKey} onclick={submitCreate}>
				{#if actionLoading}
					<Spinner size={16} />
					<span>Создаем</span>
				{:else}
					<span>Создать ключ</span>
				{/if}
			</Button>
			<Button disabled={actionLoading} onclick={() => createModal?.close()}>Отмена</Button>
		</div>
	</div>
</Modal>

<Modal
	bind:this={secretModal}
	title="Сохраните секретный ключ"
	titleId="api-key-secret-title"
	panelWidth="min(42rem, calc(100vw - 2rem))"
	closeOnBackdrop={false}
	showCloseButton={false}
>
	<div class="secret-flow">
		<p>
			Ключ для «{createdSecretName}» больше нельзя будет посмотреть в Хостамате. Скопируйте его
			сейчас и сохраните в менеджере секретов.
		</p>

		<div class="secret-box">
			<code>{createdSecret}</code>
			<CopyButton
				textToCopy={createdSecret}
				label="Скопировать секретный ключ"
				title="Скопировать"
			/>
		</div>

		<label class="acknowledgement">
			<input type="checkbox" bind:checked={secretAcknowledged} />
			<span>Я сохранил ключ и понимаю, что повторно увидеть секрет нельзя.</span>
		</label>

		<div class="modal-actions">
			<Button
				emphasis="solid"
				disabled={!secretAcknowledged}
				onclick={() => {
					createdSecret = '';
					secretModal?.close();
				}}>Готово</Button
			>
		</div>
	</div>
</Modal>

<Modal
	bind:this={editModal}
	title="Переименовать API-ключ"
	titleId="edit-api-key-title"
	initialFocusSelector="#edit-api-key-name"
	closeLabel="Закрыть форму переименования"
>
	<div class="modal-form narrow">
		<Input
			id="edit-api-key-name"
			label="Название"
			maxlength={maxNameLength}
			errorMessage={editError}
			bind:value={editName}
		/>
		<p class="modal-note">После создания можно менять только название ключа.</p>
		<div class="modal-actions">
			<Button emphasis="solid" disabled={actionLoading} onclick={submitEdit}>Сохранить</Button>
			<Button disabled={actionLoading} onclick={() => editModal?.close()}>Отмена</Button>
		</div>
	</div>
</Modal>

<Modal
	bind:this={revokeModal}
	tone="danger"
	title="Отозвать API-ключ?"
	titleId="revoke-api-key-title"
	closeLabel="Закрыть подтверждение отзыва"
>
	<div class="modal-form narrow">
		<p>
			Ключ «{revokingKey?.name}» перестанет работать сразу после отзыва. Отозванные ключи скрываются
			из списка.
		</p>
		{#if revokeError}
			<Alert variant="error" message={revokeError} />
		{/if}
		<div class="modal-actions">
			<Button tone="danger" emphasis="solid" disabled={actionLoading} onclick={submitRevoke}>
				Отозвать
			</Button>
			<Button disabled={actionLoading} onclick={() => revokeModal?.close()}>Отмена</Button>
		</div>
	</div>
</Modal>

<style>
	.dashboard-page {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 0;
	}

	.headline {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-xl);
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
		gap: var(--space-md);
		flex-wrap: wrap;
		padding-top: var(--space-sm);
	}

	.key-counter {
		font-size: 0.82rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		white-space: nowrap;
	}

	.headline-actions {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		flex-wrap: wrap;
		justify-content: flex-end;
	}

	.api-note {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		flex-wrap: wrap;
		padding: 0.65rem 0.8rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: color-mix(in hsl, var(--color-background-object) 84%, transparent);
		color: var(--color-text-secondary);
		font-size: 0.86rem;
	}

	code {
		font-family: var(--font-mono);
		font-size: 0.78rem;
		color: var(--color-text-regular);
	}

	.key-list {
		display: flex;
		flex-direction: column;
		min-width: 0;
		border-top: 1px solid var(--color-border);
	}

	.loading-state,
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-sm);
		min-height: 11rem;
		text-align: center;
		border-bottom: 1px solid var(--color-border);
	}

	.empty-state p {
		max-width: 46ch;
	}

	.key-table {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.key-table-head,
	.key-row {
		display: grid;
		grid-template-columns: minmax(13rem, 1.2fr) minmax(13rem, 1.1fr) minmax(14rem, 1fr) 6.5rem auto;
		gap: var(--space-lg);
		align-items: center;
	}

	.key-table-head {
		padding: 0.75rem 0;
		border-bottom: 1px solid var(--color-border);
		color: var(--color-text-tertiary);
		font-size: 0.72rem;
		font-weight: 700;
		letter-spacing: 0.04em;
		text-transform: uppercase;
	}

	.key-row {
		padding: var(--space-lg) 0;
		border-bottom: 1px solid var(--color-border);
	}

	.key-main {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		min-width: 0;
	}

	.key-name-line {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		min-width: 0;
		flex-wrap: wrap;
	}

	.key-name-line strong {
		max-width: 24ch;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.key-name-line code {
		padding: 0.12rem 0.35rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		background: var(--color-neutral-surface);
	}

	.key-id {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		color: var(--color-text-tertiary);
		font-family: var(--font-mono);
		font-size: 0.72rem;
	}

	.scope-list {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
		flex-wrap: wrap;
		min-width: 0;
	}

	.scope-chip,
	.status-chip {
		display: inline-flex;
		align-items: center;
		width: fit-content;
		border: 1px solid var(--color-border);
		border-radius: 999px;
		padding: 0.18rem 0.48rem;
		font-size: 0.72rem;
		font-weight: 700;
		line-height: 1.3;
		white-space: nowrap;
	}

	.scope-chip {
		background: var(--color-background-object);
		color: var(--color-text-secondary);
		font-family: var(--font-mono);
		font-weight: 500;
	}

	.status-chip.active {
		border-color: var(--color-success-border);
		background: var(--color-success-surface);
		color: var(--color-success-text);
	}

	.status-chip.expired {
		border-color: var(--color-border-strong);
		background: var(--color-neutral-surface);
		color: var(--color-text-secondary);
	}

	.status-chip.revoked {
		border-color: var(--color-danger-border);
		background: var(--color-danger-surface);
		color: var(--color-danger-text);
	}

	.date-stack {
		display: grid;
		gap: 0.1rem;
		min-width: 0;
		font-size: 0.8rem;
		color: var(--color-text-secondary);
	}

	.date-stack span {
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.date-label {
		display: inline-block;
		width: 5.8rem;
		color: var(--color-text-tertiary);
	}

	.state-cell {
		display: flex;
		align-items: center;
	}

	.key-actions {
		display: flex;
		justify-content: flex-end;
	}

	.actions-heading {
		text-align: right;
	}

	.modal-form,
	.secret-flow {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		width: 100%;
		min-width: 0;
		max-width: 100%;
	}

	.modal-form.narrow {
		max-width: 28rem;
	}

	.field-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}

	.field-heading {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: var(--space-md);
		min-width: 0;
		flex-wrap: wrap;
	}

	.field-heading h3 {
		font-size: 1rem;
	}

	.field-heading span,
	.modal-note {
		color: var(--color-text-secondary);
		font-size: 0.82rem;
		min-width: 0;
	}

	.preset-grid,
	.expiration-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(min(10rem, 100%), 1fr));
		gap: var(--space-sm);
	}

	.preset-button,
	.expiration-button {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.2rem;
		min-width: 0;
		min-height: 4.35rem;
		padding: 0.65rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-background-object);
		color: var(--color-text-regular);
		text-align: left;
		cursor: pointer;
		transition:
			background-color 140ms var(--ease-out-quart),
			border-color 140ms var(--ease-out-quart),
			box-shadow 140ms var(--ease-out-quart),
			transform 120ms var(--ease-out-quart);
	}

	.preset-button:hover,
	.expiration-button:hover,
	.expiration-button.selected {
		border-color: var(--color-accent-border);
		background: var(--color-accent-surface);
	}

	.preset-button:hover,
	.expiration-button:hover {
		transform: translateY(-1px);
	}

	.preset-button strong,
	.expiration-button strong {
		font-size: 0.82rem;
		line-height: 1.2;
		overflow-wrap: anywhere;
	}

	.preset-button span,
	.expiration-button span {
		color: var(--color-text-secondary);
		font-size: 0.72rem;
		line-height: 1.35;
		overflow-wrap: anywhere;
	}

	.scope-checkboxes {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--space-sm);
		padding: var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: color-mix(in hsl, var(--color-background-object) 82%, transparent);
	}

	.scope-checkboxes.invalid {
		border-color: var(--color-danger-border);
		background: var(--color-danger-surface);
	}

	.scope-option,
	.acknowledgement {
		display: flex;
		align-items: flex-start;
		gap: var(--space-sm);
		min-width: 0;
		cursor: pointer;
	}

	.scope-option input,
	.acknowledgement input {
		margin-top: 0.18rem;
		accent-color: var(--color-accent-strong);
	}

	.scope-option span {
		display: flex;
		flex-direction: column;
		gap: 0.05rem;
		min-width: 0;
	}

	.scope-option strong {
		font-family: var(--font-mono);
		font-size: 0.78rem;
		font-weight: 500;
		overflow-wrap: anywhere;
	}

	.scope-option small {
		color: var(--color-text-secondary);
		font-size: 0.74rem;
		overflow-wrap: anywhere;
	}

	.field-error {
		color: var(--color-danger-text);
		font-size: 0.82rem;
		font-weight: 600;
	}

	.modal-actions {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		flex-wrap: wrap;
	}

	.secret-box {
		display: grid;
		grid-template-columns: minmax(0, 1fr) auto;
		gap: var(--space-sm);
		align-items: center;
		padding: var(--space-md);
		border: 1px solid var(--color-accent-border);
		border-radius: var(--radius-md);
		background: var(--color-accent-surface);
	}

	.secret-box code {
		overflow-wrap: anywhere;
	}

	@media (max-width: 1180px) {
		.key-table-head {
			display: none;
		}

		.key-row {
			grid-template-columns: minmax(0, 1fr) auto;
			align-items: start;
		}

		.scope-list,
		.date-stack,
		.state-cell {
			grid-column: 1 / -1;
		}

		.key-actions {
			grid-column: 2;
			grid-row: 1;
		}
	}

	@media (max-width: 760px) {
		.headline {
			flex-direction: column;
		}

		.headline-actions {
			justify-content: flex-start;
		}

		.preset-grid,
		.expiration-grid,
		.scope-checkboxes {
			grid-template-columns: minmax(0, 1fr);
		}

		.key-row {
			grid-template-columns: minmax(0, 1fr);
		}

		.key-actions {
			grid-column: 1;
			grid-row: auto;
			justify-content: flex-start;
		}

		.secret-box {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
