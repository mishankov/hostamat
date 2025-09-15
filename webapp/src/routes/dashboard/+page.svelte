<script lang="ts">
  import Plus from "phosphor-svelte/lib/Plus";

	import Separator from "../../lib/components/Separator.svelte";
	import SiteContent from "./SiteContent.svelte";
	import Card from "../../lib/components/Card.svelte";

  import type { PageProps } from './$types';

  import { type Site } from './types';
	import { invalidateAll } from "$app/navigation";

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
        <Card style="height: 100px;">
          <SiteContent site={site} />
        </Card>
      {/each}
    

    <Card floating style="height: 100px;">
      <!-- TODO: should prob fix  -->
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <div class="add-site" role="button" tabindex="0" onclick={async () => {
        const response = await fetch('/api/createSite', {
              method: 'POST',
            })

            await invalidateAll()
            sites = await data.sites;
      }}>
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
    
    width: clamp(500px, 75%, 1200px);

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
