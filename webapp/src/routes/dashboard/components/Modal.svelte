<script lang="ts">
	import type { Snippet } from "svelte";

  let {
    children,
    kind = "primary"
  } : {
    children?: Snippet,
    kind?: "primary" | "danger"
  } = $props()

  let dialog = $state<HTMLDialogElement>();
  export const open = () => dialog?.showModal()
</script>

<dialog bind:this={dialog} class={kind} onclick={(e: MouseEvent) => {
    if (e.target === dialog) {
      dialog?.close()
    }
  }}>
  <div class="modal-content">
    {@render children?.()}
  </div>
</dialog>

<style>
  dialog {
    top: 50%;
    left: 50%;
    transform: translateX(-50%) translateY(-50%);
    
    border: var(--border-regular);
    border-radius: var(--border-radius-regular);

    padding: 0;

    background-color: var(--color-background-object);
    color: var(--color-text);

    transition: var(--transition);
  }

  .primary {
    border-color: var(--color-border);
  }

  .danger {
    border-color: var(--color-danger);
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

  div.modal-content {
    padding: 20px;
  }
</style>
