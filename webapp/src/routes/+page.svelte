<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main>
	<section class="hero">
		<div class="hero-content">
			<div class="hero-badge">
				<Badge text="Альфа-тестирование" />
			</div>

			<h1 class="hero-title">
				<span class="title-line">Хостинг</span>
				<span class="title-line accent">статических</span>
				<span class="title-line">сайтов</span>
			</h1>

			<p class="hero-subtitle">
				Загрузите свой сайт за секунды. Без настройки серверов,
				без сложных деплоев — просто работает.
			</p>

			<div class="hero-actions">
				<Button link="/register" kind="accent" filled>
					<span>Начать бесплатно</span>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<line x1="5" y1="12" x2="19" y2="12" />
						<polyline points="12 5 19 12 12 19" />
					</svg>
				</Button>
				<Button link="/about">Узнать больше</Button>
			</div>

			<div class="hero-stats">
				<div class="stat">
					<span class="stat-value">5</span>
					<span class="stat-label">бесплатных сайтов</span>
				</div>
				<div class="stat-divider"></div>
				<div class="stat">
					<span class="stat-value">0</span>
					<span class="stat-label">руб/мес</span>
				</div>
				<div class="stat-divider"></div>
				<div class="stat">
					<span class="stat-value mono">&lt;1мин</span>
					<span class="stat-label">на деплой</span>
				</div>
			</div>
		</div>

		<div class="hero-visual">
			<div class="terminal">
				<div class="terminal-header">
					<div class="terminal-buttons">
						<span class="btn-close"></span>
						<span class="btn-minimize"></span>
						<span class="btn-maximize"></span>
					</div>
					<span class="terminal-title">hostamat.ru</span>
				</div>
				<div class="terminal-body">
					<div class="terminal-line">
						<span class="prompt">$</span>
						<span class="command">upload</span>
						<span class="arg">my-site.zip</span>
					</div>
					<div class="terminal-line output">
						<span class="success">Uploading...</span>
					</div>
					<div class="terminal-line output">
						<span class="success">Site deployed!</span>
					</div>
					<div class="terminal-line output url">
						<span class="accent">https://my-site.hostamat.ru</span>
					</div>
					<div class="terminal-line">
						<span class="prompt">$</span>
						<span class="cursor"></span>
					</div>
				</div>
			</div>
		</div>
	</section>

	<section class="features">
		<div class="features-grid">
			<Card floating border>
				<div class="feature-card">
					<div class="feature-icon">
						<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
						</svg>
					</div>
					<h3>Мгновенный деплой</h3>
					<p>Загрузите файлы или архив — сайт появится онлайн сразу же</p>
				</div>
			</Card>

			<Card floating border>
				<div class="feature-card">
					<div class="feature-icon">
						<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
							<line x1="3" y1="9" x2="21" y2="9" />
							<line x1="9" y1="21" x2="9" y2="9" />
						</svg>
					</div>
					<h3>Простой интерфейс</h3>
					<p>Никаких сложных настроек — только то, что действительно нужно</p>
				</div>
			</Card>

			<Card floating border>
				<div class="feature-card">
					<div class="feature-icon">
						<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
						</svg>
					</div>
					<h3>HTTPS из коробки</h3>
					<p>Бесплатный SSL-сертификат для каждого сайта автоматически</p>
				</div>
			</Card>

			<Card floating border kind="action">
				<div class="feature-card highlight">
					<div class="feature-icon">
						<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
							<circle cx="12" cy="7" r="4" />
						</svg>
					</div>
					<h3>5 бесплатных сайтов</h3>
					<p>На время альфа-тестирования и навсегда после него</p>
					<Badge text="Альфа-тестирование" />
				</div>
			</Card>
		</div>
	</section>
</main>

