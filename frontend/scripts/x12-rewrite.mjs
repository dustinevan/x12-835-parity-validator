// Helpers for the viewers' SSR checks: rewrite a worked example in other
// delimiters, and reduce a parsed loop tree to a list of strings that can be
// compared across two spellings of the same file.

// Spellings of one file a reader might paste. `after` is written after every
// segment. The bar and angle bracket, the plus and less-than sign, appear in
// none of the example files, so a rewrite never collides with their data.
export const DELIMITER_VARIANTS = [
  { name: 'one long line: | > ~', element: '|', component: '>', after: '~' },
  { name: 'LF terminator: + <', element: '+', component: '<', after: '\n' },
  { name: 'CRLF line ends: | >', element: '|', component: '>', after: '\r\n' },
  { name: 'standard * : ~, CRLF after each, BOM and blank lines around', element: '*', component: ':', after: '~\r\n', before: '﻿\r\n\r\n', trailer: '\r\n\r\n' },
]

// `splitSegments` is the viewer's own, read from its parse module; it gives
// each segment's id and elements with the component separator as ':'. ISA16
// (element index 15) is the component separator itself, so it is written as
// the variant's own.
export function rewriteDelimiters(file, splitSegments, variant) {
  for (const ch of [variant.element, variant.component]) {
    if (ch !== '*' && ch !== ':' && file.includes(ch)) {
      throw new Error(`the example already contains ${JSON.stringify(ch)}, so a rewrite with it would be ambiguous`)
    }
  }
  const text = splitSegments(file)
    .map((s) => {
      const elements = s.elements.map((e, i) => (s.id === 'ISA' && i === 15 ? variant.component : e.split(':').join(variant.component)))
      return [s.id, ...elements].join(variant.element) + variant.after
    })
    .join('')
  return (variant.before ?? '') + text + (variant.trailer ?? '')
}

// Every box and row of a tree, in order, as one string each: path, kind, and
// the segment's id and elements. ISA16 is left out (it is the separator).
export function shape(root) {
  const out = []
  const seg = (s) => `${s.id}*${(s.id === 'ISA' ? s.elements.slice(0, 15) : s.elements).join('*')}`
  const walk = (box) => {
    out.push(`${box.path} ${box.boxType} ${seg(box.opening)}`)
    for (const node of box.children) {
      if (node.kind === 'box') walk(node.box)
      else out.push(`${box.path} row ${seg(node.segment)}`)
    }
    if (box.closing) out.push(`${box.path} close ${seg(box.closing)}`)
  }
  walk(root)
  return out
}
