<script lang="ts">
  let {
    child
  } = $props()
</script>

<div class="card">
  {@render child()}
</div>

<style>
  .card {
    border: 0.5px solid #bbb;
    border-radius: 5px;
    padding: 10px;

    background-color: #1a1a1a;

    transition: all 300ms cubic-bezier(.23, 1, 0.32, 1);
  }

  .card:hover {
    box-shadow: rgba(0, 0, 0, 0.25) 0 8px 15px;
    transform: scale(1.01);
  }
</style>
