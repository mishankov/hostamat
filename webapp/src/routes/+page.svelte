<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main class="page-shell">
	<section class="page-content landing">
		<div class="hero card-like">
			<Badge text="Альфа-тестирование" />
			<h1>Хостинг статических сайтов без лишней настройки</h1>
			<p>
				Загрузите архив или файлы, получите поддомен с HTTPS и обновляйте проект через аккуратный
				дашборд.
			</p>
			<div class="hero-actions">
				<Button link="/register" filled>Запустить свой сайт</Button>
				<Button link="/about" kind="action">Узнать о платформе</Button>
			</div>
		</div>

		<div class="steps">
			<Card floating border>
				<div class="step">
					<span>01</span>
					<h2>Создайте аккаунт</h2>
					<p>Регистрация занимает меньше минуты.</p>
				</div>
			</Card>
			<Card floating border>
				<div class="step">
					<span>02</span>
					<h2>Загрузите файлы</h2>
					<p>Поддерживаются отдельные файлы и ZIP-архив.</p>
				</div>
			</Card>
			<Card floating border>
				<div class="step">
					<span>03</span>
					<h2>Опубликуйте ссылку</h2>
					<p>Сайт сразу доступен по HTTPS-поддомену.</p>
				</div>
			</Card>
			<Card floating border kind="action">
				<div class="step">
					<span>04</span>
					<h2>До 5 сайтов бесплатно</h2>
					<p>Во время альфы квота доступна каждому пользователю.</p>
				</div>
			</Card>
		</div>
	</section>
</main>

<style>
	.landing {
		gap: 1.6rem;
		padding-bottom: 2rem;
	}

	.card-like {
		padding: clamp(1.2rem, 2.5vw, 2.2rem);
		border: var(--border-regular);
		border-radius: calc(var(--border-radius-regular) + 2px);
		background: color-mix(in hsl, var(--color-background-object) 86%, white 14%);
		box-shadow: var(--shadow-soft);
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	h1 {
		font-size: clamp(2rem, 5vw, 3.6rem);
		line-height: 1.08;
		max-width: 18ch;
	}

	.hero > p {
		max-width: 55ch;
		font-size: 1.05rem;
	}

	.hero-actions {
		display: flex;
		flex-wrap: wrap;
		gap: 0.8rem;
		padding-top: 0.5rem;
	}

	.steps {
		display: grid;
		grid-template-columns: repeat(2, minmax(220px, 1fr));
		gap: 0.9rem;
	}

	.step {
		min-height: 170px;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
	}

	.step span {
		font-family: var(--font-display);
		font-size: 0.9rem;
		letter-spacing: 0.15em;
		color: var(--color-accent-60);
	}

	.step h2 {
		font-size: 1.4rem;
	}

	@media (max-width: 800px) {
		.steps {
			grid-template-columns: 1fr;
		}
	}
</style>
