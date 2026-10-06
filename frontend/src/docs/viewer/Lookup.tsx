import { useEffect, useMemo, useRef, useState, type ReactNode } from 'react'
import { Tag, type TagColor } from '../../components/Tag'
import { CodeLink, NAME_LINE } from './Code'

// What an element is, what it is for, the values it takes, and how Medicare
// fills it, each sentence cited: the entry the panel opens under a row.
// The 835 rows carry one (frontend/src/data/x12About835.ts); a transaction
// whose rows carry none stays a one-line list.
export interface LookupAbout {
  is: string
  purpose?: string
  values?: string
  medicare?: string
  sources: string[]
}

// What the validator does with an element: the state the Compared fields
// page gives it, with the sentence behind it (frontend/src/data/
// lookupRows835.ts validatorNote).
export interface LookupValidatorNote {
  state: string
  why: string
  detail?: string
}

// One positional element of a transaction (CLP04, CLM05-1, HI01-2…), with
// the name a public document gives it and the page that names it. Shared
// shape for the 835, 837P, and 837I row builders (frontend/src/data/
// lookupRows835.ts, lookupRows837p.ts, lookupRows837i.ts), each built only
// from that transaction's own sourced names file — nothing here is invented.
export interface LookupElementRow {
  ref: string // the element reference: CLP04, CLM05-1, HI01-2
  name?: string // its public name; absent when no public document names it
  segment: string // the bare segment id this element belongs to: CLP, CLM, HI
  segmentName?: string // the segment's public name at this use (a loop/box title)
  tag?: string // 835 only: the segment with its qualifier at this use (DTM*232, NM1*QC); a search finds the row by it
  where: string // the part of the transaction it sits in
  source: string // the citation
  echo?: string // 835 only: the public text saying this element carries the value submitted on the 837
  use?: string // 835 only: the key of this use in x12About835.ts
  about?: LookupAbout // the entry, when one is written
  segmentAbout?: LookupAbout // the entry for the segment at this use
  validator?: LookupValidatorNote // what the validator does with it
}

// One coded value of one element, named by a public document.
export interface LookupCodeRow {
  segment: string
  element: string // the element this code is a value of, e.g. "CLM05-3"
  code: string
  name?: string
  full?: string // the maintainer's complete text, when longer than `name` (CARC, RARC)
  where: string
  source: string
  // The one element row this code belongs under, when the row builder knows
  // it (most codes): the exact `segmentName` of that row, so two different
  // boxes that reuse the same bare element ref (REF01 at a dozen different
  // loops, NM101 at every NM1) don't bleed into each other. Left out for a
  // code that is genuinely reused across more than one box with no single
  // owner (CARC/RARC, a role code valid at every NM1) — it then attaches to
  // every element row with a matching ref in the segment.
  segmentName?: string
  // What the code means beyond its name, where a document says more.
  about?: { text: string; sources: string[] }
}

interface GroupedElement {
  row: LookupElementRow
  codes: LookupCodeRow[]
}

interface Group {
  segment: string
  heading: string
  elements: GroupedElement[]
  orphanCodes: LookupCodeRow[] // a coded value whose element ref matches nothing in `elements` (should be rare)
}

// Searches ignore case and the separators people type between a segment and
// its qualifier, so "nm1 85", "NM1*85", and "nm185" find the same row.
function normalize(text: string) {
  return text.toLowerCase().replace(/[*\s\-·,…()]/g, '')
}

function matches(words: string[], ...fields: (string | undefined)[]) {
  const text = fields.filter(Boolean).join(' ').toLowerCase()
  const compact = normalize(text)
  return words.every((word) => text.includes(word) || compact.includes(normalize(word)))
}

