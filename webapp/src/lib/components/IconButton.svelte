<script lang="ts">
	import type { Snippet } from 'svelte';

	type IconButtonProps = {
		label: string;
		children?: Snippet;
		onclick?: (event: MouseEvent) => void;
		class?: string;
		tone?: 'neutral' | 'danger';
		kind?: 'primary' | 'danger';
		grouped?: boolean;
		title?: string;
		type?: 'button' | 'submit' | 'reset';
		disabled?: boolean;
		[key: string]: unknown;
	};

	let {
		label,
		children,
		onclick,
		class: className = '',
		tone = undefined,
		kind = undefined,
		grouped = false,
		title = '',
		type = 'button',
		disabled = false,
		...restProps
	}: IconButtonProps = $props();

	const resolvedTone = $derived<'neutral' | 'danger'>(
		tone || (kind === 'danger' ? 'danger' : 'neutral')
	);
</script>

<button
	{type}
	{disabled}
	class={[resolvedTone, grouped && 'grouped', className]}
	aria-label={label}
	title={title || label}
	{onclick}
	{...restProps}
>
	{@render children?.()}
</button>

<style>
	button {
		all: unset;
		display: flex;
		align-items: center;
		justify-content: center;

		padding: 6px;
		border-radius: 0;
		border: var(--border-regular);
		background: var(--color-background-object);
		color: var(--color-text-regular);

		transition:
			color 260ms ease-out,
			background-color 260ms ease-out,
			border-color 260ms ease-out;
	}

	button.grouped {
		border: none;
		border-radius: 0;
	}

	button:hover:enabled {
		transition-duration: 90ms;
		background: var(--color-accent-surface);
		cursor: pointer;
		transform: none;
	}

	button.danger:hover:enabled {
		transition-duration: 90ms;
		background: var(--color-danger-surface);
		border-color: var(--color-danger-border);
		color: var(--color-danger-text);
		cursor: pointer;
	}

	button:disabled {
		background-color: var(--color-disabled-surface);
		border-color: var(--color-disabled-border);
		color: var(--color-disabled-text);
		cursor: not-allowed;
	}
</style>
