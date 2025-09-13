<script lang="ts">
  type InputProps = {
    id: string
    type?: string
    label?: string
    value?: string
    ref?: any
    oninput?: (e: Event) => void
  }

  let {
    id,
    type = "text",
    label = "",
    value = $bindable(""),
    ref = $bindable(null),
    oninput = undefined,
  }: InputProps = $props()
</script>

{#if label}
  <label for={id}>{label}</label>
{/if}
<input {id} {type} bind:value={value} bind:this={ref} />

<style>
  input {
    padding: 12px 20px;

    background-color: var(--color-primary-dim);
    color: var(--color-primary);

    border-width: 1px;
    border-style: solid;
    border-radius: 15px;
    border-color: var(--color-border);

    transition: var(--transition);
  }

  input:focus {
    outline: none;
    border-color: var(--color-primary);
  }

  label {
    background-color: var(--color-primary-dim);
    color: var(--color-primary);
  }
</style>
