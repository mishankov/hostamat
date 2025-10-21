<script lang="ts">
	import { Check } from "$lib/icons";
  import Card from "./Card.svelte";

  let {
    children = null,
    header = ''
  } = $props()
</script>

<div class="success">
  <Card kind="confirm" style="height: fit-content; width: 100%; padding: 10px;">
    <div class="header">
      <Check size={25} color="var(--color-confirm)" />
      <h4>{header}</h4>
    </div>
    {#if children}
      <div class="content">
        {@render children?.()}
      </div>
    {/if}
  </Card>
</div>


<style>
  .success {
    width: 100%;
  }

  .header {
    display: flex;
    flex-direction: row;
    gap: 10px;
    align-items: center;
  }

  .content {
    display: flex;
    align-items: start;
    gap: 10px;
    
    padding-top: 10px;
  }
</style>