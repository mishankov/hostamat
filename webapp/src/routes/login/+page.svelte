<script>
	import Button from "$lib/components/Button.svelte";
	import Card from "../dashboard/components/Card.svelte";
</script>

<main>
  <Card>
    <div class="login-form">
      <h1>Привет! 👋</h1>

      <input type="text" placeholder="Логин" />
      <input type="password" placeholder="Пароль" />
      
      <Button kind="action">Войти</Button>
      <Button link="/register">Зарегистрироваться</Button>
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
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 10px;

    padding: 20px;
  }
</style>
