import { useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { VIEWER_837P_EXAMPLES } from '../../data/viewer837p'
import { CODES_837P, ELEMENTS_837P } from '../../data/lookupRows837p'
import { SEGMENT_NAME } from '../../data/x12Names837p'
import { BoxView, type BoxStyle } from '../viewer/BoxView'
import { describeCompact } from '../viewer/display'
import { LookupButton, LookupFrame, type LookupData } from '../viewer/LookupFrame'
import { useLookup } from '../viewer/useLookup'
import { Row } from '../viewer/Parts'
import { RawPanel } from '../viewer/RawEditor'
import { allBoxPaths, ancestorPaths, findBox, type Box } from '../viewer/tree'
import { useEditableRaw, useSyncBoxState } from '../viewer/useEditableRaw'
import { verifyExamplesOrThrow } from '../viewer837p/check'
import { boxClass, openingClass } from '../viewer837p/colors'
import { describeClosing, describeOpening, describeRow } from '../viewer837p/describe'
import { parseLoopTree, seedOpenPaths, type BoxType } from '../viewer837p/parse'

// A balance check on every example's raw file — SE01 against the true
// segment count, CLM02 against the sum of every SV102 — runs once here in
// dev, so a bad example throws immediately instead of mis-rendering
// quietly (frontend/src/docs/viewer837p/check.ts).
if (import.meta.env.DEV) verifyExamplesOrThrow(VIEWER_837P_EXAMPLES)

const STYLE: BoxStyle<BoxType> = { boxClass, openingClass }

// The three describe functions BoxView draws with; the raw-text editor runs them
// over a new tree first, so text that cannot be drawn is refused (not thrown).
const DESCRIBE = { describeOpening, describeRow, describeClosing }

const LOOKUP_837P: LookupData = { title: '837P field lookup', elements: ELEMENTS_837P, codes: CODES_837P, segmentNames: SEGMENT_NAME }

// Every element's code is a link into the 837P lookup, which opens beside the
// boxes (viewer/LookupFrame.tsx; lookup=1 and q= in the address).
export function Viewer837P() {
  const [params, setParams] = useSearchParams()
  const exampleId = params.get('example') ?? VIEWER_837P_EXAMPLES[0].id
  const example = VIEWER_837P_EXAMPLES.find((e) => e.id === exampleId) ?? VIEWER_837P_EXAMPLES[0]

  const raw = useEditableRaw(example, parseLoopTree, DESCRIBE)
  const root = raw.root
  const [openPaths, setOpenPaths] = useState<Set<string>>(() => allBoxPaths(root))
  const [maximized, setMaximized] = useState<string | null>(null)
  const [hoveredRaw, setHoveredRaw] = useState<string | null>(null)
  const lookup = useLookup()

  useSyncBoxState(root, example.id, setOpenPaths, setMaximized)

  useEffect(() => {
    if (!maximized) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setMaximized(null)
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [maximized])

  function toggleOpen(path: string) {
    setOpenPaths((prev) => {
      const next = new Set(prev)
      if (next.has(path)) next.delete(path)
      else next.add(path)
      return next
    })
  }

  function expandAll() {
    setOpenPaths(allBoxPaths(root))
  }

  function collapseAll() {
    setOpenPaths(seedOpenPaths(root, example.open))
  }

  function selectExample(id: string) {
    const next = new URLSearchParams(params)
    next.set('example', id)
    setParams(next)
  }

  const groups = useMemo(() => {
    const map = new Map<string, typeof VIEWER_837P_EXAMPLES>()
    for (const ex of VIEWER_837P_EXAMPLES) map.set(ex.group, [...(map.get(ex.group) ?? []), ex])
    return [...map.entries()]
  }, [])

  const trail: Box<BoxType>[] = maximized
    ? ancestorPaths(maximized)
        .map((p) => findBox(root, p))
        .filter((b): b is Box<BoxType> => b !== undefined)
    : []
  const maximizedBox = maximized ? findBox(root, maximized) : undefined

  return (
    <LookupFrame lookup={lookup} data={LOOKUP_837P}>
      <div className="mx-auto max-w-[100rem] px-10 py-12">
        <PageTitle>837P viewer</PageTitle>

        {/* The example list sits beside the boxes once the page is 32rem wide
            (the list's 15rem, the gap, the padding, and the 10rem the toolbar
            needs), and above them under that, so a narrow page (the lookup
            panel open on a small window) never scrolls sideways. */}
        <div className="flex flex-col gap-8 @min-[32rem]:flex-row">
          <aside className="shrink-0 @min-[32rem]:w-60">
            <nav className="sticky top-12 space-y-5">
              {groups.map(([group, examples]) => (
                <div key={group}>
                  <div className="mb-1.5 font-mono text-[11px] uppercase tracking-wide text-ink-2">{group}</div>
                  <ul className="space-y-0.5">
                    {examples.map((ex) => (
                      <li key={ex.id}>
                        <button
                          type="button"
                          onClick={() => selectExample(ex.id)}
                          aria-current={ex.id === example.id ? 'true' : undefined}
                          className={`block w-full border-l-2 px-2.5 py-1 text-left text-sm ${
                            ex.id === example.id
                              ? 'border-accent font-medium text-ink'
                              : 'border-transparent text-ink-2 hover:bg-surface hover:text-ink'
                          }`}
                        >
                          {ex.title}
                          {ex.id === example.id && raw.edited && <span className="ml-1.5 font-normal text-ink-2">(edited)</span>}
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              ))}
            </nav>
          </aside>

          <div className="min-w-0 flex-1">
            <div className="mb-4 flex flex-wrap items-center justify-end gap-3">
              <div className="flex shrink-0 gap-3 text-xs">
                <button type="button" onClick={expandAll} className="text-ink-2 underline decoration-dotted decoration-line underline-offset-2 hover:text-ink">
                  Expand all
                </button>
                <button type="button" onClick={collapseAll} className="text-ink-2 underline decoration-dotted decoration-line underline-offset-2 hover:text-ink">
                  Collapse all
                </button>
              </div>
              <LookupButton lookup={lookup} />
            </div>

            {maximizedBox && (
              <div className="mb-3 flex items-center justify-between gap-4">
                <div className="min-w-0 flex-1 overflow-hidden text-ellipsis whitespace-nowrap font-mono text-xs text-ink-2">
                  {trail.map((box) => (
                    <span key={box.path}>
                      <button type="button" onClick={() => setMaximized(box.path)} className="hover:text-ink hover:underline">
                        <Row parts={describeCompact(describeOpening(box))} />
                      </button>
                      <span className="mx-1.5">›</span>
                    </span>
                  ))}
                </div>
                <button
                  type="button"
                  onClick={() => setMaximized(null)}
                  className="shrink-0 border border-line px-2 py-1 text-xs text-ink-2 hover:text-ink"
                >
                  Restore ✕
                </button>
              </div>
            )}

            <BoxView
              box={maximizedBox ?? root}
              openPaths={openPaths}
              onToggle={toggleOpen}
              onMaximize={setMaximized}
              onHoverSegment={setHoveredRaw}
              topLevel={!!maximizedBox}
              style={STYLE}
              describeOpening={describeOpening}
              describeRow={describeRow}
              describeClosing={describeClosing}
            />

            <RawPanel raw={raw} hoveredRaw={hoveredRaw} />
          </div>
        </div>
      </div>
    </LookupFrame>
  )
}
