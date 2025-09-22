<script lang="ts">
  let {
    children,
    onclick,
    ...restProps
  } = $props()
</script>

<button {onclick} {...restProps}>
  {@render children()}
</button>

<style>
  button {
    all: unset;
    display: flex;
    align-items: center;
    justify-content: center;

    padding: 5px;
  }
</style>
