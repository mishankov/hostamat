<script lang="ts">
	import { Check } from "$lib/icons";
  import Card from "./Card.svelte";

  let {
    message = ""
  } = $props()
</script>

{#if message}
  <div class="success">
    <Card kind="confirm" style="height: 100%; width: 100%;">
      <div class="content">
        <Check size={25}/>
        <div class="success-text">{message}</div>
      </div>
    </Card>
  </div>
{/if}

<style>
  .success {
    width: 100%;
    height: 50px;
  }

  .content {
    height: 100%;
    display: flex;
    align-items: center;

    padding-left: 10px;
  }

  .success-text {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    
    padding: 10px;
  }
</style>