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

		if (domain.isApex) {
			return {
				heading: 'Как настроить корневой домен',
				intro:
					'Под корневым доменом понимается адрес вроде example.com без поддомена. Для него нужен TXT на владение и A/AAAA на входные адреса Hostamat.',
				steps: [
					'Откройте DNS-зону домена у регистратора или DNS-провайдера.',
					'Создайте TXT-запись для подтверждения владения. Если провайдер не принимает полное имя корня, используйте @ вместо полного домена.',
					'Создайте A и, если указано ниже, AAAA-записи ровно со значениями из блока «Маршрутизация».',
					'Если у корневого домена уже есть старые A, AAAA, ALIAS или ANAME-записи, удалите или замените их. Иначе трафик уйдёт не в Hostamat.',
					'Поставьте TTL 300 секунд, если провайдер просит указать TTL вручную. Это ускорит первую проверку.',
					'Подождите, пока записи распространятся, затем нажмите «Запустить проверку».'
				],
				notes: [
					'Для корневого домена не используйте CNAME, если ваш DNS-провайдер не поддерживает flattening на корне.',
					'Если панель DNS автоматически добавляет имя зоны, не дублируйте корневой домен вручную.'
				]
			};
		}

		return {
			heading: 'Как настроить поддомен',
			intro:
				'Для адресов вида docs.example.com нужен TXT на владение и CNAME на платформенный резервный адрес сайта.',
			steps: [
				'Откройте DNS-зону родительского домена у DNS-провайдера.',
				'Создайте TXT-запись из блока «Подтверждение владения». Если провайдер автоматически дописывает имя зоны, обычно достаточно ввести только левую часть поддомена.',
				'Создайте CNAME-запись из блока «Маршрутизация» и направьте её на резервный домен сайта.',
				'Удалите старые A, AAAA или другой CNAME для этого же поддомена. У одного имени не должно быть конфликтующих записей.',
				'Поставьте TTL 300 секунд, если TTL задаётся вручную.',
				'После обновления DNS нажмите «Запустить проверку».'
			],
			notes: [
				'Для поддомена используйте именно CNAME, а не A-запись на IP Hostamat.',
				'Не добавляйте https:// и не копируйте конечную точку с лишней точкой в конце.'
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
			return 'исправить DNS';
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
	const statusDescription = $derived.by(() => {
		if (domain.type === 'generated') {
			return 'Резервный платформенный адрес остаётся доступным даже после подключения собственного домена.';
		}

		if (domain.status === 'active') {
			return 'Домен подтверждён и уже включён в активную схему маршрутизации Hostamat.';
		}

		if (domain.status === 'verified') {
			return 'DNS уже сошёлся. Повторная проверка нужна только если активация в прокси не завершилась.';
		}

		if (domain.status === 'verification_failed') {
			return 'Проверка увидела несоответствие. Исправьте записи ниже и запустите проверку снова.';
		}

		return 'Запись добавлена. Сначала внесите TXT и записи маршрутизации, затем подтвердите домен вручную.';
	});
	const timestampLabel = $derived.by(() => {
		if (domain.activatedAt) {
			return `активен с ${new Date(domain.activatedAt).toLocaleDateString('ru-RU')}`;
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

<div class="domain-card-content">
	<div class="domain-card-header">
		<div class="domain-card-title-block">
			<span class="card-kicker"
				>{domain.type === 'generated' ? 'платформенный адрес' : 'пользовательский домен'}</span
			>
			<div class="domain-card-title">
				{#if canOpen}
					<Link href={`https://${domain.hostname}`} target="_blank" alwaysUnderline
						>{domain.hostname}</Link
					>
				{:else}
					<span class="hostname">{domain.hostname}</span>
				{/if}
			</div>
		</div>
		<div class="status-stack">
			<span class={['status-chip', statusTone]}>{statusLabel}</span>
			<span class="timestamp-label">{timestampLabel}</span>
		</div>
	</div>

	<p class="status-description">{statusDescription}</p>

	<div class="signal-grid">
		<div class="signal-cell">
			<span class="signal-label">тип</span>
			<strong>{domain.isApex ? 'корневой домен' : 'поддомен'}</strong>
		</div>
		<div class="signal-cell">
			<span class="signal-label">резервный адрес</span>
			<Link href={`https://${domain.siteDomain}`} target="_blank" alwaysUnderline
				>{domain.siteDomain}</Link
			>
		</div>
		<div class="signal-cell signal-cell-wide">
			<span class="signal-label">сайт</span>
			<strong>{domain.siteName}</strong>
		</div>
	</div>

	{#if logMessage}
		<div class="incident-log">
			<span class="incident-label">результат проверки</span>
			<p>{logMessage}</p>
		</div>
	{/if}

	{#if isCustom && domain.instructions}
		<Details summary="Инструкции по DNS">
			<div class="records-list">
				{#if setupGuide}
					<div class="guide-block">
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

				<div class="record-block">
					<span class="record-label">подтверждение владения</span>
					<div class="record-meta">{domain.instructions.ownership.type}</div>
					<div class="record-value">{domain.instructions.ownership.name}</div>
					<div class="record-value">{domain.instructions.ownership.value}</div>
				</div>

				{#each domain.instructions.routing as record (`${record.type}-${record.name}-${record.value}`)}
					<div class="record-block">
						<span class="record-label">маршрутизация</span>
						<div class="record-meta">{record.type}</div>
						<div class="record-value">{record.name}</div>
						<div class="record-value">{record.value}</div>
					</div>
				{/each}
			</div>
		</Details>
	{/if}

	{#if isCustom}
		<div class="domain-toolbar">
			<Button disabled={loading} emphasis="solid" onclick={handleVerify}>
				{domain.status === 'active' ? 'Проверить повторно' : 'Запустить проверку'}
			</Button>
			<Button
				tone="danger"
				disabled={loading}
				onclick={() => {
					confirmingDelete = !confirmingDelete;
					actionError = '';
				}}
			>
				{confirmingDelete ? 'Отменить удаление' : 'Удалить'}
			</Button>
		</div>

		{#if confirmingDelete}
			<div class="delete-guard">
				<div class="delete-copy">
					<span class="incident-label">отключение домена</span>
					<p>Удаление снимет домен с маршрутизации, но не затронет сайт и его резервный адрес.</p>
				</div>
				<Button tone="danger" emphasis="solid" disabled={loading} onclick={handleDelete}>
					Подтвердить удаление
				</Button>
			</div>
		{/if}
	{/if}
</div>

<style>
	.domain-card-content {
		display: flex;
		flex-direction: column;
		gap: 0.9rem;
		min-width: 0;
		padding: var(--space-lg) 0;
	}

	.domain-card-header {
		display: flex;
		align-items: start;
		justify-content: space-between;
		gap: 0.9rem;
	}

	.domain-card-title-block {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
		min-width: 0;
	}

	.card-kicker,
	.signal-label,
	.timestamp-label,
	.incident-label,
	.record-label,
	.record-meta {
		font-size: 0.68rem;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: var(--color-text-tertiary);
	}

	.domain-card-title {
		min-width: 0;
	}

	.hostname {
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 700;
		line-height: 1.15;
		color: var(--color-text-regular);
		overflow-wrap: anywhere;
	}

	.status-stack {
		display: flex;
		flex-direction: column;
		align-items: end;
		gap: 0.35rem;
	}

	.signal-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 0.65rem;
	}

	.signal-cell {
		display: flex;
		flex-direction: column;
		gap: 0.26rem;
		padding: 0.72rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-background-object);
		min-width: 0;
	}

	.signal-cell strong,
	.signal-cell :global(a) {
		font-size: 0.88rem;
		line-height: 1.35;
		overflow-wrap: anywhere;
	}

	.signal-cell-wide {
		grid-column: 1 / -1;
	}

	.status-chip {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.22rem 0.65rem;
		border-radius: 999px;
		border: 1px solid var(--chip-border);
		color: var(--chip-text);
		background: var(--chip-background);
		font-size: 0.68rem;
		font-weight: 700;
		letter-spacing: 0.03em;
		white-space: nowrap;
		text-transform: uppercase;
	}

	.status-chip.success {
		--chip-border: var(--color-success-border);
		--chip-text: var(--color-success-text);
		--chip-background: var(--color-success-surface);
	}

	.status-chip.pending {
		--chip-border: var(--color-accent-border);
		--chip-text: var(--color-accent-text);
		--chip-background: var(--color-accent-surface);
	}

	.status-chip.danger {
		--chip-border: var(--color-danger-border);
		--chip-text: var(--color-danger-text);
		--chip-background: var(--color-danger-surface);
	}

	.status-description {
		color: var(--color-text-secondary);
		max-width: 52ch;
		font-size: 0.92rem;
	}

	.incident-log {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		padding: 0.8rem;
		border-radius: var(--radius-md);
		border: 1px solid color-mix(in hsl, var(--color-danger-border) 50%, transparent);
		background: color-mix(in hsl, var(--color-danger-surface) 92%, transparent);
	}

	.incident-log p {
		color: var(--color-danger-text);
		font-size: 0.88rem;
	}

	.records-list {
		display: flex;
		flex-direction: column;
		gap: 0.65rem;
		padding-top: 0.25rem;
	}

	.guide-block {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		padding: 0.9rem;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: var(--color-background-object);
	}

	.guide-block h3 {
		font-size: 0.95rem;
		font-weight: 700;
	}

	.guide-steps {
		display: flex;
		flex-direction: column;
		gap: 0.45rem;
		padding-left: 1.15rem;
		color: var(--color-text-secondary);
	}

	.guide-notes {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		padding-top: 0.3rem;
		border-top: 1px solid var(--color-border);
	}

	.guide-notes p {
		font-size: 0.84rem;
		color: var(--color-text-tertiary);
	}

	.record-block {
		display: grid;
		grid-template-columns: minmax(0, 1fr);
		gap: 0.18rem;
		padding: 0.75rem;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: var(--color-background-object);
	}

	.record-label {
		margin-bottom: 0.08rem;
	}

	.record-value {
		overflow-wrap: anywhere;
		font-family: var(--font-mono, monospace);
		font-size: 0.84rem;
	}

	.domain-toolbar {
		display: flex;
		gap: 0.6rem;
		flex-wrap: wrap;
		margin-top: auto;
	}

	.delete-guard {
		display: flex;
		flex-wrap: wrap;
		align-items: end;
		justify-content: space-between;
		gap: 0.8rem;
		padding: 0.9rem;
		border-radius: var(--radius-md);
		border: 1px solid color-mix(in hsl, var(--color-danger-border) 45%, transparent);
		background: color-mix(in hsl, var(--color-danger-surface) 90%, transparent);
	}

	.delete-copy {
		display: flex;
		flex-direction: column;
		gap: 0.24rem;
		max-width: 34ch;
	}

	.delete-copy p {
		color: var(--color-danger-text);
		font-size: 0.88rem;
	}

	@media (max-width: 640px) {
		.domain-card-header,
		.delete-guard {
			flex-direction: column;
			align-items: start;
		}

		.status-stack {
			align-items: start;
		}

		.signal-grid {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
