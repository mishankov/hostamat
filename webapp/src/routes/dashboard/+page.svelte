<script lang="ts">
  import Plus from "phosphor-svelte/lib/Plus";

	import Separator from "./components/Separator.svelte";
  import Card from "./components/Card.svelte";
	import SiteContent from "./SiteContent.svelte";

  let sites = $state([
    {
      id: "site-1",
      url: "abc-xy-z.hostamat.ru",
      date: "2025-01-01",
      version: "v22"
    }, {
      id: "site-2",
      url: "xyz-ab-c.hostamat.ru",
      date: "2025-02-01",
      version: "v3"
    }
  ])
</script>

<div class="dashboard">
  <h1>Твои сайты</h1>

  <div class="separator">
     <Separator/>
  </div>
 
  <div class="site-list">
    {#each sites as site (site.id)}
      <Card>
        <SiteContent site={site} />
      </Card>
    {/each}

    <Card>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <div class="add-site" role="button" tabindex="0" onclick={() => sites.push({
        id: `site-${sites.length + 1}`,
        url: "new-site.hostamat.ru",
        date: "2025-03-01",
        version: "v1"
      })}>
        <p class="add-site-icon"><Plus /></p>
        <p class="add-site-text">Добавить сайт</p>
      </div>
    </Card>
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
