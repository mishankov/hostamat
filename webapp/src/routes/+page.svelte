<script>
	import { Badge, Button, Card } from '$lib/components';
</script>

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
		/*Full height - header*/
		height: calc(100% - 76px);
		width: 100vw;

		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		padding: 20px;
	}

	.landing {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: start;
		gap: 20px;
	}

	h1 {
		font-size: 50px;
	}

	.header {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.header-actions {
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
		gap: 20px;

		padding-top: 20px;
	}

	.actions {
		width: 100%;
		display: flex;
		flex-wrap: wrap;
		flex-direction: row;
		gap: 20px;
	}

	.action-card {
		width: 300px;
		height: 150px;

		display: flex;
		flex-direction: column;
		justify-content: start;
		align-items: start;
		gap: 10px;

		padding: 20px;
	}
</style>
