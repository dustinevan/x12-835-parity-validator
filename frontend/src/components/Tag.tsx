import type { ReactNode } from 'react'

export type TagColor = 'blue' | 'gray' | 'green' | 'yellow' | 'red'

const COLORS: Record<TagColor, string> = {
  blue: 'bg-accent/15 text-link',
  gray: 'bg-surface text-ink-2',
  green: 'bg-ok/15 text-ok',
  yellow: 'bg-warn/20 text-ink',
  red: 'bg-danger/15 text-danger',
}

interface TagProps {
  color?: TagColor
  children: ReactNode
}

export function Tag({ color = 'gray', children }: TagProps) {
  return (
    <span
      className={`inline-block whitespace-nowrap px-2 py-1 font-mono text-xs leading-none ${COLORS[color]}`}
    >
      {children}
    </span>
  )
}