// Groups elements by their bare segment, in the order each segment first
// appears (the file's own order), with every coded value filed under the
// specific element position it is a value of. A code whose element matches
// no row (should not happen with a complete row builder) is kept so nothing
// is silently dropped.
function buildGroups(elements: LookupElementRow[], codes: LookupCodeRow[], segmentNames: Partial<Record<string, string>>): Group[] {
  const order: string[] = []
  const bySegment = new Map<string, LookupElementRow[]>()
  for (const row of elements) {
    if (!bySegment.has(row.segment)) {
      bySegment.set(row.segment, [])
      order.push(row.segment)
    }
    bySegment.get(row.segment)!.push(row)
  }
  // A code naming the one box it belongs to (`segmentName`) attaches only to
  // the element row(s) of that exact box, so two boxes that reuse the same
  // bare element ref (REF01 at a dozen different loops, NM101 at every NM1)
  // don't bleed into each other. A code with no `segmentName` (deliberately
  // reused across boxes — CARC/RARC, a role code valid at every NM1)
  // attaches to every row sharing its element ref.
  const owningRows = (c: LookupCodeRow, segRows: LookupElementRow[]): LookupElementRow[] =>
    segRows.filter((r) => r.ref === c.element && (c.segmentName === undefined || r.segmentName === c.segmentName))

  return order.map((segment) => {
    const segRows = bySegment.get(segment)!
    const segCodes = codes.filter((c) => c.segment === segment)
    const orphanCodes = segCodes.filter((c) => owningRows(c, segRows).length === 0)
    const name = segmentNames[segment]
    return {
      segment,
      heading: name ? `${segment} · ${name}` : segment,
      elements: segRows.map((row) => ({ row, codes: segCodes.filter((c) => owningRows(c, segRows).includes(row)) })),
      orphanCodes,
    }
  })
}

// The key an element row is remembered by while its entry is open: the
// 835's use key, or the row's own place for a transaction without one.
export const rowKey = (row: LookupElementRow) => row.use ?? `${row.where}|${row.segmentName ?? ''}|${row.ref}`

// A coded value: the code, with its public name beneath, as every code on the
// docs site stands as a label (Code.tsx). It is a value of the element above
// it, not an element, so it is not a link; hovering shows the maintainer's
// full text where there is one.
function CodeRow({ row }: { row: LookupCodeRow }) {
  return (
    <div className="border-b border-line py-1.5 pl-3 last:border-b-0" title={row.full}>
      <span className="font-mono text-[13px] text-ink">{row.code}</span>
      {row.name && <span className={NAME_LINE}>{row.name}</span>}
      {row.about && <div className="mt-0.5 text-[12px] leading-4 text-ink">{row.about.text}</div>}
      <div className="mt-0.5 text-[11px] leading-4 text-ink-2">{row.about ? [row.source, ...row.about.sources].filter((s, i, all) => all.indexOf(s) === i).join('; ') : row.source}</div>
    </div>
  )
}

// The color of the validator tag: compared is blue, left out and outside
// the claim are gray, as on the Compared fields page; a qualifier that is
// part of the tag is gray too; the gap state is yellow.
const STATE_COLOR: Record<string, TagColor> = {
  compared: 'blue',
  'compared, watched': 'blue',
  'left out': 'gray',
  'outside the claim': 'gray',
  'part of the tag': 'gray',
  unlisted: 'yellow',
}

function Section({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="mt-2">
      <div className="text-[11px] font-semibold uppercase tracking-wide text-ink-2">{title}</div>
      <div className="mt-0.5 text-[13px] leading-5 text-ink">{children}</div>
    </div>
  )
}

// The entry under a row: the segment it sits in, what the element is, what
// it is for, its values, how Medicare fills it, what the validator does with
// it, and the documents behind each. Exported so a check can render every
// row's entry on the server.
export function LookupEntry({ row }: { row: LookupElementRow }) {
  const { about, segmentAbout, validator } = row
  const sources = about ? about.sources : []
  return (
    <div data-entry={rowKey(row)} className="mt-2 border-l-2 border-line pl-3 pb-1">
      {segmentAbout && (
        <Section title={`The ${row.segment} segment here`}>
          <p>{segmentAbout.is}</p>
          {segmentAbout.purpose && <p className="mt-1">{segmentAbout.purpose}</p>}
          {segmentAbout.medicare && <p className="mt-1">{segmentAbout.medicare}</p>}
        </Section>
      )}
      {about ? (
        <>
          <Section title="What it is">{about.is}</Section>
          {about.purpose && <Section title="What it is for">{about.purpose}</Section>}
          {about.values && <Section title="Values">{about.values}</Section>}
          {about.medicare && <Section title="How Medicare fills it">{about.medicare}</Section>}
        </>
      ) : (
        <Section title="What it is">{row.name ?? 'No public name on file.'}</Section>
      )}
      {validator && (
        <Section title="The validator">
          <div className="flex flex-wrap items-center gap-2">
            <span data-validator={validator.state}>
              <Tag color={STATE_COLOR[validator.state] ?? 'gray'}>{validator.state}</Tag>
            </span>
            <span className="text-[12px] leading-4 text-ink-2">{validator.why}</span>
          </div>
          {validator.detail && <div className="mt-1 font-mono text-[11px] leading-4 text-ink-2">{validator.detail}</div>}
        </Section>
      )}
      {row.echo && <Section title="From the 837">{row.echo}</Section>}
      <Section title="Sources">
        <ul className="list-disc pl-4 text-[11px] leading-4 text-ink-2">
          {[row.source, ...sources, ...(segmentAbout?.sources ?? [])]
            .filter((s, i, all) => all.indexOf(s) === i)
            .map((s) => (
              <li key={s}>{s}</li>
            ))}
        </ul>
      </Section>
    </div>
  )
}

