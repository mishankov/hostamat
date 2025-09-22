<script>
	import Button from "$lib/components/Button.svelte";
</script>


<main>
	<h1>ХОСТАМАТ</h1>
	<h2>Хостинг статических сайтов</h2>

	<Button link="/dashboard">На дашборд</Button>
</main>

<style>
	main {
		height: calc(100% - 75px);
		width: 100vw;

		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 10px;
	}
</style>