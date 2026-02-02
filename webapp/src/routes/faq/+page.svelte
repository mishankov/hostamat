<svelte:head>
	<title>FAQ| Хостамат</title>
</svelte:head>

<main>
	<div class="faq">
		<h1>FAQ</h1>

		<h2>
			Я загрузил сайт, но по адресу https://xxxx-yyyy-1234.hostamat.ru/ возвращается страница с
			ошибкой 404.
		</h2>

		<p>Хостамат следуем нескольким правилам при поиске файлов для отображения:</p>

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
				найдено, возвращаеться ошибка 404
			</li>
		</ul>
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
