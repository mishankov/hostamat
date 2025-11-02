<script>
	import { Link } from '$lib/components';
</script>

<main>
	<div class="about">
		<h1>О сервисе</h1>
		<p>
			Хостамат (<Link href="hostamat.ru" target="_blank" alwaysUnderline>hostamat.ru</Link>) —
			сервис для простого хостинга статических сайтов. Для каждого сайта выделяется домен третьего
			уровня и настроенным TLS, то есть каждый сайт доступен по HTTPS.
		</p>

		<h2>Доступные варианты загрузки файлов</h2>

		<p>
			Сайты можно загружать напрямую на сайте сервиса в виде отдельных файлов или .zip архива. Для
			загрузки доступны следующие форматы файлов:
		</p>

		<ul>
			<li>html</li>
			<li>css</li>
			<li>js</li>
			<li>png</li>
			<li>jpg</li>
			<li>jpeg</li>
			<li>svg</li>
			<li>zip</li>
		</ul>
	</div>
</main>

<style>
	main {
		width: 100%;

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
	}

	.about {
		display: flex;
		flex-direction: column;
		gap: 20px;

		width: clamp(500px, 75%, 1200px);

		padding: 20px;
	}

	p {
		line-height: 20px;
	}
</style>
