import type { ReactNode } from 'react'
import { LookupPanel, type LookupCodeRow, type LookupElementRow } from './Lookup'
import { OpenInLookup, type LookupState } from './useLookup'

// A page's lookup panel and where it sits, shared by the three viewers, the
// 835 compare page, and the Compared fields page. The panel's state is the
// address's lookup= and q= (useLookup.ts).
//
// The panel sits beside the page's content, never over it: a column at the
// right, as tall as the window, that stays in view while the page scrolls and
// scrolls its own list. Opening it narrows the content's column by the
// panel's width. A page whose content is capped (each page's max-w) loses only
// margin while the column is wider than the cap, and reflows once it is
// narrower: the content's column is a size container, so a page sets its
// breakpoints on its own width (@min-[…]: classes) rather than the window's,
// and they hold with the panel open or closed (the compare page's
// breakpoints are its table's, tree's, and text panel's widths). The top
// padding leaves the theme toggle's corner clear.

// What a page's lookup lists: a transaction's element and code rows, and its
// segment names for the group headings.
export interface LookupData {
  title: string
  elements: LookupElementRow[]
  codes: LookupCodeRow[]
  segmentNames: Partial<Record<string, string>>
}

// The button that opens and closes the panel, in a page's own controls.
export function LookupButton({ lookup }: { lookup: LookupState }) {
  return (
    <button
      type="button"
      onClick={() => lookup.setOpen(!lookup.open)}
      aria-pressed={lookup.open}
      className={`shrink-0 border border-line px-3 py-1.5 text-xs ${lookup.open ? 'bg-ink text-bg' : 'text-ink-2 hover:text-ink'}`}
    >
      Lookup
    </button>
  )
}

// The page's content, and the panel beside it while it is open. Every code
// link inside opens this panel.
export function LookupFrame({ lookup, data, children }: { lookup: LookupState; data: LookupData; children: ReactNode }) {
  return (
    <OpenInLookup.Provider value={lookup.show}>
      <div className="flex items-start">
        <div className="@container min-w-0 flex-1">{children}</div>
        {lookup.open && (
          <div className="sticky top-0 h-svh w-[26rem] shrink-0 border-l border-line bg-bg pt-14">
            <LookupPanel
              title={data.title}
              elements={data.elements}
              codes={data.codes}
              segmentNames={data.segmentNames}
              query={lookup.query}
              onQueryChange={lookup.setQuery}
              onClose={() => lookup.setOpen(false)}
            />
          </div>
        )}
      </div>
    </OpenInLookup.Provider>
  )
}
