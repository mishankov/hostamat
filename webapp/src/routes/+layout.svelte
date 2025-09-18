<script lang="ts">
	import { Button } from '$lib/components';
  import type { LayoutProps } from './$types';

  let { data, children }: LayoutProps = $props();

  const logout = async () => {
    const resp = await fetch('/api/auth/logout', { method: 'POST' });
    window.location.reload();
  }
</script>

<header>
  <a href="/">ХОСТАМАТ</a>

  <a href="/about">О сервисе</a>

  {#if data.isAuthenticated == true}
    <Button onclick={logout}>Выйти</Button>
  {:else}
    <Button link="/login">Войти</Button>
  {/if}
</header>

{@render children()}

<style>
  header {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    gap: 10px;

    padding: 10px;
  }
</style>