// Loads the app's MAC map through Vite's SSR module graph, the way the other
// checks do, and checks with no browser and no server:
//   - the fixture (src/app/fixtures/macs.json, the /api/macs response
//     body, saved from the backend's handler) holds 136 distinct numbers in
//     16 jurisdictions, in CMS's order;
//   - the panel's view (src/app/MacTree.tsx MacTreeView), rendered on the
//     server from the fixture, draws the JM subtree exactly as the drawing
//     below: the lines are read back from the rendered markup in document
//     order, with their depths, and drawn with the drawing's branches;
//   - every jurisdiction is drawn, the HHH numbers are marked, and the help
//     term, the search box, and the legend are there;
//   - the search (src/app/macLines.ts searchTree) keeps the lines it
//     should for a number, a state, a contractor, two words, and nonsense;
//   - the address's macs=1 opens the panel's button (aria-pressed) and is
//     carried by the selection's links (src/app/selection.tsx).
// To refresh the fixture from a running service:
//   curl -s localhost:8089/api/macs > src/app/fixtures/macs.json
// Run from frontend/:
//   node scripts/ssr-check-mactree.mjs
import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { StaticRouter } from 'react-router'
import { readFile } from 'node:fs/promises'
import { createServer } from 'vite'

// The drawing the MAC map follows, for one jurisdiction. The contractor's
// name is written as CMS prints it.
const JM_DRAWING = `JM  Palmetto GBA, LLC
├── Part A  (institutional claims, FISS)
│   ├── roll-up 11001  "SC and the HHH claims"
│   │   └── 11201  South Carolina
│   ├── roll-up 11003  "VA and WV"
│   │   ├── 11301  Virginia
│   │   └── 11401  West Virginia
│   ├── 11501  North Carolina                   (no roll-up named)
│   └── 11004  home health and hospice, Region C
└── Part B  (professional claims, MCS)
    ├── 11202  South Carolina
    ├── 11302  Virginia
    ├── 11402  West Virginia
    └── 11502  North Carolina`

const ORDER = ['J5', 'J6', 'J8', 'J15', 'JE', 'JF', 'JH', 'JJ', 'JK', 'JL', 'JM', 'JN', 'JA', 'JB', 'JC', 'JD']

