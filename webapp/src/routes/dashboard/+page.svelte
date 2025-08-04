<script lang="ts">
  import Plus from "phosphor-svelte/lib/Plus";

	import Separator from "./components/Separator.svelte";
	import SiteContent from "./SiteContent.svelte";
	import Card from "./components/Card.svelte";

  import type { PageProps } from './$types';

  import { type Site } from './types';

	let { data }: PageProps = $props();
  let sites = $state<Site[]>([]);

  (async () => {
    sites = await data.sites;
  })()
</script>

<div class="dashboard">
  <h1>Твои сайты</h1>

  <div class="separator">
     <Separator/>
  </div>
 
  <div class="site-list">
    {#await data.sites}
      <p>Loading...</p>
    {:then _}
      {#each sites as site (site.id)}
        <Card floating>
          <SiteContent site={site} />
        </Card>
      {/each}
    

    <Card floating>
      <!-- TODO: should prob fix  -->
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <div class="add-site" role="button" tabindex="0" onclick={() => sites.push({
        id: `site-${sites.length + 1}`,
        domain: "new-site.hostamat.ru",
        updatedAt: "2025-03-01",
        version: "v1"
      })}>
        <p class="add-site-icon"><Plus /></p>
        <p class="add-site-text">Добавить сайт</p>
      </div>
    </Card>

    {:catch error}
      <p>Error loading sites: {error.message}</p>
    {/await}
  </div>
</div>

<style>
  .dashboard {
    display: flex;
    flex-direction: column;
    gap: 10px;
    
    width: clamp(500px, 50%, 1200px);

    border: 0.5px solid #555;
    border-radius: 15px;

    padding: 20px;
  }

  .separator {
    padding-bottom: 20px;
  }

  .site-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
  
  .add-site {
    height: 100%;

    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    padding-left: 20px;
    padding-right: 20px;
  }

  .add-site-icon {
    font-size: 20px;
    user-select: none;
  }

  .add-site-text {
    color: #bbb;
    user-select: none;
  }
</style>
