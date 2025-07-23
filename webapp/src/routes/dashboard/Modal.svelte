<script lang="ts">
	import type { Snippet } from "svelte";

  let {
    ref = $bindable(),
    content,
  } : {
    ref?: HTMLDialogElement,
    content?: Snippet<[HTMLDialogElement]>,
  } = $props()
</script>

<dialog bind:this={ref}>
  {@render content?.(ref)}
</dialog>

<style scoped>
  dialog {
    border: 0.5px solid #555;
    border-radius: 5px;

    padding: 20px;

    background-color: #1a1a1a;
    color: #fff;

    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }
</style>
