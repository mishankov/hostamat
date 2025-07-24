<script lang="ts">
  import { Button } from "bits-ui";

  type ButtonKind = "primary" | "confirm" | "danger" | "action"

  type ButtonProps = {
    text?: string
    kind?: ButtonKind
    disabled?: boolean
    ref?: any
    child?: any
    onclick?: (e: Event) => void
  }

  let {
    kind = "primary",
    text = "",
    child = null,
    disabled = false,
    ref = $bindable(null),
    onclick = undefined,
  }: ButtonProps = $props()
</script>

<button {disabled} bind:this={ref} class="{kind}" onclick={onclick}>
  {#if child != null}
    {@render child()}
  {:else}
    {text}
  {/if}
</button>

<style>
  button {
    appearance: none;
    width: fit-content;

    color: #ffffff;
    background-color: #1A1A1A;

    border-width: 1px;
    border-style: solid;
    border-radius: 15px;

    padding: 12px 20px;

    text-align: center;
    user-select: none;

    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }

  .primary {
    background-color: #1A1A1A;
    color: #ffffff;
    border-color: #ffffff;
  }

  .confirm {
    background-color: hsla(110, 100%, 50%, 0.25);
    border-color: hsla(110, 100%, 50%, 1);
    color: hsla(110, 100%, 50%, 1);
  }

  .danger {
    background-color: hsla(338, 100%, 50%, 0.25);
    color: hsla(338, 100%, 50%, 1);
    border-color: hsla(338, 100%, 50%, 1);
  }

  .action {
    background-color: hsla(200, 100%, 50%, 0.25);
    color: hsla(200, 100%, 50%, 1);
    border-color: hsla(200, 100%, 50%, 1);
  }

  button:hover:enabled {
    box-shadow: rgba(0, 0, 0, 0.25) 0 8px 15px;
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

