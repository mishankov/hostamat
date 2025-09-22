<script lang="ts">
  type ButtonKind = "primary" | "confirm" | "danger" | "action"

  type ButtonProps = {
    kind?: ButtonKind
    disabled?: boolean
    link?: string
    ref?: any
    children?: any
    onclick?: (e: Event) => void
  }

  let {
    kind = "primary",
    disabled = false,
    link = "",
    children = undefined,
    ref = $bindable(null),
    onclick = undefined,
  }: ButtonProps = $props()
</script>

{#if link}
  <a href={link}>
    <button {disabled} bind:this={ref} class="{kind}" onclick={onclick}>
      {@render children?.()}
    </button>
  </a>
{:else}
  <button {disabled} bind:this={ref} class="{kind}" onclick={onclick}>
    {@render children?.()}
  </button>
{/if}

<style>
  button {
    appearance: none;
    width: fit-content;

    border-width: 1px;
    border-style: solid;
    border-radius: 15px;

    padding: 12px 20px;

    text-align: center;
    user-select: none;

    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }

  .primary {
    background-color: var(--color-primary-dim);
    color: var(--color-primary);
    border-color: var(--color-primary);
  }

  .confirm {
    background-color: var(--color-confirm-dim);
    color: var(--color-confirm);
    border-color: var(--color-confirm);
  }

  .danger {
    background-color: var(--color-danger-dim);
    color: var(--color-danger);
    border-color: var(--color-danger);
  }

  .action {
    background-color: var(--color-action-dim);
    color: var(--color-action);
    border-color: var(--color-action);
  }

  button:hover:enabled {
    box-shadow: var(--shadow);
    transform: translateY(-1px);
  }

  button:active:enabled {
    box-shadow: none;
    transform: translateY(+1px);
  }

  button:disabled {
    background-color: #555;
    color: #fff;
    border-color: #fff;
  }
</style>

