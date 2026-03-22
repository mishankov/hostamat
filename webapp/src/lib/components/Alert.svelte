<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Check, I, X } from '$lib/icons';
	import Card from './Card.svelte';

	type AlertVariant = 'error' | 'info' | 'success';
	type CardTone = 'danger' | 'neutral' | 'success';

	type AlertProps = {
		variant?: AlertVariant;
		title?: string;
		message?: string;
		children?: Snippet;
	};

	let { variant = 'info', title = '', message = '', children = undefined }: AlertProps = $props();

	const cardTone = $derived.by<CardTone>(() => {
		if (variant === 'error') {
			return 'danger';
		}

		if (variant === 'success') {
			return 'success';
		}

		return 'neutral';
	});

	const iconColor = $derived.by(() => {
		if (variant === 'error') {
			return 'var(--color-danger-50)';
		}

		if (variant === 'success') {
		return 'var(--color-success-border)';
		}

		return 'var(--color-accent-border)';
	});

	const role = $derived(variant === 'error' ? 'alert' : 'status');
</script>

<div class="alert" {role}>
	<Card tone={cardTone}>
		<div class="alert-card-body">
			<div class="header">
				{#if variant === 'error'}
					<X size={25} color={iconColor} />
				{:else if variant === 'success'}
					<Check size={25} color={iconColor} />
				{:else}
					<I size={25} color={iconColor} />
				{/if}
				{#if title}
					<h4>{title}</h4>
				{:else if message}
					<p>{message}</p>
				{/if}
			</div>
			{#if children}
				<div class="content">
					{@render children?.()}
				</div>
			{:else if message && title}
				<div class="content">{message}</div>
			{/if}
		</div>
	</Card>
</div>

<style>
	.alert {
		width: 100%;
	}

	.alert-card-body {
		height: fit-content;
		padding: 10px;
	}

	.header {
		display: flex;
		gap: 10px;
		align-items: center;
	}

	.content {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 10px;
		padding-top: 10px;
	}
</style>
