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
    border-top: var(--border-regular);
  }

  .vertical {
    border-left: var(--border-regular);
  }
</style>
