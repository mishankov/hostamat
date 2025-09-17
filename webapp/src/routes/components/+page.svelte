<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import Card from "$lib/components/Card.svelte";
	import Error from "$lib/components/Error.svelte";
	import Input from "$lib/components/Input.svelte";
	import Separator from "$lib/components/Separator.svelte";
  import { Copy, Download, Plus, Trash, Upload } from "$lib/icons";
</script>

<div style="display:flex; flex-direction:column; gap:20px; padding:10px;">
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

  <h1>Input</h1>
  <div style="width: 300px;">
    <Input id="input" label="Input" />
  </div>

  <h1>Cards</h1>

  <h2>Regular cards</h2>

  <div style="display: flex; gap: 10px">
    <Card>
      <div style="padding: 10px;">
        Card content
      </div>
    </Card>

    <Card kind="confirm">
      <div style="padding: 10px;">
        Card content
      </div>
    </Card>

    <Card kind="danger">
      <div style="padding: 10px;">
        Card content
      </div>
    </Card>

    <Card kind="action">
      <div style="padding: 10px;">
        Card content
      </div>
    </Card>
  </div>
  
  <h2>Floating card</h2>

  <Card floating>
    <div style="padding: 10px;">
      Floating card
    </div>
  </Card>

  <h2>Error card</h2>
  
  <div style="width: 300px;">
    <Error error="Error message" />
  </div>

  <h1>Separator</h1>

  <Separator />

  <h2>Vertical</h2>

  <div style="height: 100px; padding: 5px;">
    <Separator orientation="vertical" />
  </div>

  <h1>Icons</h1>
  <div style="display: flex; gap: 10px">
    <div><h2>Plus</h2><Plus /></div>
    <div><h2>Copy</h2> <Copy /></div>
    <div><h2>Trash</h2> <Trash /></div>
    <div><h2>Upload</h2> <Upload /></div>
    <div><h2>Download</h2> <Download /></div>
  </div>
</div>


<style>
  h1 {
    padding-top: 20px;
  }

  h2 {
    padding-top: 10px;
  }
</style>