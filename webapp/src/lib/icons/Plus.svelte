<script lang="ts">
  let {
    size = 100,
  } = $props()
</script>

<svg 
  width={size}
  viewBox="-50 -50 100 100"
>
  <line
    x1="0" y1="-30"
    x2="0" y2="30"
    stroke="var(--color-primary)"
    stroke-width="5"
  />

  <line
    x1="-30" y1="0"
    x2="30" y2="0"
    stroke="var(--color-primary)"
    stroke-width="5"
  />
</svg>

<style>
  line {
    transition: all 500ms;
  }

  svg:hover line {
    transform: rotate(90deg) ;
  }
</style>