// Checks the pieces behind the viewers' editable raw panel, with no browser:
//  1. The 837I parser reads its delimiters from the ISA: every example,
//     rewritten in four other delimiter spellings (see x12-rewrite.mjs), gives
//     the same segments and the same tree. (ssr-check-837p.mjs does the 837P.)
//  2. The delimiter reader refuses text with no usable ISA, with a message.
//  3. tryParse (what the panel runs on every keystroke) accepts every example
//     of all three viewers, refuses broken text with a message instead of
//     throwing, and, over many random edits of every example, never throws and
//     never accepts a tree that the boxes cannot draw (each accepted tree is
//     rendered to HTML, every box open).
//  4. Each of the three pages renders on the server with the panel in its read
//     view: an Edit button, no Reset and no "(edited)" (the text is the
//     example's own), the line saying pasted text stays in the page, and no
//     error line. Each page has no Codes/Names toggle, every cell's code is a
//     link into the lookup, and with lookup=1&q=ISA06 the lookup panel opens
//     in its own column beside the content, its search holding ISA06 and
//     showing that row. The panel is also rendered from props in its other
//     states: textarea and Done; Reset, Edit, and the parser error; a
//     one-line file shown one segment per line with its own terminator.
//  5. Informational: whether the 835 parser reads the one-long-line spelling
//     yet.
// Run from frontend/:
//   node scripts/ssr-check-viewer-edit.mjs
import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { StaticRouter } from 'react-router'
import { createServer } from 'vite'
import { DELIMITER_VARIANTS, rewriteDelimiters, shape } from './x12-rewrite.mjs'

