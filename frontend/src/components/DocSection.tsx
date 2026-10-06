import type { ReactNode } from 'react'

interface DocSectionProps {
  id: string
  title: string
  children: ReactNode
  className?: string
}

// One titled prose section: a heading (deep-linkable by id) plus body copy.
export function DocSection({ id, title, children, className }: DocSectionProps) {
  return (
    <section
      id={id}
      className={`scroll-mt-8 border-t border-line py-10 first:border-t-0 first:pt-0 ${className ?? ''}`}
    >
      <h2 className="mb-4 text-[22px] font-semibold text-ink">
        <a href={`#${id}`} className="hover:text-link">
          {title}
        </a>
      </h2>
      <div className="space-y-4 text-[15px] leading-7 text-ink-2 [&_code]:bg-surface [&_code]:px-1.5 [&_code]:py-0.5 [&_code]:font-mono [&_code]:text-[13px] [&_code]:text-ink [&_strong]:text-ink">
        {children}
      </div>
    </section>
  )
}
