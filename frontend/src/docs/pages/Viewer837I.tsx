import { useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { VIEWER_837I_EXAMPLES } from '../../data/viewer837i'
import { CODES_837I, ELEMENTS_837I } from '../../data/lookupRows837i'
import { SEGMENT_NAME } from '../../data/x12Names837i'
import { BoxView, type BoxStyle } from '../viewer/BoxView'
import { LookupButton, LookupFrame, type LookupData } from '../viewer/LookupFrame'
import { useLookup } from '../viewer/useLookup'
import { Row } from '../viewer/Parts'
import { RawPanel } from '../viewer/RawEditor'
import { allBoxPaths, ancestorPaths, findBox, type Box } from '../viewer/tree'
import { useEditableRaw, useSyncBoxState } from '../viewer/useEditableRaw'
import { verifyExamplesOrThrow } from '../viewer837i/check'
import { describeClosing, describeCompact, describeOpening, describeRow } from '../viewer837i/describe'
import { parseLoopTree, seedOpenPaths, type BoxType } from '../viewer837i/parse'

// A balance check on every example's raw file — SE01 against the true
// segment count, each claim's CLM02 against the sum of its own lines'
// SV203 — runs once here in dev, so a bad example throws immediately
// instead of mis-rendering quietly.
if (import.meta.env.DEV) verifyExamplesOrThrow(VIEWER_837I_EXAMPLES)

// Box tints follow the 835 viewer's own rule (frontend/src/docs/pages/
// Viewer835.tsx; docs/design/835-viewer.md "Color"): no hue change
// with depth, no green (reserved elsewhere for an 835 element that echoes
// an 837 submission), no amber or red (those read as "difference" and
// "warning" elsewhere in this project). The 837I has far more box types
// than the 835 (roughly thirty, once every provider and other-payer loop is
// counted), so rather than classing each one by name, every box type is
// grouped into one of three structural levels and the level carries the
// style: the envelope (ISA/GS/ST) stays neutral; the repeating "backbone"
// loops a claim is built from (the HL hierarchy, a claim, a service line,
// an other-subscriber loop) get the same neutral bordered-and-tinted-panel
// treatment the 835 gives its LX provider-summary box; every detail loop
// below that (every 2010/2310/2330/2410/2420/2430 party or sub-loop) gets
// the 835's one marked style — the page's own left accent at one constant
// strength, with a faint ink tint on its opening row only.
const BACKBONE_BOX_TYPES = new Set<BoxType>(['BILLING-HL', 'SUBSCRIBER-HL', 'PATIENT-HL', 'CLAIM', 'OTHER-SUBSCRIBER', 'LINE'])
const ENVELOPE_BOX_TYPES = new Set<BoxType>(['ISA', 'GS', 'ST'])

function boxClass(boxType: BoxType): string {
  if (ENVELOPE_BOX_TYPES.has(boxType)) return 'border border-line'
  if (BACKBONE_BOX_TYPES.has(boxType)) return 'border border-line bg-surface/50'
  return 'border border-line border-l-2 border-l-series-1'
}

// The opening row of a detail box is the one place its tint appears; the
// box's body stays on the page background so nested boxes never stack
// tints.
function openingClass(boxType: BoxType): string {
  return !ENVELOPE_BOX_TYPES.has(boxType) && !BACKBONE_BOX_TYPES.has(boxType) ? 'bg-ink/4' : ''
}

const STYLE: BoxStyle<BoxType> = { boxClass, openingClass }

// The three describe functions BoxView draws with; the raw-text editor runs them
// over a new tree first, so text that cannot be drawn is refused (not thrown).
const DESCRIBE = { describeOpening, describeRow, describeClosing }

const LOOKUP_837I: LookupData = { title: '837I field lookup', elements: ELEMENTS_837I, codes: CODES_837I, segmentNames: SEGMENT_NAME }

// Every element's code is a link into the 837I lookup, which opens beside the
// boxes (viewer/LookupFrame.tsx; lookup=1 and q= in the address).
export function Viewer837I() {
  const [params, setParams] = useSearchParams()
  const exampleId = params.get('example') ?? VIEWER_837I_EXAMPLES[0].id
  const example = VIEWER_837I_EXAMPLES.find((e) => e.id === exampleId) ?? VIEWER_837I_EXAMPLES[0]

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
    const map = new Map<string, typeof VIEWER_837I_EXAMPLES>()
    for (const ex of VIEWER_837I_EXAMPLES) map.set(ex.group, [...(map.get(ex.group) ?? []), ex])
    return [...map.entries()]
  }, [])

  const trail: Box<BoxType>[] = maximized
    ? ancestorPaths(maximized)
        .map((p) => findBox(root, p))
        .filter((b): b is Box<BoxType> => b !== undefined)
    : []
  const maximizedBox = maximized ? findBox(root, maximized) : undefined

  return (
    <LookupFrame lookup={lookup} data={LOOKUP_837I}>
      <div className="mx-auto max-w-[100rem] px-10 py-12">
        <PageTitle>837I viewer</PageTitle>

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
