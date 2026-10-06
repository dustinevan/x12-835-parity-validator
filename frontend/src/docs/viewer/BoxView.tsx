import type { KeyboardEvent, ReactNode } from 'react'
import type { SegmentDisplay } from './display'
import { Strip } from './Parts'
import { childBoxCount, type Box, type Segment } from './tree'

// How one transaction's boxes look — borders, a left-accent hue, and which
// opening row carries the subtle background tint (see each transaction's
// own colors.ts for the actual classes and the reasoning behind them:
// frontend/src/docs/viewer835/colors.ts, frontend/src/docs/viewer837p/colors.ts).
export interface BoxStyle<BoxType extends string> {
  boxClass: (boxType: BoxType) => string
  openingClass: (boxType: BoxType) => string
  // Extra classes for a plain segment row, keyed by its enclosing box's
  // type — the 835's payment-header rows directly under ST, for instance.
  // Optional; a transaction with no such rule omits it.
  rowClass?: (hostBoxType: BoxType) => string | undefined
}

// Renders one loop box and, recursively, its children — the expand/collapse
// chrome, the maximize button, and hover-to-highlight-in-the-raw-file wiring
// are all transaction-agnostic; what differs per transaction is how a box
// looks (`style`) and how its opening/closing/plain rows describe
// themselves (`describeOpening`/`describeRow`/`describeClosing`, each
// transaction's own describe.ts). `renderExtra` lets a transaction draw
// something below a maximized box's own rows (the 835's per-line economics
// bar, frontend/src/docs/viewer835/Bar.tsx); a transaction with nothing
// like that passes nothing.
export function BoxView<BoxType extends string>({
  box,
  openPaths,
  onToggle,
  onMaximize,
  onHoverSegment,
  topLevel,
  style,
  describeOpening,
  describeRow,
  describeClosing,
  renderExtra,
}: {
  box: Box<BoxType>
  openPaths: Set<string>
  onToggle: (path: string) => void
  onMaximize: (path: string) => void
  // Raised with a segment's raw text on hover, and with null on leave, so
  // the raw-file panel below can highlight the matching line — optional,
  // since this view also renders standalone (the maximized box).
  onHoverSegment?: (raw: string | null) => void
  topLevel?: boolean
  style: BoxStyle<BoxType>
  describeOpening: (box: Box<BoxType>) => SegmentDisplay
  describeRow: (segment: Segment, host: BoxType) => SegmentDisplay
  describeClosing: (box: Box<BoxType>) => SegmentDisplay
  renderExtra?: (box: Box<BoxType>) => ReactNode
}) {
  const isOpen = topLevel === true || openPaths.has(box.path)
  // Only envelope boxes (ISA/GS/ST in every transaction this project draws)
  // carry a closing segment of their own — the same generic signal
  // frontend/src/docs/viewer/tree.ts's envelopePaths uses.
  const isEnvelope = box.closing !== undefined
  const count = childBoxCount(box)
  const extra = topLevel ? renderExtra?.(box) : undefined
  const hoverProps = (raw: string) => ({
    onMouseEnter: () => onHoverSegment?.(raw),
    onMouseLeave: () => onHoverSegment?.(null),
  })

  // Enter or Space opens or closes the box only when the row itself has the
  // focus; on a control inside the row (a code's lookup link, a coded
  // value's pop-up, the maximize button) the key belongs to that control.
  function handleKey(e: KeyboardEvent<HTMLDivElement>) {
    if (e.target !== e.currentTarget) return
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      onToggle(box.path)
    }
  }

  return (
    <div className={`my-1 ${style.boxClass(box.boxType)}`}>
      <div
        role={topLevel ? undefined : 'button'}
        tabIndex={topLevel ? undefined : 0}
        aria-expanded={topLevel ? undefined : isOpen}
        onClick={topLevel ? undefined : () => onToggle(box.path)}
        onKeyDown={topLevel ? undefined : handleKey}
        {...hoverProps(box.opening.raw)}
        className={`flex min-w-0 items-start justify-between gap-3 px-3 py-1.5 ${style.openingClass(box.boxType)} ${topLevel ? '' : 'cursor-pointer hover:bg-surface'}`}
      >
        <div className="min-w-0 flex-1">
          <Strip display={describeOpening(box)} />
          {!isOpen && count > 0 && <span className="ml-1 tracking-tight text-ink-2">{'▮'.repeat(count)}</span>}
        </div>
        {!topLevel && (
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation()
              onMaximize(box.path)
            }}
            aria-label="Maximize this box"
            title="Maximize"
            className="shrink-0 text-ink-2 hover:text-ink"
          >
            ⤢
          </button>
        )}
      </div>

      <div
        className="grid transition-[grid-template-rows] duration-300 ease-out motion-reduce:transition-none"
        style={{ gridTemplateRows: isOpen ? '1fr' : '0fr' }}
      >
        <div className="min-h-0 overflow-hidden">
          <div className="space-y-1 py-1 pl-5 pr-3">
            {box.children.map((node, i) =>
              node.kind === 'box' ? (
                <BoxView
                  key={node.box.path}
                  box={node.box}
                  openPaths={openPaths}
                  onToggle={onToggle}
                  onMaximize={onMaximize}
                  onHoverSegment={onHoverSegment}
                  style={style}
                  describeOpening={describeOpening}
                  describeRow={describeRow}
                  describeClosing={describeClosing}
                  renderExtra={renderExtra}
                />
              ) : (
                <div key={i} className={style.rowClass?.(box.boxType)} {...hoverProps(node.segment.raw)}>
                  <Strip display={describeRow(node.segment, box.boxType)} />
                </div>
              ),
            )}
          </div>
          {extra && <div className="px-3 pb-3 pt-1">{extra}</div>}
          {isEnvelope && box.closing && (
            <div className="border-t border-line px-3 py-1.5" {...hoverProps(box.closing.raw)}>
              <Strip display={describeClosing(box)} />
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
