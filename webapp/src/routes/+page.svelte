<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Хостамат</title>
</svelte:head>

<main>
	<div class="landing">
		<div class="header">
			<h1>ХОСТАМАТ</h1>
			<h2>Хостинг статических сайтов</h2>
			<Badge text="Альфа-тестирование" />

			<div class="header-actions">
				<Button link="/register" filled>Загрузить свой сайт</Button>
			</div>
		</div>

		<div class="actions">
			<Card floating>
				<div class="action-card">
					<h2>Создай аккаунт</h2>
					<p>Регистрация займёт меньше минуты</p>
				</div>
			</Card>

			<Card floating>
				<div class="action-card">
					<h2>Загрузи свой сайт</h2>
					<p>Можно загрузить отдельные файлы или архив</p>
				</div>
			</Card>

			<Card floating>
				<div class="action-card">
					<h2>Получи ссылку</h2>
					<p>Которую сможет посмотреть кто угодно</p>
				</div>
			</Card>

			<Card floating>
				<div class="action-card">
					<h2>5 бесплатных сайтов</h2>
					<p>На время альфа-тестирования и после него</p>
					<Badge text="Альфа-тестирование" />
				</div>
			</Card>
		</div>
	</div>
</main>

<style>
	main {
		min-height: calc(100vh - 85px);
		width: 100vw;

		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		padding: 40px 20px;
		position: relative;
		overflow: hidden;
	}

	main::before {
		content: '';
		position: absolute;
		top: 20%;
		right: -10%;
		width: 600px;
		height: 600px;
		background: radial-gradient(circle, var(--color-accent-20) 0%, transparent 70%);
		opacity: 0.3;
		pointer-events: none;
		animation: float 20s ease-in-out infinite;
	}

	@keyframes float {
		0%, 100% {
			transform: translate(0, 0) rotate(0deg);
		}
		33% {
			transform: translate(30px, -30px) rotate(5deg);
		}
		66% {
			transform: translate(-20px, 20px) rotate(-5deg);
		}
	}

	.landing {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: start;
		gap: 48px;
		max-width: 1400px;
		width: 100%;
		z-index: 1;
	}

	h1 {
		font-size: clamp(48px, 8vw, 96px);
		line-height: 0.95;
		margin: 0;
		background: linear-gradient(135deg, var(--color-text-regular), var(--color-accent-50));
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
		position: relative;
		animation: slideInLeft 0.6s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes slideInLeft {
		from {
			transform: translateX(-50px);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	h2 {
		font-size: clamp(20px, 3vw, 28px);
		font-weight: 600;
		color: var(--color-text-secondary);
		margin-top: -12px;
		animation: slideInLeft 0.6s cubic-bezier(0.4, 0, 0.2, 1) 0.1s backwards;
	}

	.header {
		display: flex;
		flex-direction: column;
		gap: 20px;
		position: relative;
	}

	.header::after {
		content: '';
		position: absolute;
		left: -20px;
		top: 0;
		bottom: 0;
		width: 6px;
		background: linear-gradient(180deg, var(--color-accent-50), var(--color-accent-70));
		border-radius: 3px;
	}

	.header-actions {
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
		gap: 16px;
		padding-top: 12px;
		animation: slideInLeft 0.6s cubic-bezier(0.4, 0, 0.2, 1) 0.2s backwards;
	}

	.actions {
		width: 100%;
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 24px;
		animation: fadeInUp 0.8s cubic-bezier(0.4, 0, 0.2, 1) 0.3s backwards;
	}

	@keyframes fadeInUp {
		from {
			transform: translateY(30px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	.action-card {
		min-height: 180px;
		display: flex;
		flex-direction: column;
		justify-content: start;
		align-items: start;
		gap: 12px;
		padding: 28px;
		position: relative;
	}

	.action-card h2 {
		font-size: 22px;
		margin: 0;
		color: var(--color-text-regular);
		animation: none;
	}

	.action-card p {
		font-size: 14px;
		line-height: 1.6;
		color: var(--color-text-secondary);
	}

	.actions > :nth-child(1) {
		animation-delay: 0.35s;
	}

	.actions > :nth-child(2) {
		animation-delay: 0.4s;
	}

	.actions > :nth-child(3) {
		animation-delay: 0.45s;
	}

	.actions > :nth-child(4) {
		animation-delay: 0.5s;
	}
</style>
