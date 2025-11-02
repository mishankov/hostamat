<script lang="ts">
	type InputProps = {
		id: string;
		type?: string;
		label?: string;
		value?: string;
		ref?: any;
		oninput?: (e: Event) => void;
	};

	let {
		id,
		type = 'text',
		label = '',
		value = $bindable(''),
		ref = $bindable(null),
		oninput = undefined
	}: InputProps = $props();
</script>

<div class="input-container">
	<!-- Empty placeholder to utilize :placeholder-shown pseudo class -->
	<input {id} {type} bind:value bind:this={ref} placeholder=" " />
	{#if label}
		<label for={id}>{label}</label>
	{/if}
</div>

<style>
	.input-container {
		width: 100%;

		position: relative;
		display: inline-block;
	}

	input {
		width: 100%;

		padding: 12px 20px;

		background-color: var(--color-primary-dim);
		color: var(--color-primary);

		border-width: 1px;
		border-style: solid;
		border-radius: 15px;
		border-color: var(--color-border);

		transition: var(--transition);
	}

	input:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	label {
		position: absolute;
		left: 10px;
		top: 50%;
		transform: translateY(-50%);

		padding: 0 4px;

		background-color: transparent;
		color: #555;

		pointer-events: none;

		transition: var(--transition);
	}

	input:focus + label,
	input:not(:placeholder-shown) + label {
		top: 0;
		transform: translateY(-50%) translateX(-10%) scale(0.85);
		background-color: var(--color-primary-dim);
		color: var(--color-primary);
	}
</style>
