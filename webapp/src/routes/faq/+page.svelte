<script lang="ts">
	import { resolve } from '$app/paths';
	import { Card, Link } from '$lib/components';
	import { useThemeController } from '$lib/theme/controller.svelte';

	const theme = useThemeController();
</script>

<svelte:head>
	<title>FAQ | Хостамат</title>
</svelte:head>

<main class="page-shell">
	<div class={['page-content', 'faq', theme.current === 'dark' && 'theme-dark']}>
		<h1>FAQ</h1>

		<Card border>
			<section class="qa">
				<h2>После загрузки сайта возвращается 404</h2>
				<p>
					Чаще всего файл <span class="path">index.html</span> отсутствует или лежит не в корне
					архива, а внутри вложенной папки. Убедитесь, что архив содержит
					<span class="path">index.html</span>.
				</p>
				<p>Правила поиска страниц:</p>
				<ul>
					<li>
						Для <span class="path">/</span> нужен файл <span class="path">index.html</span> в корне архива
						или загруженный отдельно.
					</li>
					<li>
						Для <span class="path">/aaa</span> сервис проверяет
						<span class="path">/aaa</span>, <span class="path">/aaa.html</span> и
						<span class="path">/aaa/index.html</span>.
					</li>
				</ul>
			</section>
		</Card>

		<Card border>
			<section class="qa">
				<h2>Какой максимальный размер файла?</h2>
				<p>До 20 МБ на архив или отдельный файл.</p>
			</section>
		</Card>

		<Card border>
			<section class="qa">
				<h2>Как долго хранится сайт?</h2>
				<p>Бессрочно, пока вы не удалите его вручную.</p>
			</section>
		</Card>

		<Card border>
			<section class="qa">
				<h2>Можно обновить или удалить сайт?</h2>
				<p>
					Да, все действия доступны в <Link href={resolve('/dashboard/sites')}>Дашборде</Link>.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="qa">
				<h2>Какие форматы файлов поддерживаются?</h2>
				<p>
					Список форматов есть на странице <Link href={resolve('/about')}>«О сервисе»</Link>.
					Поддерживаются только архивы в формате <span class="path">.zip</span>.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="qa">
				<h2>Что означает «альфа-тестирование»?</h2>
				<p>
					Сервис активно развивается: некоторые функции могут работать нестабильно. Нам важна
					обратная связь через <Link
						href="https://forms.yandex.ru/u/6907bbfc49af47f3cae7650a"
						target="_blank">форму</Link
					>.
				</p>
			</section>
		</Card>
	</div>
</main>

<style>
	.faq {
		gap: 0.8rem;
		padding-bottom: 2rem;
	}

	.qa {
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
	}

	.qa h2 {
		font-size: clamp(1.15rem, 1.2vw, 1.35rem);
		line-height: 1.24;
	}

	.qa ul {
		padding-left: 1.2rem;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	.path {
		font-family:
			ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
			monospace;
		display: inline-block;
		font-size: 0.93em;
		font-weight: 600;
		line-height: 1.3;
		color: var(--color-primary-95);
		background: var(--color-accent-surface);
		border-radius: 0;
		padding: 0.12em 0.34em;
		margin: 0.14em 0;
		vertical-align: baseline;
		white-space: normal;
		overflow-wrap: anywhere;
		word-break: break-word;
		box-decoration-break: clone;
		-webkit-box-decoration-break: clone;
		transition: background-color 120ms ease;
	}

	.path:hover {
		background: var(--color-background-object);
	}

	.faq.theme-dark .path {
		background: var(--color-accent-surface);
	}

	.faq.theme-dark .path:hover {
		background: var(--color-background-object);
	}
</style>