// An element: its code, a link that searches the panel for it (so a click on
// CLP04 here narrows the list to CLP04, as a click on CLP04 on the page does),
// with its public name beneath; then where it sits, and its entry when open.
function ElementRow({ ge, open, onToggle }: { ge: GroupedElement; open: boolean; onToggle: () => void }) {
  const { row, codes } = ge
  const expandable = row.about !== undefined || row.validator !== undefined
  return (
    <div data-row={rowKey(row)} className="border-b border-line py-2 last:border-b-0">
      <div className="flex items-start gap-1.5">
        {expandable ? (
          <button type="button" onClick={onToggle} aria-expanded={open} aria-label={`${open ? 'Close' : 'Open'} the entry for ${row.ref}`}
            className="-ml-1 flex h-5 w-4 shrink-0 items-center justify-center text-ink-2 hover:text-ink">
            {open ? '▾' : '▸'}
          </button>
        ) : (
          <span className="-ml-1 w-4 shrink-0" />
        )}
        <div className="min-w-0 flex-1">
          <div className="flex items-baseline justify-between gap-2 text-[13px]">
            <CodeLink code={row.ref} />
            {row.echo && (
              <span title={row.echo} aria-label="From the 837" className="shrink-0 text-[12px] text-ink">
                ✓
              </span>
            )}
          </div>
          {row.name && <span className={NAME_LINE}>{row.name}</span>}
          <div className="mt-0.5 text-[11px] leading-4 text-ink-2">
            {row.segmentName && <span className="text-ink">{row.segmentName}</span>}
            {row.segmentName && ' · '}
            {row.where}
          </div>
          {!open && <div className="mt-0.5 text-[11px] leading-4 text-ink-2">{row.source}</div>}
          {open && <LookupEntry row={row} />}
          {codes.length > 0 && <div className="mt-1.5 border-l border-line">{codes.map((c, i) => <CodeRow key={`${c.code}-${i}`} row={c} />)}</div>}
        </div>
      </div>
    </div>
  )
}

export interface LookupPanelProps {
  title: string
  elements: LookupElementRow[]
  codes: LookupCodeRow[]
  segmentNames: Partial<Record<string, string>>
  query: string
  onQueryChange: (query: string) => void
  onClose: () => void
  footer?: ReactNode
}

