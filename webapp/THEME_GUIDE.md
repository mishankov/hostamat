# Theme System Guide

## Overview

The webapp now supports both light and dark themes with a cohesive design language across both modes.

## Using the Theme Toggle

Users can switch between light and dark themes using the theme toggle button in the header (sun/moon icon). The preference is saved in localStorage and persists across sessions.

## Theme Colors Reference

### Light Theme (Default)
- **Background**: Warm beige/cream tones
- **Text**: Dark brown/black
- **Accent**: Vibrant orange
- **Atmosphere**: Professional, warm, welcoming

### Dark Theme
- **Background**: Deep blue-black
- **Text**: Cool white
- **Accent**: Same vibrant orange (consistency)
- **Atmosphere**: Focused, modern, developer-friendly

## Color Variables

Both themes use the same variable names, making it easy to reference colors in components:

```css
/* Primary palette (background/base colors) */
--color-primary-5 to --color-primary-95

/* Accent palette (orange - for CTAs, highlights) */
--color-accent-5 to --color-accent-95

/* Inactive/gray palette (for disabled states, secondary text) */
--color-inactive-5 to --color-inactive-95

/* Semantic palettes */
--color-danger-5 to --color-danger-95  /* Red - errors, destructive actions */
--color-action-5 to --color-action-95  /* Cyan - informational */
--color-confirm-5 to --color-confirm-95 /* Green - success, confirmations */
```

## Semantic Variables

Use these for consistent theming:

```css
--color-background         /* Main background */
--color-background-object  /* Card/object backgrounds */
--color-text-regular       /* Primary text color */
--color-text-secondary     /* Secondary/muted text */
--color-border             /* Default border color */
```

## Adding Theme Support to New Components

### 1. Use Semantic Variables

```svelte
<style>
  .my-component {
    background: var(--color-background);
    color: var(--color-text-regular);
    border: var(--border-regular);
  }
</style>
```

### 2. Testing Both Themes

- Toggle between light and dark using the header button
- Check contrast ratios (especially for text)
- Verify hover/focus states work in both themes
- Test all semantic colors (danger, success, action, confirm)

### 3. Color Intensity Guidelines

Numbers indicate lightness/darkness:
- **5-20**: Darkest (backgrounds in light theme, highlights in dark theme)
- **40-60**: Medium (borders, muted elements)
- **80-95**: Lightest (highlights in light theme, backgrounds in dark theme)

In light theme: Lower numbers = darker
In dark theme: Lower numbers = lighter (inverted)

## Typography

```css
/* Headings: Lexend (bold, geometric) */
font-family: 'Lexend', sans-serif;
font-weight: 800;
text-transform: uppercase;
letter-spacing: -0.02em;

/* Body: JetBrains Mono (terminal-inspired) */
font-family: 'JetBrains Mono', monospace;
font-weight: 500;
letter-spacing: -0.01em;
```

## Animation Standards

```css
/* Standard transition */
transition: var(--transition); /* 200ms cubic-bezier(0.4, 0, 0.2, 1) */

/* Transform animations */
transform: translateY(-2px); /* Hover lift */
transform: translateY(2px);  /* Active press */

/* Entrance animations */
@keyframes fadeInUp {
  from { transform: translateY(30px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
```

## Border & Radius

```css
border: var(--border-regular);          /* 2px solid */
border-radius: var(--border-radius-regular); /* 2px - minimal, brutalist */
```

## Shadows

```css
/* Brutalist shadow (offset, no blur) */
box-shadow: 0 4px 0 0 var(--color-accent-70);

/* Subtle depth shadow */
box-shadow: var(--shadow);
```

## Best Practices

1. **Always use CSS variables** - Never hardcode colors
2. **Test both themes** - Check every component in light and dark
3. **Maintain contrast** - Use appropriate color intensities
4. **Be consistent** - Follow established patterns
5. **Semantic meaning** - Use danger for errors, confirm for success, etc.

## Grain Texture

The grain texture overlay is applied globally via `body::before`. It:
- Adds organic quality to the interface
- Animates subtly for movement
- Opacity controlled by `--grain-opacity` variable
- Works in both themes

## Implementation Example

```svelte
<script>
  import { ThemeToggle } from '$lib/components';
</script>

<div class="container">
  <ThemeToggle />

  <div class="card">
    <h2>My Component</h2>
    <p>Content with proper theming</p>
    <button class="accent-button">Action</button>
  </div>
</div>

<style>
  .container {
    background: var(--color-background);
    color: var(--color-text-regular);
    padding: 20px;
  }

  .card {
    background: var(--color-primary-10);
    border: var(--border-regular);
    border-radius: var(--border-radius-regular);
    padding: 24px;
  }

  h2 {
    font-family: 'Lexend', sans-serif;
    font-weight: 800;
    text-transform: uppercase;
    color: var(--color-text-regular);
  }

  .accent-button {
    background: var(--color-accent-50);
    color: white;
    border: none;
    padding: 12px 24px;
    border-radius: var(--border-radius-regular);
    box-shadow: 0 4px 0 0 var(--color-accent-70);
    transition: var(--transition);
  }

  .accent-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 0 0 var(--color-accent-70);
  }

  .accent-button:active {
    transform: translateY(2px);
    box-shadow: 0 2px 0 0 var(--color-accent-70);
  }
</style>
```

---

**Note**: The theme system is designed to be robust and maintainable. Always use the provided variables and follow the established patterns for consistency across the application.
