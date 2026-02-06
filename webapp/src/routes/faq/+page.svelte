<script>
	import { resolve } from '$app/paths';
	import { Card, Link } from '$lib/components';
</script>

<svelte:head>
	<title>FAQ | Хостамат</title>
</svelte:head>

<main class="page-shell">
	<div class="page-content faq">
		<h1>FAQ</h1>

		<Card border>
			<section class="qa">
				<h2>После загрузки сайта возвращается 404</h2>
				<p>
					Чаще всего файл <span class="path">index.html</span> лежит не в корне архива, а внутри
					вложенной папки. Убедитесь, что архив содержит файлы напрямую.
				</p>
				<p>Правила поиска страниц:</p>
				<ul>
					<li>
						Для <span class="path">/</span> нужен файл <span class="path">index.html</span> в корне
						архива или загруженный отдельно.
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
					Да, все действия доступны в <Link href={resolve('/dashboard')}>Дашборде</Link>.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="qa">
				<h2>Какие форматы файлов поддерживаются?</h2>
				<p>
					Список доступен на странице <Link href={resolve('/about')}>«О сервисе»</Link>. Архивы
					принимаются только в формате <span class="path">.zip</span>.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="qa">
				<h2>Что означает «альфа-тестирование»?</h2>
				<p>
					Сервис активно развивается: часть функций может меняться, а стабильность не всегда финальная.
					Нам важна обратная связь через <Link
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
		font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
			monospace;
		color: var(--color-primary-90);
	}
</style>
