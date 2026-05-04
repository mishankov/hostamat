<script lang="ts">
	import { page } from '$app/state';
	import { createFeedback } from '$lib/api-client';
	import { Alert, Button, Input, Modal } from '$lib/components';
	import { Check, I, Spinner } from '$lib/icons';

	const maxMessageLength = 5000;

	let { email = '' }: { email?: string } = $props();

	let modal = $state<Modal>();
	let message = $state('');
	let errorMessage = $state('');
	let submitError = $state('');
	let sent = $state(false);
	let loading = $state(false);

	const currentUrl = $derived(`${page.url.pathname}${page.url.search}`);
	const remainingCharacters = $derived(maxMessageLength - [...message].length);
	const emailUnavailable = $derived(email.trim() === '');

	const openModal = () => {
		errorMessage = '';
		submitError = '';
		sent = false;
		modal?.open();
	};

	const validate = () => {
		const trimmed = message.trim();
		if (!trimmed) {
			errorMessage = 'Напишите сообщение';
			return false;
		}

		if ([...trimmed].length > maxMessageLength) {
			errorMessage = 'Сообщение должно быть не длиннее 5000 символов';
			return false;
		}

		if (emailUnavailable) {
			errorMessage = 'В аккаунте не найден email для ответа';
			return false;
		}

		errorMessage = '';
		return true;
	};

	const send = async () => {
		submitError = '';
		sent = false;

		if (!validate()) {
			return;
		}

		let error: string | null = null;

		try {
			loading = true;
			const result = await createFeedback(message.trim(), currentUrl);
			error = result.error;
		} catch {
			error = 'Не удалось отправить сообщение';
		} finally {
			loading = false;
		}

		if (error) {
			submitError = error;
			return;
		}

		message = '';
		sent = true;
	};
</script>

<div class="feedback-widget">
	<button
		type="button"
		class="feedback-trigger"
		aria-label="Открыть форму обратной связи"
		onclick={openModal}
	>
		<I size={19} color="currentColor" />
		<span>Обратная связь</span>
	</button>
</div>

<Modal
	bind:this={modal}
	title="Обратная связь"
	panelWidth="min(32rem, calc(100vw - 1rem))"
	titleId="feedback-title"
	descriptionId="feedback-description"
	initialFocusSelector="#feedback-message"
>
	<form
		class="feedback-form"
		onsubmit={(event) => {
			event.preventDefault();
			send();
		}}
	>
		<p id="feedback-description" class="intro">
			Расскажите, что сломалось, чего не хватает или где интерфейс сбивает с толку. Или что
			понравилось
		</p>

		<Input
			id="feedback-email"
			label="Email для ответа"
			value={email}
			readonly
			invalid={emailUnavailable}
			errorMessage={emailUnavailable ? 'Email не найден в аккаунте' : ''}
		/>

		<div class="textarea-field">
			<div class="textarea-label-row">
				<label for="feedback-message">Сообщение</label>
				<span class={['counter', remainingCharacters < 0 && 'over']}>{remainingCharacters}</span>
			</div>
			<textarea
				id="feedback-message"
				bind:value={message}
				class={errorMessage ? 'invalid' : ''}
				rows="7"
				maxlength={maxMessageLength + 200}
				placeholder="Например: не понял, как подключить домен..."
				aria-invalid={errorMessage ? 'true' : undefined}
				aria-describedby={errorMessage ? 'feedback-message-error' : undefined}
				disabled={loading}
				oninput={() => {
					errorMessage = '';
					submitError = '';
				}}
			></textarea>
			{#if errorMessage}
				<span id="feedback-message-error" class="field-error">{errorMessage}</span>
			{/if}
		</div>

		{#if submitError}
			<Alert variant="error" message={submitError} />
		{/if}

		{#if sent}
			<Alert variant="success" message="Сообщение отправлено! Мы ответим в ближайшее время" />
		{/if}

		<div class="actions">
			<Button type="submit" emphasis="solid" disabled={loading || emailUnavailable}>
				{#if loading}
					<Spinner size={17} color="currentColor" />
					Отправляем
				{:else if sent}
					<Check size={17} color="currentColor" />
					Отправлено
				{:else}
					Отправить
				{/if}
			</Button>
			<Button disabled={loading} onclick={() => modal?.close()}>Закрыть</Button>
		</div>
	</form>
</Modal>

<style>
	.feedback-widget {
		width: 100%;
	}

	.feedback-trigger {
		display: inline-flex;
		align-items: center;
		justify-content: flex-start;
		gap: var(--space-sm);
		width: 100%;
		min-height: 2.7rem;
		padding: 0.68rem 0.85rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-background-object);
		color: var(--color-text-regular);
		box-shadow: none;
		font-size: 0.92rem;
		font-weight: 700;
		line-height: 1.2;
		cursor: pointer;
		transition:
			background-color 180ms var(--ease-out-quart),
			border-color 180ms var(--ease-out-quart),
			box-shadow 180ms var(--ease-out-quart),
			color 180ms var(--ease-out-quart),
			transform 120ms var(--ease-out-quart);
	}

	.feedback-trigger:hover {
		transform: translateX(2px);
		background: var(--color-neutral-surface);
		border-color: var(--color-border);
		box-shadow: var(--shadow-sm);
	}

	.feedback-trigger:active {
		transform: translateX(1px) scale(0.99);
	}

	.feedback-trigger:focus-visible {
		outline: none;
		box-shadow: var(--focus-ring);
	}

	.feedback-trigger :global(svg) {
		flex: 0 0 auto;
		color: var(--color-text-secondary);
	}

	.feedback-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.intro {
		margin: 0;
		color: var(--color-text-secondary);
		font-size: 0.95rem;
		line-height: 1.55;
	}

	.textarea-field {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}

	.textarea-label-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-md);
	}

	label {
		font-size: 0.8125rem;
		font-weight: 600;
		color: var(--color-text-secondary);
		user-select: none;
	}

	.counter {
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--color-text-tertiary);
		font-variant-numeric: tabular-nums;
	}

	.counter.over {
		color: var(--color-danger-text);
	}

	textarea {
		width: 100%;
		min-height: 11rem;
		resize: vertical;
		padding: 0.75rem 0.9rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-background-object);
		color: var(--color-text-regular);
		font: inherit;
		font-size: 0.95rem;
		line-height: 1.45;
		transition:
			border-color 180ms var(--ease-out-quart),
			box-shadow 180ms var(--ease-out-quart),
			background-color 180ms var(--ease-out-quart);
	}

	textarea:focus {
		outline: none;
		border-color: var(--color-accent-border);
		box-shadow:
			var(--focus-ring),
			0 0 0 1px var(--color-accent-border);
	}

	textarea::placeholder {
		color: var(--color-disabled-text);
	}

	textarea.invalid {
		border-color: var(--color-danger-strong);
		background: var(--color-danger-surface);
		box-shadow: 0 0 0 1px color-mix(in hsl, var(--color-danger-strong) 24%, transparent);
	}

	.field-error {
		color: var(--color-danger-text);
		font-size: 0.8125rem;
		font-weight: 500;
	}

	.actions {
		display: flex;
		justify-content: flex-start;
		gap: var(--space-sm);
	}

	@media (max-width: 700px) {
		.actions {
			flex-direction: column;
			align-items: stretch;
		}

		.actions :global(.button) {
			width: 100%;
		}
	}
</style>
