import { Navigate, NavLink, Route, Routes, useLocation } from 'react-router'
import { ThemeToggle } from '../components/ThemeToggle'
import { APP_URL } from './appUrl'
import { ComparedFields } from './pages/ComparedFields'
import { Reference835 } from './pages/Reference835'
import { Viewer835 } from './pages/Viewer835'
import { Compare835 } from './pages/Compare835'
import { Viewer837P } from './pages/Viewer837P'
import { Viewer837I } from './pages/Viewer837I'
import { Format835Flatfile } from './pages/formats/Format835Flatfile'
import { FormatTransmittal907 } from './pages/formats/FormatTransmittal907'
import { FormatTransmittal579 } from './pages/formats/FormatTransmittal579'
import { FormatCedi837p } from './pages/formats/FormatCedi837p'
import { FormatCedi835 } from './pages/formats/FormatCedi835'

interface NavItem {
  path: string
  label: string
  exact?: boolean
  sections?: { id: string; label: string }[]
  children?: NavItem[]
}

// The local docs site: viewers, the compare page, compared fields, file formats, and the reference.
// The review pages are the validator app (vite.config.ts), served by Go.
const NAV: NavItem[] = [
  { path: '/835-viewer', label: '835 viewer' },
  { path: '/835-compare', label: '835 compare' },
  { path: '/837p-viewer', label: '837P viewer' },
  { path: '/837i-viewer', label: '837I viewer' },
  {
    path: '/formats',
    label: 'File formats',
    children: [
      { path: '/formats/835-flat-file', label: 'CMS 835 flat file' },
      { path: '/formats/transmittal-907', label: 'Transmittal 907 (2011)' },
      { path: '/formats/transmittal-579', label: 'Transmittal 579 (2009)' },
      { path: '/formats/cedi-837p', label: 'CEDI 837P guide' },
      { path: '/formats/cedi-835', label: 'CEDI 835 guide' },
    ],
  },
  {
    path: '/835-reference',
    label: '835 Reference',
    sections: [
      { id: 'what-it-is', label: 'What an 835 is' },
      { id: 'money', label: 'The money' },
      { id: 'identifiers', label: 'Identifiers' },
      { id: 'codes', label: 'Codes and qualifiers' },
      { id: 'terms', label: 'Acronyms' },
      { id: 'sources', label: 'Sources' },
    ],
  },
  {
    path: '/compared-fields',
    label: 'Compared fields',
    sections: [
      { id: 'tree', label: 'In file order' },
      { id: 'left-out', label: 'Left out' },
    ],
  },
]

// The 835 viewer was first published at /835-nested; that address, with its
// ?example= query, still opens the viewer.
function RedirectToViewer() {
  const location = useLocation()
  return <Navigate to={{ pathname: '/835-viewer', search: location.search }} replace />
}

// The 835 compare page was first published at /835-list; that address, with
// its ref=, a=, b=, hide=, and differs= query, still opens the compare page.
function RedirectToCompare() {
  const location = useLocation()
  return <Navigate to={{ pathname: '/835-compare', search: location.search }} replace />
}

// /lookup/835 and /lookup/837, the addresses of the standalone field
// lookups, open the viewers' lookup panels. /lookup/837 covered both 837
// types, so it opens the 837P viewer.
function RedirectTo835Lookup() {
  return <Navigate to="/835-viewer?lookup=1" replace />
}
function RedirectTo837Lookup() {
  return <Navigate to="/837p-viewer?lookup=1" replace />
}

function NavTree({ items, depth }: { items: NavItem[]; depth: number }) {
  const location = useLocation()
  return (
    <ul className={depth === 0 ? 'space-y-0.5' : 'mt-0.5 space-y-0.5 border-l border-line pl-3'}>
      {items.map((item) => (
        <li key={item.path}>
          <NavLink
            to={item.path}
            end={item.exact}
            className={({ isActive }) =>
              `block border-l-2 px-2.5 py-1 text-sm ${
                isActive
                  ? 'border-accent font-medium text-ink'
                  : 'border-transparent text-ink-2 hover:bg-surface hover:text-ink'
              }`
            }
          >
            {item.label}
          </NavLink>
          {item.sections && location.pathname === item.path && (
            <ul className="mb-1 mt-0.5">
              {item.sections.map((s) => (
                <li key={s.id}>
                  <a
                    href={`#${s.id}`}
                    className="block py-[3px] pl-7 pr-2 font-mono text-[12px] text-ink-2 hover:text-ink"
                  >
                    {s.label}
                  </a>
                </li>
              ))}
            </ul>
          )}
          {item.children && <NavTree items={item.children} depth={depth + 1} />}
        </li>
      ))}
    </ul>
  )
}

export default function App() {
  return (
    <div className="flex min-h-svh">
      <ThemeToggle />
      <nav className="sticky top-0 h-svh w-72 shrink-0 overflow-y-auto border-r border-line px-3 py-6">
        <div className="mb-7 flex items-center gap-2 px-2 font-mono text-[12px] font-medium tracking-[0.04em] text-ink">
          <span className="h-2.5 w-2.5 shrink-0 bg-accent" />
          x12-835-parity-validator
        </div>
        <a href={APP_URL} className="mb-4 block border-l-2 border-transparent px-2.5 py-1 text-sm text-link hover:bg-surface">
          Open the validator app ↗
        </a>
        <NavTree items={NAV} depth={0} />
      </nav>
      <main className="min-w-0 flex-1">
        <Routes>
          <Route path="/835-reference" element={<Reference835 />} />
          <Route path="/835-viewer" element={<Viewer835 />} />
          <Route path="/835-compare" element={<Compare835 />} />
          <Route path="/837p-viewer" element={<Viewer837P />} />
          <Route path="/835-nested" element={<RedirectToViewer />} />
          <Route path="/835-list" element={<RedirectToCompare />} />
          <Route path="/837i-viewer" element={<Viewer837I />} />
          <Route path="/lookup/835" element={<RedirectTo835Lookup />} />
          <Route path="/lookup/837" element={<RedirectTo837Lookup />} />
          <Route path="/formats" element={<Format835Flatfile />} />
          <Route path="/formats/835-flat-file" element={<Format835Flatfile />} />
          <Route path="/formats/transmittal-907" element={<FormatTransmittal907 />} />
          <Route path="/formats/transmittal-579" element={<FormatTransmittal579 />} />
          <Route path="/formats/cedi-837p" element={<FormatCedi837p />} />
          <Route path="/formats/cedi-835" element={<FormatCedi835 />} />
          <Route path="/" element={<Navigate to="/835-viewer" replace />} />
          <Route path="/compared-fields" element={<ComparedFields />} />
        </Routes>
      </main>
    </div>
  )
}
