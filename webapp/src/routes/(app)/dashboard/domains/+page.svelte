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
	let generatedDomains = $derived(domains.filter((domain) => domain.type === 'generated'));
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

	const refresh = async () => {
		loading = true;
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

<div class="domains-page">
	<div class="headline">
		<h1>Домены</h1>
		<div class="headline-meta">
			<p>Подключайте пользовательские домены и управляйте адресами сайтов в одном списке.</p>
			{#if !loading}
				<span class="domain-counter">доменов: {domains.length}</span>
			{/if}
		</div>
	</div>

	{#if errorMessage}
		<Alert variant="error" message={errorMessage} />
	{/if}

	{#if loading}
		<div class="domains-list">
			<div class="loading-state" aria-live="polite">
				<Spinner size={50} />
			</div>
		</div>
	{:else if sites.length === 0}
		<div class="domains-list">
			<div class="empty-state">
				<Alert variant="info" title="Сначала создайте сайт">
					Пользовательский домен можно подключить только после создания сайта. Резервный адрес
					появится автоматически.
				</Alert>
			</div>
		</div>
	{:else}
		<section class="domain-section">
			<div class="section-head">
				<div>
					<h2>Пользовательские домены</h2>
					<p>Ваши собственные адреса с проверкой DNS, подтверждением владения и управлением.</p>
				</div>
				<span class="section-counter">{customDomains.length}</span>
			</div>

			<div class="domains-list">
				{#each customDomains as domain (domain.id)}
					<div class="domain-card-shell">
						<DomainCard {domain} onrefresh={refresh} />
					</div>
				{/each}

				<AddDomainCard {sites} onrefresh={refresh} />
			</div>
		</section>

		<section class="domain-section">
			<div class="section-head">
				<div>
					<h2>Платформенные адреса</h2>
					<p>Резервные домены Hostamat для каждого сайта. Они остаются доступны всегда.</p>
				</div>
				<span class="section-counter">{generatedDomains.length}</span>
			</div>

			<div class="domains-list">
				{#each generatedDomains as domain (domain.id)}
					<div class="domain-card-shell">
						<DomainCard {domain} onrefresh={refresh} />
					</div>
				{/each}
			</div>
		</section>
	{/if}
</div>

<style>
	.domains-page {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 0;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		padding-bottom: var(--space-lg);
		border-bottom: 1px solid var(--color-border);
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

	.domain-section {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.section-head {
		display: flex;
		align-items: end;
		justify-content: space-between;
		gap: var(--space-md);
		flex-wrap: wrap;
	}

	.section-head h2 {
		font-size: 1.15rem;
		font-weight: 700;
	}

	.section-head p {
		max-width: 60ch;
		color: var(--color-text-secondary);
		font-size: 0.92rem;
	}

	.section-counter {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: 2.5rem;
		padding: 0.3rem 0.7rem;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: var(--color-background-object);
		color: var(--color-text-tertiary);
		font-size: 0.82rem;
		font-weight: 600;
	}

	.domains-list {
		--domain-tile-width: 25rem;

		display: flex;
		flex-wrap: wrap;
		gap: var(--space-lg);
		align-items: stretch;
	}

	.loading-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 10.625rem;
		width: var(--domain-tile-width);
	}

	.empty-state {
		width: min(34rem, 100%);
	}

	.domain-card-shell {
		display: flex;
		min-height: 10.625rem;
		width: var(--domain-tile-width);
		flex: 0 0 var(--domain-tile-width);
		min-width: 0;
	}

	.domain-card-shell :global(.card) {
		height: 100%;
		border-radius: var(--radius-lg);
	}

	@media (max-width: 820px) {
		.domains-list {
			--domain-tile-width: 100%;
		}

		.empty-state {
			width: 100%;
		}
	}
</style>
