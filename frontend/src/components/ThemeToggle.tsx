import { useState } from 'react'

type Theme = 'light' | 'dark'

const STORAGE_KEY = 'validator-theme'

// Current effective mode: an explicit data-theme attribute (set by the
// flash-free inline script in index.html, or a prior click) wins; with no
// stored choice, follow the system preference.
function effectiveTheme(): Theme {
  const stored = document.documentElement.getAttribute('data-theme')
  if (stored === 'light' || stored === 'dark') return stored
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

// Fixed to the top-right corner by default; `inline` places it in a layout,
// such as the app's row of corner controls.
export function ThemeToggle({ inline = false }: { inline?: boolean }) {
  const [theme, setTheme] = useState<Theme>(effectiveTheme)
  const next: Theme = theme === 'dark' ? 'light' : 'dark'

  function handleClick() {
    document.documentElement.setAttribute('data-theme', next)
    localStorage.setItem(STORAGE_KEY, next)
    setTheme(next)
  }

  return (
    <button
      type="button"
      onClick={handleClick}
      aria-label={`Switch to ${next} theme`}
      className={`${inline ? 'shrink-0' : 'fixed right-4 top-4 z-50'} border border-line bg-bg px-2 py-1 font-mono text-[11px] text-ink-2 hover:text-ink`}
    >
      {next === 'dark' ? '☾ dark' : '☀ light'}
    </button>
  )
}
