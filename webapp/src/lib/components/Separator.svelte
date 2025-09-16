<script lang="ts">
  let {
    orientation = "horizontal"
  } : {
    orientation?: "horizontal" | "vertical"
  } = $props()
</script>

<div class="separator {orientation}"></div>

<style>
  .separator {
    border: 0;
  }

  .horizontal {
    width: 100%;
    border-top: var(--border-regular);
  }

  .vertical {
    height: 100%;
    border-left: var(--border-regular);
  }
</style>
