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
		gap: 5px;
	}

	label {
		background-color: transparent;
		font-size: 14px;
		color: var(--color-inactive-70);
		pointer-events: none;
		transition: var(--transition);
		user-select: none;
	}

	input {
		width: 100%;

		padding: 12px 10px;

		font-size: 16px;

		background-color: var(--color-primary-7);
		color: var(--color-primary);

		border-width: 1px;
		border-style: solid;
		border-radius: var(--border-radius-regular);
		border-color: var(--color-inactive-50);

		transition: var(--transition);
	}

	input:focus {
		outline: none;
		border-color: var(--color-primary-90);
	}

	input::placeholder {
		color: var(--color-inactive-30);
	}

	input.invalid {
		border-color: var(--color-danger-50);
		background: var(--color-danger-5);
	}

	span.error-message {
		color: var(--color-danger-50);
		font-size: 14px;
	}
</style>
