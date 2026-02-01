<script lang="ts">
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';

	let theme = $state('light');

	onMount(() => {
		const savedTheme = localStorage.getItem('theme') || 'light';
		theme = savedTheme;
		document.documentElement.setAttribute('data-theme', savedTheme);
	});

	function toggleTheme() {
		const newTheme = theme === 'light' ? 'dark' : 'light';
		theme = newTheme;
		if (browser) {
			localStorage.setItem('theme', newTheme);
			document.documentElement.setAttribute('data-theme', newTheme);
		}
	}
</script>

<button onclick={toggleTheme} class="theme-toggle" aria-label="Toggle theme">
	{#if theme === 'light'}
		<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
			<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
		</svg>
	{:else}
		<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
			<circle cx="12" cy="12" r="5"></circle>
			<line x1="12" y1="1" x2="12" y2="3"></line>
			<line x1="12" y1="21" x2="12" y2="23"></line>
			<line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
			<line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
			<line x1="1" y1="12" x2="3" y2="12"></line>
			<line x1="21" y1="12" x2="23" y2="12"></line>
			<line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
			<line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
		</svg>
	{/if}
</button>

<style>
	.theme-toggle {
		appearance: none;
		background: transparent;
		border: var(--border-regular);
		border-radius: var(--border-radius-regular);
		padding: 10px;
		cursor: pointer;
		color: var(--color-text-regular);
		transition: var(--transition);
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
	}

	.theme-toggle:hover {
		background: var(--color-primary-10);
		transform: translateY(-2px);
		box-shadow: var(--shadow);
	}

	.theme-toggle:active {
		transform: translateY(0);
		box-shadow: none;
	}

	svg {
		display: block;
	}
</style>
