<script lang="ts">
	import Copy from "phosphor-svelte/lib/Copy";
	import Download from "phosphor-svelte/lib/Download";
	import Trash from "phosphor-svelte/lib/Trash";
	import Upload from "phosphor-svelte/lib/Upload";
	import Modal from "./components/Modal.svelte";
	import Button from "$lib/components/Button.svelte";

  let {
    site
  } = $props()

  let uploadModal = $state<Modal>()
  let deleteModal = $state<Modal>()
</script>

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
    <div class="site-toolbar-action" title="Загрузить сайт">
      <Upload onclick={() => uploadModal?.open()} />
      <Modal bind:this={uploadModal} kind="primary">
        {#snippet content()}
          <div class="upload-modal">
            <p>Загрузить сайт</p>
            <form method="dialog">
              <Button text="Загрузить" kind="primary" />
            </form>
          </div>
        {/snippet}
      </Modal>
    </div>

    <div class="site-toolbar-action" title="Скачать сайт">
      <Download />
    </div>  

    <div class="site-toolbar-action" title="Удалить сайт">
      <Trash onclick={() => deleteModal?.open()} />
      <Modal bind:this={deleteModal} kind="danger">
        {#snippet content()}
          <div class="delete-modal">
            <p>Вы уверены, что хотите удалить сайт {site.url}?</p>
            <form method="dialog">
              <Button text="Отмена" kind="primary" />
              <Button text="Удалить" kind="danger" />
            </form>
          </div>
        {/snippet}
      </Modal>
    </div> 
  </div>
</div>

<style>
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

  .upload-modal {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .delete-modal {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
</style>