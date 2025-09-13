<script lang="ts">
	import Card from "./Card.svelte";

  let {
    error = ""
  } = $props()
</script>

{#if error}
  <div class="error">
    <Card kind="danger" style="height: 100%; width: 100%;">
      <div class="error-text">{error}</div>
    </Card>
  </div>
{/if}

<style>
  .error {
    width: 100%;
    height: 50px;
  }

  .error-text {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    
    padding: 10px;
  }
</style>