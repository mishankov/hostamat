<script>
	import { Badge, Link } from '$lib/components';
</script>

<svelte:head>
	<title>О сервисе | Хостамат</title>
</svelte:head>

<main>
	<div class="about">
		<header class="page-header">
			<h1>О сервисе</h1>
			<p class="intro">
				<Link href="https://hostamat.ru" target="_blank" alwaysUnderline>Хостамат</Link> — сервис для простого хостинга статических сайтов.
			</p>
		</header>

		<div class="features">
			<section class="feature-section">
				<div class="feature-icon">01</div>
				<div class="feature-content">
					<h2>Простота использования</h2>
					<p>
						Для каждого сайта выделяется домен третьего уровня с настроенным TLS, поэтому каждый сайт
						доступен по HTTPS сразу после загрузки.
					</p>
				</div>
			</section>

			<section class="feature-section">
				<div class="feature-icon">02</div>
				<div class="feature-content">
					<h2>Гибкая загрузка</h2>
					<p>
						Сайты можно загружать напрямую на сайте сервиса в виде отдельных файлов или .zip архива.
					</p>
				</div>
			</section>

			<section class="feature-section formats">
				<div class="feature-icon">03</div>
				<div class="feature-content">
					<h2>Поддерживаемые форматы</h2>

					<div class="format-group">
						<h3>Исходный код</h3>
						<div class="file-types">
							<Badge text="html" kind="muted" />
							<Badge text="css" kind="muted" />
							<Badge text="js" kind="muted" />
						</div>
					</div>

					<div class="format-group">
						<h3>Изображения</h3>
						<div class="file-types">
							<Badge text="png" kind="muted" />
							<Badge text="jpg" kind="muted" />
							<Badge text="jpeg" kind="muted" />
							<Badge text="svg" kind="muted" />
							<Badge text="gif" kind="muted" />
							<Badge text="ico" kind="muted" />
							<Badge text="webp" kind="muted" />
						</div>
					</div>

					<div class="format-group">
						<h3>Шрифты</h3>
						<div class="file-types">
							<Badge text="woff2" kind="muted" />
							<Badge text="woff" kind="muted" />
							<Badge text="ttf" kind="muted" />
							<Badge text="eot" kind="muted" />
						</div>
					</div>
				</div>
			</section>

			<section class="feature-section">
				<div class="feature-icon">?</div>
				<div class="feature-content">
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
			</section>
		</div>
	</div>
</main>

<style>
	main {
		width: 100%;
		min-height: calc(100vh - 140px);

		display: flex;
		justify-content: center;

		padding: 0 24px;
	}

	.about {
		width: 100%;
		max-width: 800px;

		display: flex;
		flex-direction: column;
		gap: 48px;

		padding: 40px 0 60px;
	}

	.page-header {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.page-header h1 {
		margin: 0;
	}

	.intro {
		font-size: 1.125rem;
		line-height: 1.7;
		color: var(--color-text-secondary);
		max-width: 600px;
		margin: 0;
	}

	.features {
		display: flex;
		flex-direction: column;
		gap: 40px;
	}

	.feature-section {
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 24px;
		align-items: start;
	}

	.feature-icon {
		font-family: var(--font-display);
		font-size: 1.5rem;
		font-weight: 800;
		color: var(--color-accent-50);
		opacity: 0.5;

		width: 48px;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;

		background: var(--color-background-object);
		border: 1px solid var(--color-border-subtle);
		border-radius: var(--border-radius-regular);
	}

	.feature-content {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.feature-content h2 {
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 700;
		margin: 0;
	}

	.feature-content p {
		font-size: 0.9375rem;
		line-height: 1.7;
		color: var(--color-text-secondary);
		margin: 0;
	}

	.feature-content h3 {
		font-family: var(--font-mono);
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.03em;
		margin: 0;
	}

	.formats .feature-content {
		gap: 20px;
	}

	.format-group {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.file-types {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}

	@media (max-width: 600px) {
		main {
			padding: 0 16px;
		}

		.about {
			padding: 24px 0 40px;
			gap: 32px;
		}

		.feature-section {
			grid-template-columns: 1fr;
			gap: 16px;
		}

		.feature-icon {
			width: 40px;
			height: 40px;
			font-size: 1.25rem;
		}
	}
</style>
