<script lang="ts">
	import { navigateAfterAuthChange } from '$lib/auth/navigation';
	import { Alert, Button, Input } from '$lib/components';
	import { onMount } from 'svelte';

	let login = $state('');
	let password = $state('');
	let error = $state('');
	let loginButtonDisabled = $derived(!(login.length > 0 && password.length > 0));

	const loginAction = async () => {
		error = '';
		const res = await fetch('/api/auth/login', {
			method: 'POST',
			body: JSON.stringify({
				login,
				password
			})
		});

		if (res.ok) {
			await navigateAfterAuthChange('/dashboard');
			return;
		}

		if (res.status === 401) {
			error = 'Неверный логин или пароль';
		} else {
			error = 'Ошибка сервера';
		}
	};

	const onSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		await loginAction();
	};

	onMount(async () => {
		console.log('onMount called');

		// TODO: fix this hack
		// await new Promise((resolve) => setTimeout(resolve, 1000));

		window.YaAuthSuggest.init(
			{
				client_id: 'af7d560e82f349a5be6a69ec17e82a70',
				response_type: 'token',
				redirect_uri: 'http://localhost:5173/oauth-redirect'
			},
			'http://localhost', // TODO: different for prod
			{
				view: 'button',
				parentId: 'oauthContainer',
				buttonSize: 'm',
				buttonView: 'main',
				buttonTheme: 'light',
				buttonBorderRadius: '9',
				buttonIcon: 'ya'
			}
		)
			.then(({ handler }) => {
				console.log('got handler');
				handler();
			})
			.then((data) => console.log('Сообщение с токеном', data))
			.catch((error) => console.log('Обработка ошибки', error));

		console.log('onMount finished');
	});
</script>

<svelte:head>
	<title>Вход | Хостамат</title>

	<!-- <script
		src="https://yastatic.net/s3/passport-sdk/autofill/v1/sdk-suggest-with-polyfills-latest.js"
	>
	</script> -->
</svelte:head>

<h1>Вход</h1>
<p class="lead">Вернитесь к вашим сайтам и публикациям.</p>

{#if error}
	<Alert variant="error" message={error} />
{/if}

<div id="oauthContainer"></div>

<style>
	.lead {
		font-size: 0.95rem;
		margin-top: -0.4rem;
		margin-bottom: 0.4rem;
	}

	div#oauthContainer {
		width: 100%;
	}
</style>
