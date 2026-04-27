<script lang="ts">
	import { Button, IconButton, Link, Modal, Select, Toolbar } from '$lib/components';
	import { Check, Edit, Trash } from '$lib/icons';
	import { deleteDomain, updateDomain, verifyDomain } from '$lib/api-client';
	import type { DomainRecord, Site } from '$lib/types';

	interface Props {
		domain: DomainRecord;
		sites: Site[];
		onrefresh: () => Promise<void>;
	}

	let { domain, sites, onrefresh = async () => {} }: Props = $props();

	let loading = $state(false);
	let actionError = $state('');
	let deleteModal = $state<Modal>();
	let editModal = $state<Modal>();
	let selectedSiteId = $state('');
	let editError = $state('');
	let dnsOpen = $state(false);

	const deleteTitleId = $derived(`delete-domain-title-${domain.id}`);
	const deleteDescId = $derived(`delete-domain-description-${domain.id}`);
	const deleteButtonId = $derived(`delete-domain-confirm-${domain.id}`);
	const editTitleId = $derived(`edit-domain-title-${domain.id}`);
	const editSelectId = $derived(`edit-domain-select-${domain.id}`);

	const isCustom = $derived(domain.type === 'custom');
	const canOpen = $derived(domain.status === 'active');
	const errorText = $derived(actionError || domain.lastVerificationError || '');
	const setupGuide = $derived.by(() => {
		if (!isCustom || !domain.instructions) {
			return null;
		}

		return {
			heading: 'Настройка DNS',
			intro:
				'Добавьте две DNS-записи: TXT подтверждает владение, CNAME направляет поддомен на сайт в Хостамате.',
			steps: [
				'Откройте DNS-зону домена у вашего DNS-провайдера.',
				'Скопируйте TXT-запись из блока ниже. Она нужна только для проверки владения.',
				'Скопируйте CNAME-запись из блока ниже. Она должна указывать на DNS-адрес Хостамата.',
				'Если для этого поддомена уже есть другие записи, удалите их перед проверкой.',
				'Если провайдер просит TTL, поставьте 300 секунд.',
				'Когда DNS обновится, нажмите «Проверить».'
			],
			notes: [
				'В поле имени обычно можно вводить полное значение из таблицы. Если панель сама дописывает домен зоны, оставьте только левую часть.',
				'Не добавляйте https://, слеши и лишнюю точку в конце.'
			]
		};
	});
	const statusLabel = $derived.by(() => {
		if (domain.type === 'generated') {
			return 'резерв';
		}

		if (domain.status === 'active') {
			return 'активен';
		}

		if (domain.status === 'verified') {
			return 'проверен';
		}

		if (domain.status === 'verification_failed') {
			return 'ошибка DNS';
		}

		return 'ожидает';
	});
	const statusTone = $derived.by(() => {
		if (domain.type === 'generated' || domain.status === 'active') {
			return 'success';
		}

		if (domain.status === 'verification_failed') {
			return 'danger';
		}

		return 'pending';
	});
	const statusDetail = $derived.by(() => {
		if (errorText) {
			return errorText;
		}

		if (domain.type === 'generated') {
			return 'резервный адрес';
		}

		if (domain.status === 'active' && domain.activatedAt) {
			return `с ${new Date(domain.activatedAt).toLocaleDateString('ru-RU')}`;
		}

		if (domain.status === 'verified') {
			return 'ожидает активации';
		}

		if (domain.lastCheckedAt) {
			return 'проверка не пройдена';
		}

		return 'добавьте DNS-записи';
	});
	const detailHasError = $derived(Boolean(errorText));

	const handleVerify = async () => {
		loading = true;
		actionError = '';

		const { data, error } = await verifyDomain(domain.id);
		if (error) {
			actionError = error;
			loading = false;
			return;
		}

		if (data) {
			await onrefresh();
		}

		loading = false;
	};

	const handleDelete = async () => {
		loading = true;
		actionError = '';

		const { data, error } = await deleteDomain(domain.id);
		if (error) {
			actionError = error;
			loading = false;
			deleteModal?.close();
			return;
		}

		if (data) {
			deleteModal?.close();
			await onrefresh();
		}

		loading = false;
	};

	const handleEdit = async () => {
		loading = true;
		editError = '';

		const { data, error } = await updateDomain(domain.id, selectedSiteId);
		if (error) {
			editError = error;
			loading = false;
			return;
		}

		if (data) {
			editModal?.close();
			await onrefresh();
		}

		loading = false;
	};
</script>

