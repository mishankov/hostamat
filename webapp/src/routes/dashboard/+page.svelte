<script>
	import Separator from "./Separator.svelte";
  import Card from "./Card.svelte";

  let sites = [
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
  ]
</script>

<div class="dashboard">
  <h1>Твои сайты</h1>

  <div class="separator">
     <Separator/>
  </div>
 
  <div class="site-list">
    {#each sites as site}
      <Card>
        {#snippet child()}
          <div class="site-content">
            <a href="https://{site.url}" target="_blank" class="site-url">{site.url}</a>
            <p class="site-data">
              Дата последней загрузки: {site.date}
            </p>
            <p class="site-data">
              Версия: {site.version}
            </p>

            <div class="site-toolbar">

            </div>
          </div>
        {/snippet}
      </Card>
    {/each}

    <Card>
        {#snippet child()}
          <div class="add-site">
            <p class="add-site-icon">+</p>
            <p class="add-site-text">Добавить сайт</p>
          </div>
        {/snippet}
    </Card>
  </div>
</div>

<style>
  .dashboard {
    display: flex;
    flex-direction: column;
    gap: 10px;

    max-width: 1200px;

    border: 0.5px solid #555;
    border-radius: 15px;

    padding: 20px;
  }

  .separator {
    padding-bottom: 20px;
  }

  .site-list {
    display: flex;
    gap: 10px;
  }
  
  .site-content {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }

  a.site-url {
    font-size: 20px;
    color: #fff;
    text-decoration: none;
    text-decoration: underline;

    padding-bottom: 10px;
  }

  p.site-data {
    font-size: 12px;
    color: #bbb;
  }
  
  .add-site {
    height: 100%;
    width: 100%;

    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
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
