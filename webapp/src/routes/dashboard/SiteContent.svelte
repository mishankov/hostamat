<script lang="ts">
	import { Button, Link, Modal } from "$lib/components";
  import  { Copy, Download, Trash, Upload } from "$lib/icons";

	import FileUploader from "./components/FileUploader.svelte";
	import IconButton from "./components/IconButton.svelte";
  import { type Site } from './types';

  interface Props {
    site: Site
  }

  let {
    site
  }: Props = $props()

  let uploadModal = $state<Modal>()
  let deleteModal = $state<Modal>()
</script>

<div class="site-content">
  <div class="site-url-container">
    <Link href="https://{site.domain}" target="_blank" alwaysUnderline>{site.domain}</Link> 
    <IconButton title="Скопировать URL" onclick={()=>{navigator.clipboard.writeText(`https://${site.domain}`)}}>
      <Copy size={20} />
    </IconButton>
  </div>

  <div class="site-toolbar">
    <IconButton onclick={() => uploadModal?.open()} title="Загрузить сайт">
      <Upload size={20} />
    </IconButton>
    <Modal bind:this={uploadModal} kind="primary">
      <div class="upload-modal">
        <h3>Загрузить файлы для {site.domain}</h3>
        <FileUploader multiple={true} uploadUrl="/api/uploadFile" additionalHeaders={{
          'H-Site-Id': site.id
        }}/>
      </div>
    </Modal>

    <IconButton title="Скачать сайт" onclick={()=>{}}>
      <Download size={20} />
    </IconButton> 

    <IconButton title="Удалить сайт" onclick={() => deleteModal?.open()}>
      <Trash size={20} />
    </IconButton>
    <Modal bind:this={deleteModal} kind="danger">
      <div class="delete-modal">
        <p>Вы уверены, что хотите удалить сайт {site.domain}?</p>
        <form method="dialog">
          <Button kind="primary">Отмена</Button>
          <Button kind="danger" onclick={async () => {
            const response = await fetch('/api/deleteSite', {
              method: 'POST',
              headers: {
                'H-Site-Id': site.id
              }
            })

            window.location.reload()
          }}>Удалить</Button>
        </form>
      </div>
    </Modal>
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
    gap: 3px;

    padding-bottom: 10px;
  }

  .site-toolbar {
    width: fit-content;
    align-self: flex-end;
    padding-top: 10px;

    display: flex;
    gap: 5px;
  }

  .upload-modal {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .delete-modal {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .delete-modal > form {
    display: flex;
    gap: 10px;
  }
</style>