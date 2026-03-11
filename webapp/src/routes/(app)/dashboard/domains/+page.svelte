<script lang="ts">
	import { Badge, Card } from '$lib/components';
</script>

<svelte:head>
	<title>Домены | Хостамат</title>
</svelte:head>

<div class="domains-page">
	<div class="headline">
		<div class="headline-row">
			<h1>Домены</h1>
			<Badge text="Soon" />
		</div>
		<p>Просматривайте и управляете доменами и поддоменами</p>
	</div>

	<Card border cornerBorder class="domains-card">
		<div class="domains-card-content">
			<h2>Как это работает сейчас</h2>
			<p>
				Каждый статический сайт в Хостамат уже получает собственный поддомен автоматически после
				создания.
			</p>
			<p>
				Когда раздел доменов будет запущен, здесь появится управление подключением ваших доменов к
				существующим сайтам.
			</p>
		</div>
	</Card>
</div>

<style>
	.domains-page {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		min-width: 0;
	}

	.headline {
		display: flex;
		flex-direction: column;
		gap: 0.45rem;
	}

	.headline-row {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.75rem;
	}

	:global(.domains-card) {
		box-shadow: var(--shadow-sm);
	}

	.domains-card-content {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		padding: 1.15rem;
		max-width: 66ch;
	}

	.domains-card-content h2 {
		font-size: clamp(1.15rem, 1.7vw, 1.45rem);
	}
</style>
