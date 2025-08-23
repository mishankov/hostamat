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

    border-radius: 15px;
    padding: 5px;

    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }

  button:hover {
    background-color: #555;
  }

  button:active {
    background-color: #333;
  }
</style>
