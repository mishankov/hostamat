<script>
	import { Badge, Card, Link } from '$lib/components';
</script>

<svelte:head>
	<title>О сервисе | Хостамат</title>
</svelte:head>

<main>
	<div class="about">
		<header class="about-header">
			<h1>О сервисе</h1>
			<p class="about-intro">
				Хостамат (<Link href="https://hostamat.ru" target="_blank" alwaysUnderline>hostamat.ru</Link>)
				— сервис для простого хостинга статических сайтов.
			</p>
		</header>

		<div class="features-grid">
			<Card border>
				<div class="feature-card">
					<div class="feature-icon">
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
						</svg>
					</div>
					<h2>Простота использования</h2>
					<p>
						Для каждого сайта выделяется домен третьего уровня с настроенным TLS, поэтому каждый сайт
						доступен по HTTPS сразу после загрузки.
					</p>
				</div>
			</Card>

			<Card border>
				<div class="feature-card">
					<div class="feature-icon">
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
							<polyline points="17 8 12 3 7 8" />
							<line x1="12" y1="3" x2="12" y2="15" />
						</svg>
					</div>
					<h2>Гибкая загрузка</h2>
					<p>
						Сайты можно загружать напрямую на сайте сервиса в виде отдельных файлов или .zip архива.
					</p>
				</div>
			</Card>

			<Card border>
				<div class="feature-card formats-card">
					<div class="feature-icon">
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
							<polyline points="14 2 14 8 20 8" />
							<line x1="16" y1="13" x2="8" y2="13" />
							<line x1="16" y1="17" x2="8" y2="17" />
						</svg>
					</div>
					<h2>Поддерживаемые форматы</h2>

					<div class="format-section">
						<h3>Исходный код</h3>
						<div class="file-types">
							<Badge text="html" kind="primary" />
							<Badge text="css" kind="primary" />
							<Badge text="js" kind="primary" />
						</div>
					</div>

					<div class="format-section">
						<h3>Изображения</h3>
						<div class="file-types">
							<Badge text="png" kind="action" />
							<Badge text="jpg" kind="action" />
							<Badge text="jpeg" kind="action" />
							<Badge text="svg" kind="action" />
							<Badge text="gif" kind="action" />
							<Badge text="ico" kind="action" />
							<Badge text="webp" kind="action" />
						</div>
					</div>

					<div class="format-section">
						<h3>Шрифты</h3>
						<div class="file-types">
							<Badge text="woff2" kind="confirm" />
							<Badge text="woff" kind="confirm" />
							<Badge text="ttf" kind="confirm" />
							<Badge text="eot" kind="confirm" />
						</div>
					</div>
				</div>
			</Card>

			<Card border>
				<div class="feature-card">
					<div class="feature-icon">
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" />
						</svg>
					</div>
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
			</Card>
		</div>
	</div>
</main>

<style>
	main {
		width: 100%;
		min-height: calc(100vh - 140px);

		display: flex;
		justify-content: center;
		padding: 40px 24px;
	}

	.about {
		display: flex;
		flex-direction: column;
		gap: 32px;

		width: 100%;
		max-width: 900px;
	}

	.about-header {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.about-header h1 {
		font-size: 32px;
		font-weight: 700;
	}

	.about-intro {
		font-size: 18px;
		line-height: 1.6;
		color: var(--color-text-secondary);
		max-width: 600px;
	}

	.features-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 20px;
	}

	.feature-card {
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding: 24px;
	}

	.feature-icon {
		width: 44px;
		height: 44px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-accent-10);
		border-radius: var(--border-radius-regular);
		color: var(--color-accent-50);
	}

	.feature-card h2 {
		font-size: 18px;
		font-weight: 600;
		margin: 0;
	}

	.feature-card p {
		font-size: 15px;
		line-height: 1.6;
		color: var(--color-text-secondary);
		margin: 0;
	}

	.formats-card {
		grid-column: span 2;
	}

	.format-section {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding-top: 8px;
	}

	.format-section h3 {
		font-size: 13px;
		font-weight: 600;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.file-types {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}

	@media (max-width: 700px) {
		main {
			padding: 24px 16px;
		}

		.about-header h1 {
			font-size: 24px;
		}

		.formats-card {
			grid-column: span 1;
		}
	}
</style>
