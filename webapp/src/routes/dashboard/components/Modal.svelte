<script lang="ts">
	import type { Snippet } from "svelte";

  let {
    content,
    kind = "primary"
  } : {
    content?: Snippet,
    kind?: "primary" | "danger"
  } = $props()

  let dialog = $state<HTMLDialogElement>();
  export const open = () => dialog?.showModal()
</script>

<dialog bind:this={dialog} class={kind}>
  {@render content?.()}
</dialog>

<style>
  dialog {
    top: 50%;
    left: 50%;
    transform: translateX(-50%) translateY(-50%);
    
    border-width: 0.5px;
    border-style: solid;
    border-radius: 5px;

    padding: 20px;

    background-color: #1a1a1a;
    color: #fff;

    transition: all 200ms cubic-bezier(.23, 1, 0.32, 1);
  }

  .primary {
    border-color: #555;
  }

  .danger {
    border-color: hsla(338, 100%, 50%, 1);
  }

  dialog[open] {
    animation: fadein 300ms ease-in forwards;
  }

  @keyframes fadein {
    0% {
      opacity: 0;
      transform: translateX(-50%) translateY(-50%) scale(.95);
    }
    100% {
      opacity: 1;
      transform: translateX(-50%) translateY(-50%) scale(1);
    }
  }

  ::backdrop {
    background-color: rgba(0, 0, 0, 0.5);
  }
</style>
