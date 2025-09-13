<script>
	import Button from "$lib/components/Button.svelte";
	import Card from "$lib/components/Card.svelte";
	import Input from "$lib/components/Input.svelte";
	import { redirect } from "@sveltejs/kit";

  let login = $state("")
  let password = $state("")
  let error = $state("Тест (удалить)")
  let loginButtonDisabled = $derived(!(login.length > 0 && password.length > 0))

  const loginAction = async () => {
    const res = await fetch("/api/login", {
      method: "POST",
      body: JSON.stringify({
        login,
        password,
      }),
    })

    if (res.ok) {
      window.location.href = "/dashboard";
    }

    if (res.status === 401) {
      error = "Неверный логин или пароль";
    } else {
      error = "Ошибка сервера";
    }
  }
</script>

<main>
  <Card>
    <div class="login-form">
      <h1>Привет! 👋</h1>

      {#if error || true}
        <div class="error">
          <Card kind="danger" style="height: 100%; width: 100%;">
            <div class="error-text">{error}</div>
          </Card>
        </div>
      {/if}

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
    height: 100vh;
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

  .error {
    width: 100%;
    height: 50px;
  }

  .error-text {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    
    padding: 10px;
  }

  .actions {
    width: 100%;

    display: flex;
    justify-content: space-between;
    gap: 10px;

    padding-top: 10px;
  }
</style>
