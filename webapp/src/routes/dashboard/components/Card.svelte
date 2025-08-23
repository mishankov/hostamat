<script lang="ts">
  let {
    children,
    floating = false,
    ...restProps
  } = $props()
</script>

<div class={["card", floating && "floating"]} {...restProps}>
  {@render children()}
</div>

<style>
  .card {
    width: fit-content;

    border: var(--border-regular);
    border-radius: var(--border-radius-regular);

    background-color: var(--color-background-object);

    transition: var(--transition);
  }

  .floating:hover {
    box-shadow: var(--shadow);
    transform: scale(1.01);
  }
</style>
