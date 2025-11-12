<script>
	import { Badge, Card, Link } from '$lib/components';
</script>

<svelte:head>
	<title>О сервисе | Хостамат</title>
</svelte:head>

<main>
	<div class="about">
		<h1>О сервисе</h1>
		<p>
			Хостамат (<Link href="https://hostamat.ru" target="_blank" alwaysUnderline>hostamat.ru</Link>)
			— сервис для простого хостинга статических сайтов.
		</p>

		<div class="cards">
			<div class="card">
				<h2>Простота использования</h2>
				<p>
					Для каждого сайта выделяется домен третьего уровня с настроенным TLS, поэтому каждый сайт
					доступен по HTTPS сразу после загрузки.
				</p>
			</div>

			<div class="card">
				<h2>Гибкая загрузка</h2>
				<p>
					Сайты можно загружать напрямую на сайте сервиса в виде отдельных файлов или .zip архива.
				</p>
			</div>

			<div class="card">
				<h2>Поддерживаемые форматы</h2>
				<div class="file-types">
					<Badge text="HTML" />
					<Badge text="CSS" />
					<Badge text="JavaScript" />
					<Badge text="PNG" />
					<Badge text="JPG" />
					<Badge text="JPEG" />
					<Badge text="SVG" />
					<Badge text="ZIP" />
				</div>
			</div>

			<div class="card">
				<h2>Обратная связь</h2>
				<p>
					Для сообщения о возникших ошибках, вопросах или желании оставить какую-либо другую
					обратную связь можно воспользоваться <Link
						href="https://forms.yandex.ru/u/6907bbfc49af47f3cae7650a"
						target="_blank"
						alwaysUnderline>формой</Link
					>
				</p>
			</div>
		</div>
	</div>
</main>

<style>
	main {
		width: 100%;

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
	}

	.about {
		display: flex;
		flex-direction: column;
		gap: 10px;

		width: clamp(500px, 75%, 1200px);

		padding: 20px;
	}

	p {
		font-size: 1.25rem;
		line-height: 1.6;
		color: var(--color-text-secondary);
		max-width: 600px;
	}

	.cards {
		display: flex;
		flex-direction: column;
		gap: 32px;
	}

	.card {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.card h2 {
		padding-top: 10px;
	}

	.file-types {
		display: flex;
		gap: 12px;
		max-width: 600px;
	}
</style>
