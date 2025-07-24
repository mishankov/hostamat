<script lang="ts">
  let {
    children,
    ...restProps
  } = $props()
</script>

<div class="card" {...restProps}>
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
</style>