// A searchable, grouped lookup of every element and coded value a
// transaction's public documents name, shown as a panel beside a page's
// content (LookupFrame.tsx), not a standalone page: the three viewers, the
// 835 compare page, and the Compared fields page. A row shows its code with
// the public name beneath, as every code on the docs site does (Code.tsx).
//
// A row stays short until it is opened; an opened row shows its entry
// (LookupEntry). Which rows are open is local to the panel, not in the
// address and not in storage: the address already names the panel and the
// search (lookup=, q=), and an open entry is a reading position, like the
// expanded groups of the compare page's tree, not something a link needs to
// reproduce. When the search narrows the list to one element, that one opens
// by itself, so typing CLP09, or clicking CLP09 on the page, shows CLP09's
// entry without a second click. A new search starts the list at its top.
export function LookupPanel({ title, elements, codes, segmentNames, query, onQueryChange, onClose, footer }: LookupPanelProps) {
  const namedElements = elements.filter((e) => e.name).length
  const namedCodes = codes.filter((c) => c.name).length
  const [opened, setOpened] = useState<Set<string>>(() => new Set())
  const toggle = (key: string) =>
    setOpened((prev) => {
      const next = new Set(prev)
      if (!next.delete(key)) next.add(key)
      return next
    })
  const list = useRef<HTMLDivElement>(null)
  useEffect(() => {
    list.current?.scrollTo({ top: 0 })
  }, [query])

  // An element row is also found by its tag (DTM*232), and a coded value by
  // its segment and code written the same way, so a click on DTM*232 or
  // NM1*QC lists that use's elements and the qualifier code under them.
  const words = query.trim().toLowerCase().split(/\s+/).filter(Boolean)
  const elementMatches = (r: LookupElementRow) => matches(words, r.ref, r.name, r.segment, r.segmentName, r.tag, r.where, r.source, r.about?.is)
  const codeMatches = (c: LookupCodeRow) => matches(words, c.element, c.code, `${c.segment}*${c.code}`, c.name, c.full, c.where, c.source, c.segmentName, c.about?.text)

  const groups = useMemo(() => buildGroups(elements, codes, segmentNames), [elements, codes, segmentNames])

  const visibleGroups = groups
    .map((group) => {
      const visibleElements = group.elements
        .map((ge) => {
          const visibleCodes = words.length === 0 ? ge.codes : ge.codes.filter(codeMatches)
          const visible = words.length === 0 || elementMatches(ge.row) || visibleCodes.length > 0
          return visible ? { row: ge.row, codes: visibleCodes } : null
        })
        .filter((ge): ge is GroupedElement => ge !== null)
      const orphanCodes = words.length === 0 ? group.orphanCodes : group.orphanCodes.filter(codeMatches)
      return { ...group, elements: visibleElements, orphanCodes }
    })
    .filter((group) => group.elements.length > 0 || group.orphanCodes.length > 0)

  const shownElements = visibleGroups.reduce((n, g) => n + g.elements.length, 0)
  const shownCodes = visibleGroups.reduce((n, g) => n + g.elements.reduce((m, ge) => m + ge.codes.length, 0) + g.orphanCodes.length, 0)
  const only = shownElements === 1 ? visibleGroups[0]?.elements[0]?.row : undefined
  const isOpen = (row: LookupElementRow) => opened.has(rowKey(row)) || (only !== undefined && rowKey(only) === rowKey(row))

  return (
    <div role="complementary" aria-label={title} className="flex h-full flex-col bg-bg">
      <div className="flex items-start justify-between gap-3 border-b border-line px-4 pb-3 pt-4">
        <div>
          <h2 className="font-serif text-lg text-ink">{title}</h2>
          <p className="mt-1 text-[11px] leading-4 text-ink-2">
            {namedElements} of {elements.length} element rows have a public name · {namedCodes} of {codes.length} coded-value rows have a public name
          </p>
        </div>
        <button type="button" onClick={onClose} aria-label="Close lookup" className="shrink-0 border border-line px-2 py-1 text-xs text-ink-2 hover:text-ink">
          Close ✕
        </button>
      </div>

      <div className="border-b border-line px-4 py-3">
        <label htmlFor="lookup-panel-search" className="mb-1.5 block text-xs font-semibold text-ink">
          Find an element, code, or word
        </label>
        <input
          id="lookup-panel-search"
          type="search"
          value={query}
          onChange={(event) => onQueryChange(event.target.value)}
          placeholder="Try CLP04, NM1*85, or diagnosis"
          className="w-full border border-line bg-[var(--color-bg)] px-3 py-2 font-mono text-sm text-ink focus:outline-2 focus:outline-accent"
        />
        <p className="mt-1.5 text-[11px] text-ink-2" role="status">
          {shownElements} elements · {shownCodes} codes
        </p>
      </div>

      <div ref={list} className="flex-1 overflow-y-auto px-4 py-2">
        {visibleGroups.length === 0 && <p className="py-6 text-sm text-ink-2">Nothing matches “{query}”.</p>}
        {visibleGroups.map((group) => (
          <section key={group.segment} className="py-3">
            <h3 className="mb-1 font-mono text-[11px] uppercase tracking-wide text-ink-2">{group.heading}</h3>
            <div>
              {group.elements.map((ge, i) => (
                <ElementRow key={`${ge.row.ref}-${i}`} ge={ge} open={isOpen(ge.row)} onToggle={() => toggle(rowKey(ge.row))} />
              ))}
              {group.orphanCodes.length > 0 && (
                <div className="mt-1 border-l border-line pl-1">
                  {group.orphanCodes.map((c, i) => (
                    <CodeRow key={`${c.element}-${c.code}-${i}`} row={c} />
                  ))}
                </div>
              )}
            </div>
          </section>
        ))}
      </div>

      {footer && <div className="border-t border-line px-4 py-3 text-[11px] leading-5 text-ink-2">{footer}</div>}
    </div>
  )
}
