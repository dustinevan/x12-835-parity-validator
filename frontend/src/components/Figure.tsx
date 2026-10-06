import type { ReactNode } from 'react'

interface FigureProps {
  number: number
  caption: string
  children: ReactNode
}

// A numbered figure: bordered, horizontally scrollable when wide, captioned.
// Figures are inline SVG drawn with theme tokens so they work in both themes.
export function Figure({ number, caption, children }: FigureProps) {
  return (
    <figure className="my-6 border border-line bg-bg">
      <div className="overflow-x-auto p-4">{children}</div>
      <figcaption className="border-t border-line px-4 py-2 font-mono text-[12px] text-ink-2">
        Figure {number}. {caption}
      </figcaption>
    </figure>
  )
}
