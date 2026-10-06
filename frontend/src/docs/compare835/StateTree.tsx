import { useState } from 'react'
import { Tag, type TagColor } from '../../components/Tag'
import { CodeLink, NAME_LINE } from '../viewer/Code'
import { leavesOf, lookupOf, type TreeNode } from './tree'

// The element tree of the 835 compare page's selector (ListTree.tsx), drawn
// the same way: each part, its groups, and their elements, one under the
// other, with the same triangles, the same labels (a group's or an element's
// code, a link into the page's lookup, with its public name beneath; a part's
// title), and the node's key on hover. Where that selector has a box that
// shows or hides an element, this has a tag that says what the comparison
// does with it, and a part or group shows a tally of the tags below it, so a
// closed one still tells where it stands. The tree and the tags come from the
// page; nothing here decides a state.
//
// Every node starts open; which are closed is local to the page, not in the
// address and not in storage.

// A leaf's state, the color of its tag, and the text shown on hover (why
// the state is what it is).
export interface Mark<S extends string> {
  state: S
  color: TagColor
  title?: string
}

// The tag column's width: the widest tag, "compared, watched", with room.
const TAG_SLOT = 'w-[10.5rem] shrink-0'

// The ids of the nodes that have something under them.
const expandableIds = (node: TreeNode): string[] => (node.children.length > 0 ? [node.id, ...node.children.flatMap(expandableIds)] : [])

export function StateTree<S extends string>({ tree, markOf, tallyOf }: {
  tree: TreeNode[]
  markOf: (leaf: TreeNode) => Mark<S>
  tallyOf: (states: S[]) => string
}) {
  const [closed, setClosed] = useState<Set<string>>(() => new Set())
  const flip = (id: string) => setClosed((prev) => {
    const next = new Set(prev)
    if (!next.delete(id)) next.add(id)
    return next
  })

  const marks = new Map<string, Mark<S>>()
  for (const leaf of tree.flatMap(leavesOf)) marks.set(leaf.id, markOf(leaf))

  const entry = (node: TreeNode, depth: number) => {
    const indent = { paddingLeft: `${depth * 0.875 + 0.25}rem` }
    const mark = marks.get(node.id)
    if (node.children.length === 0 && mark) {
      return (
        <li key={node.id} data-leaf={node.id}>
          <div className="flex items-start" style={indent}>
            <span className="w-5 shrink-0" />
            <div className="min-w-0 flex-1 py-0.5 pr-2">
              <div className="flex min-h-6 items-center gap-3">
                <CodeLink code={node.codes} search={lookupOf(node)} title={node.id} className="min-w-0 [overflow-wrap:anywhere]" />
                {/* A dotted leader, so the eye follows the row from its code to its tag. */}
                <span aria-hidden="true" className="mb-1.5 min-w-3 flex-1 self-end border-b border-dotted border-line" />
                <span data-state={mark.state} title={mark.title} className={TAG_SLOT}>
                  <Tag color={mark.color}>{mark.state}</Tag>
                </span>
              </div>
              {/* The element's public name beneath its code, the sourced name the compare page's rows carry; an element with no public name shows the code alone. */}
              {node.names !== node.codes && <span className={NAME_LINE}>{node.names}</span>}
            </div>
          </div>
        </li>
      )
    }
    const open = !closed.has(node.id)
    const tally = tallyOf(leavesOf(node).flatMap((leaf) => marks.get(leaf.id)?.state ?? []))
    return (
      <li key={node.id} data-group={node.id}>
        <div className="flex items-start" style={indent}>
          <button type="button" onClick={() => flip(node.id)} aria-expanded={open} aria-label={`${open ? 'Collapse' : 'Expand'} ${node.codes}`}
            className="flex h-6 w-5 shrink-0 items-center justify-center text-ink-2 hover:text-ink">
            {open ? '▾' : '▸'}
          </button>
          {/* A part's title opens or closes it; a group's code opens the lookup, and its triangle opens or closes it. */}
          <div className="min-w-0 flex-1 py-1 pr-2">
            <div className="flex flex-wrap items-baseline gap-x-3">
              {depth === 0 ? (
                <button type="button" onClick={() => flip(node.id)} aria-expanded={open} title={node.id}
                  className="min-w-0 text-left font-semibold [overflow-wrap:anywhere] text-ink">
                  {node.codes}
                </button>
              ) : (
                <CodeLink code={node.codes} search={lookupOf(node)} title={node.id} className="min-w-0 [overflow-wrap:anywhere]" />
              )}
              <span data-tally={tally} className="font-mono text-[11px] text-ink-2">{tally}</span>
            </div>
            {node.names !== node.codes && <span className={NAME_LINE}>{node.names}</span>}
          </div>
        </div>
        {open && <ul>{node.children.map((child) => entry(child, depth + 1))}</ul>}
      </li>
    )
  }

  return (
    <div className="max-w-[32rem] border border-line text-xs">
      <div className="flex items-center gap-3 border-b border-line px-3 py-2">
        <span className="font-semibold text-ink">Elements by state</span>
        <span className="ml-auto flex gap-3">
          <button type="button" onClick={() => setClosed(new Set())} className="text-link hover:underline">Expand all</button>
          <button type="button" onClick={() => setClosed(new Set(tree.flatMap(expandableIds)))} className="text-link hover:underline">Collapse all</button>
        </span>
      </div>
      <ul className="py-1">{tree.map((node) => entry(node, 0))}</ul>
    </div>
  )
}
