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
		gap: 6px;
	}

	label {
		background-color: transparent;
		font-size: 13px;
		font-weight: 600;
		color: var(--color-text-secondary);
		pointer-events: none;
		transition: var(--transition);
		user-select: none;
	}

	input {
		width: 100%;

		padding: 0.74rem 0.82rem;

		font-size: 16px;

		background-color: var(--color-background-object);
		color: var(--color-text-regular);

		border-width: 1px;
		border-style: solid;
		border-radius: calc(var(--border-radius-regular) - 2px);
		border-color: color-mix(in hsl, var(--color-border) 86%, var(--color-primary-20));

		transition: var(--transition);
	}

	input:focus {
		outline: none;
		border-color: var(--color-accent-70);
		box-shadow: var(--focus-ring);
	}

	input::placeholder {
		color: var(--color-inactive-50);
	}

	input.invalid {
		border-color: var(--color-danger-50);
		background: var(--color-danger-10);
	}

	span.error-message {
		color: var(--color-danger-50);
		font-size: 13px;
	}
</style>
