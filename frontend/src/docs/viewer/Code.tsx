import { useContext } from 'react'
import { OpenInLookup } from './useLookup'

// How the docs site shows an X12 code, the same on every page that has a
// lookup panel:
//   - A code is a link into the page's lookup: a click opens the panel with
//     the code in its search (the panel is LookupFrame.tsx; the address's
//     lookup= and q= are useLookup.ts). It is monospace in the link tone,
//     underlined on hover.
//   - A code's public name sits beneath it, small and dim, the one minor
//     line (NAME_LINE). The compare page's row key set the style
//     (compare835/ListGrid.tsx).
//   - A code that labels a value (a viewer cell) has the value beside it on
//     the code's line (Parts.tsx).
//
// A code here is an element reference (CLP04, SVC01, CLM05-1), a segment ID
// (CLP), or a segment with its qualifier (DTM*232, NM1*QC). A value that is
// itself a code (CLP02's 1, a CARC) is a value, not a link: a search for a
// value finds no element. It shows its name on hover instead (Term.tsx).

// The public name beneath a code.
export const NAME_LINE = 'block text-[11px] font-normal leading-tight text-ink-2 [overflow-wrap:anywhere]'

// One code as a link. `search` is what the lookup is searched for when it
// differs from the code shown (the compare page's 'MOA03–MOA07' searches
// MOA03); `title` is the text a hover shows, where a page has one. The click
// stays on the code: a code inside a box's opening row opens the lookup, not
// the box.
export function CodeLink({ code, search, title, className }: { code: string; search?: string; title?: string; className?: string }) {
  const open = useContext(OpenInLookup)
  const term = search ?? code
  if (open === null) {
    return (
      <span title={title} className={`font-mono ${className ?? ''}`}>
        {code}
      </span>
    )
  }
  return (
    <button
      type="button"
      title={title}
      data-code={term}
      onClick={(e) => {
        e.stopPropagation()
        open(term)
      }}
      className={`cursor-pointer text-left font-mono text-link underline-offset-2 hover:underline focus-visible:underline ${className ?? ''}`}
    >
      {code}
    </button>
  )
}

// A code standing alone as a label: the code, then its public name beneath
// when it has one. The code's size comes from the element it sits in.
export function CodeName({ code, name, search, title }: { code: string; name?: string; search?: string; title?: string }) {
  return (
    <>
      <CodeLink code={code} search={search} title={title} />
      {name !== undefined && name !== code && <span className={NAME_LINE}>{name}</span>}
    </>
  )
}
