// The tree of every element the worked examples send, for the Compared fields
// page: the union of the trees the 835 compare page builds (tree.ts treeOf)
// for each example alone, so each part, group, and element appears once, in
// the order that page lists them (rows.ts: the payment, payer, payee, provider
// summary, each claim then its service lines, then the provider adjustments;
// a claim's and a service line's groups in rows.ts ORDER). The envelope is not
// in it, as it is not in that page's tree (file.ts).
//
// One difference from that page's tree. A segment that carries a qualifier
// (a date kind, a reference kind, an amount kind, and so on) has two elements
// there, DTM01 and DTM02, shared by every qualifier, so an element cannot
// stand one way for DTM*232 and another for DTM*050. Here a group of such a
// segment lists its qualified tags instead ('DTM*232', 'DTM*050'), in the
// order the examples send them, with the key 'claim.DTM.DTM*232'. NM1 is
// already one group per role on both pages. Every other node is the same,
// with the same key and label.
//
// Nothing in this file uses React or reads storage or the address.
import {
  CLAIM_AMT_NAME,
  CLAIM_PER_NAME,
  CLAIM_QTY_NAME,
  CLAIM_REF_NAME,
  DTM_NAME,
  HEADER_REF_NAME,
  LQ_NAME,
  PAYER_PER_NAME,
  REF_PAYEE_EL_NAME,
  REF_PAYER_EL_NAME,
  SERVICE_AMT_NAME,
  SERVICE_REF_NAME,
} from '../../data/x12Names835'
import type { ParsedFile } from './file'
import { COLS, listOf, type ListSection } from './rows'
import { treeOf, type TreeNode } from './tree'

// The public name of a qualified tag ('DTM*232', 'REF*6R') in a scope, from
// the sourced names in x12Names835.ts, the same names the compare page's rows
// carry; a tag those names do not cover returns the tag itself, so the leaf
// shows its code alone. Nothing is invented here.
function nameOfTag(scope: string, tag: string): string {
  const star = tag.indexOf('*')
  if (star < 0) return tag
  const segment = tag.slice(0, star)
  const kind = tag.slice(star + 1)
  switch (segment) {
    case 'DTM':
      return DTM_NAME[kind] ?? tag
    case 'REF':
      if (scope === 'payment') return HEADER_REF_NAME[kind] ?? tag
      if (scope === 'payer') return REF_PAYER_EL_NAME[1]
      if (scope === 'payee') return REF_PAYEE_EL_NAME[1]
      if (scope === 'claim') return CLAIM_REF_NAME[kind] ?? tag
      if (scope === 'line') return SERVICE_REF_NAME[kind] ?? tag
      return tag
    case 'AMT':
      if (scope === 'claim') return CLAIM_AMT_NAME
      if (scope === 'line') return SERVICE_AMT_NAME
      return tag
    case 'QTY':
      return scope === 'claim' ? CLAIM_QTY_NAME : tag
    case 'LQ':
      return LQ_NAME[kind] ?? tag
    case 'PER':
      if (scope === 'payer') return PAYER_PER_NAME[kind] ?? tag
      if (scope === 'claim') return CLAIM_PER_NAME[kind] ?? tag
      return tag
    default:
      return tag
  }
}

// The segments whose first element names a kind. NM1 is left out: its tree
// groups are already one per role.
const TAGGED = new Set(['DTM', 'REF', 'AMT', 'QTY', 'LQ', 'PER'])

// The qualified tags of each tagged group in a list, by the group's key
// ('claim.DTM'), in the order the list's rows reach them. The tag is read from
// the segment a row came from, as rows.ts names it: 'DTM*232', or 'DTM' alone
// when the segment's first element is empty.
function tagsOf(sections: ListSection[]): Map<string, string[]> {
  const tags = new Map<string, string[]>()
  for (const row of sections.flatMap((s) => s.rows)) {
    if (!TAGGED.has(row.group)) continue
    const segment = COLS.map((c) => row.from[c]?.[0]).find((s) => s !== undefined)
    if (!segment) continue
    const kind = segment.elements[0] ?? ''
    const tag = kind ? `${segment.id}*${kind}` : segment.id
    const id = `${row.scope}.${row.group}`
    const list = tags.get(id) ?? []
    if (!list.includes(tag)) list.push(tag)
    tags.set(id, list)
  }
  return tags
}

// One example's tree: the compare page's, with each tagged group's elements
// replaced by its tags.
function treeOfExample(file: ParsedFile): TreeNode[] {
  const sections = listOf({ ref: file })
  const tags = tagsOf(sections)
  return treeOf(sections).map((part) => ({
    ...part,
    children: part.children.map((group) => {
      const list = tags.get(group.id)
      if (!list) return group
      return { ...group, children: list.map((tag): TreeNode => ({ id: `${group.id}.${tag}`, codes: tag, names: nameOfTag(part.id, tag), children: [] })) }
    }),
  }))
}

// The ids of each list merged into one, by the rule rows.ts mergeOrder uses
// for its columns (that function is not exported): the first list's order; an
// id only a later list has goes after the id before it there; one with no id
// before it there goes before the first later id the merged list has, else
// last. Every list keeps its own order.
function mergeOrder(lists: string[][]): string[] {
  const merged: string[] = []
  for (const list of lists) {
    let last = -1
    list.forEach((id, i) => {
      const at = merged.indexOf(id)
      if (at >= 0) {
        last = at
        return
      }
      if (last < 0) {
        const next = list.slice(i + 1).map((later) => merged.indexOf(later)).find((index) => index >= 0)
        last = (next ?? merged.length) - 1
      }
      merged.splice(last + 1, 0, id)
      last += 1
    })
  }
  return merged
}

// Nodes with one id are one node; its children are the union of theirs.
function unionOf(lists: TreeNode[][]): TreeNode[] {
  const first = new Map<string, TreeNode>()
  const below = new Map<string, TreeNode[][]>()
  for (const list of lists) {
    for (const node of list) {
      if (!first.has(node.id)) first.set(node.id, node)
      below.set(node.id, [...(below.get(node.id) ?? []), node.children])
    }
  }
  return mergeOrder(lists.map((list) => list.map((node) => node.id))).map((id) => ({
    ...first.get(id)!,
    children: unionOf(below.get(id) ?? []),
  }))
}

// A claim's and a service line's groups in the order rows.ts lists them (its
// ORDER, which is not exported): the head segment's group first, then the
// groups of these segments in this order, a segment not named last. Groups of
// one segment keep the order the examples give them. The other parts keep the
// order the files send, which the merge above gives.
const ORDER: Record<string, { head: string; then: string[] }> = {
  claim: { head: 'CLP', then: ['NM1', 'MOA', 'MIA', 'DTM', 'REF', 'PER', 'AMT', 'QTY', 'CAS'] },
  line: { head: 'SVC', then: ['DTM', 'CAS', 'REF', 'AMT', 'LQ'] },
}

function inListOrder(part: TreeNode): TreeNode {
  const order = ORDER[part.id]
  if (!order) return part
  const rank = (group: TreeNode) => {
    const segment = group.codes.split('*')[0]
    if (segment === order.head) return -1
    const at = order.then.indexOf(segment)
    return at < 0 ? order.then.length : at
  }
  return { ...part, children: [...part.children].sort((x, y) => rank(x) - rank(y)) }
}

// The union tree of the given files, one example each.
export function unionTreeOf(files: ParsedFile[]): TreeNode[] {
  return unionOf(files.map(treeOfExample)).map(inListOrder)
}
