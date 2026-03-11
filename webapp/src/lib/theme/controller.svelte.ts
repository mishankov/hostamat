import { createContext } from 'svelte';

export type Theme = 'light' | 'dark';

const THEME_STORAGE_KEY = 'hostamat-theme';

const resolveDocumentTheme = (): Theme => {
	if (typeof document === 'undefined') {
		return 'light';
	}

	return document.documentElement.getAttribute('data-theme') === 'dark' ? 'dark' : 'light';
};

export class ThemeController {
	current = $state<Theme>('light');

	constructor(initialTheme: Theme = 'light') {
		this.current = initialTheme;
	}

	toggle() {
		this.current = this.current === 'light' ? 'dark' : 'light';
		this.sync();
	}

	initializeFromBrowser() {
		this.current = this.resolvePreferredTheme();
		this.sync();
	}

	private resolvePreferredTheme(): Theme {
		const fallbackTheme = this.resolveSystemTheme();

		try {
			const savedTheme = localStorage.getItem(THEME_STORAGE_KEY);
			if (savedTheme === 'dark' || savedTheme === 'light') {
				return savedTheme;
			}
		} catch {
			return fallbackTheme;
		}

		return fallbackTheme;
	}

	private resolveSystemTheme(): Theme {
		if (typeof window === 'undefined') {
			return resolveDocumentTheme();
		}

		return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
	}

	private sync() {
		if (typeof document !== 'undefined') {
			document.documentElement.setAttribute('data-theme', this.current);
		}

		try {
			localStorage.setItem(THEME_STORAGE_KEY, this.current);
		} catch {
			// Ignore storage write failures.
		}
	}
}

const [useThemeController, setThemeController] = createContext<ThemeController>();

export const provideThemeController = (controller: ThemeController) => setThemeController(controller);
export { useThemeController, resolveDocumentTheme };
