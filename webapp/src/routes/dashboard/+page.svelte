<script>
  import { Toolbar, Tooltip } from "bits-ui";

  import Upload from "phosphor-svelte/lib/Upload";
  import Download from "phosphor-svelte/lib/Download";
  import Trash from "phosphor-svelte/lib/Trash";
  import Copy from "phosphor-svelte/lib/Copy";

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
            <div class="site-url-container">
              <a href="https://{site.url}" target="_blank" class="site-url">{site.url}</a> 
              <Copy />
            </div>
            <p class="site-data">
              Дата последней загрузки: {site.date}
            </p>
            <p class="site-data">
              Версия: {site.version}
            </p>

            <div class="site-toolbar">
              <Toolbar.Root>
                {#snippet child()}
                  <Toolbar.Button>
                    {#snippet child({ props })}
                      <div class="site-toolbar-action" title="Загрузить сайт" {...props}>
                        <Upload />
                      </div>
                    {/snippet}
                  </Toolbar.Button>

                  <Toolbar.Button>
                    {#snippet child({ props })}
                      <div class="site-toolbar-action" title="Скачать сайт" {...props}>
                        <Download />
                      </div>  
                    {/snippet}
                  </Toolbar.Button>

                  <Toolbar.Button>
                    {#snippet child({ props })}
                      <div class="site-toolbar-action" title="Удалить сайт" {...props}>
                        <Trash />
                      </div> 
                    {/snippet}
                  </Toolbar.Button>
                {/snippet}
              </Toolbar.Root>
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

    padding-top: 20px;
    padding-left: 20px;
    padding-right: 20px;
    padding-bottom: 10px;
  }

  .site-url-container {
    display: flex;
    align-items: center;
    gap: 5px;

    padding-bottom: 10px;
  }

  a.site-url {
    font-size: 20px;
    color: #fff;
    text-decoration: none;
    text-decoration: underline;
  }

  p.site-data {
    font-size: 12px;
    color: #bbb;
  }

  .site-toolbar {
    width: fit-content;
    align-self: flex-end;
    padding-top: 10px;

    display: flex;
    gap: 5px;
  }

  .site-toolbar-action {
    display: flex;
    align-items: center;
    justify-content: center;

    border-radius: 15px;
    padding: 5px;

    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }

  .site-toolbar-action:hover {
    background-color: #555;
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
