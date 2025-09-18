<script lang="ts">
	import { Button, Card } from '$lib/components';
  import type { LayoutProps } from './$types';

  let { data, children }: LayoutProps = $props();

  const logout = async () => {
    const resp = await fetch('/api/auth/logout', { method: 'POST' });
    window.location.reload();
  }
</script>

<header>
  <Card>
    <div class="header-card">
      <a href="/">ХОСТАМАТ</a>

      <a href="/about">О сервисе</a>

      {#if data.isAuthenticated == true}
        <Button onclick={logout}>Выйти</Button>
      {:else}
        <Button link="/login">Войти</Button>
      {/if}
    </div>
  </Card>
</header>

{@render children()}

<style>
  header {
    width: 100%;
    height: fit-content;

    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
  }

  .header-card {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    gap: 10px;

    padding: 10px;

    height: 75px;
  }
</style>