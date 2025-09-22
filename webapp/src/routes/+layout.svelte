<script lang="ts">
	import { Button, Card, Link } from '$lib/components';
  import type { LayoutProps } from './$types';

  let { data, children }: LayoutProps = $props();

  const logout = async () => {
    const resp = await fetch('/api/auth/logout', { method: 'POST' });
    window.location.reload();
  }
</script>

<header>
  <div class="header-card">
    <div class="nav-links">
      <Link href="/"><span class="logo">ХОСТАМАТ</span></Link>
      <Link href="/about">О сервисе</Link>
    </div>

    {#if data.isAuthenticated == true}
      <Button onclick={logout}>Выйти</Button>
    {:else}
      <Button link="/login">Войти</Button>
    {/if}
  </div>
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

    background-color: var(--color-background-object);
    border-bottom: var(--border-regular);
  }

  .header-card {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    gap: 10px;

    padding: 20px;

    height: 75px;
    width: clamp(500px, 75%, 1200px);
  }

  .nav-links {
    display: flex;
    flex-direction: row;
    align-items: end;
    gap: 20px;
  }

  span.logo {
    font-size: larger;
    font-weight: bold;
  }
</style>