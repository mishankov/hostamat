<script>
	import Button from "$lib/components/Button.svelte";
	import Card from "$lib/components/Card.svelte";
	import Error from "$lib/components/Error.svelte";
	import Input from "$lib/components/Input.svelte";

  let login = $state("")
  let password = $state("")
  let error = $state("")
  let loginButtonDisabled = $derived(!(login.length > 0 && password.length > 0))

  const loginAction = async () => {
    error = "";
    const res = await fetch("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({
        login,
        password,
      }),
    })

    if (res.ok) {
      window.location.href = "/dashboard";
      return;
    }

    if (res.status === 401) {
      error = "Неверный логин или пароль";
    } else {
      error = "Ошибка сервера";
    }
  }
</script>

<svelte:window on:keydown={(e) => {
  if (e.key === "Enter" && !loginButtonDisabled) {
    loginAction();
  }
}} />

<main>
  <Card>
    <div class="login-form">
      <h1>Вход</h1>

      <Error {error} />

      <Input id="login" type="text" label="Логин" bind:value={login} />
      <Input id="password" type="password" label="Пароль" bind:value={password} />
      
      <div class="actions">
        <Button kind="action" onclick={loginAction} disabled={loginButtonDisabled}>Войти</Button>
        <Button link="/register">Создать аккаунт</Button>
      </div>
    </div>  
  </Card>
</main>

<style>
  main {
    height: calc(100% - 75px);
    width: 100vw;

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .login-form {
    width: 300px;

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    
    gap: 15px;

    padding: 20px;
  }

  .actions {
    width: 100%;

    display: flex;
    justify-content: space-between;
    gap: 10px;

    padding-top: 10px;
  }
</style>
