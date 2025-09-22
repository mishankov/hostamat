<script lang="ts">
  let {
    size = 100,
  } = $props()
</script>

<svg 
  width={size}
  viewBox="-50 -50 100 100"
>

  <g class="cap">
    <polyline
      points="
        -10 -20
        -5 -30
        5 -30
        10 -20
      "
      stroke="currentColor"
      stroke-width="5"
      stroke-linecap="round"
      stroke-linejoin="round"
      fill="none"
    />

    <line
      x1="-30" y1="-20"
      x2="30" y2="-20"
      stroke="currentColor"
      stroke-width="5"
      stroke-linecap="round"
      stroke-linejoin="round"
      fill="none"
    />
   </g>

  <line 
    x1="-8" y1="-10"
    x2="-8" y2="20"
    stroke="currentColor"
    stroke-width="5"
    stroke-linecap="round"
    stroke-linejoin="round"
  />

  <line 
    x1="8" y1="-10"
    x2="8" y2="20"
    stroke="currentColor"
    stroke-width="5"
    stroke-linecap="round"
    stroke-linejoin="round"
  />

  <polyline class="can"
    points="
      -25 -20
      -20 30
      20 30
      25 -20
    "
    stroke="currentColor"
    stroke-width="5"
    stroke-linecap="round"
    stroke-linejoin="round"
    fill="none"
  />
</svg>

<style>
  .cap {
    transform-box: fill-box;
    transform-origin: center;
    transition: all 500ms;
  }

  svg:hover .cap {
    transform: translateY(-10px) translateX(2px) rotate(8deg);
  }

  svg:active .cap {
    transform: translateY(-15px) translateX(4px) rotate(15deg);
  }
</style>