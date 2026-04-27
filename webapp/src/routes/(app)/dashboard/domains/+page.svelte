<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';

	import { Alert } from '$lib/components';
	import { Spinner } from '$lib/icons';
	import type { DomainRecord, Site } from '$lib/types';

	import AddDomainCard from './AddDomainCard.svelte';
	import DomainCard from './DomainCard.svelte';

	let { data }: PageProps = $props();

	let sites = $state.raw<Site[]>([]);
	let domains = $state.raw<DomainRecord[]>([]);
	let loading = $state(true);
	let errorMessage = $state('');
	let customDomains = $derived(domains.filter((domain) => domain.type === 'custom'));

	const updateState = async () => {
		const [sitesResult, domainsResult] = await Promise.all([data.sitesResp, data.domainsResp]);

		if (sitesResult.data) {
			sites = sitesResult.data.sites ?? [];
		} else {
			sites = [];
		}

		if (domainsResult.data) {
			domains = domainsResult.data.domains ?? [];
			errorMessage = '';
		} else {
			domains = [];
			errorMessage = domainsResult.error;
		}

		if (sitesResult.error) {
			errorMessage = errorMessage ? `${sitesResult.error}. ${errorMessage}` : sitesResult.error;
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

	updateState();
</script>

<svelte:head>
	<title>Домены | Хостамат</title>
	<meta
		name="description"
		content="Подключайте собственные домены, проверяйте DNS и управляйте адресами сайтов."
	/>
</svelte:head>

<div class="dashboard-page">
	<div class="headline">
		<h1>Домены</h1>
		<div class="headline-meta">
			<p>Подключайте свои домены</p>
			{#if !loading}
				<span class="domain-counter">доменов: {customDomains.length}</span>
			{/if}
		</div>
	</div>

	{#if errorMessage}
		<Alert variant="error" message={errorMessage} />
	{/if}

	{#if loading}
		<div class="domain-list">
			<div class="loading-state" aria-live="polite">
				<Spinner size={50} />
			</div>
		</div>
	{:else if sites.length === 0}
		<div class="domain-list">
			<div class="limit-info-wrap">
				<Alert variant="info" title="Сначала создайте сайт">
					Пользовательский домен можно подключить только после создания сайта. Резервный адрес
					появится автоматически.
				</Alert>
			</div>
		</div>
	{:else if customDomains.length > 0 || true}
		<div class="domain-list">
			{#each customDomains as domain (domain.id)}
				<div class="domain-row">
					<DomainCard {domain} {sites} onrefresh={refresh} />
				</div>
			{/each}
		</div>
		<div class="add-action">
			<AddDomainCard {sites} onrefresh={refresh} />
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

	.domain-counter {
		font-family: var(--font-body);
		font-size: 0.82rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		white-space: nowrap;
	}

	.domain-list {
		display: flex;
		flex-direction: column;
		min-width: 0;
		border-top: 1px solid var(--color-border);
	}

	.domain-row {
		border-bottom: 1px solid var(--color-border);
		min-width: 0;
	}

	.loading-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 10.625rem;
	}

	.limit-info-wrap {
		min-height: 6.875rem;
		display: flex;
		align-items: center;
		padding: var(--space-lg) 0;
	}

	.add-action {
		padding-top: var(--space-sm);
	}

	@media (max-width: 820px) {
		.limit-info-wrap {
			min-height: auto;
		}
	}
</style>
