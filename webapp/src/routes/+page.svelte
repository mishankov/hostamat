<script lang="ts">
	import { Badge, Button, Card } from '$lib/components';
	import { useThemeController } from '$lib/theme/controller.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const theme = useThemeController();
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main class="page-shell">
	<section class={['page-content', 'landing', theme.current === 'dark' && 'theme-dark']}>
		<Card surface="chrome" class="hero-shell">
			<div class="hero-card">
				<div class="hero-main">
					<div class="hero-copy">
						<Badge text="Альфа-тестирование" />
						<h1>Хостинг статических сайтов без лишней настройки</h1>
						<p class="hero-description">
							Загрузите файлы или ZIP-архив, получите HTTPS-поддомен и управляйте проектом в удобной
							панели.
						</p>
						<div class="hero-actions">
							{#if data.isAuthenticated}
								<Button href="/dashboard/sites" emphasis="solid">Открыть дашборд</Button>
							{:else}
								<Button href="/dashboard/sites" emphasis="solid">Запустить свой сайт</Button>
							{/if}
							<Button href="/about">Узнать о платформе</Button>
						</div>
					</div>
					<div class="hero-ascii">
						<div class="terminal-window" aria-hidden="true">
							<div class="terminal-titlebar">
								<div class="terminal-controls">
									<span class="terminal-control is-close"></span>
									<span class="terminal-control is-minimize"></span>
									<span class="terminal-control is-zoom"></span>
								</div>
								<p class="terminal-title">Terminal — deploy@hostamat</p>
							</div>
							<div class="terminal-body">
								<div class="hero-terminal">
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
					</div>
				</div>
			</div>
		</Card>

		<div class="steps">
			<Card class="step-card">
				<div class="step">
					<span>01</span>
					<h2>Войдите через Яндекс</h2>
					<p>Авторизация занимает пару кликов.</p>
				</div>
			</Card>
			<Card class="step-card">
				<div class="step">
					<span>02</span>
					<h2>Загрузите файлы</h2>
					<p>Можно загрузить отдельные файлы или ZIP-архив.</p>
				</div>
			</Card>
			<Card class="step-card">
				<div class="step">
					<span>03</span>
					<h2>Поделитесь ссылкой</h2>
					<p>Сайт сразу доступен по HTTPS-поддомену.</p>
				</div>
			</Card>
			<Card class="step-card step-card-featured">
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

	.hero-card {
		padding: clamp(1.2rem, 2.5vw, 2.2rem);
	}

	:global(.hero-shell) {
		box-shadow: var(--shadow-sm);
	}

	.hero-main {
		display: grid;
		grid-template-columns: minmax(500px, 1.2fr) minmax(320px, 0.95fr);
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
		inline-size: 100%;
	}

	.terminal-window {
		--terminal-window-border: color-mix(
			in hsl,
			var(--color-border) 82%,
			var(--color-accent-30) 18%
		);
		--terminal-window-bg: color-mix(
			in hsl,
			var(--color-background-object) 97%,
			var(--color-background)
		);
		--terminal-window-shadow: 0 10px 20px color-mix(in hsl, var(--color-accent-30) 10%, transparent);
		--terminal-titlebar-bg-top: color-mix(
			in hsl,
			var(--color-background-object) 94%,
			var(--color-accent-20) 6%
		);
		--terminal-titlebar-bg-bottom: color-mix(
			in hsl,
			var(--color-background-object) 97%,
			var(--color-background)
		);
		--terminal-title-color: color-mix(in hsl, var(--color-accent-80) 72%, var(--color-primary-80));
		--terminal-body-bg: color-mix(
			in hsl,
			var(--color-background-object) 95%,
			var(--color-accent-10) 5%
		);
		--terminal-text: color-mix(in hsl, var(--color-accent-90) 78%, var(--color-primary-90));
		--terminal-text-dim: color-mix(in hsl, var(--color-accent-70) 82%, var(--color-primary-80));
		--terminal-line-color: color-mix(in hsl, var(--color-accent-60) 60%, transparent);
		--terminal-brand-color: color-mix(in hsl, var(--color-accent-80) 84%, var(--color-background));
		--terminal-brand-stroke: color-mix(in hsl, var(--color-accent-60) 28%, transparent);
		--terminal-status: color-mix(in hsl, var(--color-confirm-60) 70%, var(--color-accent-70));
		--terminal-scan: color-mix(in hsl, var(--color-accent-40) 38%, transparent);
		--terminal-control-close: #ff5f57;
		--terminal-control-minimize: #febc2e;
		--terminal-control-zoom: #28c840;
		position: relative;
		overflow: hidden;
		border-radius: 14px;
		border: 1px solid var(--terminal-window-border);
		background: var(--terminal-window-bg);
		box-shadow: var(--terminal-window-shadow);
	}

	.terminal-titlebar {
		position: relative;
		display: flex;
		align-items: center;
		height: 2rem;
		padding: 0 0.72rem;
		border-bottom: 1px solid var(--terminal-window-border);
		background: linear-gradient(
			180deg,
			var(--terminal-titlebar-bg-top) 0%,
			var(--terminal-titlebar-bg-bottom) 100%
		);
	}

	.terminal-controls {
		display: flex;
		align-items: center;
		gap: 0.38rem;
	}

	.terminal-control {
		display: inline-block;
		inline-size: 0.72rem;
		block-size: 0.72rem;
		border-radius: 50%;
		box-shadow: 0 0 0 1px color-mix(in hsl, var(--terminal-window-border) 82%, transparent);
	}

	.terminal-control.is-close {
		background: var(--terminal-control-close);
	}

	.terminal-control.is-minimize {
		background: var(--terminal-control-minimize);
	}

	.terminal-control.is-zoom {
		background: var(--terminal-control-zoom);
	}

	.terminal-title {
		position: absolute;
		left: 50%;
		top: 50%;
		max-inline-size: calc(100% - 7rem);
		transform: translate(-50%, -50%);
		margin: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-family:
			ui-sans-serif,
			-apple-system,
			system-ui,
			sans-serif;
		font-size: 0.68rem;
		font-weight: 520;
		letter-spacing: 0.01em;
		color: var(--terminal-title-color);
	}

	.terminal-body {
		position: relative;
		overflow: hidden;
		background: var(--terminal-body-bg);
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
		color: var(--terminal-text);
		text-shadow: none;
	}

	.terminal-body::after {
		content: '';
		position: absolute;
		inset-inline: 0;
		top: -22%;
		height: 18%;
		pointer-events: none;
		background: linear-gradient(180deg, transparent 0%, var(--terminal-scan) 50%, transparent 100%);
		opacity: 0;
		animation: asciiScanPass 1100ms ease-out 580ms 1 both;
	}

	.terminal-head {
		display: flex;
		align-items: center;
		gap: 0.48rem;
	}

	.terminal-head-label {
		color: var(--terminal-text);
		font-weight: 600;
		white-space: nowrap;
	}

	.terminal-head-rule {
		flex: 1;
		height: 1px;
		background: repeating-linear-gradient(
			90deg,
			var(--terminal-line-color) 0 8px,
			transparent 8px 12px
		);
	}

	.terminal-brand {
		margin: 0;
		padding-top: 0.04em;
		font-family: var(--font-body);
		font-size: clamp(1.82rem, 4.3vw, 4.3rem);
		font-weight: 700;
		line-height: 0.88;
		letter-spacing: 0.05em;
		color: var(--terminal-brand-color);
		text-shadow: none;
		-webkit-text-stroke: 1px var(--terminal-brand-stroke);
	}

	.terminal-divider {
		margin: 0.1rem 0 0.05rem;
		border-top: 1px dashed var(--terminal-line-color);
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
		color: var(--terminal-text-dim);
	}

	.terminal-step-text {
		overflow-wrap: anywhere;
	}

	.terminal-step-fill {
		border-bottom: 1px dotted var(--terminal-line-color);
		transform: translateY(-0.08rem);
	}

	.terminal-step-status {
		color: var(--terminal-status);
	}

	.terminal-meta-row {
		margin: 0;
		display: flex;
		flex-wrap: wrap;
		gap: 0 0.36rem;
	}

	.terminal-meta-key {
		color: var(--terminal-text-dim);
	}

	.landing.theme-dark .terminal-window {
		--terminal-window-border: color-mix(
			in hsl,
			var(--color-border) 64%,
			var(--color-accent-20) 36%
		);
		--terminal-window-bg: color-mix(
			in hsl,
			var(--color-background-object) 88%,
			var(--color-primary-10) 12%
		);
		--terminal-window-shadow: 0 14px 28px color-mix(in hsl, var(--color-accent-30) 18%, transparent);
		--terminal-titlebar-bg-top: color-mix(
			in hsl,
			var(--color-background-object) 80%,
			var(--color-primary-20) 20%
		);
		--terminal-titlebar-bg-bottom: color-mix(
			in hsl,
			var(--color-background-object) 86%,
			var(--color-primary-10) 14%
		);
		--terminal-title-color: color-mix(in hsl, var(--color-accent-70) 80%, var(--color-primary-70));
		--terminal-body-bg: color-mix(in hsl, var(--color-background-object) 84%, black 16%);
		--terminal-text: color-mix(in hsl, var(--color-accent-80) 84%, var(--color-background));
		--terminal-text-dim: color-mix(in hsl, var(--color-accent-60) 84%, var(--color-primary-70));
		--terminal-line-color: color-mix(in hsl, var(--color-accent-50) 62%, transparent);
		--terminal-brand-color: color-mix(in hsl, var(--color-accent-80) 90%, var(--color-background));
		--terminal-brand-stroke: color-mix(in hsl, var(--color-accent-70) 38%, transparent);
		--terminal-status: color-mix(in hsl, var(--color-confirm-60) 64%, var(--color-accent-70));
		--terminal-scan: color-mix(in hsl, var(--color-accent-40) 42%, transparent);
		--terminal-control-close: #ff5f57;
		--terminal-control-minimize: #febc2e;
		--terminal-control-zoom: #28c840;
	}

	.terminal-meta-value {
		overflow-wrap: anywhere;
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

	:global(.step-card) {
		box-shadow: var(--shadow-sm);
		transition:
			transform 180ms ease-out,
			box-shadow 180ms ease-out;
	}

	:global(.step-card:hover) {
		box-shadow: var(--shadow-lg);
		transform: translateY(-2px);
	}

	:global(.step-card-featured) {
		background: var(--color-accent-surface);
		border-color: var(--color-accent-border);
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
		.terminal-titlebar {
			height: 1.8rem;
			padding-inline: 0.58rem;
		}

		.terminal-control {
			inline-size: 0.66rem;
			block-size: 0.66rem;
		}

		.terminal-title {
			max-inline-size: calc(100% - 6rem);
			font-size: 0.63rem;
		}

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
		.terminal-body::after,
		.terminal-line {
			animation: none;
		}

		.terminal-body::after {
			opacity: 0;
		}

		.terminal-line {
			opacity: 1;
			transform: none;
		}
	}
</style>
