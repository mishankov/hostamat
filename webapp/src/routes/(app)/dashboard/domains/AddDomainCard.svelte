<script lang="ts">
	import { Button, Card, Input, Modal } from '$lib/components';
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

		if (!hostname.trim()) {
			errorMessage = 'Домен обязателен';
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

<div class="add-domain-card-shell">
	<Card class="add-domain-card" style="height: 100%;">
		<button
			type="button"
			class="add-domain"
			onclick={() => {
				resetState();
				modal?.open();
			}}
		>
			{#if loading}
				<div class="prompt-line">
					<span class="prompt-symbol">+</span>
					<span class="prompt-text">подключаем домен</span>
					<Spinner size={16} />
				</div>
			{:else}
				<div class="prompt-line">
					<span class="prompt-symbol">+</span>
					<span class="prompt-text">подключить домен</span>
				</div>
			{/if}

			<p class="prompt-description">
				Привяжите адрес к сайту, а DNS-инструкции появятся прямо в карточке домена.
			</p>
		</button>
	</Card>
</div>

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
			label="Домен"
			placeholder="example.com"
			{errorMessage}
			bind:value={hostname}
		/>

		<p class="modal-note">
			После добавления в карточке домена появятся точные записи `TXT`, `CNAME`, `A` или `AAAA`.
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
	.add-domain-card-shell {
		display: flex;
		min-height: 10.625rem;
		width: 25rem;
		flex: 0 0 25rem;
	}

	.add-domain {
		all: unset;
		box-sizing: border-box;
		width: 100%;
		height: 100%;
		min-height: 10.625rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-md);
		padding: var(--space-lg) 1.3rem;
		cursor: pointer;
		text-align: center;
		border-radius: var(--radius-lg);
	}

	.prompt-line {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		font-family: var(--font-body);
		font-size: 0.9rem;
		font-weight: 500;
	}

	.prompt-symbol {
		color: var(--color-accent-strong);
		font-weight: 700;
		font-size: 1.1rem;
	}

	.prompt-text {
		color: var(--color-text-secondary);
	}

	.prompt-description {
		max-width: 28ch;
		color: var(--color-text-secondary);
		font-size: 0.88rem;
		line-height: 1.5;
	}

	:global(.add-domain-card) {
		border-style: dashed;
		border-color: var(--color-border);
		border-radius: var(--radius-lg);
		background: transparent;
		transition:
			border-color 180ms ease-out,
			background 180ms ease-out,
			box-shadow 180ms ease-out,
			transform 180ms ease-out;
	}

	:global(.add-domain-card:hover) {
		border-color: var(--color-accent-border);
		background: var(--color-accent-surface);
		box-shadow: var(--shadow-sm);
		transform: translateY(-2px);
	}

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
		background-color: var(--color-background-object);
		color: var(--color-text-regular);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		transition: var(--transition);
	}

	.field-group select:focus {
		outline: none;
		border-color: var(--color-accent-border);
		box-shadow: var(--focus-ring);
	}

	.modal-note {
		max-width: 42ch;
		color: var(--color-text-secondary);
		font-size: 0.88rem;
		line-height: 1.5;
	}

	.modal-actions {
		display: flex;
		gap: var(--space-lg);
		align-items: center;
		flex-wrap: wrap;
	}

	.action-spinner {
		display: inline-flex;
		align-items: center;
	}

	@media (max-width: 820px) {
		.add-domain-card-shell {
			width: 100%;
			flex-basis: 100%;
		}
	}
</style>
