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
		...restProps
	} = $props();
</script>

<div class="input-container">
	{#if label}
		<label for={id}>{label}</label>
	{/if}
	<input
		{id}
		class={[(invalid || errorMessage) && 'invalid']}
		{type}
		bind:value
		bind:this={ref}
		{placeholder}
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
		gap: 8px;
	}

	label {
		background-color: transparent;
		font-size: 13px;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.03em;
		color: var(--color-text-regular);
		transition: var(--transition);
		user-select: none;
	}

	input {
		width: 100%;
		padding: 14px 16px;
		font-size: 15px;
		font-family: 'JetBrains Mono', monospace;
		font-weight: 500;

		background-color: var(--color-primary-5);
		color: var(--color-text-regular);

		border: var(--border-regular);
		border-radius: var(--border-radius-regular);
		border-color: var(--color-primary-20);

		transition: var(--transition);
		box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	input:focus {
		outline: none;
		border-color: var(--color-accent-50);
		background-color: var(--color-primary-7);
		box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.05), 0 0 0 3px var(--color-accent-10);
	}

	input::placeholder {
		color: var(--color-inactive-40);
		font-style: italic;
	}

	input.invalid {
		border-color: var(--color-danger-50);
		background: var(--color-danger-5);
	}

	input.invalid:focus {
		box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.05), 0 0 0 3px var(--color-danger-10);
	}

	span.error-message {
		color: var(--color-danger-60);
		font-size: 12px;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.02em;
	}
</style>
