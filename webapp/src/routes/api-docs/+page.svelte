<script lang="ts">
	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { Link } from '$lib/components';
	import openapiSpec from '$lib/api-docs/openapi.yaml?raw';

	const openapiHref = resolve('/api-docs/openapi.yaml');
	const tutorialHref = resolve('/api-docs/tutorial');
	const scalarScriptUrl = 'https://cdn.jsdelivr.net/npm/@scalar/api-reference';

	type ScalarGlobal = {
		createApiReference: (
			selector: string,
			configuration: Record<string, unknown>
		) => ScalarReferenceInstance;
	};

	type ScalarReferenceInstance = {
		destroy: () => void;
	};

	type ScalarWindow = Window & {
		Scalar?: ScalarGlobal;
		__hostamatScalarReference?: ScalarReferenceInstance;
	};

	const scalarWindow = () => window as ScalarWindow;
	const scalarGlobal = () => scalarWindow().Scalar;

	const loadScalarScript = () => {
		return new Promise<void>((resolveScript, rejectScript) => {
			const existingScript = document.querySelector<HTMLScriptElement>(
				`script[src="${scalarScriptUrl}"]`
			);

			if (existingScript) {
				if (scalarGlobal()) {
					resolveScript();
					return;
				}

				existingScript.addEventListener('load', () => resolveScript(), { once: true });
				existingScript.addEventListener('error', rejectScript, { once: true });
				return;
			}

			const script = document.createElement('script');
			script.src = scalarScriptUrl;
			script.async = true;
			script.addEventListener('load', () => resolveScript(), { once: true });
			script.addEventListener('error', rejectScript, { once: true });
			document.head.append(script);
		});
	};

	let scalarError = $state('');

	const clearScalarContainer = () => {
		document.getElementById('api-reference')?.replaceChildren();
	};

	const destroyPreviousScalar = () => {
		scalarWindow().__hostamatScalarReference?.destroy();
		scalarWindow().__hostamatScalarReference = undefined;
		document.dispatchEvent(new CustomEvent('scalar:destroy-references'));
		clearScalarContainer();
	};

	onMount(() => {
		if (!browser) {
			return;
		}

		let cancelled = false;

		const mountScalar = async () => {
			try {
				await loadScalarScript();
				if (cancelled) {
					return;
				}

				const scalar = scalarGlobal();
				if (!scalar) {
					throw new Error('Scalar API Reference was not initialized by the script.');
				}

				destroyPreviousScalar();
				scalarWindow().__hostamatScalarReference = scalar.createApiReference('#api-reference', {
					content: openapiSpec,
					theme: 'default',
					layout: 'modern',
					defaultHttpClient: {
						targetKey: 'shell',
						clientKey: 'curl'
					}
				});
			} catch {
				scalarError = 'Не удалось загрузить интерактивную документацию. Откройте OpenAPI YAML.';
			}
		};

		mountScalar();

		return () => {
			cancelled = true;
			destroyPreviousScalar();
		};
	});
</script>

<svelte:head>
	<title>Public API | Хостамат</title>
	<meta
		name="description"
		content="Интерактивная документация Public API Хостамат для сайтов, загрузок и доменов."
	/>
</svelte:head>

<main class="api-docs-shell">
	<div class="docs-topbar">
		<div>
			<span class="eyebrow">Public API</span>
			<h1>Документация API</h1>
		</div>
		<nav aria-label="Навигация по API документации">
			<Link href={tutorialHref}>Быстрый старт</Link>
			<Link href={openapiHref}>OpenAPI YAML</Link>
		</nav>
	</div>

	{#if scalarError}
		<div class="scalar-error" role="status">
			<p>{scalarError}</p>
			<Link href={openapiHref} alwaysUnderline>Открыть спецификацию</Link>
		</div>
	{/if}

	<div id="api-reference" class="api-reference" aria-label="Документация Public API"></div>
</main>

<style>
	.api-docs-shell {
		min-height: calc(100vh - 4rem);
		background: var(--color-background);
	}

	.docs-topbar {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: var(--space-lg);
		padding: var(--space-lg) var(--space-xl);
		border-bottom: 1px solid var(--color-border);
		background: var(--color-background-object);
	}

	.docs-topbar h1 {
		margin: 0;
		font-size: clamp(1.65rem, 3vw, 2.35rem);
		font-weight: 800;
	}

	.eyebrow {
		color: var(--color-text-tertiary);
		font-family: var(--font-mono);
		font-size: 0.75rem;
		text-transform: uppercase;
	}

	.docs-topbar nav {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		flex-wrap: wrap;
	}

	.scalar-error {
		margin: var(--space-lg) var(--space-xl) 0;
		padding: var(--space-md);
		border: 1px solid var(--color-danger-border);
		border-radius: var(--radius-md);
		background: var(--color-danger-surface);
		color: var(--color-danger-text);
	}

	.api-reference {
		min-height: calc(100vh - 9rem);
	}

	@media (max-width: 760px) {
		.docs-topbar {
			align-items: flex-start;
			flex-direction: column;
			padding: var(--space-lg);
		}
	}
</style>
