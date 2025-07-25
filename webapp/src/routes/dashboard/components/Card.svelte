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

    border: 0.5px solid #555;
    border-radius: 5px;

    background-color: #1a1a1a;

    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }

  .floating:hover {
    box-shadow: rgba(0, 0, 0, 0.25) 0 8px 15px;
    transform: scale(1.01);
  }
</style>
