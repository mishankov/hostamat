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

	let focused = $state(false);
</script>

<div class="input-container" class:focused class:has-value={value?.length > 0}>
	{#if label}
		<label for={id}>{label}</label>
	{/if}
	<div class="input-wrapper" class:invalid={invalid || errorMessage}>
		<input
			{id}
			class:invalid={invalid || errorMessage}
			{type}
			bind:value
			bind:this={ref}
			{placeholder}
			onfocus={() => focused = true}
			onblur={() => focused = false}
			{...restProps}
		/>
		<div class="input-border"></div>
	</div>
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
		font-family: var(--font-mono);
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		letter-spacing: 0.02em;
		text-transform: uppercase;
		transition: color var(--transition-fast);
		user-select: none;
	}

	.input-container.focused label {
		color: var(--color-accent-50);
	}

	.input-wrapper {
		position: relative;
		width: 100%;
	}

	input {
		width: 100%;
		padding: 14px 16px;

		font-family: var(--font-mono);
		font-size: 0.9375rem;
		font-weight: 400;

		background-color: var(--color-background-object);
		color: var(--color-text-primary);

		border: 1px solid var(--color-border);
		border-radius: var(--border-radius-regular);

		transition: all var(--transition-fast);
	}

	/* Animated border highlight */
	.input-border {
		position: absolute;
		inset: 0;
		border-radius: var(--border-radius-regular);
		pointer-events: none;

		background: linear-gradient(135deg, var(--color-accent-50), var(--color-accent-60));
		opacity: 0;
		transition: opacity var(--transition-fast);

		mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
		mask-composite: exclude;
		-webkit-mask-composite: xor;
		padding: 2px;
	}

	input:focus {
		outline: none;
		background-color: var(--color-background-elevated);
		border-color: transparent;
	}

	input:focus + .input-border {
		opacity: 1;
	}

	input:hover:not(:focus) {
		border-color: var(--color-border-strong);
	}

	input::placeholder {
		color: var(--color-text-muted);
		opacity: 0.6;
	}

	/* Invalid state */
	.input-wrapper.invalid input {
		border-color: var(--color-danger-40);
		background-color: var(--color-danger-5);
	}

	.input-wrapper.invalid .input-border {
		background: linear-gradient(135deg, var(--color-danger-50), var(--color-danger-60));
	}

	.input-wrapper.invalid input:focus {
		border-color: transparent;
	}

	span.error-message {
		display: flex;
		align-items: center;
		gap: 6px;
		color: var(--color-danger-60);
		font-size: 0.8125rem;
		font-weight: 500;
	}

	span.error-message::before {
		content: '!';
		display: flex;
		align-items: center;
		justify-content: center;
		width: 16px;
		height: 16px;
		font-size: 0.6875rem;
		font-weight: 700;
		background: var(--color-danger-50);
		color: var(--color-danger-95);
		border-radius: 50%;
	}
</style>
