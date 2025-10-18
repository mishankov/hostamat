<script lang="ts">
	import { X } from "$lib/icons";
import Card from "./Card.svelte";

  let {
    error = ""
  } = $props()
</script>

{#if error}
  <div class="error">
    <Card kind="danger" style="height: 100%; width: 100%;">
      <div class="content">
        <X size={20} />
        <div class="error-text">{error}</div>
      </div>
    </Card>
  </div>
{/if}

<style>
  .error {
    width: 100%;
    height: 50px;
  }

  .content {
    height: 100%;
    display: flex;
    align-items: center;

    padding-left: 10px;
  }

  .error-text {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    
    padding: 10px;
  }
</style>