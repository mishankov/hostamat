<script>
	import { Card, Link } from '$lib/components';
</script>

<main>
	<div class="about">
		<div class="header">
			<h1>О сервисе</h1>
			<p class="lead">
				Хостамат (<Link href="hostamat.ru" target="_blank" alwaysUnderline>hostamat.ru</Link>) —
				сервис для простого хостинга статических сайтов.
			</p>
		</div>

		<div class="cards">
			<Card floating>
				<div class="card">
					<h2>🚀 Простота использования</h2>
					<p>
						Для каждого сайта выделяется домен третьего уровня с настроенным TLS, поэтому каждый
						сайт доступен по HTTPS сразу после загрузки.
					</p>
				</div>
			</Card>

			<Card floating>
				<div class="card">
					<h2>📁 Гибкая загрузка</h2>
					<p>
						Сайты можно загружать напрямую на сайте сервиса в виде отдельных файлов или .zip архива.
					</p>
				</div>
			</Card>

			<Card floating>
				<div class="card">
					<h2>📄 Поддерживаемые форматы</h2>
					<div class="file-types">
						<div class="file-type">HTML</div>
						<div class="file-type">CSS</div>
						<div class="file-type">JavaScript</div>
						<div class="file-type">PNG</div>
						<div class="file-type">JPG</div>
						<div class="file-type">JPEG</div>
						<div class="file-type">SVG</div>
						<div class="file-type">ZIP</div>
					</div>
				</div>
			</Card>

			<Card floating>
				<div class="card">
					<p>
						Для сообщения о возникших ошибках, вопросах или желании оставить какую-либо другую
						обратную связь можно воспользоваться <Link
							href="https://forms.yandex.ru/u/6907bbfc49af47f3cae7650a"
							target="_blank"
							alwaysUnderline>формой</Link
						>
					</p>
				</div>
			</Card>
		</div>
	</div>
</main>

<style>
	main {
		width: 100%;
		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: flex-start;
		min-height: calc(100vh - 200px);
		padding: 40px 20px;
	}

	.about {
		display: flex;
		flex-direction: column;
		gap: 40px;
		width: clamp(500px, 75%, 800px);
	}

	.header {
		text-align: center;
		margin-bottom: 20px;
	}

	h1 {
		font-size: 2.5rem;
		font-weight: 700;
		margin-bottom: 16px;
	}

	.lead {
		font-size: 1.25rem;
		line-height: 1.6;
		color: var(--color-text-secondary);
		max-width: 600px;
		margin: 0 auto;
	}

	.cards {
		display: flex;
		flex-direction: column;
		gap: 32px;
	}

	.card {
		padding: 24px;

		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.card p {
		line-height: 1.6;
		color: var(--color-text-secondary);
		font-size: 1rem;
	}

	.file-types {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
		gap: 12px;
		margin-top: 16px;
	}

	.file-type {
		background: var(--color-primary-dim);
		border: var(--border-regular);
		border-radius: var(--border-radius-regular);
		padding: 8px 12px;
		text-align: center;
		font-size: 0.9rem;
		font-weight: 500;
		transition: var(--transition);
	}
</style>
