<script lang="ts">
	import { resolve } from '$app/paths';
	import { Link } from '$lib/components';
	import tutorialMarkdown from '$lib/api-docs/tutorial.md?raw';

	const referenceHref = resolve('/api-docs');
	const openapiHref = resolve('/api-docs/openapi.yaml');

	const escapeHTML = (value: string) =>
		value
			.replaceAll('&', '&amp;')
			.replaceAll('<', '&lt;')
			.replaceAll('>', '&gt;')
			.replaceAll('"', '&quot;');

	const renderInline = (value: string) =>
		escapeHTML(value).replace(/`([^`]+)`/g, '<code>$1</code>');

	const closeList = (chunks: string[], listType: 'ol' | 'ul' | null) => {
		if (listType) {
			chunks.push(`</${listType}>`);
		}

		return null;
	};

	const renderMarkdown = (markdown: string) => {
		const chunks: string[] = [];
		const lines = markdown.split('\n');
		let inCodeBlock = false;
		let listType: 'ol' | 'ul' | null = null;
		let codeLanguage = '';

		for (const line of lines) {
			if (line.startsWith('```')) {
				if (inCodeBlock) {
					chunks.push('</code></pre>');
					inCodeBlock = false;
					continue;
				}

				listType = closeList(chunks, listType);
				codeLanguage = line.slice(3).trim();
				chunks.push(`<pre><code${codeLanguage ? ` class="language-${escapeHTML(codeLanguage)}"` : ''}>`);
				inCodeBlock = true;
				continue;
			}

			if (inCodeBlock) {
				chunks.push(`${escapeHTML(line)}\n`);
				continue;
			}

			if (!line.trim()) {
				listType = closeList(chunks, listType);
				continue;
			}

			const heading = /^(#{1,3})\s+(.+)$/.exec(line);
			if (heading) {
				listType = closeList(chunks, listType);
				const level = heading[1].length;
				chunks.push(`<h${level}>${renderInline(heading[2])}</h${level}>`);
				continue;
			}

			const ordered = /^\d+\.\s+(.+)$/.exec(line);
			if (ordered) {
				if (listType !== 'ol') {
					listType = closeList(chunks, listType);
					chunks.push('<ol>');
					listType = 'ol';
				}
				chunks.push(`<li>${renderInline(ordered[1])}</li>`);
				continue;
			}

			const unordered = /^\s*-\s+(.+)$/.exec(line);
			if (unordered) {
				if (listType !== 'ul') {
					listType = closeList(chunks, listType);
					chunks.push('<ul>');
					listType = 'ul';
				}
				chunks.push(`<li>${renderInline(unordered[1])}</li>`);
				continue;
			}

			listType = closeList(chunks, listType);
			chunks.push(`<p>${renderInline(line)}</p>`);
		}

		closeList(chunks, listType);
		return chunks.join('\n');
	};

	const tutorialHtml = renderMarkdown(tutorialMarkdown);
</script>

<svelte:head>
	<title>Быстрый старт Public API | Хостамат</title>
	<meta name="description" content="Руководство по использованию Public API Хостамата." />
</svelte:head>

<main class="tutorial-shell">
	<header class="tutorial-header">
		<div>
			<span class="eyebrow">Public API</span>
			<h1>Быстрый старт</h1>
			<p>Создание ключа, публикация zip-архива и подключение пользовательского домена.</p>
		</div>
		<nav aria-label="Навигация по API документации">
			<Link href={referenceHref}>Интерактивная документация</Link>
			<Link href={openapiHref}>OpenAPI YAML</Link>
		</nav>
	</header>

	<article class="tutorial-content">
		{@html tutorialHtml}
	</article>
</main>

<style>
	.tutorial-shell {
		min-height: calc(100vh - 4rem);
		background: var(--color-background);
	}

	.tutorial-header {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: var(--space-lg);
		padding: var(--space-xl);
		border-bottom: 1px solid var(--color-border);
		background: var(--color-background-object);
	}

	.tutorial-header h1 {
		margin: 0;
		font-size: clamp(1.8rem, 4vw, 2.8rem);
		font-weight: 800;
	}

	.tutorial-header p {
		max-width: 60ch;
		color: var(--color-text-secondary);
	}

	.eyebrow {
		color: var(--color-text-tertiary);
		font-family: var(--font-mono);
		font-size: 0.75rem;
		text-transform: uppercase;
	}

	.tutorial-header nav {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		flex-wrap: wrap;
	}

	.tutorial-content {
		width: min(52rem, calc(100vw - 2rem));
		margin: 0 auto;
		padding: var(--space-2xl) 0;
		color: var(--color-text-regular);
	}

	.tutorial-content :global(h1) {
		display: none;
	}

	.tutorial-content :global(h2) {
		margin: var(--space-2xl) 0 var(--space-md);
		font-size: 1.45rem;
	}

	.tutorial-content :global(h3) {
		margin: var(--space-xl) 0 var(--space-sm);
		font-size: 1.05rem;
	}

	.tutorial-content :global(p),
	.tutorial-content :global(li) {
		color: var(--color-text-secondary);
		line-height: 1.7;
	}

	.tutorial-content :global(ul),
	.tutorial-content :global(ol) {
		padding-left: 1.35rem;
	}

	.tutorial-content :global(code) {
		font-family: var(--font-mono);
		font-size: 0.84em;
	}

	.tutorial-content :global(p code),
	.tutorial-content :global(li code) {
		padding: 0.1rem 0.32rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		background: var(--color-neutral-surface);
		color: var(--color-text-regular);
	}

	.tutorial-content :global(pre) {
		overflow: auto;
		padding: var(--space-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-neutral-surface);
	}

	.tutorial-content :global(pre code) {
		color: var(--color-text-regular);
		white-space: pre;
	}

	@media (max-width: 760px) {
		.tutorial-header {
			align-items: flex-start;
			flex-direction: column;
			padding: var(--space-lg);
		}
	}
</style>
