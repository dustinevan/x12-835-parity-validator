import { useMemo } from 'react'
import { splitX12, terminatorSuffix } from './x12split'

interface Line {
  text: string // the segment as the file sent it, without its terminator
  segment: boolean // false for a line of text that is not a readable segment
}

// One line per segment, split on the file's own delimiters (so a file that
// arrives on one line, or with CRLF or another terminator, still reads one
// segment per line). Text with no readable ISA is shown as the lines it has
// and never throws, since the panel also shows text mid-edit.
function linesOf(file: string): { lines: Line[]; suffix: string } {
  try {
    const { segments, delimiters } = splitX12(file)
    return { lines: segments.map((s) => ({ text: s.raw, segment: true })), suffix: terminatorSuffix(delimiters) }
  } catch {
    const lines = file.split(/\r?\n/)
    while (lines.length > 0 && lines[lines.length - 1].trim() === '') lines.pop()
    return { lines: lines.map((text) => ({ text, segment: false })), suffix: '' }
  }
}

// The raw X12 file as read-only text: one segment per line, nothing added (no
// line numbers, no annotations) beyond restoring each segment's own trailing
// terminator, which the file itself holds — splitting trims it only to
// parse. Long segments scroll horizontally per line rather than wrapping,
// since a wrapped line reads as more lines than the file has.
// Hovering a box row in the viewer above highlights its line here
// (`hoveredRaw` is that row's raw segment text).
export function RawFile({ file, hoveredRaw }: { file: string; hoveredRaw: string | null }) {
  const { lines, suffix } = useMemo(() => linesOf(file), [file])
  return (
    <div className="overflow-x-auto border border-line bg-surface p-4 font-mono text-[12px] text-ink">
      {lines.map((line, i) => (
        <div key={i} className={`whitespace-nowrap px-1 ${line.segment && line.text === hoveredRaw ? 'bg-line' : ''}`}>
          {line.text}
          {line.segment ? suffix : ''}
        </div>
      ))}
    </div>
  )
}
