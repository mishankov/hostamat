<script lang="ts">
  import type { PageProps } from './$types';
	import { invalidateAll } from "$app/navigation";
  
	import { Card, Error } from "$lib/components";
	import { Spinner, Plus } from "$lib/icons";
  
  import { type Site } from './types';
	import SiteContent from "./SiteContent.svelte";

	let { data }: PageProps = $props();
  let sites = $state<Site[]>([]);

  (async () => {
    sites = await data.sites;
  })()
</script>

<main>
  <div class="dashboard">
    <h1>Твои сайты</h1>
  
    <div class="site-list">
      {#await data.sites}
        <Spinner />
      {:then _}
        {#each sites as site (site.id)}
          <Card style="height: 110px;">
            <SiteContent site={site} />
          </Card>
        {/each}
      

      <Card floating style="height: 110px;">
        <button class="add-site" onclick={async () => {
          const response = await fetch('/api/createSite', {
                method: 'POST',
              })

              await invalidateAll()
              sites = await data.sites;
        }}>
          <p class="add-site-icon"><Plus size={40} /></p>
          <p class="add-site-text">Добавить сайт</p>
        </button>
      </Card>

      {:catch error}
        <Error message="Ошибка при загрузке сайтов: {error.message}" />
      {/await}
    </div>
  </div>
</main>


<style>
  main {
    width: 100%;

    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
  }
  
  .dashboard {
    display: flex;
    flex-direction: column;
    gap: 10px;
    
    width: clamp(500px, 75%, 1200px);

    padding: 20px;
  }

  .site-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
  
  .add-site {
    all: unset;

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
