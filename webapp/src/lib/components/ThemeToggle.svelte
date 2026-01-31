<script lang="ts">
	import { browser } from '$app/environment';

	let theme = $state<'light' | 'dark'>('dark');

	if (browser) {
		const stored = localStorage.getItem('theme');
		if (stored === 'light' || stored === 'dark') {
			theme = stored;
		} else if (window.matchMedia('(prefers-color-scheme: light)').matches) {
			theme = 'light';
		}
		document.documentElement.setAttribute('data-theme', theme);
	}

	const toggleTheme = () => {
		theme = theme === 'dark' ? 'light' : 'dark';
		if (browser) {
			localStorage.setItem('theme', theme);
			document.documentElement.setAttribute('data-theme', theme);
		}
	};
</script>

<button
	class="theme-toggle"
	onclick={toggleTheme}
	title={theme === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'}
	aria-label={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
>
	<div class="toggle-track">
		<div class="toggle-thumb" class:light={theme === 'light'}>
			{#if theme === 'dark'}
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
				</svg>
			{:else}
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<circle cx="12" cy="12" r="5" />
					<line x1="12" y1="1" x2="12" y2="3" />
					<line x1="12" y1="21" x2="12" y2="23" />
					<line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
					<line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
					<line x1="1" y1="12" x2="3" y2="12" />
					<line x1="21" y1="12" x2="23" y2="12" />
					<line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
					<line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
				</svg>
			{/if}
		</div>
	</div>
</button>

<style>
	.theme-toggle {
		all: unset;
		cursor: pointer;
		display: flex;
		align-items: center;
	}

	.toggle-track {
		width: 48px;
		height: 26px;
		background: var(--color-primary-15);
		border-radius: 13px;
		padding: 3px;
		transition: var(--transition);
		border: 1px solid var(--glass-border);
	}

	.toggle-track:hover {
		background: var(--color-primary-20);
	}

	.toggle-thumb {
		width: 20px;
		height: 20px;
		background: var(--color-accent-50);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-primary-5);
		transition: var(--transition-bounce);
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
	}

	.toggle-thumb.light {
		transform: translateX(22px);
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
	}

	.theme-toggle:focus-visible .toggle-track {
		outline: 2px solid var(--color-accent-50);
		outline-offset: 2px;
	}
</style>
