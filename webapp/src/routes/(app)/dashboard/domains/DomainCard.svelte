<script lang="ts">
	import { Button, Details, Link } from '$lib/components';
	import { deleteDomain, verifyDomain } from '$lib/api-client';
	import type { DomainRecord } from '$lib/types';

	interface Props {
		domain: DomainRecord;
		onrefresh: () => Promise<void>;
	}

	let { domain, onrefresh = async () => {} }: Props = $props();

	let loading = $state(false);
	let actionError = $state('');
	let confirmingDelete = $state(false);

	const isCustom = $derived(domain.type === 'custom');
	const canOpen = $derived(domain.status === 'active');
	const logMessage = $derived(actionError || domain.lastVerificationError || '');
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
	const timestampLabel = $derived.by(() => {
		if (domain.activatedAt) {
			return `с ${new Date(domain.activatedAt).toLocaleDateString('ru-RU')}`;
		}

		if (domain.lastCheckedAt) {
			return `проверен ${new Date(domain.lastCheckedAt).toLocaleDateString('ru-RU')}`;
		}

		return 'ещё не проверен';
	});

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
			confirmingDelete = false;
			return;
		}

		if (data) {
			confirmingDelete = false;
			await onrefresh();
		}

		loading = false;
	};
</script>

<div class="domain-card">
	<div class="domain-main">
		<div class="domain-data">
			<div class="data-row">
				<span class="data-label">{domain.type === 'generated' ? 'резерв' : 'домен'}:</span>
				{#if canOpen}
					<Link href={`https://${domain.hostname}`} target="_blank" alwaysUnderline>
						{domain.hostname}
					</Link>
				{:else}
					<span class="hostname">{domain.hostname}</span>
				{/if}
			</div>
			<div class="domain-meta">
				<span class="status-chip {statusTone}">
					<span class="status-pip" aria-hidden="true"></span>
					{statusLabel}
				</span>
				<span class="meta-text">{timestampLabel}</span>
				<span class="meta-text">{domain.siteName}</span>
			</div>
		</div>

		{#if isCustom}
			<div class="domain-toolbar">
				<Button disabled={loading} emphasis="solid" onclick={handleVerify}>
					{domain.status === 'active' ? 'Проверить' : 'Проверить'}
				</Button>
				<Button
					tone="danger"
					disabled={loading}
					onclick={() => {
						confirmingDelete = !confirmingDelete;
						actionError = '';
					}}
				>
					{confirmingDelete ? 'Отменить' : 'Удалить'}
				</Button>
			</div>
		{/if}
	</div>

	{#if logMessage}
		<div class="domain-log">
			<span class="data-label">ошибка:</span>
			<span class="log-text">{logMessage}</span>
		</div>
	{/if}

	{#if isCustom && domain.instructions}
		<Details summary="Как настроить DNS">
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
		</Details>
	{/if}

	{#if isCustom && confirmingDelete}
		<div class="delete-guard">
			<div class="delete-guard-text">
				<span class="data-label">отключение домена</span>
				<p>Удаление снимет домен с маршрутизации, но не затронет сайт и его резервный адрес.</p>
			</div>
			<Button tone="danger" emphasis="solid" disabled={loading} onclick={handleDelete}>
				Подтвердить удаление
			</Button>
		</div>
	{/if}
</div>

<style>
	.domain-card {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		min-width: 0;
		padding: var(--space-lg) 0;
	}

	.domain-main {
		display: grid;
		grid-template-columns: 1fr auto;
		align-items: center;
		gap: var(--space-lg);
		min-width: 0;
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

	.hostname {
		font-family: var(--font-mono);
		font-size: 0.9rem;
		font-weight: 500;
		color: var(--color-text-regular);
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
		gap: 0.3rem;
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

	.status-pip {
		display: block;
		width: 5px;
		height: 5px;
		border-radius: 1px;
		background: currentColor;
	}

	.meta-text {
		font-size: 0.82rem;
		color: var(--color-text-secondary);
		white-space: nowrap;
	}

	.domain-toolbar {
		display: flex;
		gap: var(--space-sm);
		flex-wrap: wrap;
		flex-shrink: 0;
	}

	.domain-log {
		display: flex;
		align-items: baseline;
		gap: var(--space-sm);
		flex-wrap: wrap;
		padding: 0.45rem 0.6rem;
		border-radius: var(--radius-sm);
		border: 1px solid color-mix(in oklch, var(--color-danger-border) 40%, transparent);
		background: color-mix(in oklch, var(--color-danger-surface) 94%, transparent);
	}

	.log-text {
		font-size: 0.84rem;
		color: var(--color-danger-text);
		overflow-wrap: anywhere;
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

	.delete-guard {
		display: flex;
		flex-wrap: wrap;
		align-items: end;
		justify-content: space-between;
		gap: var(--space-md);
		padding: var(--space-md);
		border-radius: var(--radius-sm);
		border: 1px solid color-mix(in oklch, var(--color-danger-border) 40%, transparent);
		background: color-mix(in oklch, var(--color-danger-surface) 92%, transparent);
	}

	.delete-guard-text {
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
		max-width: 36ch;
	}

	.delete-guard-text p {
		color: var(--color-danger-text);
		font-size: 0.84rem;
		line-height: 1.4;
	}

	@media (max-width: 640px) {
		.domain-main {
			grid-template-columns: 1fr;
			gap: var(--space-md);
		}

		.domain-toolbar {
			justify-content: flex-start;
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
