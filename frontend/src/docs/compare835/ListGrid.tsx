import type { ReactNode } from 'react'
import { isAmount, isZero, signedMoney, subtract } from '../../format'
import { CodeLink, NAME_LINE } from '../viewer/Code'
import type { Part } from '../viewer/display'
import { Term } from '../viewer/Term'
import { COLS, isMarked, type Col, type Label, type ListRow, type ListSection } from './rows'

// The grid of the list: one row per element, a key column, then the
// reference and the two compared files. Amber marks a compared cell whose
// value differs from the reference's (rows.ts isMarked; a dollar amount within
// five cents of the reference's is the same amount and is not marked): a
// tinted cell with a bold value and, for a dollar amount, the signed
// difference (compared minus reference). A faint · is an element the file
// does not carry, or carries empty.

const HINTED = 'text-left underline decoration-dotted decoration-line underline-offset-2 hover:text-ink hover:decoration-ink-2'

// A label's code: a link into the page's lookup (viewer/Code.tsx) when it is
// an element's code, plain monospace when it is not (a CAS row's CO-45).
function KeyCode({ label }: { label: Label }) {
  if (label.lookup === undefined) return <span className="font-mono">{label.codes}</span>
  return <CodeLink code={label.codes} search={label.lookup} title={label.hint} />
}

// A value as the file carries it. A coded value shows what it names on hover
// (the name in the viewer's tables; for a CARC or RARC, the full text too).
function Value({ text, part }: { text: string; part?: Part }) {
  if (part?.kind === 'label' && part.name) {
    return <Term trigger={text} triggerClassName={HINTED} title={part.name} body={part.full} />
  }
  if (part?.kind === 'joined') {
    return (
      <>
        {part.parts.map((p, i) => (
          <span key={i}>
            {i > 0 && <span className="text-ink-2">:</span>}
            {p.kind === 'label' && p.name
              ? <Term trigger={p.code} triggerClassName={HINTED} title={p.name} body={p.full} />
              : p.kind === 'label' ? p.code : p.text}
          </span>
        ))}
      </>
    )
  }
  return <>{text}</>
}

function ListCell({ section, row, col }: { section: ListSection; row: ListRow; col: Col }) {
  const value = row.values[col]
  const marked = isMarked(section, row, col)
  const cell = `border-l-4 py-1.5 pl-2 align-top ${marked ? 'border-warn bg-warn/15' : 'border-transparent'}`
  if (value === undefined) return <td className={cell}><span className="font-mono text-ink-2/50">·</span></td>
  const reference = row.values.ref
  const delta = marked && row.amount && reference !== undefined && isAmount(reference) && isAmount(value) ? subtract(value, reference) : undefined
  return (
    <td className={cell}>
      <span className="flex flex-wrap items-center gap-x-2 gap-y-0.5">
        <span className={`font-mono text-xs [overflow-wrap:anywhere] ${marked ? 'font-semibold' : ''} text-ink`}>
          <Value text={value} part={row.parts[col]} />
        </span>
        {delta !== undefined && !isZero(delta) && (
          <span className="bg-warn px-1.5 font-mono text-[11px] font-semibold text-black">{signedMoney(delta)}</span>
        )}
      </span>
    </td>
  )
}

// The key: the code is the main text, the public name small beneath it
// (both at once, code major, name minor, with no mode to switch); every
// code on the docs site stands as a label this way (viewer/Code.tsx). A row's
// context (the segment's qualifier, DTM*232, or its place, PLB (2nd)) stays
// inline after the code, a link too, its name on hover. A marked row's key is
// heavier.
// What a page adds to a row: a note after the key (the validator app's claim
// screen puts the compared set's state there) and a class on the row.
export interface RowDecoration {
  note?: ReactNode
  className?: string
}

function KeyCell({ section, row, note }: { section: ListSection; row: ListRow; note?: ReactNode }) {
  const active = COLS.some((c) => isMarked(section, row, c))
  return (
    <td className={`py-1.5 pl-3 align-top [overflow-wrap:anywhere] ${active ? 'font-medium text-ink' : 'text-ink-2'}`}>
      <span className="text-xs">
        <KeyCode label={row.key} />
      </span>
      {row.context && (
        <span className="ml-2 text-[11px] font-normal">
          <KeyCode label={row.context} />
        </span>
      )}
      {row.key.names !== row.key.codes && <span className={NAME_LINE}>{row.key.names}</span>}
      {note}
    </td>
  )
}

// `head` is the content of each column's header cell: the example it shows.
// `decorate`, when a page gives one, adds a note and a class to each row.
export function ListGrid({ sections, head, decorate }: {
  sections: ListSection[]
  head: Record<Col, ReactNode>
  decorate?: (section: ListSection, row: ListRow) => RowDecoration | undefined
}) {
  const th = 'sticky top-0 z-10 bg-bg py-2 pl-3 pr-2 shadow-[inset_0_-2px_0_0_var(--text)]'
  return (
    // The table is only as wide as its content — a key column and three value
    // columns, 52rem in all (16rem and 12rem each) — so a row's values sit
    // beside its key and the eye reads across without following lines. The
    // page's spare width stays outside the table.
    // Where the page is narrower than that (a narrow window, or the lookup
    // panel open beside it), the four columns narrow in the same proportions,
    // and the page never scrolls sideways.
    <table className="w-[min(100%,52rem)] table-fixed text-sm">
      <thead>
        <tr className="text-left align-bottom text-xs text-ink">
          <th className={`${th} w-4/13`}>Element</th>
          {COLS.map((col) => <th key={col} className={`${th} w-3/13`}>{head[col]}</th>)}
        </tr>
      </thead>
      {sections.map((section) => (
        <tbody key={section.id}>
          <tr>
            <td colSpan={COLS.length + 1} className="pt-5 pb-1.5 pl-3 font-mono text-[11px] font-semibold uppercase tracking-wider text-ink-2">
              {section.title}
              {section.note && <span className="ml-3 font-normal normal-case tracking-normal">{section.note}</span>}
            </td>
          </tr>
          {section.rows.map((row, i) => {
            const deco = decorate?.(section, row)
            return (
              <tr key={i} className={`border-b border-line ${deco?.className ?? ''}`}>
                <KeyCell section={section} row={row} note={deco?.note} />
                {COLS.map((col) => <ListCell key={col} section={section} row={row} col={col} />)}
              </tr>
            )
          })}
        </tbody>
      ))}
    </table>
  )
}
