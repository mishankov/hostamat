<script lang="ts">
	let {
		id,
		invalid = false,
		errorMessage = '',
		type = 'text',
		label = '',
		placeholder = ' ',
		value = $bindable(''),
		ref = $bindable(null),
		readonly = false,
		tabindex = undefined,
		...restProps
	} = $props();
</script>

<div class="input-container">
	{#if label}
		<label for={id}>{label}</label>
	{/if}
	<input
		{id}
		class={[(invalid || errorMessage) && 'invalid', readonly && 'readonly']}
		{type}
		bind:value
		bind:this={ref}
		{placeholder}
		{readonly}
		tabindex={readonly ? -1 : tabindex}
		aria-readonly={readonly ? 'true' : undefined}
		{...restProps}
	/>
	{#if errorMessage}
		<span class="error-message">{errorMessage}</span>
	{/if}
</div>

<style>
	.input-container {
		width: 100%;

		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}

	label {
		background-color: transparent;
		font-size: 0.8125rem;
		font-weight: 600;
		color: var(--color-text-secondary);
		pointer-events: none;
		transition: var(--transition);
		user-select: none;
	}

	input {
		width: 100%;

		padding: 0.65rem 0.9rem;

		font-size: 0.95rem;

		background-color: var(--color-background-object);
		color: var(--color-text-regular);

		border-width: 1px;
		border-style: solid;
		border-radius: var(--radius-md);
		border-color: var(--color-border);

		transition:
			border-color 180ms var(--ease-out-quart),
			box-shadow 180ms var(--ease-out-quart),
			background-color 180ms var(--ease-out-quart);
	}

	input:focus {
		outline: none;
		border-color: var(--color-accent-border);
		box-shadow:
			var(--focus-ring),
			0 0 0 1px var(--color-accent-border);
	}

	input.readonly {
		background: var(--color-neutral-surface);
		border-color: color-mix(in hsl, var(--color-border) 74%, var(--color-disabled-border));
		box-shadow: inset 0 0 0 1px color-mix(in hsl, var(--color-border) 55%, transparent);
		color: var(--color-text-secondary);
		caret-color: transparent;
		cursor: default;
	}

	input.readonly:focus {
		border-color: color-mix(in hsl, var(--color-border) 74%, var(--color-disabled-border));
		box-shadow: inset 0 0 0 1px color-mix(in hsl, var(--color-border) 55%, transparent);
	}

	input::placeholder {
		color: var(--color-disabled-text);
	}

	input.invalid {
		border-color: var(--color-danger-strong);
		background: var(--color-danger-surface);
		box-shadow: 0 0 0 1px color-mix(in hsl, var(--color-danger-strong) 24%, transparent);
	}

	input.invalid:focus {
		border-color: var(--color-danger-strong);
		box-shadow:
			0 0 0 3px color-mix(in hsl, var(--color-danger-strong) 20%, transparent),
			0 0 0 1px var(--color-danger-strong);
	}

	span.error-message {
		color: var(--color-danger-text);
		font-size: 0.8125rem;
		font-weight: 500;
	}
</style>
