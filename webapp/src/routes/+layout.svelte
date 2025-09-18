<script lang="ts">
	import { Button, Card, Link } from '$lib/components';
	import Separator from '$lib/components/Separator.svelte';
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
      <Link href="/"><span class="logo">ХОСТАМАТ</span></Link>
      <Link href="/about">О сервисе</Link>

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
    width: clamp(500px, 50vw, 1200px);
  }

  span.logo {
    font-size: larger;
    font-weight: bold;
  }
</style>