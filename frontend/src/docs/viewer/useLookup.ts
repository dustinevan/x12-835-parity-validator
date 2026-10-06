import { createContext, useEffect } from 'react'
import { useSearchParams } from 'react-router'

// The lookup panel's state, kept in the address: lookup=1 while it is open and
// q= its search, so a link reproduces the panel. show(code) is what a click on
// a code does: open the panel with that code in the search. The panel and
// where it sits are LookupFrame.tsx; a code's link is Code.tsx.
export interface LookupState {
  open: boolean
  query: string
  setOpen: (open: boolean) => void
  setQuery: (query: string) => void
  show: (code: string) => void
}

// Opens the page's lookup panel searched for a code. LookupFrame provides it;
// with no frame above it, a code is plain monospace text (Code.tsx).
export const OpenInLookup = createContext<((code: string) => void) | null>(null)

export function useLookup(): LookupState {
  const [params, setParams] = useSearchParams()
  const open = params.get('lookup') === '1'
  const query = params.get('q') ?? ''

  const change = (apply: (next: URLSearchParams) => void) =>
    setParams((prev) => {
      const next = new URLSearchParams(prev)
      apply(next)
      return next
    })
  const setOpen = (on: boolean) =>
    change((next) => {
      if (on) next.set('lookup', '1')
      else {
        next.delete('lookup')
        next.delete('q')
      }
    })
  const setQuery = (text: string) =>
    change((next) => {
      if (text) next.set('q', text)
      else next.delete('q')
    })
  const show = (code: string) =>
    change((next) => {
      next.set('lookup', '1')
      next.set('q', code)
    })

  // Escape closes the panel. setParams changes with the address, so the
  // listener is set again with it and never writes an address the page has
  // left.
  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key !== 'Escape') return
      setParams((prev) => {
        const next = new URLSearchParams(prev)
        next.delete('lookup')
        next.delete('q')
        return next
      })
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [open, setParams])

  return { open, query, setOpen, setQuery, show }
}
