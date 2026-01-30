<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main>
	<div class="landing">
		<section class="hero">
			<div class="hero-content">
				<Badge text="Альфа-тестирование" />
				<h1 class="hero-title">
					<span class="title-line">Запусти свой</span>
					<span class="title-line accent">статический сайт</span>
					<span class="title-line">за минуту</span>
				</h1>
				<p class="hero-subtitle">
					Хостинг без настроек. Загрузи файлы — получи ссылку. HTTPS включён.
				</p>
				<div class="hero-actions">
					<Button link="/register" filled>Загрузить свой сайт</Button>
					<Button link="/about">Узнать больше</Button>
				</div>
			</div>

			<div class="hero-visual">
				<div class="code-window">
					<div class="window-header">
						<div class="window-dots">
							<span class="dot red"></span>
							<span class="dot yellow"></span>
							<span class="dot green"></span>
						</div>
						<span class="window-title">terminal</span>
					</div>
					<div class="window-content">
						<div class="code-line">
							<span class="prompt">$</span>
							<span class="command">upload</span>
							<span class="arg">./my-site</span>
						</div>
						<div class="code-line output">
							<span class="symbol">→</span>
							<span>Загрузка...</span>
							<span class="loader"></span>
						</div>
						<div class="code-line output success">
							<span class="symbol">✓</span>
							<span>Готово!</span>
						</div>
						<div class="code-line output link">
							<span class="symbol">◈</span>
							<span class="url">https://my-site.hostamat.ru</span>
						</div>
					</div>
				</div>
			</div>
		</section>

		<section class="features">
			<h2 class="section-title">Как это работает</h2>
			<div class="feature-grid">
				<Card floating border>
					<div class="feature-card">
						<div class="feature-number">01</div>
						<h3>Создай аккаунт</h3>
						<p>Регистрация займёт меньше минуты. Никаких лишних вопросов.</p>
					</div>
				</Card>

				<Card floating border>
					<div class="feature-card">
						<div class="feature-number">02</div>
						<h3>Загрузи сайт</h3>
						<p>Перетащи файлы или загрузи ZIP-архив. Мы всё распакуем.</p>
					</div>
				</Card>

				<Card floating border>
					<div class="feature-card">
						<div class="feature-number">03</div>
						<h3>Получи ссылку</h3>
						<p>Твой сайт доступен по HTTPS сразу после загрузки.</p>
					</div>
				</Card>

				<Card floating border>
					<div class="feature-card highlight">
						<div class="feature-number">∞</div>
						<h3>5 бесплатных сайтов</h3>
						<p>На время альфа-тестирования и навсегда после него.</p>
						<Badge text="Бесплатно" kind="confirm" />
					</div>
				</Card>
			</div>
		</section>
	</div>
</main>

<style>
	main {
		width: 100%;
		min-height: calc(100vh - 140px);

		display: flex;
		flex-direction: column;
		align-items: center;

		padding: 0 24px;
	}

	.landing {
		width: 100%;
		max-width: 1200px;

		display: flex;
		flex-direction: column;
		gap: 100px;

		padding: 60px 0 80px;
	}

	/* Hero Section */
	.hero {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 60px;
		align-items: center;
	}

	.hero-content {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.hero-title {
		font-size: clamp(2.5rem, 5vw, 3.5rem);
		line-height: 1.1;
		margin: 0;
	}

	.title-line {
		display: block;
	}

	.title-line.accent {
		color: var(--color-accent-50);
		background: linear-gradient(135deg, var(--color-accent-50), var(--color-accent-70));
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.hero-subtitle {
		font-size: 1.125rem;
		color: var(--color-text-secondary);
		line-height: 1.6;
		max-width: 440px;
	}

	.hero-actions {
		display: flex;
		flex-wrap: wrap;
		gap: 16px;
		padding-top: 8px;
	}

	/* Code Window Visual */
	.hero-visual {
		display: flex;
		justify-content: center;
		perspective: 1000px;
	}

	.code-window {
		width: 100%;
		max-width: 420px;

		background: var(--color-primary-7);
		border: 1px solid var(--color-border);
		border-radius: var(--border-radius-lg);

		box-shadow: var(--shadow-lg);
		transform: rotateY(-5deg) rotateX(2deg);
		transition: transform var(--transition-slow);
	}

	.code-window:hover {
		transform: rotateY(0deg) rotateX(0deg);
	}

	.window-header {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 14px 16px;
		border-bottom: 1px solid var(--color-border-subtle);
	}

	.window-dots {
		display: flex;
		gap: 8px;
	}

	.dot {
		width: 12px;
		height: 12px;
		border-radius: 50%;
	}

	.dot.red { background: hsl(0, 65%, 55%); }
	.dot.yellow { background: hsl(45, 80%, 55%); }
	.dot.green { background: hsl(120, 50%, 50%); }

	.window-title {
		font-size: 0.75rem;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.window-content {
		padding: 20px;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.code-line {
		display: flex;
		align-items: center;
		gap: 10px;
		font-family: var(--font-mono);
		font-size: 0.9375rem;
	}

	.prompt {
		color: var(--color-accent-50);
		font-weight: 600;
	}

	.command {
		color: var(--color-text-primary);
	}

	.arg {
		color: var(--color-text-secondary);
	}

	.code-line.output {
		color: var(--color-text-secondary);
		padding-left: 4px;
	}

	.code-line.output .symbol {
		width: 14px;
	}

	.code-line.success {
		color: var(--color-confirm-60);
	}

	.code-line.link {
		color: var(--color-accent-60);
	}

	.code-line .url {
		text-decoration: underline;
		text-underline-offset: 3px;
	}

	.loader {
		width: 10px;
		height: 10px;
		border: 2px solid var(--color-accent-30);
		border-top-color: var(--color-accent-50);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	/* Features Section */
	.features {
		display: flex;
		flex-direction: column;
		gap: 40px;
	}

	.section-title {
		text-align: center;
		color: var(--color-text-secondary);
		font-weight: 500;
		letter-spacing: 0.02em;
	}

	.feature-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
		gap: 20px;
	}

	.feature-card {
		padding: 28px;
		display: flex;
		flex-direction: column;
		gap: 12px;
		height: 100%;
	}

	.feature-number {
		font-family: var(--font-display);
		font-size: 2rem;
		font-weight: 800;
		color: var(--color-accent-40);
		opacity: 0.6;
		line-height: 1;
	}

	.feature-card h3 {
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 700;
		margin: 0;
	}

	.feature-card p {
		font-size: 0.9375rem;
		color: var(--color-text-secondary);
		line-height: 1.6;
		margin: 0;
		flex: 1;
	}

	.feature-card.highlight .feature-number {
		color: var(--color-confirm-50);
	}

	/* Responsive */
	@media (max-width: 900px) {
		.hero {
			grid-template-columns: 1fr;
			gap: 48px;
		}

		.hero-visual {
			order: -1;
		}

		.code-window {
			transform: none;
			max-width: 100%;
		}

		.landing {
			gap: 80px;
			padding: 40px 0 60px;
		}
	}

	@media (max-width: 600px) {
		main {
			padding: 0 16px;
		}

		.hero-title {
			font-size: 2rem;
		}

		.feature-card {
			padding: 24px;
		}
	}
</style>
