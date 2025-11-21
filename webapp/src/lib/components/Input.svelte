<script lang="ts">
	type InputProps = {
		id: string;
		type?: string;
		label?: string;
		placeholder?: string;
		value?: string;
		ref?: any;
	};

	let {
		id,
		type = 'text',
		label = '',
		placeholder = ' ',
		value = $bindable(''),
		ref = $bindable(null)
	}: InputProps = $props();
</script>

<div class="input-container">
	{#if label}
		<label for={id}>{label}</label>
	{/if}
	<input {id} {type} bind:value bind:this={ref} {placeholder} />
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
		color: var(--color-primary-90);
		pointer-events: none;
		transition: var(--transition);
		user-select: none;
	}

	input {
		width: 100%;

		padding: 12px 10px;

		background-color: var(--color-primary-7);
		color: var(--color-primary);

		border-width: 1px;
		border-style: solid;
		border-radius: var(--border-radius-regular);
		border-color: var(--color-primary-30);

		transition: var(--transition);
	}

	input:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	input::placeholder {
		color: var(--color-primary-30);
	}
</style>
