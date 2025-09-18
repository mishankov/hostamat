<script>
	import Button from "$lib/components/Button.svelte";
	import Card from "$lib/components/Card.svelte";
	import Error from "$lib/components/Error.svelte";
	import Input from "$lib/components/Input.svelte";

  let login = $state("")
  let password = $state("")
  let error = $state("")
  let registerButtonDisabled = $derived(!(login.length > 0 && password.length > 0))

  const registerAction = async () => {
    error = "";
    const res = await fetch("/api/auth/register", {
      method: "POST",
      body: JSON.stringify({
        login,
        password,
      }),
    })

    if (res.ok) {
      window.location.href = "/login";
      return;
    }

    error = "Ошибка регистрации";
  }
</script>

<svelte:window on:keydown={(e) => {
  if (e.key === "Enter" && !registerButtonDisabled) {
    registerAction();
  }
}} />

<main>
  <Card>
    <div class="form">
      <h1>Регистрация</h1>

      <Error {error} />

      <Input id="login" type="text" label="Логин" bind:value={login} />
      <Input id="password" type="password" label="Пароль" bind:value={password} />
      
      <div class="actions">
        <Button kind="action" onclick={registerAction} disabled={registerButtonDisabled}>Создать</Button>
        <Button link="/login">Войти</Button>
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

  .form {
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
