<script lang="ts">
  type CardKind = "primary" | "confirm" | "danger" | "action"
  let {
    children,
    floating = false,
    kind = "primary",
    ...restProps
  } = $props()
</script>

<div class={["card", floating && "floating", kind]} {...restProps}>
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

  .primary {
    border-color: var(--color-border);
  }

  .confirm {
    border-color: var(--color-confirm);
    background-color: var(--color-confirm-dim);
  }

  .danger {
    border-color: var(--color-danger);
    background-color: var(--color-danger-dim);
  }

  .action {
    border-color: var(--color-action);
    background-color: var(--color-action-dim);
  }

  .floating:hover {
    box-shadow: var(--shadow);
    transform: scale(1.01);
  }
</style>
