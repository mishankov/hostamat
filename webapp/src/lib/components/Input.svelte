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

<div class="input-container" class:focused class:has-value={value.length > 0} class:invalid={invalid || errorMessage}>
	{#if label}
		<label for={id}>{label}</label>
	{/if}
	<div class="input-wrapper">
		<input
			{id}
			class={[(invalid || errorMessage) && 'invalid']}
			{type}
			bind:value
			bind:this={ref}
			{placeholder}
			onfocus={() => focused = true}
			onblur={() => focused = false}
			{...restProps}
		/>
		<div class="input-glow"></div>
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
		gap: 6px;
	}

	label {
		font-size: 13px;
		font-weight: 500;
		color: var(--color-text-secondary);
		transition: var(--transition);
		user-select: none;
	}

	.input-container.focused label,
	.input-container.has-value label {
		color: var(--color-accent-60);
	}

	.input-container.invalid label {
		color: var(--color-danger-60);
	}

	.input-wrapper {
		position: relative;
	}

	input {
		width: 100%;
		padding: 12px 16px;

		font-family: 'Outfit', sans-serif;
		font-size: 15px;

		background: var(--color-background-elevated);
		color: var(--color-text-regular);

		border: 1px solid var(--color-border);
		border-radius: var(--border-radius-regular);

		transition: var(--transition);
	}

	input:hover {
		border-color: var(--color-primary-30);
	}

	input:focus {
		outline: none;
		border-color: var(--color-accent-50);
		background: var(--color-background);
	}

	/* Glow effect on focus */
	.input-glow {
		position: absolute;
		inset: -1px;
		border-radius: calc(var(--border-radius-regular) + 1px);
		background: transparent;
		opacity: 0;
		transition: var(--transition);
		pointer-events: none;
		z-index: -1;
	}

	.input-container.focused .input-glow {
		opacity: 1;
		box-shadow: var(--glow-accent);
	}

	input::placeholder {
		color: var(--color-inactive-40);
	}

	input.invalid {
		border-color: var(--color-danger-50);
		background: var(--color-danger-5);
	}

	.input-container.invalid .input-glow {
		box-shadow: var(--glow-danger);
	}

	.input-container.invalid.focused .input-glow {
		opacity: 1;
	}

	span.error-message {
		color: var(--color-danger-60);
		font-size: 13px;
		display: flex;
		align-items: center;
		gap: 4px;
	}

	span.error-message::before {
		content: '!';
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 14px;
		height: 14px;
		font-size: 10px;
		font-weight: 700;
		background: var(--color-danger-50);
		color: white;
		border-radius: 50%;
	}
</style>
