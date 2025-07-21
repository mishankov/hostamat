<script lang="ts">
	import Button from "$lib/components/Button.svelte";
</script>

<div style="display:flex; flex-direction:column; gap:10px; padding:10px;">
  <h1>Buttons</h1>

  <Button>
    {#snippet child()}
      Snippet
    {/snippet}
  </Button>

  <h2>Button kinds</h2>

  <div style="display: flex; gap: 10px">
    <Button text="Primary" kind="primary" />
    <Button text="Confirm" kind="confirm" />
    <Button text="Decline" kind="decline" />
    <Button text="Action" kind="action" />
  </div>

  <h2>On click event</h2>

  <Button text="Click me!"  onclick={(event) => alert("Hey!")} />
</div>
