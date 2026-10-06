// The nodes the page's selector lists, derived from the rows so it always
// matches them: a part (the payment, a claim, a service line, …), its groups
// of segments, and each group's elements. A node's id is its path, the part
// first ('claim', 'claim.CLP', 'claim.CLP.CLP03'; 'line.REF.REF02'), and the
// address's hide= lists the ids of the nodes unchecked. A part and its groups
// are the same in every claim and every service line, so one box hides a
// group or an element in all of them.
import { SCOPES, SCOPE_TITLE, type ListRow, type ListSection } from './rows'

export interface TreeNode {
  id: string
  codes: string // the label: a part's title, a group's code, an element's code
  names: string // the public name shown beneath the code (the code again where there is none)
  children: TreeNode[]
}

export const rowPath = (row: ListRow): string[] => [row.scope, `${row.scope}.${row.group}`, `${row.scope}.${row.group}.${row.leaf}`]

// The code a group or element node searches the lookup for: the last part of
// its key ('claim.NM1*QC' searches NM1*QC; 'claim.MOA.MOA03', shown as
// MOA03–MOA07, searches MOA03; 'claim.DTM.DTM*232' searches DTM*232).
export const lookupOf = (node: TreeNode): string => node.id.slice(node.id.lastIndexOf('.') + 1)

const distinct = (values: (string | undefined)[]) => values.filter((v, i, all): v is string => v !== undefined && all.indexOf(v) === i)

export function treeOf(sections: ListSection[]): TreeNode[] {
  return SCOPES.flatMap((scope) => {
    const rows = sections.filter((s) => s.scope === scope).flatMap((s) => s.rows)
    if (rows.length === 0) return []
    const groups = distinct(rows.map((r) => r.group))
    const children = groups.map((group): TreeNode => {
      const own = rows.filter((r) => r.group === group)
      const id = `${scope}.${group}`
      const leaves = distinct(own.map((r) => r.leaf)).map((leaf): TreeNode => {
        const same = own.filter((r) => r.leaf === leaf)
        return {
          id: `${id}.${leaf}`,
          codes: same[0].leafCodes,
          names: distinct(same.map((r) => r.name)).join(' · ') || same[0].leafCodes,
          children: [],
        }
      })
      return { id, codes: group, names: distinct(own.map((r) => r.groupName)).join(' · ') || group, children: leaves }
    })
    return [{ id: scope, codes: SCOPE_TITLE[scope], names: SCOPE_TITLE[scope], children }]
  })
}

// ---------------------------------------------------------------- the checkboxes

// A node's box. Its leaves decide it: fully shown, fully hidden, or partly.
// A node is hidden when its own id, or the id of a node above it, is in the
// hidden set (hide=).

export type Check = 'on' | 'off' | 'mixed'

export const leavesOf = (node: TreeNode): TreeNode[] => (node.children.length ? node.children.flatMap(leavesOf) : [node])

// 'claim.CLP.CLP03' is under 'claim.CLP', which is under 'claim'.
const pathOf = (id: string) => {
  const parts = id.split('.')
  return parts.map((_, i) => parts.slice(0, i + 1).join('.'))
}

const isHidden = (id: string, hidden: Set<string>) => pathOf(id).some((p) => hidden.has(p))

export function checkOf(node: TreeNode, hidden: Set<string>): Check {
  const leaves = leavesOf(node)
  const off = leaves.filter((l) => isHidden(l.id, hidden)).length
  return off === 0 ? 'on' : off === leaves.length ? 'off' : 'mixed'
}

export const idsOf = (node: TreeNode): string[] => [node.id, ...node.children.flatMap(idsOf)]

// The shortest list of ids that hides exactly the leaves not in `shown`: a
// node is listed when every leaf under it is hidden, else its children are
// looked at.
function hiddenFor(node: TreeNode, shown: Set<string>, out: Set<string>) {
  const leaves = leavesOf(node)
  const count = leaves.filter((l) => shown.has(l.id)).length
  if (count === leaves.length) return
  if (count === 0) {
    out.add(node.id)
    return
  }
  for (const child of node.children) hiddenFor(child, shown, out)
}

// The shortest hidden set that keeps exactly the leaves `keep` accepts, by the
// same rule a click uses. Every id in it is a node of this tree. The built-in
// presets (presets.ts) are made with it.
export function hiddenKeeping(tree: TreeNode[], keep: (leaf: TreeNode) => boolean): string[] {
  const shown = new Set(tree.flatMap(leavesOf).filter(keep).map((leaf) => leaf.id))
  const out = new Set<string>()
  for (const root of tree) hiddenFor(root, shown, out)
  return [...out]
}

// The hidden ids after a click on `node`: a click on a fully checked node
// hides it, and a click on an unchecked or partly checked one shows it. Hidden
// ids the tree has no node for (from a link made on other files) are kept.
export function toggled(tree: TreeNode[], hidden: Set<string>, node: TreeNode): Set<string> {
  const show = checkOf(node, hidden) !== 'on'
  const under = new Set(leavesOf(node).map((l) => l.id))
  const shown = new Set<string>()
  for (const root of tree) {
    for (const leaf of leavesOf(root)) if (under.has(leaf.id) ? show : !isHidden(leaf.id, hidden)) shown.add(leaf.id)
  }
  const known = new Set(tree.flatMap(idsOf))
  const next = new Set([...hidden].filter((id) => !known.has(id)))
  for (const root of tree) hiddenFor(root, shown, next)
  return next
}
