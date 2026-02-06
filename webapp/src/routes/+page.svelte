<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main class="page-shell landing-shell">
	<div class="landing-grid-bg" aria-hidden="true"></div>
	<section class="page-content landing">
		<Card border cornerBorder class="card-like">
			<Badge text="Альфа-тестирование" />
			<h1>Хостинг статических сайтов без лишней настройки</h1>
			<p class="hero-description">
				Загрузите архив или файлы, получите поддомен с HTTPS и обновляйте проект через аккуратный
				дашборд.
			</p>
			<div class="hero-actions">
				<Button link="/register" filled>Запустить свой сайт</Button>
				<Button link="/about" kind="action">Узнать о платформе</Button>
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
					<p>Поддерживаются отдельные файлы и ZIP-архив.</p>
				</div>
			</Card>
			<Card floating border cornerBorder>
				<div class="step">
					<span>03</span>
					<h2>Опубликуйте ссылку</h2>
					<p>Сайт сразу доступен по HTTPS-поддомену.</p>
				</div>
			</Card>
			<Card floating border cornerBorder kind="action">
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
	.landing-shell {
		position: relative;
		isolation: isolate;
		overflow: hidden;

		--landing-grid-size: 28px;
		--landing-grid-major-size: 112px;
		--landing-grid-line: color-mix(in hsl, var(--color-accent-60) 8%, transparent);
		--landing-grid-major-line: color-mix(in hsl, var(--color-accent-70) 11%, transparent);
		--landing-glow-strong: color-mix(in hsl, var(--color-accent-40) 20%, transparent);
		--landing-glow-soft: color-mix(in hsl, var(--color-accent-30) 11%, transparent);
		--landing-fade: color-mix(in hsl, var(--color-background) 95%, transparent);
	}

	:global(:root[data-theme='dark']) .landing-shell {
		--landing-grid-line: color-mix(in hsl, var(--color-accent-70) 10%, transparent);
		--landing-grid-major-line: color-mix(in hsl, var(--color-accent-80) 14%, transparent);
		--landing-glow-strong: color-mix(in hsl, var(--color-accent-60) 14%, transparent);
		--landing-glow-soft: color-mix(in hsl, var(--color-accent-50) 8%, transparent);
		--landing-fade: color-mix(in hsl, var(--color-background) 92%, transparent);
	}

	.landing-grid-bg {
		position: fixed;
		inset: 0;
		pointer-events: none;
		z-index: 0;
		background:
			radial-gradient(68% 55% at 86% 7%, var(--landing-glow-strong) 0%, transparent 69%),
			radial-gradient(52% 42% at 11% 92%, var(--landing-glow-soft) 0%, transparent 74%),
			linear-gradient(180deg, var(--landing-fade) 0%, transparent 32%, transparent 100%),
			repeating-linear-gradient(
				0deg,
				transparent 0 calc(var(--landing-grid-major-size) - 1px),
				var(--landing-grid-major-line) calc(var(--landing-grid-major-size) - 1px)
					var(--landing-grid-major-size)
			),
			repeating-linear-gradient(
				90deg,
				transparent 0 calc(var(--landing-grid-major-size) - 1px),
				var(--landing-grid-major-line) calc(var(--landing-grid-major-size) - 1px)
					var(--landing-grid-major-size)
			),
			repeating-linear-gradient(
				0deg,
				transparent 0 calc(var(--landing-grid-size) - 1px),
				var(--landing-grid-line) calc(var(--landing-grid-size) - 1px) var(--landing-grid-size)
			),
			repeating-linear-gradient(
				90deg,
				transparent 0 calc(var(--landing-grid-size) - 1px),
				var(--landing-grid-line) calc(var(--landing-grid-size) - 1px) var(--landing-grid-size)
			);
	}

	.landing {
		position: relative;
		z-index: 1;
		gap: 1.6rem;
		padding-bottom: 2rem;
	}

	.landing :global(.card) {
		background: color-mix(in hsl, var(--color-background-object) 96%, var(--color-background));
		box-shadow: 0 10px 28px hsl(214 35% 12% / 0.09);
	}

	:global(:root[data-theme='dark']) .landing :global(.card) {
		background: color-mix(in hsl, var(--color-background-object) 94%, var(--color-primary-10));
		box-shadow: 0 14px 36px hsl(212 60% 3% / 0.42);
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
		.landing-shell {
			--landing-grid-size: 34px;
			--landing-grid-major-size: 136px;
			--landing-grid-line: color-mix(in hsl, var(--color-accent-60) 6%, transparent);
			--landing-grid-major-line: color-mix(in hsl, var(--color-accent-70) 8%, transparent);
			--landing-glow-strong: color-mix(in hsl, var(--color-accent-40) 14%, transparent);
			--landing-glow-soft: color-mix(in hsl, var(--color-accent-30) 8%, transparent);
		}

		:global(:root[data-theme='dark']) .landing-shell {
			--landing-grid-line: color-mix(in hsl, var(--color-accent-70) 7%, transparent);
			--landing-grid-major-line: color-mix(in hsl, var(--color-accent-80) 10%, transparent);
			--landing-glow-strong: color-mix(in hsl, var(--color-accent-60) 10%, transparent);
			--landing-glow-soft: color-mix(in hsl, var(--color-accent-50) 6%, transparent);
		}

		.steps {
			grid-template-columns: 1fr;
		}
	}
</style>
