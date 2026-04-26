<script lang="ts">
	import { Badge, Button } from '$lib/components';
	import { useThemeController } from '$lib/theme/controller.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const theme = useThemeController();
</script>

<svelte:head>
	<title>Хостамат</title>
	<meta name="description" content="Хостинг статических сайтов с автоматическими HTTPS-поддоменами. Загрузите файлы — получите готовый сайт." />
</svelte:head>

<main class="page-shell">
	<section class={['page-content', 'landing', theme.current === 'dark' && 'theme-dark']}>
		<div class="hero">
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

		<div class="steps">
			<div class="step">
				<span class="step-number">01</span>
				<div class="step-body">
					<h2>Войдите через Яндекс</h2>
					<p>Авторизация занимает пару кликов.</p>
				</div>
			</div>
			<div class="step">
				<span class="step-number">02</span>
				<div class="step-body">
					<h2>Загрузите файлы</h2>
					<p>Можно загрузить отдельные файлы или ZIP-архив.</p>
				</div>
			</div>
			<div class="step">
				<span class="step-number">03</span>
				<div class="step-body">
					<h2>Поделитесь ссылкой</h2>
					<p>Сайт сразу доступен по HTTPS-поддомену.</p>
				</div>
			</div>
			<div class="step step-featured">
				<span class="step-number">04</span>
				<div class="step-body">
					<h2>До 5 сайтов бесплатно</h2>
					<p>Во время альфа-тестирования до 5 сайтов доступны каждому пользователю.</p>
				</div>
			</div>
		</div>
	</section>
</main>

<style>
	.landing {
		position: relative;
		z-index: 1;
		gap: var(--space-2xl);
		padding-bottom: var(--space-2xl);
	}

	.hero {
		border: var(--border-regular);
		background: var(--color-background-object);
	}

	.hero-main {
		display: grid;
		grid-template-columns: minmax(31.25rem, 1.2fr) minmax(20rem, 0.95fr);
		align-items: center;
		gap: clamp(1.5rem, 3vw, 3rem);
		padding: clamp(var(--space-xl), 3vw, var(--space-2xl));
	}

	.hero-copy {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 0;
	}

	h1 {
		font-size: clamp(2.2rem, 5vw, 3.6rem);
		line-height: 1.1;
		max-width: 18ch;
		font-weight: 800;
	}

	.hero-description {
		max-width: 44ch;
		font-size: 1.05rem;
		line-height: 1.6;
	}

	.hero-ascii {
		margin-top: 0.1rem;
		position: relative;
		inline-size: 100%;
	}

	.terminal-window {
		--terminal-window-border: color-mix(
			in hsl,
			var(--color-border) 70%,
			var(--color-accent-30) 30%
		);
		--terminal-window-bg: color-mix(
			in hsl,
			var(--color-background-object) 98%,
			var(--color-background)
		);
		--terminal-window-shadow: var(--shadow-lg);
		--terminal-titlebar-bg-top: color-mix(
			in hsl,
			var(--color-background-object) 96%,
			var(--color-accent-20) 4%
		);
		--terminal-titlebar-bg-bottom: color-mix(
			in hsl,
			var(--color-background-object) 98%,
			var(--color-background)
		);
		--terminal-title-color: color-mix(in hsl, var(--color-accent-80) 70%, var(--color-primary-80));
		--terminal-body-bg: color-mix(
			in hsl,
			var(--color-background-object) 97%,
			var(--color-accent-10) 3%
		);
		--terminal-text: color-mix(in hsl, var(--color-accent-90) 75%, var(--color-primary-90));
		--terminal-text-dim: color-mix(in hsl, var(--color-accent-70) 80%, var(--color-primary-80));
		--terminal-line-color: color-mix(in hsl, var(--color-accent-60) 55%, transparent);
		--terminal-brand-color: color-mix(in hsl, var(--color-accent-80) 82%, var(--color-background));
		--terminal-brand-stroke: color-mix(in hsl, var(--color-accent-60) 24%, transparent);
		--terminal-status: color-mix(in hsl, var(--color-confirm-60) 65%, var(--color-accent-70));
		--terminal-scan: color-mix(in hsl, var(--color-accent-40) 32%, transparent);
		--terminal-control-close: #ff5f57;
		--terminal-control-minimize: #febc2e;
		--terminal-control-zoom: #28c840;
		position: relative;
		overflow: hidden;
		border-radius: var(--radius-lg);
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
			var(--color-border) 60%,
			var(--color-accent-20) 40%
		);
		--terminal-window-bg: color-mix(
			in hsl,
			var(--color-background-object) 90%,
			var(--color-primary-10) 10%
		);
		--terminal-window-shadow: var(--shadow-xl);
		--terminal-titlebar-bg-top: color-mix(
			in hsl,
			var(--color-background-object) 85%,
			var(--color-primary-20) 15%
		);
		--terminal-titlebar-bg-bottom: color-mix(
			in hsl,
			var(--color-background-object) 90%,
			var(--color-primary-10) 10%
		);
		--terminal-title-color: color-mix(in hsl, var(--color-accent-70) 75%, var(--color-primary-70));
		--terminal-body-bg: color-mix(in hsl, var(--color-background-object) 88%, black 12%);
		--terminal-text: color-mix(in hsl, var(--color-accent-80) 82%, var(--color-background));
		--terminal-text-dim: color-mix(in hsl, var(--color-accent-60) 80%, var(--color-primary-70));
		--terminal-line-color: color-mix(in hsl, var(--color-accent-50) 55%, transparent);
		--terminal-brand-color: color-mix(in hsl, var(--color-accent-80) 88%, var(--color-background));
		--terminal-brand-stroke: color-mix(in hsl, var(--color-accent-70) 32%, transparent);
		--terminal-status: color-mix(in hsl, var(--color-confirm-60) 60%, var(--color-accent-70));
		--terminal-scan: color-mix(in hsl, var(--color-accent-40) 36%, transparent);
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
		gap: var(--space-md);
		padding-top: var(--space-sm);
	}

	.steps {
		display: grid;
		grid-template-columns: repeat(4, minmax(0, 1fr));
		border-top: var(--border-regular);
	}

	.step {
		padding: var(--space-xl);
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		border-left: var(--border-regular);
	}

	.step:first-child {
		border-left: none;
	}

	.step-featured {
		background: var(--color-accent-surface);
	}

	.step-number {
		font-family: var(--font-display);
		font-size: 0.8rem;
		font-weight: 700;
		letter-spacing: 0.1em;
		color: var(--color-accent-strong);
	}

	.step-body {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}

	.step-body h2 {
		font-size: 1.3rem;
		font-weight: 700;
	}

	.step-body p {
		max-width: 32ch;
		font-size: 0.92rem;
		line-height: 1.5;
	}

	@media (max-width: 920px) {
		.hero-main {
			grid-template-columns: 1fr;
			gap: var(--space-md);
		}

		.hero-terminal {
			font-size: clamp(0.56rem, 1.5vw, 0.78rem);
			line-height: 1.04;
		}

		.terminal-brand {
			font-size: clamp(1.64rem, 8vw, 3rem);
		}
	}

	@media (max-width: 800px) {
		.steps {
			grid-template-columns: 1fr;
		}

		.step {
			border-left: none;
			border-top: var(--border-regular);
			padding: var(--space-lg);
		}

		.step:first-child {
			border-top: none;
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
