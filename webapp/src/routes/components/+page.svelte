<script lang="ts">
	import { Alert, Button, Card, Input, Link, Separator } from '$lib/components';
	import { Copy, Download, Plus, Spinner, Trash, Upload, X } from '$lib/icons';
	import Check from '$lib/icons/Check.svelte';
</script>

<svelte:head>
	<title>Компоненты | Хостамат</title>
</svelte:head>

<div style="display:flex; flex-direction:column; gap:20px; padding:10px;">
	<h1>Buttons</h1>

	<Button>Snippet</Button>

	<h2>Button kinds</h2>

	<div style="display: flex; gap: 10px">
		<Button kind="primary">Primary</Button>
		<Button kind="confirm">Confirm</Button>
		<Button kind="danger">Danger</Button>
		<Button kind="action">Action</Button>
	</div>

	<h3>Filled</h3>

	<div style="display: flex; gap: 10px">
		<Button kind="primary" filled>Primary</Button>
		<Button kind="confirm" filled>Confirm</Button>
		<Button kind="danger" filled>Danger</Button>
		<Button kind="action" filled>Action</Button>
	</div>

	<h2>On click event</h2>

	<Button onclick={(event) => alert('Hey!')}>Click me!</Button>

	<h1>Input</h1>
	<div style="width: 300px;">
		<Input id="input" label="Input" />
	</div>

	<h1>Link</h1>
	<Link href="/components">Link</Link>
	<Link href="/components" alwaysUnderline>Always underline link</Link>

	<h1>Cards</h1>

	<h2>Regular cards</h2>

	<div style="display: flex; gap: 10px">
		<Card>
			<div style="padding: 10px;">Card content</div>
		</Card>

		<Card kind="confirm">
			<div style="padding: 10px;">Card content</div>
		</Card>

		<Card kind="danger">
			<div style="padding: 10px;">Card content</div>
		</Card>

		<Card kind="action">
			<div style="padding: 10px;">Card content</div>
		</Card>
	</div>

	<h2>Floating card</h2>

	<div style="width: 300px;">
		<Card floating>
			<div style="padding: 10px;">Floating card</div>
		</Card>
	</div>

	<h2>Error card</h2>

	<div style="width: 300px;">
		<Alert variant="error" title="Error message" />
	</div>

	<h2>Success card</h2>

	<div style="width: 300px;">
		<Alert variant="success" title="Success message" />
	</div>

	<h1>Separator</h1>

	<Separator />

	<h2>Vertical</h2>

	<div style="height: 100px; padding: 5px;">
		<Separator orientation="vertical" />
	</div>

	<h1>Icons</h1>
	<div style="display: flex; gap: 10px">
		<div>
			<h2>Plus</h2>
			<Plus />
		</div>
		<div>
			<h2>Copy</h2>
			<Copy />
		</div>
		<div>
			<h2>Trash</h2>
			<Trash />
		</div>
		<div>
			<h2>Upload</h2>
			<Upload />
		</div>
		<div>
			<h2>Download</h2>
			<Download />
		</div>
		<div>
			<h2>X</h2>
			<X animated />
		</div>
		<div>
			<h2>Check</h2>
			<Check />
		</div>
		<div>
			<h2>Spinner</h2>
			<Spinner />
		</div>
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
