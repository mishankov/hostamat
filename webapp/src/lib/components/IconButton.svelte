<script lang="ts">
	type IconButtonProps = {
		label: string;
		children?: any;
		onclick?: (event: Event) => void;
		class?: string;
		kind?: 'primary' | 'danger';
		title?: string;
		type?: 'button' | 'submit' | 'reset';
		disabled?: boolean;
		[key: string]: any;
	};

	let {
		label,
		children,
		onclick,
		class: className = '',
		kind = 'primary',
		title = '',
		type = 'button',
		disabled = false,
		...restProps
	}: IconButtonProps = $props();
</script>

<button
	{type}
	{disabled}
	class={[kind, className]}
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
		border-radius: var(--border-radius-regular);
		border: var(--border-regular);
		background: var(--color-background-object);

		transition:
			color 260ms ease-out,
			background-color 260ms ease-out,
			border-color 260ms ease-out;
	}

	button:hover:enabled {
		transition-duration: 90ms;
		background: var(--color-accent-20);
		cursor: pointer;
		transform: none;
	}

	button.danger:hover:enabled {
		transition-duration: 90ms;
		background: var(--color-danger-20);
		cursor: pointer;
	}

	button:disabled {
		background-color: var(--color-inactive-90);
		border-color: var(--color-inactive-70);
		color: var(--color-inactive-50);
		cursor: not-allowed;
	}
</style>
