<script>
	import { Badge, Card, Link } from '$lib/components';
</script>

<svelte:head>
	<title>О сервисе | Хостамат</title>
</svelte:head>

<main class="page-shell">
	<div class="page-content about">
		<h1>О сервисе</h1>
		<p class="intro">
			Хостамат (<Link href="https://hostamat.ru" target="_blank" alwaysUnderline>hostamat.ru</Link>)
			помогает быстро публиковать статические сайты.
		</p>

		<div class="cards">
			<Card border>
				<section class="card-section">
					<h2>Простота публикации</h2>
					<p>
						Для каждого сайта автоматически создаётся домен третьего уровня с TLS, поэтому HTTPS
						работает сразу после загрузки.
					</p>
				</section>
			</Card>

			<Card border>
				<section class="card-section">
					<h2>Гибкая загрузка</h2>
					<p>Вы можете загрузить отдельные файлы или ZIP-архив.</p>
				</section>
			</Card>

			<Card border>
				<section class="card-section">
					<h2>Поддерживаемые форматы</h2>
					<h3>Исходный код</h3>
					<div class="file-types">
						<Badge text="html" />
						<Badge text="css" />
						<Badge text="js" />
					</div>
					<h3>Изображения</h3>
					<div class="file-types">
						<Badge text="png" />
						<Badge text="jpg" />
						<Badge text="jpeg" />
						<Badge text="svg" />
						<Badge text="gif" />
						<Badge text="ico" />
						<Badge text="webp" />
					</div>
					<h3>Шрифты</h3>
					<div class="file-types">
						<Badge text="woff2" />
						<Badge text="woff" />
						<Badge text="ttf" />
						<Badge text="eot" />
					</div>
				</section>
			</Card>

			<Card border>
				<section class="card-section">
					<h2>Обратная связь</h2>
					<p>
						Если заметили баг или хотите предложить идею, отправьте сообщение через <Link
							href="https://forms.yandex.ru/u/6907bbfc49af47f3cae7650a"
							target="_blank"
							alwaysUnderline>форму</Link
						>.
					</p>
				</section>
			</Card>
		</div>
	</div>
</main>

<style>
	.about {
		gap: 1rem;
		padding-bottom: 2rem;
	}

	.intro {
		max-width: 66ch;
		font-size: 1.05rem;
	}

	.cards {
		display: grid;
		grid-template-columns: repeat(2, minmax(240px, 1fr));
		gap: 0.9rem;
	}

	.card-section {
		height: 100%;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.7rem;
	}

	.card-section h2 {
		font-size: clamp(1.2rem, 1.4vw, 1.45rem);
		line-height: 1.25;
	}

	.card-section h3 {
		font-size: clamp(1.02rem, 1.2vw, 1.18rem);
		line-height: 1.3;
		margin: 0.35rem 0 0;
	}

	.file-types {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}

	@media (max-width: 860px) {
		.cards {
			grid-template-columns: 1fr;
		}
	}
</style>
