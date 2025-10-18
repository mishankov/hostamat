<script lang="ts">
  let {
    size = 100,
    color = "var(--color-primary)"
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
    stroke={color}
    stroke-width="7"
    stroke-dasharray="60 188"
    stroke-dashoffset="0"
    stroke-linecap="round"
    fill="none"
  />
</svg>

<style>
  circle {
    animation: rotate 1.4s linear infinite;
  }

  @keyframes rotate {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>