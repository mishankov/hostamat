<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main class="page-shell">
	<section class="page-content landing">
		<Card border cornerBorder class="card-like">
			<Badge text="Альфа-тестирование" />
			<h1>Хостинг статических сайтов без лишней настройки</h1>
			<p class="hero-description">
				Загрузите файлы или ZIP-архив, получите HTTPS-поддомен и управляйте проектом в удобной
				панели.
			</p>
			<div class="hero-actions">
				<Button href="/register" filled>Запустить свой сайт</Button>
				<Button href="/about">Узнать о платформе</Button>
			</div>
		</Card>

		<div class="steps">
			<Card floating border cornerBorder>
				<div class="step">
					<span>01</span>
					<h2>Создайте аккаунт</h2>
					<p>Регистрация занимает меньше минуты.</p>
				</div>
			</Card>
			<Card floating border cornerBorder>
				<div class="step">
					<span>02</span>
					<h2>Загрузите файлы</h2>
					<p>Можно загрузить отдельные файлы или ZIP-архив.</p>
				</div>
			</Card>
			<Card floating border cornerBorder>
				<div class="step">
					<span>03</span>
					<h2>Поделитесь ссылкой</h2>
					<p>Сайт сразу доступен по HTTPS-поддомену.</p>
				</div>
			</Card>
			<Card floating border cornerBorder kind="action">
				<div class="step">
					<span>04</span>
					<h2>До 5 сайтов бесплатно</h2>
					<p>Во время альфа-тестирования до 5 сайтов доступны каждому пользователю.</p>
				</div>
			</Card>
		</div>
	</section>
</main>

<style>
	.landing {
		position: relative;
		z-index: 1;
		gap: 1.6rem;
		padding-bottom: 2rem;
	}

	.landing :global(.card) {
		background: color-mix(in hsl, var(--color-background-object) 96%, var(--color-background));
		box-shadow: var(--shadow-sm);
	}

	:global(:root[data-theme='dark']) .landing :global(.card) {
		background: color-mix(in hsl, var(--color-background-object) 94%, var(--color-primary-10));
		box-shadow: var(--shadow-md);
	}

	:global(.card-like) {
		padding: clamp(1.2rem, 2.5vw, 2.2rem);
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	h1 {
		font-size: clamp(2rem, 5vw, 3.6rem);
		line-height: 1.08;
		max-width: 18ch;
	}

	.hero-description {
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