// Quotes and runs of spaces are a matter of type: compare without them.
const plain = (text) => text.replace(/[“”]/g, '"').replace(/ +/g, ' ').trim()
const unescape = (text) => text.replace(/&quot;/g, '"').replace(/&#x27;/g, "'").replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&amp;/g, '&')
const count = (html, text) => html.split(text).length - 1

// Draws lines given in document order with their depths, as the drawing
// does: the depth-0 line, then each line under it with its branch.
function draw(lines) {
  const root = { text: lines[0].text, children: [] }
  const stack = [root]
  for (const line of lines.slice(1)) {
    const node = { text: line.text, children: [] }
    stack.length = line.depth
    stack[line.depth - 1].children.push(node)
    stack.push(node)
  }
  const out = [root.text]
  const drawUnder = (node, indent) => {
    node.children.forEach((child, i) => {
      const last = i === node.children.length - 1
      out.push(`${indent}${last ? '└── ' : '├── '}${child.text}`)
      drawUnder(child, indent + (last ? '    ' : '│   '))
    })
  }
  drawUnder(root, '')
  return out.join('\n')
}

// The ids in a part of the response, roll-ups and the numbers under them.
function idsOf(data) {
  const ids = []
  for (const j of data.jurisdictions) {
    for (const p of j.parts) {
      for (const r of p.roll_ups) ids.push(r.id, ...r.workloads.map((w) => w.id))
      for (const list of [p.workloads, p.hhh, p.other]) ids.push(...list.map((w) => w.id))
    }
  }
  return ids
}

// The ids a list of lines holds.
function linesIds(nodes) {
  const ids = []
  const add = (node) => {
    if (node.id !== undefined) ids.push(node.id)
    for (const child of node.children) add(child)
  }
  for (const node of nodes) add(node)
  return ids
}

async function main() {
  const server = await createServer({
    configFile: false,
    root: process.cwd(),
    plugins: [],
    server: { middlewareMode: true },
    appType: 'custom',
    logLevel: 'warn',
  })

  try {
    const { MacTreeView, MacMapButton, useMacMap } = await server.ssrLoadModule('/src/app/MacTree.tsx')
    const lines = await server.ssrLoadModule('/src/app/macLines.ts')
    const { SelectionProvider, useSelection } = await server.ssrLoadModule('/src/app/selection.tsx')
    const data = JSON.parse(await readFile(new URL('../src/app/fixtures/macs.json', import.meta.url), 'utf8'))

    const failures = []
    const fail = (message) => failures.push(message)

    // ---- the fixture
    const ids = idsOf(data)
    if (ids.length !== 136) fail(`the fixture holds ${ids.length} numbers, not 136`)
    if (new Set(ids).size !== ids.length) fail('a number appears more than once in the fixture')
    const names = data.jurisdictions.map((j) => j.jurisdiction)
    if (JSON.stringify(names) !== JSON.stringify(ORDER)) fail(`the jurisdictions are ${names.join(', ')}, not CMS's order`)

    // ---- the panel from the fixture
    const render = (element, location = '/?macs=1') => renderToStaticMarkup(createElement(StaticRouter, { location }, element))
    const html = render(createElement(MacTreeView, { data, onClose: () => {} }))
    const drawn = [...html.matchAll(/data-node="([^"]*)" data-depth="(\d+)"/g)].map((m) => ({ text: unescape(m[1]), depth: Number(m[2]) }))
    const nodes = lines.macNodes(data)
    let parts = 0
    for (const j of nodes) parts += j.children.length
    if (drawn.length !== linesIds(nodes).length + nodes.length + parts) {
      fail(`the panel draws ${drawn.length} lines; the tree has ${linesIds(nodes).length} numbers, ${nodes.length} jurisdictions, and ${parts} Parts`)
    }
    const tops = drawn.map((line, i) => (line.depth === 0 ? i : -1)).filter((i) => i >= 0)
    const drawnNames = tops.map((i) => drawn[i].text.split(' ')[0])
    if (JSON.stringify(drawnNames) !== JSON.stringify(ORDER)) fail(`the panel draws the jurisdictions ${drawnNames.join(', ')}`)
    const jmAt = tops[ORDER.indexOf('JM')]
    const jmEnd = tops[ORDER.indexOf('JM') + 1] ?? drawn.length
    const jm = draw(drawn.slice(jmAt, jmEnd))
    const want = JM_DRAWING.split('\n').map(plain).join('\n')
    const got = jm.split('\n').map(plain).join('\n')
    if (got !== want) fail(`the JM subtree draws as\n${jm}\nnot\n${JM_DRAWING}`)
    if (count(html, '>HHH</span>') !== 4) fail(`${count(html, '>HHH</span>')} HHH numbers are marked, not 4 (06004, 11004, 14014, 15004)`)
    for (const text of ['aria-label="About Workloads and roll-ups"', 'type="search"', 'Find a number, state, or contractor', 'How a number reads', '136 numbers in 16 jurisdictions', 'no states named', 'Internal processing number']) {
      if (!html.includes(text)) fail(`the panel does not show “${text}”`)
    }
    for (const j of data.jurisdictions) if (!html.includes(`>${j.contractor.replace(/&/g, '&amp;')}<`)) fail(`the panel does not name ${j.jurisdiction}'s contractor`)

    // ---- the search
    const search = (query) => linesIds(lines.searchTree(nodes, query)).sort()
    const tests = [
      ['11201', ['11001', '11201']],
      ['virginia', ['11003', '11301', '11302', '11401', '11402', '15004', '18003']],
      ['palmetto', [...idsOf({ jurisdictions: data.jurisdictions.filter((j) => j.contractor.startsWith('Palmetto')) })].sort()],
      ['palmetto virginia', ['11003', '11301', '11302', '11401', '11402']],
      ['JM', [...idsOf({ jurisdictions: data.jurisdictions.filter((j) => j.jurisdiction === 'JM') })].sort()],
      ['region c', ['11004']],
      ['queens', ['13292']],
      ['zzz', []],
    ]
    for (const [query, wantIds] of tests) {
      const gotIds = search(query)
      if (JSON.stringify(gotIds) !== JSON.stringify([...wantIds].sort())) fail(`searching “${query}” keeps ${gotIds.join(', ') || 'nothing'}, not ${wantIds.join(', ') || 'nothing'}`)
    }
    if (lines.searchTree(nodes, '  ') !== nodes) fail('an empty search does not keep the whole tree')

    // ---- the address: macs=1 presses the button, and links carry it
    function Probe() {
      const map = useMacMap()
      const sel = useSelection()
      return createElement('div', null, createElement(MacMapButton, { map }), createElement('a', { href: sel.link('/claims') }, 'link'))
    }
    const probe = (location) =>
      renderToStaticMarkup(createElement(StaticRouter, { location }, createElement(SelectionProvider, null, createElement(Probe))))
    const opened = probe('/?macs=1')
    if (!opened.includes('aria-pressed="true"')) fail('macs=1 does not press the MAC map button')
    if (!opened.includes('href="/claims?macs=1"')) fail(`a link under macs=1 does not carry it: ${opened}`)
    const shut = probe('/')
    if (!shut.includes('aria-pressed="false"')) fail('the MAC map button is pressed with no macs=1')
    if (shut.includes('macs=')) fail('a link carries macs= with the panel closed')

    // ---- nothing in the map's view reads the network at render time
    for (const file of ['/src/app/macLines.ts']) {
      const code = await readFile(new URL(`..${file}`, import.meta.url), 'utf8')
      if (/\bfetch\(|\bapi\./.test(code)) fail(`${file} reads the network`)
    }

    if (failures.length > 0) {
      console.error('MAC map SSR check failed:')
      for (const f of failures) console.error(`  - ${f}`)
      process.exitCode = 1
      return
    }
    console.log(`MAC map SSR check passed: the fixture holds ${ids.length} numbers in ${names.length} jurisdictions in CMS's order; the panel draws ${drawn.length} lines, the JM subtree exactly as the drawing, ${count(html, '>HHH</span>')} HHH marks; ${tests.length} searches keep what they should; macs=1 presses the button and rides on links.`)
  } finally {
    await server.close()
  }
}

main().catch((err) => {
  console.error(err)
  process.exitCode = 1
})