<div class="domain-card">
	<div class="domain-header">
		<div class="domain-data">
			<div class="data-row">
				{#if canOpen}
					<Link href={`https://${domain.hostname}`} target="_blank" alwaysUnderline>
						{domain.hostname}
					</Link>
				{:else}
					<span class="hostname">{domain.hostname}</span>
				{/if}
				<span class="site-name">{domain.siteName}</span>
			</div>
			<div class="domain-meta">
				<span class="status-chip {statusTone}">{statusLabel}</span>
				<span class="meta-detail" class:error={detailHasError}>{statusDetail}</span>
			</div>
		</div>

		{#if isCustom && domain.instructions}
			<button
				type="button"
				class="dns-toggle"
				aria-expanded={dnsOpen}
				onclick={() => {
					dnsOpen = !dnsOpen;
				}}>DNS-записи</button
			>
		{/if}
	</div>

	{#if isCustom}
		<div class="domain-toolbar">
			<Toolbar>
				<IconButton
					grouped
					label="Проверить домен"
					title="Проверить домен"
					disabled={loading}
					onclick={handleVerify}
				>
					<Check size={25} />
				</IconButton>
				<IconButton
					grouped
					label="Изменить сайт"
					title="Изменить сайт"
					disabled={loading}
					onclick={() => {
						selectedSiteId = domain.siteId;
						editError = '';
						editModal?.open();
					}}
				>
					<Edit size={25} />
				</IconButton>
				<IconButton
					grouped
					label="Удалить домен"
					title="Удалить домен"
					tone="danger"
					disabled={loading}
					onclick={() => deleteModal?.open()}
				>
					<Trash size={25} />
				</IconButton>
			</Toolbar>
		</div>
	{/if}

	{#if isCustom && domain.instructions && dnsOpen}
		<div class="domain-details">
			<div class="dns-zone">
				{#if setupGuide}
					<div class="dns-guide">
						<h3>{setupGuide.heading}</h3>
						<p>{setupGuide.intro}</p>

						<ol class="guide-steps">
							{#each setupGuide.steps as step (`step-${step}`)}
								<li>{step}</li>
							{/each}
						</ol>

						<div class="guide-notes">
							{#each setupGuide.notes as note (`note-${note}`)}
								<p>{note}</p>
							{/each}
						</div>
					</div>
				{/if}

				<div class="zone-table">
					<div class="zone-header">
						<span>ТИП</span>
						<span>ИМЯ</span>
						<span>ЗНАЧЕНИЕ</span>
					</div>
					<div class="zone-body">
						<div class="zone-row zone-row-ownership">
							<span class="zone-type">{domain.instructions.ownership.type}</span>
							<span class="zone-name">{domain.instructions.ownership.name}</span>
							<span class="zone-value">{domain.instructions.ownership.value}</span>
						</div>
						{#each domain.instructions.routing as record (`${record.type}-${record.name}-${record.value}`)}
							<div class="zone-row">
								<span class="zone-type">{record.type}</span>
								<span class="zone-name">{record.name}</span>
								<span class="zone-value">{record.value}</span>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>
	{/if}

	<Modal
		bind:this={editModal}
		tone="neutral"
		title={`Изменить сайт для "${domain.hostname}"`}
		titleId={editTitleId}
		initialFocusSelector={'#' + editSelectId}
	>
		<div class="edit-modal">
			<Select
				id={editSelectId}
				label="Сайт"
				options={sites.map((site) => ({ value: site.id, label: `${site.name} (${site.domain})` }))}
				errorMessage={editError}
				bind:value={selectedSiteId}
			/>
			<div class="edit-actions">
				<Button emphasis="solid" disabled={loading} onclick={handleEdit}>Изменить</Button>
				<Button disabled={loading} onclick={() => editModal?.close()}>Отмена</Button>
			</div>
		</div>
	</Modal>

	<Modal
		bind:this={deleteModal}
		tone="danger"
		title="Подтверждение удаления"
		titleId={deleteTitleId}
		descriptionId={deleteDescId}
		initialFocusSelector={'#' + deleteButtonId}
	>
		<div class="delete-modal">
			<p id={deleteDescId} class="delete-modal-text">
				Вы уверены, что хотите удалить домен "{domain.hostname}"? Домен будет снят с маршрутизации,
				но сайт и его резервный адрес останутся без изменений.
			</p>
			<div class="delete-actions">
				<Button
					id={deleteButtonId}
					tone="danger"
					emphasis="solid"
					disabled={loading}
					onclick={handleDelete}>Удалить</Button
				>
				<Button tone="danger" onclick={() => deleteModal?.close()}>Отмена</Button>
			</div>
		</div>
	</Modal>
</div>

<style>
	.domain-card {
		display: grid;
		grid-template-columns: 1fr auto;
		align-items: center;
		column-gap: var(--space-lg);
		min-width: 0;
		padding: var(--space-lg) 0;
	}

	.domain-header {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		min-width: 0;
	}

	.domain-toolbar {
		align-self: center;
	}

	.domain-data {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		min-width: 0;
	}

	.domain-data :global(a) {
		font-family: var(--font-mono);
		font-size: 0.9rem;
		font-weight: 500;
	}

	.data-row {
		display: flex;
		align-items: baseline;
		gap: var(--space-sm);
		flex-wrap: wrap;
		min-width: 0;
	}

	.hostname {
		font-family: var(--font-mono);
		font-size: 0.9rem;
		font-weight: 500;
		color: var(--color-text-regular);
		overflow-wrap: anywhere;
	}

	.site-name {
		font-size: 0.82rem;
		color: var(--color-text-secondary);
		overflow-wrap: anywhere;
	}

	.domain-meta {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		flex-wrap: wrap;
	}

	.status-chip {
		display: inline-flex;
		align-items: center;
		padding: 0.12rem 0.45rem;
		border: 1px solid var(--chip-border);
		background: var(--chip-bg);
		color: var(--chip-text);
		border-radius: var(--radius-sm);
		font-size: 0.68rem;
		font-weight: 700;
		letter-spacing: 0.05em;
		text-transform: uppercase;
		white-space: nowrap;
	}

	.status-chip.success {
		--chip-border: var(--color-success-border);
		--chip-bg: var(--color-success-surface);
		--chip-text: var(--color-success-text);
	}

	.status-chip.pending {
		--chip-border: var(--color-accent-border);
		--chip-bg: var(--color-accent-surface);
		--chip-text: var(--color-accent-text);
	}

	.status-chip.danger {
		--chip-border: var(--color-danger-border);
		--chip-bg: var(--color-danger-surface);
		--chip-text: var(--color-danger-text);
	}

	.meta-detail {
		font-size: 0.82rem;
		color: var(--color-text-secondary);
		overflow-wrap: anywhere;
	}

	.meta-detail.error {
		color: var(--color-danger-text);
	}

	.dns-toggle {
		all: unset;
		width: fit-content;
		font-weight: 650;
		color: var(--color-accent-text);
		text-decoration: underline;
		text-underline-offset: 3px;
		cursor: pointer;
	}

	.dns-toggle:hover {
		color: var(--color-accent-strong);
	}

	.domain-details {
		grid-column: 1;
		min-width: 0;
	}

	.dns-zone {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		padding-top: var(--space-xs);
	}

	.dns-guide {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		padding: var(--space-md);
		border-radius: var(--radius-sm);
		border: 1px solid var(--color-border);
		background: var(--color-background-object);
	}

	.dns-guide h3 {
		font-size: 0.9rem;
		font-weight: 700;
		letter-spacing: -0.01em;
	}

	.dns-guide p {
		font-size: 0.85rem;
		line-height: 1.45;
	}

	.guide-steps {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		padding-left: 1.05rem;
		color: var(--color-text-secondary);
		font-size: 0.85rem;
		line-height: 1.45;
	}

	.guide-notes {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		padding-top: 0.35rem;
		border-top: 1px solid var(--color-border);
	}

	.guide-notes p {
		font-size: 0.8rem;
		color: var(--color-text-tertiary);
		line-height: 1.4;
	}

	.zone-table {
		display: flex;
		flex-direction: column;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		overflow: hidden;
		background: var(--color-background-object);
	}

	.zone-header {
		display: grid;
		grid-template-columns: 4rem 1fr 2fr;
		gap: var(--space-sm);
		padding: 0.45rem 0.75rem;
		background: var(--color-neutral-surface);
		border-bottom: 1px solid var(--color-border);
	}

	.zone-header span {
		font-size: 0.65rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-text-tertiary);
	}

	.zone-body {
		display: flex;
		flex-direction: column;
	}

	.zone-row {
		display: grid;
		grid-template-columns: 4rem 1fr 2fr;
		gap: var(--space-sm);
		padding: 0.5rem 0.75rem;
		border-bottom: 1px solid var(--color-border);
		font-family: var(--font-mono);
		font-size: 0.78rem;
		line-height: 1.35;
		align-items: center;
	}

	.zone-row:last-child {
		border-bottom: none;
	}

	.zone-row-ownership {
		background: color-mix(in oklch, var(--color-accent-surface) 55%, transparent);
	}

	.zone-type {
		font-weight: 600;
		color: var(--color-text-regular);
	}

	.zone-name {
		color: var(--color-text-secondary);
		overflow-wrap: anywhere;
	}

	.zone-value {
		color: var(--color-text-regular);
		overflow-wrap: anywhere;
	}

	.edit-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.edit-actions {
		display: flex;
		gap: var(--space-sm);
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.delete-modal-text {
		color: var(--color-danger-text);
	}

	.delete-actions {
		display: flex;
		gap: var(--space-sm);
	}

	@media (max-width: 640px) {
		.domain-card {
			grid-template-columns: 1fr;
			gap: var(--space-md);
		}

		.domain-toolbar {
			width: 100%;
		}

		.zone-header,
		.zone-row {
			grid-template-columns: 3rem 1fr;
			gap: var(--space-xs);
		}

		.zone-value {
			grid-column: 1 / -1;
			padding-top: 0.05rem;
		}
	}
</style>
