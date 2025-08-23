<script lang="ts">
	import SiteUploaderItem from "$lib/SiteUploaderItem.svelte";
	import { onMount } from 'svelte';

	interface Site {
		id: string;
		subdomain: string;
		domain: string;
		createdAt: string; // ISO date string from JSON
	}

	type CreateSiteResponse = Site;

	let sites: Site[] = [];

	onMount(async () => {
		await fetchSites();
	});

	async function fetchSites() {
		try {
			const response = await fetch("/api/listSites");
			sites = await response.json();
		} catch (error) {
			console.error('Failed to fetch sites:', error);
		}
	}

	async function createSite() {
		try {
			const response = await fetch('/api/createSite', {
				method: 'POST'
			});
			
			if (response.ok) {
				await fetchSites();
			}
		} catch (error) {
			console.error('Failed to create site:', error);
		}
	}
</script>

<div class="site-management">
	<div class="create-site">
		<button on:click={createSite}>Create Site</button>
	</div>

	{#each sites as site}
		<SiteUploaderItem {site} />
	{/each}
</div>

<style>
	.site-management {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.create-site {
		display: flex;
		gap: 1rem;
		margin-bottom: 2rem;
	}

</style>
