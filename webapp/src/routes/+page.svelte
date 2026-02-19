<script>
	import { Badge, Button, Card } from '$lib/components';

	const heroAsciiInnerWidth = 65;
	const frameRow = (text = '') => `|${text.padEnd(heroAsciiInnerWidth, ' ')}|`;
	const frameTop = (label = '') => {
		const prefix = `- ${label} `;
		return `+${prefix}${'-'.repeat(Math.max(0, heroAsciiInnerWidth - prefix.length))}+`;
	};
	const frameDivider = () => frameRow('-'.repeat(heroAsciiInnerWidth));
	const frameBottom = () => `+${'-'.repeat(heroAsciiInnerWidth)}+`;

	const heroAsciiLines = [
		frameTop('deploy@hostamat'),
		frameRow(' _   _  ___  ____ _____   _    __  __    _    _____'),
		frameRow('| | | |/ _ \\/ ___|_   _| / \\  |  \\/  |  / \\  |_   _|'),
		frameRow('| |_| | | | \\___ \\ | |  / _ \\ | |\\/| | / _ \\   | |'),
		frameRow('|  _  | |_| |___) || | / ___ \\| |  | |/ ___ \\  | |'),
		frameRow('|_| |_|\\___/|____/ |_|/_/   \\_\\_|  |_/_/   \\_\\ |_|'),
		frameDivider(),
		frameRow(' $ hostamat publish ./portfolio.zip'),
		frameRow(' [01] unpack assets ................................ [ok]'),
		frameRow(' [02] reserve subdomain ............................ [ok]'),
		frameRow(' [03] issue tls certificate ........................ [ok]'),
		frameRow(' [04] warm edge cache .............................. [ok]'),
		frameDivider(),
		frameRow(' domain: happy-cappy-12345.hostamat.ru'),
		frameRow(' status: live | https'),
		frameBottom()
	];
	const heroAsciiRows = heroAsciiLines;
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
				<div class="hero-ascii" aria-hidden="true">
					<pre>{#each heroAsciiRows as row, idx}<span
								class="ascii-line"
								style={`--line-index: ${idx};`}>{row}</span
							>{/each}</pre>
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

	.hero-ascii pre {
		margin: 0;
		display: block;
		padding: 0.85rem 0.95rem;
		font-family: var(--font-body);
		font-size: clamp(0.54rem, 1.05vw, 0.86rem);
		line-height: 1;
		letter-spacing: 0;
		font-kerning: none;
		font-variant-ligatures: none;
		white-space: normal;
		color: var(--color-accent-70);
		text-shadow: 0 0 0.35rem color-mix(in hsl, var(--color-accent-30) 12%, transparent);
	}

	:global(:root[data-theme='dark']) .hero-ascii pre {
		color: var(--color-accent-80);
		text-shadow: 0 0 0.45rem color-mix(in hsl, var(--color-accent-60) 14%, transparent);
	}

	.ascii-line {
		display: block;
		white-space: pre;
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

		.hero-ascii pre {
			font-size: clamp(0.5rem, 1.5vw, 0.72rem);
			line-height: 1.04;
		}
	}

	@media (max-width: 680px) {
		.hero-ascii pre {
			font-size: clamp(0.5rem, 1.95vw, 0.62rem);
			line-height: 1;
		}
	}

	@media (max-width: 800px) {
		.steps {
			grid-template-columns: 1fr;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.hero-ascii::after,
		.ascii-line {
			animation: none;
		}

		.hero-ascii::after {
			opacity: 0;
		}

		.ascii-line {
			opacity: 1;
			transform: none;
		}
	}
</style>
