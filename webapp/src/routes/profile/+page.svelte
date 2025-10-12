<script lang="ts">
	import { IconButton } from '$lib/components';
	import { Copy } from '$lib/icons';
  import type { PageProps } from './$types';

  let { data }: PageProps = $props();
</script>

<main>
  <div class="profile">
    <h1>Профиль</h1>

    <div class="profile-data">
      <div class="profile-data-row">
        <span class="label">Логин</span>
        <div class="profile-data-value">
          <span>{data.username}</span>
          <IconButton onclick={() => {navigator.clipboard.writeText(data.username)}}>
            <Copy size={15}/>
          </IconButton>
        </div>
      </div>
    </div>
  </div>
</main>

<style>
  main {
    width: 100%;

    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
  }
  
  .profile {
    display: flex;
    flex-direction: column;
    gap: 10px;
    
    width: clamp(500px, 75%, 1200px);

    padding: 20px;
  }

  .profile-data {
    padding-top: 10px;
  }

  .profile-data-row {
    display: flex;
    flex-direction: row;
    align-items: center;

    gap: 20px;
  }

  span.label {
    color: var(--color-text-secondary);
  }

  .profile-data-value {
    display: flex;
    flex-direction: row;
    align-items: center;
  }
</style>