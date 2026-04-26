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
		<span class="theme-icon" aria-hidden="true" data-icon="sun">
			<Sun size={18} color="var(--color-accent-strong)" />
		</span>
	{:else}
		<span class="theme-icon" aria-hidden="true" data-icon="moon">
			<Moon size={18} color="var(--color-accent-strong)" />
		</span>
	{/if}
</IconButton>

<style>
	.theme-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		animation: theme-icon-in 220ms var(--ease-out-quart) both;
	}

	@keyframes theme-icon-in {
		from {
			opacity: 0;
			transform: rotate(-30deg) scale(0.9);
		}
		to {
			opacity: 1;
			transform: rotate(0deg) scale(1);
		}
	}
</style>