// A small seeded generator, so a failure repeats.
function prng(seed) {
  let a = seed >>> 0
  return () => {
    a = (a + 0x6d2b79f5) >>> 0
    let t = a
    t = Math.imul(t ^ (t >>> 15), t | 1)
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61)
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

const MUTATIONS_PER_EXAMPLE = 25

async function main() {
  const server = await createServer({
    configFile: false,
    root: process.cwd(),
    plugins: [],
    server: { middlewareMode: true },
    appType: 'custom',
    logLevel: 'warn',
  })

  const failures = []
  const fail = (msg) => failures.push(msg)

  try {
    const load = (path) => server.ssrLoadModule(path)
    const { readDelimiters } = await load('/src/docs/viewer/x12split.ts')
    const { tryParse } = await load('/src/docs/viewer/useEditableRaw.ts')
    const { BoxView } = await load('/src/docs/viewer/BoxView.tsx')
    const { allBoxPaths } = await load('/src/docs/viewer/tree.ts')

    const viewers = []
    for (const v of [
      { name: '835', data: '/src/data/viewer835Examples.ts', key: 'VIEWER_835_EXAMPLES', dir: 'viewer835' },
      { name: '837P', data: '/src/data/viewer837p/index.ts', key: 'VIEWER_837P_EXAMPLES', dir: 'viewer837p' },
      { name: '837I', data: '/src/data/viewer837i/index.ts', key: 'VIEWER_837I_EXAMPLES', dir: 'viewer837i' },
    ]) {
      const parse = await load(`/src/docs/${v.dir}/parse.ts`)
      const describe = await load(`/src/docs/${v.dir}/describe.ts`)
      const examples = (await load(v.data))[v.key]
      viewers.push({ ...v, parse, describe, examples })
    }

    // 1. 837I delimiters.
    const i837 = viewers.find((v) => v.name === '837I')
    for (const ex of i837.examples) {
      const total = i837.parse.splitSegments(ex.file).length
      const root = i837.parse.parseLoopTree(ex.file)
      for (const variant of DELIMITER_VARIANTS) {
        try {
          const rewritten = rewriteDelimiters(ex.file, i837.parse.splitSegments, variant)
          const count = i837.parse.splitSegments(rewritten).length
          if (count !== total) {
            fail(`837I ${ex.id} (${variant.name}): ${count} segments read, the file has ${total}`)
            continue
          }
          if (JSON.stringify(shape(i837.parse.parseLoopTree(rewritten))) !== JSON.stringify(shape(root))) {
            fail(`837I ${ex.id} (${variant.name}): the tree differs from the standard-delimiter tree`)
          }
        } catch (err) {
          fail(`837I ${ex.id} (${variant.name}): ${err instanceof Error ? err.message : err}`)
        }
      }
    }

    // 2. The delimiter reader on text with no usable ISA.
    const isa = i837.examples[0].file.split('\n')[0]
    const badHeaders = {
      empty: '',
      'blank text': '   \n  ',
      'no ISA': 'GS*HC*A*B~',
      'ISA alone': 'ISA',
      'ISA cut short': isa.slice(0, 60),
      'no terminator after ISA16': isa.slice(0, isa.lastIndexOf('~')),
      'separators the same': 'ISA' + '*'.repeat(16) + '*',
    }
    for (const [name, text] of Object.entries(badHeaders)) {
      try {
        readDelimiters(text)
        fail(`readDelimiters accepted ${name}`)
      } catch (err) {
        if (!(err instanceof Error) || !err.message.startsWith('x12:')) fail(`readDelimiters on ${name}: unexpected error ${err}`)
      }
    }
    const d = readDelimiters(rewriteDelimiters(i837.examples[0].file, i837.parse.splitSegments, DELIMITER_VARIANTS[1])).delimiters
    if (d.element !== '+' || d.component !== '<' || d.terminator !== '\n') fail(`readDelimiters read ${JSON.stringify(d)} from the + < LF spelling`)

    // 3. tryParse on every example, on broken text, and under random edits.
    let accepted = 0
    let refused = 0
    for (const v of viewers) {
      const parse = v.parse.parseLoopTree
      const describe = v.describe
      const style = { boxClass: () => '', openingClass: () => '' }
      const render = (root) => {
        const props = {
          box: root,
          openPaths: allBoxPaths(root),
          onToggle: () => {},
          onMaximize: () => {},
          style,
          describeOpening: describe.describeOpening,
          describeRow: describe.describeRow,
          describeClosing: describe.describeClosing,
        }
        renderToStaticMarkup(createElement(BoxView, props))
      }

      for (const ex of v.examples) {
        const r = tryParse(ex.file, parse, describe)
        if (!r.ok) fail(`${v.name} ${ex.id}: tryParse refused the example itself: ${r.error}`)
      }

      const first = v.examples[0].file
      const lines = first.split('\n')
      const cases = {
        empty: '',
        'one word': 'hello',
        'no SE': lines.filter((l) => !l.startsWith('SE*')).join('\n'),
        'no GE': lines.filter((l) => !l.startsWith('GE*')).join('\n'),
        'no IEA': lines.filter((l) => !l.startsWith('IEA*')).join('\n'),
        'half the file': lines.slice(0, Math.floor(lines.length / 2)).join('\n'),
        'segment ID in lower case': first.replace('\nGS*', '\ngs*'),
      }
      const expect = { 'no SE': 'ST has no closing SE', 'no GE': 'GS has no closing GE', 'no IEA': 'ISA has no closing IEA' }
      for (const [name, text] of Object.entries(cases)) {
        let r
        try {
          r = tryParse(text, parse, describe)
        } catch (err) {
          fail(`${v.name} tryParse threw on ${name}: ${err}`)
          continue
        }
        if (r.ok) fail(`${v.name} tryParse accepted ${name}`)
        else if (expect[name] && !r.error.includes(expect[name])) fail(`${v.name} ${name}: expected "${expect[name]}", got "${r.error}"`)
      }

      // Random edits: nothing may throw; whatever is accepted must render.
      const rand = prng(20261001)
      const pick = (n) => Math.floor(rand() * n)
      const noise = ['*', '~', ':', '|', '\n', ' ', 'A', '9', '^']
      for (const ex of v.examples) {
        for (let n = 0; n < MUTATIONS_PER_EXAMPLE; n++) {
          let text = ex.file
          const edits = 1 + pick(3)
          for (let e = 0; e < edits; e++) {
            const at = pick(text.length)
            switch (pick(6)) {
              case 0:
                text = text.slice(0, at) + text.slice(at + 1)
                break
              case 1:
                text = text.slice(0, at) + noise[pick(noise.length)] + text.slice(at)
                break
              case 2: {
                const rows = text.split('\n')
                rows.splice(pick(rows.length), 1)
                text = rows.join('\n')
                break
              }
              case 3: {
                const rows = text.split('\n')
                const k = pick(rows.length)
                rows.splice(k, 0, rows[k])
                text = rows.join('\n')
                break
              }
              case 4:
                text = text.slice(0, at)
                break
              default: {
                const rows = text.split('\n')
                const k = pick(rows.length)
                rows[k] = rows[k].split('*').slice(0, 1 + pick(4)).join('*') + '~'
                text = rows.join('\n')
              }
            }
          }
          let r
          try {
            r = tryParse(text, parse, describe)
          } catch (err) {
            fail(`${v.name} ${ex.id}: tryParse threw on an edited copy: ${err}`)
            continue
          }
          if (!r.ok) {
            refused++
            if (typeof r.error !== 'string' || r.error === '') fail(`${v.name} ${ex.id}: a refusal with no message`)
            continue
          }
          accepted++
          try {
            render(r.root)
          } catch (err) {
            fail(`${v.name} ${ex.id}: an accepted edit would not render: ${err instanceof Error ? err.stack : err}`)
          }
        }
      }
    }

    // 4. The pages, rendered once, then once more with the lookup open. Every
    //    cell's code is a link into the lookup (data-code, docs/viewer/
    //    Code.tsx), there is no Codes/Names toggle, and the lookup panel sits
    //    in its own column beside the content, never fixed over it.
    for (const [name, path, exportName] of [
      ['835', '/src/docs/pages/Viewer835.tsx', 'Viewer835'],
      ['837P', '/src/docs/pages/Viewer837P.tsx', 'Viewer837P'],
      ['837I', '/src/docs/pages/Viewer837I.tsx', 'Viewer837I'],
    ]) {
      try {
        const Page = (await load(path))[exportName]
        const html = renderToStaticMarkup(createElement(StaticRouter, { location: '/' }, createElement(Page)))
        for (const want of ['Raw file', '>Edit<', 'stays in this page']) if (!html.includes(want)) fail(`${name} page: the markup lacks ${JSON.stringify(want)}`)
        for (const dont of ['(edited)', '>Reset<', '<textarea']) if (html.includes(dont)) fail(`${name} page: the unedited read view contains ${JSON.stringify(dont)}`)
        for (const dont of ['>Codes</button>', '>Names</button>']) if (html.includes(dont)) fail(`${name} page: the Codes/Names toggle is still there`)
        for (const want of ['data-code="ISA"', 'data-code="ISA01"', 'data-code="ISA06"']) if (!html.includes(want)) fail(`${name} page: no lookup link ${want}`)
        if (html.includes('role="complementary"')) fail(`${name} page: the lookup panel is open with no lookup=1`)

        const open = renderToStaticMarkup(createElement(StaticRouter, { location: '/?lookup=1&q=ISA06' }, createElement(Page)))
        const at = open.indexOf('role="complementary"')
        if (at < 0) fail(`${name} page: lookup=1 does not open the lookup panel`)
        else {
          const panel = open.slice(at)
          if (!panel.includes('value="ISA06"')) fail(`${name} page: the panel's search does not hold q=ISA06`)
          if (!panel.includes('data-row=')) fail(`${name} page: the panel searched for ISA06 shows no row`)
        }
        if (!open.includes('class="@container min-w-0 flex-1"') || !/class="sticky top-0 h-svh w-\[26rem\][^"]*"/.test(open)) fail(`${name} page: the panel is not a column beside the content`)
        if (/class="(?:[^"]* )?fixed(?: [^"]*)?"/.test(open)) fail(`${name} page: something on the page is fixed over it`)
      } catch (err) {
        fail(`${name} page did not render: ${err instanceof Error ? err.stack : err}`)
      }
    }

    // The panel itself in its other states (it is driven only by its props).
    {
      const { RawPanel } = await load('/src/docs/viewer/RawEditor.tsx')
      const text = i837.examples[0].file
      const base = { root: null, text, edited: false, editing: false, error: null, setText: () => {}, setEditing: () => {}, reset: () => {} }
      const html = (raw) => renderToStaticMarkup(createElement(RawPanel, { raw: { ...base, ...raw }, hoveredRaw: null }))
      const editing = html({ editing: true })
      if (!editing.includes('<textarea') || !editing.includes('>Done<') || editing.includes('>Reset<')) fail('RawPanel: editing and unedited should show a textarea and Done, and no Reset')
      const edited = html({ editing: false, edited: true, error: 'x12: reading the ISA: the header is incomplete' })
      if (!edited.includes('>Reset<') || !edited.includes('>Edit<') || edited.includes('<textarea')) fail('RawPanel: edited, read view should show Reset and Edit and no textarea')
      if (!edited.includes('x12: reading the ISA: the header is incomplete') || !edited.includes('last text that parsed')) fail('RawPanel: the parser error is not shown under the panel')
      const oneLine = html({ text: rewriteDelimiters(text, i837.parse.splitSegments, DELIMITER_VARIANTS[0]) })
      if ((oneLine.match(/whitespace-nowrap/g) ?? []).length !== i837.parse.splitSegments(text).length) fail('RawPanel: a one-line file is not shown one segment per line')
      if (!oneLine.includes('~</div>')) fail('RawPanel: the file\'s own terminator is not shown after each segment')
      const crlf = html({ text: rewriteDelimiters(text, i837.parse.splitSegments, DELIMITER_VARIANTS[2]) })
      if (crlf.includes('~</div>') || crlf.includes('\r')) fail('RawPanel: a line-break terminator should add no character')
    }

    // 5. Informational: the 835 parser's own delimiter handling.
    const v835 = viewers.find((v) => v.name === '835')
    const sample = v835.examples[0].file
    const spelled = rewriteDelimiters(sample, i837.parse.splitSegments, DELIMITER_VARIANTS[0])
    const r835 = tryParse(spelled, v835.parse.parseLoopTree, v835.describe)
    console.log(`835 parser on the one-long-line | > ~ spelling: ${r835.ok ? 'reads it' : `does not read it yet (${r835.error})`}`)

    if (failures.length > 0) {
      console.error(`Viewer edit SSR check FAILED (${failures.length}):`)
      for (const f of failures.slice(0, 40)) console.error(`  - ${f}`)
      process.exitCode = 1
      return
    }
    console.log(
      `Viewer edit SSR check passed: ${i837.examples.length} 837I examples read under ${DELIMITER_VARIANTS.length} delimiter spellings; ` +
        `${accepted + refused} random edits (${accepted} accepted and rendered, ${refused} refused with a message), none threw.`,
    )
  } finally {
    await server.close()
  }
}

main().catch((err) => {
  console.error(err)
  process.exitCode = 1
})