<style>
	main {
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	/* ═══════════════════════════════════════════════════════════
	   HERO SECTION
	   ═══════════════════════════════════════════════════════════ */
	.hero {
		width: 100%;
		max-width: 1400px;
		min-height: calc(100vh - 200px);
		padding: 60px 24px;

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

	.hero-badge {
		animation: slideDown 0.6s ease-out;
	}

	.hero-title {
		font-size: clamp(42px, 6vw, 72px);
		font-weight: 700;
		line-height: 1.05;
		letter-spacing: -0.03em;
		animation: slideUp 0.6s ease-out;
	}

	.title-line {
		display: block;
	}

	.title-line.accent {
		background: linear-gradient(135deg, var(--color-accent-50) 0%, var(--color-accent-70) 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.hero-subtitle {
		font-size: 18px;
		color: var(--color-text-secondary);
		max-width: 480px;
		line-height: 1.6;
		animation: slideUp 0.6s ease-out 0.1s both;
	}

	.hero-actions {
		display: flex;
		gap: 12px;
		padding-top: 8px;
		animation: slideUp 0.6s ease-out 0.2s both;
	}

	.hero-stats {
		display: flex;
		align-items: center;
		gap: 24px;
		padding-top: 24px;
		animation: fadeIn 0.8s ease-out 0.4s both;
	}

	.stat {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.stat-value {
		font-size: 28px;
		font-weight: 700;
		color: var(--color-accent-50);
	}

	.stat-value.mono {
		font-family: 'JetBrains Mono', monospace;
	}

	.stat-label {
		font-size: 13px;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.stat-divider {
		width: 1px;
		height: 40px;
		background: var(--color-border);
	}

	/* ═══════════════════════════════════════════════════════════
	   TERMINAL VISUAL
	   ═══════════════════════════════════════════════════════════ */
	.hero-visual {
		display: flex;
		justify-content: center;
		align-items: center;
		animation: scaleIn 0.8s ease-out 0.3s both;
	}

	.terminal {
		width: 100%;
		max-width: 500px;
		background: var(--color-primary-7);
		border-radius: var(--border-radius-lg);
		border: 1px solid var(--glass-border);
		box-shadow: var(--shadow-lg);
		overflow: hidden;
	}

	.terminal-header {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px 16px;
		background: var(--color-primary-10);
		border-bottom: 1px solid var(--glass-border);
	}

	.terminal-buttons {
		display: flex;
		gap: 6px;
	}

	.terminal-buttons span {
		width: 12px;
		height: 12px;
		border-radius: 50%;
	}

	.btn-close { background: var(--color-danger-50); }
	.btn-minimize { background: var(--color-accent-50); }
	.btn-maximize { background: var(--color-confirm-50); }

	.terminal-title {
		font-family: 'JetBrains Mono', monospace;
		font-size: 12px;
		color: var(--color-text-muted);
		flex: 1;
		text-align: center;
		margin-right: 54px;
	}

	.terminal-body {
		padding: 20px;
		font-family: 'JetBrains Mono', monospace;
		font-size: 14px;
		line-height: 1.8;
	}

	.terminal-line {
		display: flex;
		gap: 8px;
	}

	.terminal-line.output {
		padding-left: 18px;
	}

	.terminal-line.url {
		padding-top: 8px;
	}

	.prompt {
		color: var(--color-accent-50);
	}

	.command {
		color: var(--color-primary-80);
	}

	.arg {
		color: var(--color-text-secondary);
	}

	.success {
		color: var(--color-confirm-60);
	}

	.terminal-line .accent {
		color: var(--color-accent-60);
		text-decoration: underline;
		text-underline-offset: 3px;
	}

	.cursor {
		width: 8px;
		height: 18px;
		background: var(--color-accent-50);
		animation: pulse 1s infinite;
	}

	/* ═══════════════════════════════════════════════════════════
	   FEATURES SECTION
	   ═══════════════════════════════════════════════════════════ */
	.features {
		width: 100%;
		max-width: 1400px;
		padding: 40px 24px 80px;
	}

	.features-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 20px;
	}

	.feature-card {
		padding: 28px;
		display: flex;
		flex-direction: column;
		gap: 12px;
		min-height: 200px;
	}

	.feature-icon {
		width: 48px;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-accent-10);
		border-radius: var(--border-radius-regular);
		color: var(--color-accent-50);
		margin-bottom: 8px;
	}

	.feature-card h3 {
		font-size: 20px;
		font-weight: 600;
	}

	.feature-card p {
		color: var(--color-text-secondary);
		font-size: 15px;
		line-height: 1.5;
	}

	.feature-card.highlight .feature-icon {
		background: var(--color-action-20);
		color: var(--color-action-60);
	}

	/* ═══════════════════════════════════════════════════════════
	   RESPONSIVE
	   ═══════════════════════════════════════════════════════════ */
	@media (max-width: 900px) {
		.hero {
			grid-template-columns: 1fr;
			min-height: auto;
			padding: 40px 20px;
			gap: 40px;
		}

		.hero-visual {
			order: -1;
		}

		.terminal {
			max-width: 400px;
		}

		.hero-stats {
			flex-wrap: wrap;
		}
	}

	@media (max-width: 600px) {
		.hero-title {
			font-size: 36px;
		}

		.hero-actions {
			flex-direction: column;
		}

		.hero-stats {
			gap: 16px;
		}

		.stat-divider {
			display: none;
		}

		.stat {
			flex: 1;
			min-width: 80px;
		}
	}
</style>
