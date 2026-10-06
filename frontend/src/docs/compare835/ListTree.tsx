import { useState, type ReactNode } from 'react'
import { CodeLink, NAME_LINE } from '../viewer/Code'
import { checkOf, lookupOf, toggled, type Check, type TreeNode } from './tree'

// The selector beside the list: a file tree of the files' content, each node
// with a box. An unchecked node hides its rows; a parent checks or unchecks
// its children. The unchecked nodes are the address's hide= (their ids,
// tree.ts treeOf); which nodes are expanded is local to the page. The preset
// controls (PresetControls.tsx) come in through `top`, under the header. A
// label shows its id on hover: the key the preset text uses (presets.ts).
//
// A part (Payment, Claim, …) is a title. A group's or an element's code is a
// link into the page's lookup, with the public name beneath it, as every code
// on the docs site stands as a label (viewer/Code.tsx).

// A square box in the page's line color that never changes, with a smaller
// solid square inside it when the node is shown, a half-width one when only
// some of a parent's children are shown, and nothing when hidden; it has no
// check mark, rounded corners, or color change. A button with the checkbox
// role, so the keyboard and screen readers treat
// it as one.
const INNER: Record<Check, string> = {
  on: 'h-2 w-2 bg-series-1',
  off: '',
  mixed: 'h-2 w-1 bg-series-1',
}

function Box({ state, label, onChange }: { state: Check; label: string; onChange: () => void }) {
  return (
    <button type="button" role="checkbox" aria-checked={state === 'mixed' ? 'mixed' : state === 'on'} aria-label={label} onClick={onChange}
      className="mt-0.5 flex h-3.5 w-3.5 shrink-0 items-center justify-center border border-ink-2 bg-bg">
      {state !== 'off' && <span className={INNER[state]} />}
    </button>
  )
}

export function ListTree({ tree, hidden, onChange, top }: {
  tree: TreeNode[]
  hidden: Set<string>
  onChange: (hidden: Set<string>) => void
  top?: ReactNode
}) {
  // Every part starts open and its groups closed; a click on the arrow flips a node.
  const [flipped, setFlipped] = useState<Set<string>>(() => new Set())
  const flip = (id: string) => setFlipped((prev) => {
    const next = new Set(prev)
    if (!next.delete(id)) next.add(id)
    return next
  })

  const entry = (node: TreeNode, depth: number) => {
    const expandable = node.children.length > 0
    const expanded = (depth === 0) !== flipped.has(node.id)
    return (
      <li key={node.id}>
        <div className="flex items-start" style={{ paddingLeft: `${depth * 0.875 + 0.25}rem` }}>
          {expandable ? (
            <button type="button" onClick={() => flip(node.id)} aria-expanded={expanded} aria-label={`${expanded ? 'Collapse' : 'Expand'} ${node.codes}`}
              className="flex h-6 w-5 shrink-0 items-center justify-center text-ink-2 hover:text-ink">
              {expanded ? '▾' : '▸'}
            </button>
          ) : <span className="w-5 shrink-0" />}
          {/* The box alone toggles show/hide; a part's title expands or
              collapses it, and a code opens the lookup. */}
          <div className="flex min-w-0 flex-1 items-start gap-2 py-1 pr-2">
            <Box state={checkOf(node, hidden)} label={`Show ${node.codes}`} onChange={() => onChange(toggled(tree, hidden, node))} />
            {depth === 0 ? (
              <button type="button" onClick={() => flip(node.id)} aria-expanded={expanded} title={node.id}
                className="min-w-0 text-left font-semibold [overflow-wrap:anywhere] text-ink">
                {node.codes}
              </button>
            ) : (
              <span className="min-w-0 [overflow-wrap:anywhere]">
                <CodeLink code={node.codes} search={lookupOf(node)} title={node.id} />
                {node.names !== node.codes && <span className={NAME_LINE}>{node.names}</span>}
              </span>
            )}
          </div>
        </div>
        {expandable && expanded && <ul>{node.children.map((child) => entry(child, depth + 1))}</ul>}
      </li>
    )
  }

  // Beside the table once the page is 72rem wide (the tree's 13rem, the
  // table's 52rem, and their gap and padding; Compare835.tsx), above it under
  // that. The width is the page's own, not the window's, so it holds with
  // the lookup panel open (viewer/LookupFrame.tsx).
  return (
    <nav aria-label="Elements shown" className="w-full shrink-0 self-start border border-line text-xs @min-[72rem]:sticky @min-[72rem]:top-6 @min-[72rem]:max-h-[calc(100svh-3rem)] @min-[72rem]:w-52 @min-[72rem]:overflow-y-auto">
      <div className="flex items-center gap-3 border-b border-line px-3 py-2">
        <span className="font-semibold text-ink">Elements shown</span>
        <span className="ml-auto flex gap-3">
          <button type="button" onClick={() => onChange(new Set())} className="text-link hover:underline">All</button>
          <button type="button" onClick={() => onChange(new Set(tree.map((node) => node.id)))} className="text-link hover:underline">None</button>
        </span>
      </div>
      {top}
      <ul className="py-1">{tree.map((node) => entry(node, 0))}</ul>
    </nav>
  )
}
