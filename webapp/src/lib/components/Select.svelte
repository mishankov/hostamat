<script lang="ts">
	let {
		id,
		label = '',
		options,
		placeholder = '',
		errorMessage = '',
		value = $bindable(''),
		...restProps
	}: {
		id: string;
		label?: string;
		options: { value: string; label: string }[];
		placeholder?: string;
		errorMessage?: string;
		value?: string;
		[key: string]: unknown;
	} = $props();
</script>

<div class="select-container">
	{#if label}
		<label for={id}>{label}</label>
	{/if}
	<select {id} class={[errorMessage && 'invalid']} bind:value {...restProps}>
		{#if placeholder}
			<option value="" disabled selected>{placeholder}</option>
		{/if}
		{#each options as option (option.value)}
			<option value={option.value}>{option.label}</option>
		{/each}
	</select>
	{#if errorMessage}
		<span class="error-message">{errorMessage}</span>
	{/if}
</div>

<style>
	.select-container {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}

	label {
		font-size: 0.8125rem;
		font-weight: 600;
		color: var(--color-text-secondary);
		pointer-events: none;
		transition: var(--transition);
		user-select: none;
	}

	select {
		width: 100%;
		padding: 0.65rem 2.5rem 0.65rem 0.9rem;
		font-size: 0.95rem;
		font-family: var(--font-body);
		background-color: var(--color-background-object);
		color: var(--color-text-regular);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		transition:
			border-color 180ms var(--ease-out-quart),
			box-shadow 180ms var(--ease-out-quart),
			background-color 180ms var(--ease-out-quart);
		cursor: pointer;
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 0.75rem center;
		background-size: 1rem;
	}

	select:focus {
		outline: none;
		border-color: var(--color-accent-border);
		box-shadow: var(--focus-ring), 0 0 0 1px var(--color-accent-border);
	}

	select:has(option[value='']:checked) {
		color: var(--color-disabled-text);
	}

	select.invalid {
		border-color: var(--color-danger-strong);
		background: var(--color-danger-surface);
		box-shadow: 0 0 0 1px color-mix(in hsl, var(--color-danger-strong) 24%, transparent);
	}

	select.invalid:focus {
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
