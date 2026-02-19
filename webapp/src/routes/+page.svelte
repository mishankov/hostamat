<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main class="page-shell">
	<section class="page-content landing">
		<Card border cornerBorder class="card-like">
			<div class="hero-main">
				<div class="hero-copy">
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
				</div>
				<div class="hero-ascii">
					<div class="hero-terminal" aria-hidden="true">
						<div class="terminal-head terminal-line" style="--line-index: 0;">
							<span class="terminal-head-label">deploy@hostamat</span>
							<span class="terminal-head-rule"></span>
						</div>
						<p class="terminal-brand terminal-line" style="--line-index: 1;">HOSTAMAT</p>
						<div class="terminal-divider terminal-line" style="--line-index: 2;"></div>
						<p class="terminal-command terminal-line" style="--line-index: 3;">
							$ hostamat publish ./portfolio.zip
						</p>
						<ul class="terminal-steps">
							<li class="terminal-step terminal-line" style="--line-index: 4;">
								<span class="terminal-step-index">[01]</span>
								<span class="terminal-step-text">unpack assets</span>
								<span class="terminal-step-fill"></span>
								<span class="terminal-step-status">[ok]</span>
							</li>
							<li class="terminal-step terminal-line" style="--line-index: 5;">
								<span class="terminal-step-index">[02]</span>
								<span class="terminal-step-text">reserve subdomain</span>
								<span class="terminal-step-fill"></span>
								<span class="terminal-step-status">[ok]</span>
							</li>
							<li class="terminal-step terminal-line" style="--line-index: 6;">
								<span class="terminal-step-index">[03]</span>
								<span class="terminal-step-text">issue tls certificate</span>
								<span class="terminal-step-fill"></span>
								<span class="terminal-step-status">[ok]</span>
							</li>
							<li class="terminal-step terminal-line" style="--line-index: 7;">
								<span class="terminal-step-index">[04]</span>
								<span class="terminal-step-text">warm edge cache</span>
								<span class="terminal-step-fill"></span>
								<span class="terminal-step-status">[ok]</span>
							</li>
						</ul>
						<div class="terminal-divider terminal-line" style="--line-index: 8;"></div>
						<p class="terminal-meta-row terminal-line" style="--line-index: 9;">
							<span class="terminal-meta-key">domain:</span>
							<span class="terminal-meta-value">happy-cappy-12345.hostamat.ru</span>
						</p>
						<p class="terminal-meta-row terminal-line" style="--line-index: 10;">
							<span class="terminal-meta-key">status:</span>
							<span class="terminal-meta-value">live | https</span>
						</p>
					</div>
				</div>
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
	}

	.hero-main {
		display: grid;
		grid-template-columns: minmax(320px, 0.95fr) minmax(500px, 1.2fr);
		align-items: center;
		gap: clamp(1rem, 2.4vw, 2.4rem);
	}

	.hero-copy {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		min-width: 0;
	}

	h1 {
		font-size: clamp(2rem, 5vw, 3.6rem);
		line-height: 1.08;
		max-width: 18ch;
	}

	.hero-description {
		max-width: 44ch;
		font-size: 1.05rem;
	}

	.hero-ascii {
		margin-top: 0.1rem;
		position: relative;
		overflow: hidden;
		inline-size: 100%;
		border: var(--border-width) var(--border-style)
			color-mix(in hsl, var(--color-border) 75%, var(--color-accent-30) 25%);
		border-radius: 12px;
		background: linear-gradient(
			180deg,
			color-mix(in hsl, var(--color-background-object) 90%, var(--color-accent-20) 10%) 0%,
			color-mix(in hsl, var(--color-background-object) 94%, var(--color-primary-10) 6%) 100%
		);
		box-shadow:
			inset 0 1px 0 color-mix(in hsl, var(--color-background-object) 75%, white 25%),
			inset 0 0 0 1px color-mix(in hsl, var(--color-accent-20) 65%, transparent),
			0 10px 24px color-mix(in hsl, var(--color-accent-30) 16%, transparent);
	}

	.hero-terminal {
		position: relative;
		z-index: 1;
		display: grid;
		gap: 0.3rem;
		padding: 0.92rem 1.02rem;
		font-family: var(--font-body);
		font-size: clamp(0.58rem, 1.04vw, 0.9rem);
		line-height: 1.1;
		letter-spacing: 0;
		font-kerning: none;
		font-variant-ligatures: none;
		color: var(--color-accent-70);
		text-shadow: 0 0 0.35rem color-mix(in hsl, var(--color-accent-30) 12%, transparent);
	}

	.hero-ascii::after {
		content: '';
		position: absolute;
		inset-inline: 0;
		top: -22%;
		height: 18%;
		pointer-events: none;
		background: linear-gradient(
			180deg,
			transparent 0%,
			color-mix(in hsl, var(--color-accent-40) 64%, transparent) 50%,
			transparent 100%
		);
		opacity: 0;
		animation: asciiScanPass 1100ms ease-out 580ms 1 both;
	}

	.terminal-head {
		display: flex;
		align-items: center;
		gap: 0.48rem;
	}

	.terminal-head-label {
		color: var(--color-accent-80);
		font-weight: 600;
		white-space: nowrap;
	}

	.terminal-head-rule {
		flex: 1;
		height: 1px;
		background: repeating-linear-gradient(
			90deg,
			color-mix(in hsl, var(--color-accent-50) 72%, transparent) 0 8px,
			transparent 8px 12px
		);
	}

	.terminal-brand {
		margin: 0;
		padding-top: 0.04em;
		font-family: var(--font-display);
		font-size: clamp(1.82rem, 4.3vw, 4.3rem);
		font-weight: 700;
		line-height: 0.88;
		letter-spacing: 0.05em;
		color: color-mix(in hsl, var(--color-accent-70) 92%, var(--color-background-object));
		text-shadow:
			0 0 0.6rem color-mix(in hsl, var(--color-accent-30) 18%, transparent),
			0 0 1px color-mix(in hsl, var(--color-accent-40) 45%, transparent);
		-webkit-text-stroke: 1px color-mix(in hsl, var(--color-accent-50) 44%, transparent);
	}

	.terminal-divider {
		margin: 0.1rem 0 0.05rem;
		border-top: 1px dashed color-mix(in hsl, var(--color-accent-50) 72%, transparent);
	}

	.terminal-command {
		margin: 0;
		overflow-wrap: anywhere;
	}

	.terminal-steps {
		margin: 0;
		padding: 0;
		list-style: none;
		display: grid;
		gap: 0.18rem;
	}

	.terminal-step {
		display: grid;
		grid-template-columns: auto auto minmax(0, 1fr) auto;
		align-items: center;
		gap: 0.46rem;
		min-width: 0;
	}

	.terminal-step-index {
		color: color-mix(in hsl, var(--color-accent-60) 86%, var(--color-accent-90));
	}

	.terminal-step-text {
		overflow-wrap: anywhere;
	}

	.terminal-step-fill {
		border-bottom: 1px dotted color-mix(in hsl, var(--color-accent-50) 72%, transparent);
		transform: translateY(-0.08rem);
	}

	.terminal-step-status {
		color: color-mix(in hsl, var(--color-confirm-60) 58%, var(--color-accent-70));
	}

	.terminal-meta-row {
		margin: 0;
		display: flex;
		flex-wrap: wrap;
		gap: 0 0.36rem;
	}

	.terminal-meta-key {
		color: color-mix(in hsl, var(--color-accent-60) 86%, var(--color-accent-90));
	}

	.terminal-meta-value {
		overflow-wrap: anywhere;
	}

	:global(:root[data-theme='dark']) .hero-terminal {
		color: var(--color-accent-80);
		text-shadow: 0 0 0.45rem color-mix(in hsl, var(--color-accent-60) 14%, transparent);
	}

	:global(:root[data-theme='dark']) .terminal-brand {
		color: color-mix(in hsl, var(--color-accent-80) 92%, var(--color-background-object));
		-webkit-text-stroke: 1px color-mix(in hsl, var(--color-accent-70) 52%, transparent);
	}

	.terminal-line {
		opacity: 0;
		transform: translateY(0.22rem);
		animation: asciiLineReveal 360ms cubic-bezier(0.2, 0.75, 0.25, 1) forwards;
		animation-delay: calc(var(--line-index) * 32ms);
	}

	@keyframes asciiLineReveal {
		from {
			opacity: 0;
			transform: translateY(0.22rem);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes asciiScanPass {
		0% {
			opacity: 0;
			transform: translateY(0%);
		}
		16% {
			opacity: 0.38;
		}
		100% {
			opacity: 0;
			transform: translateY(680%);
		}
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

	@media (max-width: 920px) {
		.hero-main {
			grid-template-columns: 1fr;
			gap: 0.9rem;
		}

		.hero-terminal {
			font-size: clamp(0.56rem, 1.5vw, 0.78rem);
			line-height: 1.04;
		}

		.terminal-brand {
			font-size: clamp(1.64rem, 8vw, 3rem);
		}
	}

	@media (max-width: 680px) {
		.hero-terminal {
			font-size: clamp(0.56rem, 1.95vw, 0.68rem);
			line-height: 1;
		}

		.terminal-step {
			grid-template-columns: auto minmax(0, 1fr) auto;
		}

		.terminal-step-index {
			grid-column: 1;
		}

		.terminal-step-text {
			grid-column: 2;
		}

		.terminal-step-fill {
			display: none;
		}

		.terminal-step-status {
			grid-column: 3;
			justify-self: end;
		}
	}

	@media (max-width: 800px) {
		.steps {
			grid-template-columns: 1fr;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.hero-ascii::after,
		.terminal-line {
			animation: none;
		}

		.hero-ascii::after {
			opacity: 0;
		}

		.terminal-line {
			opacity: 1;
			transform: none;
		}
	}
</style>
