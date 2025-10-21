<script lang="ts">
  let {
    size = 100,
    color = "var(--color-primary)",
  } = $props()
</script>

<svg 
  width={size}
  viewBox="-50 -50 100 100"
>
  <circle cx="0" cy="0" r="40" stroke-width="7" stroke={color} fill="transparent"/>
  
  <line
    x1="0" y1="-20"
    x2="0" y2="-20"
    stroke={color}
    stroke-width="10"
    stroke-linecap="round"
    stroke-linejoin="round"
  />
  
  <line
    x1="0" y1="0"
    x2="0" y2="20"
    stroke={color}
    stroke-width="10"
    stroke-linecap="round"
    stroke-linejoin="round"
  />
</svg>

<style>
  line {
    transition: all 500ms;
  }

  svg:hover line.animated {
    transform: rotate(90deg) ;
  }
</style>
