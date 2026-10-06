import type { ReactNode } from 'react'

interface PageTitleProps {
  children: ReactNode
  lede?: ReactNode
  // One sentence saying what the page answers (the app reads it from
  // src/app/help.ts); always visible, directly under the title.
  purpose?: string
  // Controls beside the title, such as a link to a related page.
  actions?: ReactNode
}

export function PageTitle({ children, lede, purpose, actions }: PageTitleProps) {
  const title = <h1 className="font-serif text-4xl tracking-tight text-ink">{children}</h1>
  return (
    <header className="mb-10">
      <div className="mb-6 h-1 w-10 bg-accent" />
      {actions ? <div className="flex flex-wrap items-baseline gap-x-6 gap-y-2">{title}{actions}</div> : title}
      {purpose && <p className="mt-2 max-w-[72ch] text-sm leading-6 text-ink-2">{purpose}</p>}
      {lede && <p className="mt-4 max-w-[65ch] text-[15px] leading-7 text-ink-2">{lede}</p>}
    </header>
  )
}
