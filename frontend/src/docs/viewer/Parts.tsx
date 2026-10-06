import { CodeLink, CodeName, NAME_LINE } from './Code'
import type { Cell, Part, SegmentDisplay } from './display'
import { Term } from './Term'

// How a segment shows in every viewer (docs/design/835-viewer.md): its
// segment ID leading the row, the public name beneath it,
// then one cell per element. A cell is two lines. The first is the element's
// code and the value the file carries, on one baseline at one size: the code
// in monospace in the link tone, a link into the page's lookup (Code.tsx);
// the value in the body face at full contrast, the brightest thing in the
// cell. The second is the element's public name, small and dim, the cell's
// one minor line, across its full width. A cell is as wide as the longer of
// its two lines; a name wider than the box wraps and the cell grows taller,
// never cut. An element with no public name has no second line.

const HINTED = 'underline decoration-dotted decoration-line underline-offset-2 hover:decoration-ink-2'

// A value as the file carries it. A coded value shows what it names on hover
// (its public name, and for a CARC or RARC the maintainer's full text); one
// with no public name is its code alone. A present-but-empty element is a
// faint ·, so its position still shows honestly.
function ValueView({ part }: { part: Part }) {
  if (part.kind === 'joined') {
    return (
      <>
        {part.parts.map((p, i) => (
          <span key={i}>
            {i > 0 && <span className="text-ink-2">:</span>}
            <ValueView part={p} />
          </span>
        ))}
      </>
    )
  }
  if (part.kind === 'value') return part.text === '' ? <span className="text-ink-2">·</span> : <>{part.text}</>
  if (!part.name) return <>{part.code}</>
  const body = part.full !== undefined && part.full !== part.name ? part.full : undefined
  return <Term trigger={part.code} triggerClassName={HINTED} title={part.name} body={body} />
}

// A part as plain text, with no hover.
function textOf(part: Part): string {
  if (part.kind === 'joined') return part.parts.map((p) => textOf(p)).join(':')
  if (part.kind === 'value') return part.text === '' ? '·' : part.text
  return part.code
}

// A segment on one line, as plain text: the trail of a maximized box, where
// each step is itself a button, so nothing in it can be another control.
export function Row({ parts, className }: { parts: Part[]; className?: string }) {
  return (
    <span className={className}>
      {parts.map((part, i) => (
        <span key={i} className="mr-2 inline-block align-baseline last:mr-0">
          {textOf(part)}
        </span>
      ))}
    </span>
  )
}

function CellView({ cell }: { cell: Cell }) {
  const empty = cell.part.kind === 'value' && cell.part.text === ''
  return (
    <div className={`max-w-full shrink-0 border-r border-line px-2 py-1 last:border-r-0 ${empty ? 'opacity-50' : ''}`}>
      <div className="flex items-baseline gap-1.5 text-[13px] leading-5">
        <CodeLink code={cell.ref} className="shrink-0" />
        {cell.marker && <span className={`inline-block h-1.5 w-1.5 shrink-0 self-center rounded-full ${cell.marker}`} aria-hidden="true" />}
        <span className="min-w-0 text-ink [overflow-wrap:anywhere]">
          <ValueView part={cell.part} />
        </span>
      </div>
      {cell.name && <span className={NAME_LINE}>{cell.name}</span>}
    </div>
  )
}

// A segment's whole display: its segment ID with its name beneath (a code
// standing alone as a label, Code.tsx), then one cell per element the file
// holds, side by side in the same row. The strip only wraps onto more lines
// when it is wider than its box (flex-wrap, never a forced one-element-per-
// line layout) — long segments like ISA (16 elements) or CLP (up to 14) wrap
// there naturally.
export function Strip({ display, className }: { display: SegmentDisplay; className?: string }) {
  return (
    <div className={`flex min-w-0 flex-wrap items-stretch ${className ?? ''}`}>
      {display.lead.code !== '' && (
        <div className="mr-2 max-w-full shrink-0 py-1 text-[13px] leading-5">
          <CodeName code={display.lead.code} name={display.lead.name} />
        </div>
      )}
      {display.cells.map((c, i) => (
        <CellView key={i} cell={c} />
      ))}
    </div>
  )
}
