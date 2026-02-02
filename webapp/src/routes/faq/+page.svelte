<script>
	import { resolve } from '$app/paths';
	import { Link } from '$lib/components';
</script>

<svelte:head>
	<title>FAQ | Хостамат</title>
</svelte:head>

<main>
	<div class="faq">
		<h1>FAQ</h1>

		<h2>
			Я загрузил сайт, но по адресу https://xxxx-yyyy-1234.hostamat.ru/ возвращается страница с
			ошибкой 404.
		</h2>

		<p>
			<strong>Самая частая причина:</strong> файл
			<span class="path">index.html</span> находится не в корне архива, а внутри вложенной папки. При
			создании архива убедитесь, что вы добавляете файлы напрямую, а не папку с ними.
		</p>

		<p>Хостамат следует нескольким правилам при поиске файлов для отображения:</p>

		<ul>
			<li>
				<span class="path">/</span> (корневой путь) - для отображения корневого пути должен быть
				файл
				<span class="path">index.html</span>, либо загруженый отдельно, либо лежащий в корне
				загруженного архива
			</li>
			<li>
				<span class="path">/aaa</span> (вложенный путь) - для отображения вложенного пути идёт поиск
				сначала файла <span class="path">/aaa</span>, потом <span class="path">/aaa.html</span>, и
				после него <span class="path">/aaa/index.html</span>. Если ни одного из этих файлов не
				найдено, возвращается ошибка 404
			</li>
		</ul>

		<h2>Какой максимальный размер загружаемого файла?</h2>
		<p>Максимальный размер загружаемого архива или отдельного файла — 20 МБ.</p>

		<h2>Как долго хранится загруженный сайт?</h2>
		<p>Загруженные сайты хранятся бессрочно, пока вы не удалите их вручную.</p>

		<h2>Можно ли обновить или удалить загруженный сайт?</h2>
		<p>
			Да, вы можете обновить или удалить свой сайт на <Link href={resolve('/dashboard')}
				>Дашборде</Link
			>.
		</p>

		<h2>Какие форматы файлов поддерживаются?</h2>
		<p>
			Информация о поддерживаемых форматах доступна на странице <Link href={resolve('/about')}
				>«О сервисе»</Link
			>. Архивы принимаются только в формате <span class="path">.zip</span>.
		</p>
	</div>
</main>

<style>
	main {
		width: 100%;
		height: calc(100% - 76px);

		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: start;
	}

	.faq {
		display: flex;
		flex-direction: column;
		gap: 10px;

		width: clamp(500px, 75%, 1200px);

		padding: 20px;
	}

	h2 {
		/*font-weight: 400;*/
		padding-top: 10px;
	}

	p,
	li {
		font-size: 1.25rem;
		line-height: 1.6;
		color: var(--color-text-secondary);
		max-width: 600px;
	}

	span.path {
		font-family: monospace;
		color: var(--color-text-regular);
	}
</style>
