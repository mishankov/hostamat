<script lang="ts">
	import { Card, Modal, Input, Button } from '$lib/components';
	import { Spinner } from '$lib/icons';
	import { createSite } from '$lib/api-client';

	let { loading = $bindable(false), onsuccess = async () => {} } = $props();

	let modal = $state<Modal>();
	let newSiteName = $state('');
	let newSiteNameError = $state('');

	const resetState = () => {
		newSiteName = 'Новый сайт';
		newSiteNameError = '';
	};

	resetState();
</script>

<div class="add-site-card-shell">
	<Card class="add-site-card" style="height: 100%;">
		<button
			type="button"
			class="add-site"
			onclick={() => {
				resetState();
				modal?.open();
			}}
		>
			{#if loading}
				<div class="prompt-line">
					<span class="prompt-symbol">+</span>
					<span class="prompt-text">добавляем сайт</span>
					<Spinner size={16} />
				</div>
			{:else}
				<div class="prompt-line">
					<span class="prompt-symbol">+</span>
					<span class="prompt-text">добавить сайт</span>
				</div>
			{/if}
		</button>
	</Card>
</div>

<Modal
	bind:this={modal}
	title="Новый сайт"
	titleId="new-site-modal-title"
	initialFocusSelector="#new-site-name"
	closeLabel="Закрыть форму создания сайта"
>
	<div class="modal-body">
		<Input
			id="new-site-name"
			label="Название"
			errorMessage={newSiteNameError}
			bind:value={newSiteName}
			maxlength="25"
		></Input>
		<div class="modal-actions">
			<Button
				emphasis="solid"
				disabled={loading}
				onclick={async () => {
					if (newSiteName.length === 0) {
						newSiteNameError = 'Название не должно быть пустым';
						return;
					}

					loading = true;
					const { data, error } = await createSite(newSiteName);

					if (error) {
						newSiteNameError = error;
						loading = false;
					}

					if (data) {
						await onsuccess();
						modal?.close();
					}
				}}>Создать</Button
			>
			<Button
				disabled={loading}
				onclick={() => {
					modal?.close();
				}}>Отмена</Button
			>
			{#if loading}
				<span class="action-spinner"><Spinner size={30} /></span>
			{/if}
		</div>
	</div>
</Modal>

<style>
	.add-site-card-shell {
		display: flex;
		min-height: 10.625rem;
	}

	.add-site {
		all: unset;

		width: 100%;
		height: 100%;
		min-height: 10.625rem;
		box-sizing: border-box;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-md);

		padding: var(--space-lg) 1.3rem;
		cursor: pointer;
		border-radius: var(--radius-lg);
	}

	.prompt-line {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		font-family: var(--font-body);
		font-size: 0.9rem;
		font-weight: 500;
	}

	.prompt-symbol {
		color: var(--color-accent-strong);
		font-weight: 700;
		font-size: 1.1rem;
	}

	.prompt-text {
		color: var(--color-text-secondary);
	}

	:global(.add-site-card) {
		border-style: dashed;
		border-color: var(--color-border);
		border-radius: var(--radius-lg);
		background: transparent;
		transition:
			border-color 180ms ease-out,
			background 180ms ease-out,
			box-shadow 180ms ease-out,
			transform 180ms ease-out;
	}

	:global(.add-site-card:hover) {
		border-color: var(--color-accent-border);
		background: var(--color-accent-surface);
		box-shadow: var(--shadow-sm);
		transform: translateY(-2px);
	}

	.modal-body {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 16.875rem;
	}

	.modal-actions {
		display: flex;
		gap: var(--space-lg);
		align-items: center;
	}

	.action-spinner {
		display: inline-flex;
		align-items: center;
	}
</style>
