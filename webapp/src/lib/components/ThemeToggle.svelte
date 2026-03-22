<script lang="ts">
	import { Moon, Sun } from '$lib/icons';
	import { useThemeController } from '$lib/theme/controller.svelte';
	import IconButton from './IconButton.svelte';

	const theme = useThemeController();
	const label = $derived(
		theme.current === 'light' ? 'Включить тёмную тему' : 'Включить светлую тему'
	);
</script>

<IconButton {label} title={label} onclick={() => theme.toggle()}>
	{#if theme.current === 'light'}
		<span class="theme-icon" aria-hidden="true">
			<Sun size={18} color="var(--color-accent-strong)" />
		</span>
	{:else}
		<span class="theme-icon" aria-hidden="true">
			<Moon size={18} color="var(--color-accent-strong)" />
		</span>
	{/if}
</IconButton>

<style>
	.theme-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}
</style>
