<script lang="ts">
  let {
    size = 100,
  } = $props()
</script>

<svg 
  width={size}
  viewBox="-50 -50 100 100"
  class="spinner"
>
  <circle
    cx="0"
    cy="0"
    r="40"
    stroke="var(--color-primary)"
    stroke-width="5"
    stroke-dasharray="60 20"
    stroke-linecap="round"
    fill="none"
  />
</svg>

<style>
  .spinner {
    animation: rotate 2s linear infinite;
  }

  @keyframes rotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>