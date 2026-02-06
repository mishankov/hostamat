<script>
	import { Card, Link } from '$lib/components';
</script>

<svelte:head>
	<title>Политика конфиденциальности | Хостамат</title>
</svelte:head>

<main class="page-shell">
	<div class="page-content legal-page">
		<h1>Политика конфиденциальности</h1>

		<Card border>
			<section class="legal-section">
				<h2>1. Общие положения</h2>
				<p>
					Настоящая Политика конфиденциальности определяет порядок обработки и защиты персональных
					данных пользователей сервиса "Хостамат" (далее - Сервис).
				</p>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>2. Собираемые данные</h2>
				<p>При использовании Сервиса мы собираем следующую информацию:</p>
				<ul>
					<!-- <li>Электронная почта</li> -->
					<li>Пароль (в хэшированном виде)</li>
					<li>Информация о загруженных сайтах и файлах</li>
					<li>Технические данные (IP-адрес, браузер, время доступа)</li>
				</ul>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>3. Цели обработки</h2>
				<p>Персональные данные обрабатываются для следующих целей:</p>
				<ul>
					<li>Регистрация и авторизация пользователей</li>
					<li>Предоставление услуг хостинга статических сайтов</li>
					<li>Техническая поддержка пользователей</li>
					<li>Улучшение качества Сервиса</li>
					<li>Соблюдение требований законодательства РФ</li>
				</ul>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>4. Передача данных</h2>
				<p>
					Персональные данные пользователей не передаются третьим лицам, за исключением случаев,
					предусмотренных законодательством Российской Федерации.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>5. Безопасность данных</h2>
				<p>
					Мы принимаем необходимые организационные и технические меры для защиты персональных данных
					от несанкционированного доступа, уничтожения, изменения, блокирования, копирования,
					распространения, а также от иных неправомерных действий.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>6. Права пользователей</h2>
				<p>Пользователь имеет право:</p>
				<ul>
					<li>На доступ к своим персональным данным</li>
					<li>На исправление неточных персональных данных</li>
					<li>На блокирование и удаление своих персональных данных</li>
					<li>На отзыв согласия на обработку персональных данных</li>
				</ul>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>7. Сроки хранения данных</h2>
				<p>
					Персональные данные хранятся в течение срока, необходимого для достижения целей их
					обработки, или до момента отзыва согласия пользователем.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>8. Контактная информация</h2>
				<p>
					По вопросам, связанным с обработкой персональных данных, вы можете обратиться через <Link
						href="https://forms.yandex.ru/u/6907bbfc49af47f3cae7650a"
						target="_blank"
						alwaysUnderline>форму обратной связи</Link
					>.
				</p>
			</section>
		</Card>

		<Card border>
			<section class="legal-section">
				<h2>9. Изменения в Политике</h2>
				<p>
					Мы оставляем за собой право вносить изменения в настоящую Политику конфиденциальности.
					Новая редакция Политики вступает в силу с момента ее размещения на Сайте.
				</p>
			</section>
		</Card>
	</div>
</main>

<style>
	.legal-page {
		gap: 0.8rem;
		padding-bottom: 2rem;
	}

	.legal-page h1 {
		font-size: clamp(1.8rem, 4vw, 2.3rem);
		margin-bottom: 0.3rem;
	}

	.legal-section {
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
	}

	.legal-section h2 {
		font-size: clamp(1.2rem, 1.4vw, 1.45rem);
		line-height: 1.25;
	}

	.legal-section h3 {
		font-size: 1.05rem;
		margin-top: 0.2rem;
		margin-bottom: 0.1rem;
	}

	.legal-section p {
		line-height: 1.6;
		color: var(--color-text-secondary);
	}

	.legal-section ul,
	.legal-section ol {
		padding-left: 1.2rem;
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		line-height: 1.6;
	}

	.legal-section li {
		color: var(--color-text-secondary);
	}
</style>
