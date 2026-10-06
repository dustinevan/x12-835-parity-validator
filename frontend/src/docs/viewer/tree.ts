// Transaction-agnostic loop-tree primitives shared by every X12
// viewer (frontend/src/docs/viewer835/ for the 835, frontend/src/docs/viewer837p/
// for the 837P, and any future one). A transaction-specific `parse.ts`
// builds a tree of these shapes from its own grammar (which segment opens
// which box, in what nesting — see frontend/src/docs/viewer837p/parse.ts for
// that per-transaction loop-rule table); everything below only walks a tree
// already built, so it carries no segment- or loop-ID knowledge of any one
// transaction.
export interface Segment {
  id: string
  // elements[0] is element 01, elements[1] is element 02, and so on.
  elements: string[]
  raw: string
}

export interface Box<BoxType extends string = string> {
  // Dotted path from the root, unique within one parsed file.
  path: string
  boxType: BoxType
  opening: Segment
  // Only envelope boxes (ISA/GS/ST in every transaction this project draws)
  // have a closing segment (IEA/GE/SE); a loop box closes with a plain
  // border, no segment of its own. This is also the generic signal
  // `envelopePaths` below uses to find the envelope chain without knowing
  // any transaction's actual box-type names.
  closing?: Segment
  children: Node<BoxType>[]
}

export type Node<BoxType extends string = string> = { kind: 'segment'; segment: Segment } | { kind: 'box'; box: Box<BoxType> }

export function splitSegments(raw: string): Segment[] {
  return raw
    .split('~')
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
    .map((s) => {
      const elements = s.split('*')
      return { id: elements[0], elements: elements.slice(1), raw: s }
    })
}

export function childBoxCount<T extends string>(box: Box<T>): number {
  return box.children.filter((n) => n.kind === 'box').length
}

function childBoxes<T extends string>(box: Box<T>): Box<T>[] {
  return box.children.filter((n): n is Extract<Node<T>, { kind: 'box' }> => n.kind === 'box').map((n) => n.box)
}

// The envelope chain a box sits in: every box from the root down while each
// one still has its own closing segment (ISA/GS/ST in both the 835 and the
// 837P — every transaction's envelope, never a repeating content loop).
// Collapsing this chain to one line is never the point of an example, so
// every example's `open` list can focus on the part it demonstrates.
export function envelopePaths<T extends string>(root: Box<T>): Set<string> {
  const open = new Set<string>()
  let box: Box<T> | undefined = root
  while (box && box.closing !== undefined) {
    open.add(box.path)
    box = childBoxes(box).find((b) => b.closing !== undefined)
  }
  return open
}

// Seeds the open-path set for an example: the envelope, plus every
// requested path and its ancestors (so a deep request is actually visible).
export function seedOpenPaths<T extends string>(root: Box<T>, requested: string[] | undefined): Set<string> {
  const open = envelopePaths(root)
  for (const path of requested ?? []) {
    const parts = path.split('.')
    for (let n = 1; n <= parts.length; n++) open.add(parts.slice(0, n).join('.'))
  }
  return open
}

export function findBox<T extends string>(root: Box<T>, path: string): Box<T> | undefined {
  if (root.path === path) return root
  for (const box of childBoxes(root)) {
    const found = findBox(box, path)
    if (found) return found
  }
  return undefined
}

// Every box path in the tree, root included — every example starts with all
// of them open, as does "Expand all"; "Collapse all" returns to
// seedOpenPaths' focused set.
export function allBoxPaths<T extends string>(root: Box<T>): Set<string> {
  const all = new Set<string>([root.path])
  for (const box of childBoxes(root)) for (const path of allBoxPaths(box)) all.add(path)
  return all
}

export function ancestorPaths(path: string): string[] {
  const parts = path.split('.')
  const ancestors: string[] = []
  for (let n = 1; n < parts.length; n++) ancestors.push(parts.slice(0, n).join('.'))
  return ancestors
}
