<script lang="ts">
	import { Card } from '$lib/components';
    import type { LayoutProps } from './$types';

    let { children }: LayoutProps = $props();
</script>

<main>
  <Card>
    <div class="form">
      {@render children()}
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

  .form {
    width: 300px;

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    
    gap: 15px;

    padding: 20px;
  }
</style>

