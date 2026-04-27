<script lang="ts">
	import { Button, Input, Modal } from '$lib/components';
	import { Spinner } from '$lib/icons';
	import { createDomain } from '$lib/api-client';
	import type { Site } from '$lib/types';

	interface Props {
		sites: Site[];
		onrefresh: () => Promise<void>;
	}

	let { sites, onrefresh = async () => {} }: Props = $props();

	let modal = $state<Modal>();
	let loading = $state(false);
	let hostname = $state('');
	let selectedSiteId = $state('');
	let errorMessage = $state('');

	const validateSubdomainHostname = (value: string): string | null => {
		const normalized = value.trim().toLowerCase().replace(/\.$/, '');

		if (!normalized) {
			return 'Поддомен обязателен';
		}

		if (normalized.includes('://') || /[/:@]/.test(normalized)) {
			return 'Укажите поддомен без схемы и пути';
		}

		const labels = normalized.split('.');
		if (labels.length < 3) {
			return 'Подключите поддомен, например www.example.com или docs.example.com';
		}

		const validLabel = /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/;
		if (!labels.every((label) => validLabel.test(label))) {
			return 'Укажите корректный поддомен';
		}

		return null;
	};

	const resetState = () => {
		selectedSiteId = sites[0]?.id ?? '';
		hostname = '';
		errorMessage = '';
	};

	const submit = async () => {
		if (!selectedSiteId) {
			errorMessage = 'Сначала выберите сайт';
			return;
		}

		const validationError = validateSubdomainHostname(hostname);
		if (validationError) {
			errorMessage = validationError;
			return;
		}

		loading = true;
		errorMessage = '';

		const { data, error } = await createDomain(selectedSiteId, hostname);
		if (error) {
			errorMessage = error;
			loading = false;
			return;
		}

		if (data) {
			await onrefresh();
			loading = false;
			modal?.close();
		}
	};

	resetState();
</script>

<Button
	emphasis="solid"
	disabled={loading}
	onclick={() => {
		resetState();
		modal?.open();
	}}
>
	{#if loading}
		<Spinner size={16} />
		<span>Добавляем домен</span>
	{:else}
		<span>Добавить домен</span>
	{/if}
</Button>

<Modal
	bind:this={modal}
	title="Подключить домен"
	titleId="new-domain-modal-title"
	initialFocusSelector="#new-domain-hostname"
	closeLabel="Закрыть форму подключения домена"
>
	<div class="modal-body">
		<div class="field-group">
			<label for="new-domain-site">Сайт</label>
			<select id="new-domain-site" bind:value={selectedSiteId}>
				{#each sites as site (site.id)}
					<option value={site.id}>{site.name} ({site.domain})</option>
				{/each}
			</select>
		</div>

		<Input
			id="new-domain-hostname"
			label="Поддомен"
			placeholder="www.example.com"
			{errorMessage}
			bind:value={hostname}
		/>

		<p class="modal-note">
			После добавления в карточке поддомена появятся точные записи `TXT` и `CNAME`.
		</p>

		<div class="modal-actions">
			<Button emphasis="solid" disabled={loading} onclick={submit}>Подключить</Button>
			<Button
				disabled={loading}
				onclick={() => {
					modal?.close();
				}}>Отмена</Button
			>
			{#if loading}
				<span class="action-spinner"><Spinner size={30} /></span>
			{/if}
		</div>
	</div>
</Modal>

<style>
	.modal-body {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: min(100%, 24rem);
	}

	.field-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}

	.field-group label {
		font-size: 0.8125rem;
		font-weight: 600;
		color: var(--color-text-secondary);
	}

	.field-group select {
		width: 100%;
		padding: 0.65rem 0.9rem;
		font-size: 0.95rem;
		font-family: var(--font-body);
		background-color: var(--color-background-object);
		color: var(--color-text-regular);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		transition: var(--transition);
	}

	.field-group select:focus {
		outline: none;
		border-color: var(--color-accent-border);
		box-shadow: var(--focus-ring);
	}

	.modal-note {
		max-width: 48ch;
		color: var(--color-text-secondary);
		font-size: 0.87rem;
		line-height: 1.5;
	}

	.modal-actions {
		display: flex;
		gap: var(--space-md);
		align-items: center;
		flex-wrap: wrap;
	}

	.action-spinner {
		display: inline-flex;
		align-items: center;
	}
</style>
