<script lang="ts">
  let {
    uploadProgress = $bindable(0)
   } = $props()
</script>

<div class="progress-bar">
  <div class="progress-fill" style="width: {uploadProgress}%"></div>
</div>

<style>
  .progress-bar {
    width: 100%;
    height: 30px;
    
    background-color: var(--color-primary-dim);
    
    border: var(--border-regular);
    border-radius: var(--border-radius-regular);
  }
  
  .progress-fill {
    height: 100%;
    background-color: var(--color-primary);
    transition: var(--transition);
  }
</style>