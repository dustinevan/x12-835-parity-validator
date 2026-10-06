import { useCallback, useRef, useState } from 'react'
import { Link, NavLink, Route, Routes, useLocation } from 'react-router'
import { ThemeToggle } from '../components/ThemeToggle'
import { VendorStrip } from './VendorPicker'
import { HELP, PAGE_BY_PATH } from './help'
import { HelpDrawer } from './HelpDrawer'
import { MacMapButton, MacMapPanel, useMacMap } from './MacTree'
import { ClaimCompare } from './pages/ClaimCompare'
import { Overview } from './pages/Overview'
import { ViewDifferences } from './pages/ViewDifferences'
import { ViewGroup } from './pages/ViewGroup'
import { ViewUnmatched } from './pages/ViewUnmatched'
import { ViewWorkloads } from './pages/ViewWorkloads'
import { useSelection } from './selection'

// The validator app's shell: the sections sit in a left navigation column,
// with the two submission pickers above them; the page controls sit in the
// top right corner; the MAC map, when open, is a column at the right. The
// Overview is the app's root, and the review pages (Claims, Differences,
// Unmatched, Workloads) and the claim screen (/claim, opened from links
// inside the pages) have a path each.

const NAV = [
  { path: '/', label: 'Overview' },
  { path: '/claims', label: 'Claims' },
  { path: '/differences', label: 'Differences' },
  { path: '/unmatched', label: 'Unmatched' },
  { path: '/workloads', label: 'Workloads' },
]

function NotFound() {
  const { pathname } = useLocation()
  const sel = useSelection()
  return (
    <div className="mx-auto max-w-[72rem] px-10 py-12 font-mono text-[13px] text-ink-2">
      No page at {pathname}. <Link to={sel.link('/')} className="text-link hover:underline">Open the overview</Link>.
    </div>
  )
}

export function Shell() {
  const sel = useSelection()
  const { pathname } = useLocation()
  const page = PAGE_BY_PATH[pathname]
  const helpButton = useRef<HTMLButtonElement>(null)
  const [helpOpen, setHelpOpen] = useState(false)
  const closeHelp = useCallback((restoreFocus: boolean) => {
    setHelpOpen(false)
    if (restoreFocus) helpButton.current?.focus()
  }, [])
  const macMap = useMacMap()

  return (
    <div className="flex min-h-svh">
      <nav className="sticky top-0 flex h-svh w-[18rem] shrink-0 flex-col overflow-y-auto border-r border-line px-3 py-6">
        <Link to={sel.link('/')} className="mb-5 flex items-center gap-2 px-2 text-sm font-semibold text-ink">
          <span className="h-2.5 w-2.5 shrink-0 bg-accent" />
          835 Parity Validator
        </Link>
        <div className="mb-5 border-y border-line px-2 py-3 [&_select]:max-w-[13rem]">
          <div className="mb-2 font-mono text-[11px] uppercase tracking-wide text-ink-2">Submissions shown</div>
          <VendorStrip />
        </div>
        <ul className="space-y-0.5">
          {NAV.map((item) => (
            <li key={item.path}>
              <NavLink to={sel.link(item.path)} end={item.path === '/'}
                className={({ isActive }) => `block border-l-2 px-2.5 py-1 text-sm ${isActive ? 'border-accent font-medium text-ink' : 'border-transparent text-ink-2 hover:bg-surface hover:text-ink'}`}>
                {item.label}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
      {/* The page controls sit in the top right corner, where the docs site
          keeps them: the help for the page, the MAC map, then the theme toggle. */}
      <div className="fixed right-4 top-4 z-40 flex items-center gap-2">
        {page && (
          <button ref={helpButton} type="button" onClick={() => setHelpOpen((o) => !o)}
            aria-label="How to read this page" aria-haspopup="dialog" aria-expanded={helpOpen}
            className={`shrink-0 border bg-bg px-2 py-1 font-mono text-[11px] hover:text-ink ${helpOpen ? 'border-ink text-ink' : 'border-line text-ink-2'}`}>
            ? how to read this page
          </button>
        )}
        <MacMapButton map={macMap} />
        <ThemeToggle inline />
      </div>
      <HelpDrawer help={page && HELP[page]} open={helpOpen && page !== undefined} onClose={closeHelp} opener={helpButton} />
      {/* The page's column is a size container, so a page's @min-[…] classes
          follow its own width with the MAC map open or closed. */}
      <main className="@container min-w-0 flex-1">
        <Routes>
          <Route path="/" element={<Overview />} />
          <Route path="/claims" element={<ViewGroup />} />
          <Route path="/differences" element={<ViewDifferences />} />
          <Route path="/unmatched" element={<ViewUnmatched />} />
          <Route path="/workloads" element={<ViewWorkloads />} />
          <Route path="/claim" element={<ClaimCompare />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
      {macMap.open && <MacMapPanel onClose={() => macMap.setOpen(false)} />}
    </div>
  )
}
