<script lang="ts">
  let {
    size = 100,
    color = "var(--color-primary)"
  } = $props()
</script>

<svg 
  width={size}
  viewBox="-50 -50 100 100"
>
  <line
    x1="-30" y1="-5"
    x2="-5" y2="20"
    stroke={color}
    stroke-width="7"
    stroke-linecap="round"
    stroke-linejoin="round"
  />

  <line
    x1="-5" y1="20"
    x2="40" y2="-25"
    stroke={color}
    stroke-width="7"
    stroke-linecap="round"
    stroke-linejoin="round"
  />
</svg>
