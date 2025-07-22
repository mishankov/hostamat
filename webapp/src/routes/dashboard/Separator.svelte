<script lang="ts">
  import { Separator } from "bits-ui";

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
    border-top: 0.5px solid #555;
  }

  .vertical {
    border-left: 0.5px solid #555;
  }
</style>
