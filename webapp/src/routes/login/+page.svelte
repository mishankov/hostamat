<script>
	import Button from "$lib/components/Button.svelte";
	import Card from "$lib/components/Card.svelte";
	import Input from "$lib/components/Input.svelte";
	import { redirect } from "@sveltejs/kit";

  let login = $state("")
  let password = $state("")

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
  }
</script>

<main>
  <Card>
    <div class="login-form">
      <h1>Привет! 👋</h1>

      <Input id="login" type="text" label="Логин" bind:value={login} />
      <Input id="password" type="password" label="Пароль" bind:value={password} />
      
      <div class="actions">
        <Button kind="action" onclick={loginAction}>Войти</Button>
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
    
    gap: 10px;

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
