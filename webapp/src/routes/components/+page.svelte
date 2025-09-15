<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import Copy from "$lib/icons/Copy.svelte";
</script>

<div style="display:flex; flex-direction:column; gap:10px; padding:10px;">
  <h1>Buttons</h1>

  <Button>
    Snippet
  </Button>

  <h2>Button kinds</h2>

  <div style="display: flex; gap: 10px">
    <Button kind="primary">Primary</Button>
    <Button kind="confirm">Confirm</Button>
    <Button kind="danger">Danger</Button>
    <Button kind="action">Action</Button>
  </div>

  <h2>On click event</h2>

  <Button onclick={(event) => alert("Hey!")}>Click me!</Button>

  <h1>Icons</h1>

  <Copy />
</div>
