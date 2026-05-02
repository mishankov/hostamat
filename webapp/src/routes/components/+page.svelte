<script lang="ts">
	import { Alert, Button, Card, Input, Link, Separator } from '$lib/components';
	import {
		Check,
		Copy,
		Download,
		Edit,
		I as Info,
		Moon,
		Plus,
		Spinner,
		Sun,
		Trash,
		Upload,
		Warning,
		X
	} from '$lib/icons';
</script>

<svelte:head>
	<title>Компоненты | Хостамат</title>
</svelte:head>

<div style="display:flex; flex-direction:column; gap:20px; padding:10px;">
	<h1>Buttons</h1>

	<Button>Snippet</Button>

	<h2>Button tones</h2>

	<div style="display: flex; gap: 10px">
		<Button tone="accent">Accent</Button>
		<Button tone="success">Success</Button>
		<Button tone="danger">Danger</Button>
	</div>

	<h3>Solid</h3>

	<div style="display: flex; gap: 10px">
		<Button tone="accent" emphasis="solid">Accent</Button>
		<Button tone="success" emphasis="solid">Success</Button>
		<Button tone="danger" emphasis="solid">Danger</Button>
	</div>

	<h3>Sizes</h3>

	<div style="display: flex; align-items: center; gap: 10px">
		<Button size="xs" emphasis="solid">Extra small</Button>
		<Button size="sm" emphasis="solid">Small</Button>
		<Button size="default" emphasis="solid">Default</Button>
		<Button size="lg" emphasis="solid">Large</Button>
		<Button size="xl" emphasis="solid">Extra large</Button>
	</div>

	<h3>Disabled</h3>

	<div style="display: flex; gap: 10px">
		<Button emphasis="solid" disabled>Solid disabled</Button>
		<Button disabled>Outline disabled</Button>
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

		<Card tone="success">
			<div style="padding: 10px;">Card content</div>
		</Card>

		<Card tone="danger">
			<div style="padding: 10px;">Card content</div>
		</Card>

		<Card surface="chrome">
			<div style="padding: 10px;">Chrome card</div>
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
	<div class="icons-grid">
		<div class="icon-preview">
			<h2>Plus</h2>
			<Plus />
		</div>
		<div class="icon-preview">
			<h2>Copy</h2>
			<Copy />
		</div>
		<div class="icon-preview">
			<h2>Trash</h2>
			<Trash />
		</div>
		<div class="icon-preview">
			<h2>Upload</h2>
			<Upload />
		</div>
		<div class="icon-preview">
			<h2>Download</h2>
			<Download />
		</div>
		<div class="icon-preview">
			<h2>Edit</h2>
			<Edit />
		</div>
		<div class="icon-preview">
			<h2>X</h2>
			<X animated />
		</div>
		<div class="icon-preview">
			<h2>Check</h2>
			<Check />
		</div>
		<div class="icon-preview">
			<h2>Warning</h2>
			<Warning />
		</div>
		<div class="icon-preview">
			<h2>Info</h2>
			<Info />
		</div>
		<div class="icon-preview">
			<h2>Sun</h2>
			<Sun />
		</div>
		<div class="icon-preview">
			<h2>Moon</h2>
			<Moon />
		</div>
		<div class="icon-preview">
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

	.icons-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
		gap: 16px;
		max-width: 920px;
	}

	.icon-preview {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 8px;
		min-height: 150px;
		border: 1px solid var(--color-border);
		border-radius: var(--border-radius-regular);
	}

	.icon-preview h2 {
		padding-top: 0;
		font-size: 14px;
		font-weight: 500;
	}
</style>
