<script lang="ts">
  import { Button } from "bits-ui";

  type ButtonKind = "primary" | "confirm" | "decline" | "action"

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

    border: 2px solid #ffffff;
    border-radius: 15px;

    padding: 12px 20px;

    text-align: center;
    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }

  .confirm {
    background-color: hsla(110, 100%, 50%, 0.25);
    border: 2px solid hsla(110, 100%, 50%, 1);
    color: hsla(110, 100%, 50%, 1);
  }

  .decline {
    background-color: hsla(338, 100%, 50%, 0.25);
    border: 2px solid hsla(338, 100%, 50%, 1);
    color: hsla(338, 100%, 50%, 1);
  }

  .action {
    background-color: hsla(200, 100%, 50%, 0.25);
    border: 2px solid hsla(200, 100%, 50%, 1);
    color: hsla(200, 100%, 50%, 1);
  }

  button:hover {
    box-shadow: rgba(0, 0, 0, 0.25) 0 8px 15px;
    transform: translateY(-1px);
  }

  button:active {
    box-shadow: none;
    transform: translateY(+1px);
  }
</style>

