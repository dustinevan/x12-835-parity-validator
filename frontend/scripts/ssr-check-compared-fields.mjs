// Loads the Compared fields page through Vite's SSR module graph, the way
// ssr-check-835-compare.mjs loads the compare page, and checks with no browser
// that
//   - the page's tree ("In file order", compare835/unionTree.ts) is the union
//     of the trees the compare page builds for the examples: every element any
//     example sends has a leaf and every leaf comes from one, each id is held
//     once, and each section's groups and leaves keep the order the compare
//     page lists them in;
//   - the page, rendered on the server, draws every leaf of that tree once,
//     with exactly one state tag: the state stateOf gives it, in blue, gray,
//     or yellow only; and every part and group shows the tally of its leaves;
//   - every group's and leaf's code, and every tag of the leave-out list, is a
//     link into the 835 lookup (docs/viewer/Code.tsx), with the public name
//     beneath the code in the tree; a part's title opens and closes it; and
//     with lookup=1&q=CLP03 the lookup panel opens beside the page, its search
//     holding CLP03 and showing that element;
//   - the lede states the five-cent rule with the cents of AmountMatchCents
//     in backend/lib/core/service/compare.go, read as text, which
//     AMOUNT_MATCH_CENTS in data/comparedFields.ts must equal;
//   - the leave-out list on the page (LEFT_OUT in data/comparedFields.ts) is
//     the leave-out list of compare set 1 in backend/lib/core/service/
//     compareset.go, which this script reads as text: the same level and tag on
//     both sides, none on one side only, and each entry drawn once on the page
//     with its public name and its reason;
//   - the four watched elements are the ones the page names, and the page says
//     that they are compared and recorded like every other element;
//   - stateOf answers as the rules read for keys the examples do not send (a
//     qualifier or role no rule mentions, an unplaced segment, a payment-level
//     key); "left out" is only ever inside a claim, so the elements outside it
//     (PLB, BPR01, the N1 loops) are "outside the claim"; and no element is
//     "unlisted": inside a claim and covered by neither the compared rule nor
//     the left-out rule;
//   - the page has two sections, tree and left-out, the docs navigation
//     (docs/App.tsx) lists the same two, and nothing else is on the page;
//   - nothing in the page, the data, or the tree's code reads or writes
//     storage or the network, and the only thing the page puts in the address
//     is the lookup panel's lookup= and q=, through the lookup's own code
//     (docs/viewer/useLookup.ts), which writes no other key.
// It prints the state counts, and the left-out leaves by where their state
// comes from.
// Run from frontend/:
//   node scripts/ssr-check-compared-fields.mjs
import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { StaticRouter } from 'react-router'
import { readFile } from 'node:fs/promises'
import { createServer } from 'vite'

const TAGGED = new Set(['DTM', 'REF', 'AMT', 'QTY', 'LQ', 'PER'])
const distinct = (list) => list.filter((x, i) => list.indexOf(x) === i)
const twice = (list) => distinct(list.filter((x, i) => list.indexOf(x) !== i))
const reEscape = (text) => text.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
const count = (html, text) => html.split(text).length - 1
// Text as React writes it into markup.
const escapeHtml = (text) => text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#x27;')

// The text of some markup: tags removed, the entities React writes decoded,
// white space made single.
const textOf = (markup) =>
  markup
    .replace(/<[^>]+>/g, '')
    .replace(/&#x27;/g, "'")
    .replace(/&quot;/g, '"')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&amp;/g, '&')
    .replace(/\s+/g, ' ')
    .trim()

// Whether `small` appears in `big` in the same order, not necessarily side by side.
function isSubsequence(small, big) {
  let at = 0
  for (const item of big) if (item === small[at]) at++
  return at === small.length
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
    const { VIEWER_835_EXAMPLES } = await server.ssrLoadModule('/src/data/viewer835Examples.ts')
    const { parsedOf } = await server.ssrLoadModule('/src/docs/compare835/columns.ts')
    const { COLS, listOf } = await server.ssrLoadModule('/src/docs/compare835/rows.ts')
    const { leavesOf } = await server.ssrLoadModule('/src/docs/compare835/tree.ts')
    const { unionTreeOf } = await server.ssrLoadModule('/src/docs/compare835/unionTree.ts')
    const data = await server.ssrLoadModule('/src/data/comparedFields.ts')
    const { ComparedFields } = await server.ssrLoadModule('/src/docs/pages/ComparedFields.tsx')
    const { LEFT_OUT, NAME_ELEMENTS, WATCHED, COMPARED, LEAF_STATES, stateOf, tallyText } = data

    const examples = VIEWER_835_EXAMPLES
    if (!Array.isArray(examples) || examples.length === 0) throw new Error('VIEWER_835_EXAMPLES is empty, nothing to check')

    const failures = []
    const fail = (message) => failures.push(message)

    // ---- the union tree: what the examples send, once each, in the compare page's order
    const files = examples.map((ex) => parsedOf(ex.id))
    const tree = unionTreeOf(files)
    const everyNode = (nodes) => nodes.flatMap((n) => [n, ...everyNode(n.children)])
    const nodes = everyNode(tree)
    const leaves = tree.flatMap(leavesOf)
    const ids = nodes.map((n) => n.id)
    if (new Set(ids).size !== ids.length) fail(`the union tree holds an id twice: ${ids.filter((id, i) => ids.indexOf(id) !== i).join(', ')}`)
    if (leaves.some((l) => l.children.length > 0)) fail('a leaf of the union tree has children')
    if (tree.some((part) => part.children.length === 0 || part.children.some((g) => g.children.length === 0))) fail('a part or group of the union tree holds nothing')

    // The key a row's element has on this page, derived again here from the rows
    // and the segments they came from: the compare page's key, except that a
    // segment with a qualifier has one element per qualified tag.
    const keyOf = (row) => {
      if (!TAGGED.has(row.group)) return `${row.scope}.${row.group}.${row.leaf}`
      const segment = COLS.map((c) => row.from[c]?.[0]).find((s) => s !== undefined)
      const kind = segment.elements[0] ?? ''
      return `${row.scope}.${row.group}.${kind ? `${segment.id}*${kind}` : segment.id}`
    }
    const unionGroups = new Map() // part id -> its group ids
    const unionLeaves = new Map() // group id -> its leaf ids
    for (const part of tree) {
      unionGroups.set(part.id, part.children.map((g) => g.id))
      for (const group of part.children) unionLeaves.set(group.id, group.children.map((l) => l.id))
    }
    const fromRows = new Set()
    let rowCount = 0
    let sectionCount = 0
    for (const ex of examples) {
      for (const section of listOf({ ref: parsedOf(ex.id) })) {
        sectionCount++
        const rows = section.rows
        rowCount += rows.length
        for (const row of rows) fromRows.add(keyOf(row))
        // the groups of one section, and each group's elements, keep their order in the union
        const groups = distinct(rows.map((r) => `${r.scope}.${r.group}`))
        if (!isSubsequence(groups, unionGroups.get(section.scope) ?? [])) fail(`${ex.id}: the groups of “${section.title}” (${groups.join(', ')}) are not in the union's order`)
        for (const group of groups) {
          const own = distinct(rows.filter((r) => `${r.scope}.${r.group}` === group).map(keyOf))
          if (!isSubsequence(own, unionLeaves.get(group) ?? [])) fail(`${ex.id}: the elements of ${group} in “${section.title}” (${own.join(', ')}) are not in the union's order`)
        }
      }
    }
    const leafIds = new Set(leaves.map((l) => l.id))
    for (const key of fromRows) if (!leafIds.has(key)) fail(`an example sends ${key}, which has no leaf`)
    for (const id of leafIds) if (!fromRows.has(id)) fail(`the leaf ${id} comes from no row of any example`)
    // The order of a claim's groups and of a service line's is rows.ts's own. No
    // one example holds every pair of groups (none has both MOA and MIA), so ask
    // rows.ts directly: a claim and a line holding one segment of each group the
    // examples send, in reverse, listed by rows.ts, give the order the union must have.
    {
      const first = (segments) => {
        const seen = new Map()
        for (const s of segments) {
          const key = s.id === 'NM1' ? `NM1*${s.elements[0]}` : s.id
          if (!seen.has(key)) seen.set(key, s)
        }
        return [...seen.values()].reverse()
      }
      const base = files.find((f) => f.claims.some((c) => c.lines.length > 0))
      const baseClaim = base.claims.find((c) => c.lines.length > 0)
      const claim = {
        ...baseClaim,
        segs: first(files.flatMap((f) => f.claims.flatMap((c) => c.segs))),
        lines: [{ ...baseClaim.lines[0], segs: first(files.flatMap((f) => f.claims.flatMap((c) => c.lines.flatMap((l) => l.segs)))) }],
      }
      const sections = listOf({ ref: { ...base, claims: [claim] } })
      for (const scope of ['claim', 'line']) {
        const own = distinct(sections.filter((s) => s.scope === scope).flatMap((s) => s.rows).map((r) => r.group.split('*')[0]))
        const union = distinct(tree.find((p) => p.id === scope).children.map((g) => g.codes.split('*')[0]))
        if (JSON.stringify(own) !== JSON.stringify(union)) fail(`the segments of ${scope} in the union go ${union.join(', ')}; the compare page's rows go ${own.join(', ')}`)
        if (own.length < 5) fail(`the order check for ${scope} found only ${own.join(', ')}`)
      }
    }
    // the parts come in the compare page's order
    const PARTS = ['payment', 'payer', 'payee', 'summary', 'claim', 'line', 'adjustment']
    if (JSON.stringify(tree.map((p) => p.id)) !== JSON.stringify(PARTS)) fail(`the parts are ${tree.map((p) => p.id).join(', ')}`)

    // ---- the leave-out list against compare set 1 in the code (compareset.go, read as text)
    const read = (path) => readFile(new URL(`../src/${path}`, import.meta.url), 'utf8')
    const goSource = await readFile(new URL('../../backend/lib/core/service/compareset.go', import.meta.url), 'utf8')
    const setBlock = /var compareSet1 = CompareSet\{[\s\S]*?\n\}\n/.exec(goSource)?.[0]
    let codeList = []
    if (!setBlock) fail('compareset.go has no "var compareSet1 = CompareSet{" block')
    else {
      if (!/\bName:\s*"1"/.test(setBlock)) fail('compareSet1 in compareset.go is not named "1"')
      // The entries are written {LevelClaim, "TAG", "reason"}: the level, then the tag.
      codeList = [...setBlock.matchAll(/\{\s*Level(Claim|Line)\s*,\s*"([^"]+)"/g)].map((m) => `${m[1].toLowerCase()}:${m[2]}`)
      if (codeList.length === 0) fail('found no leave-out entries in compareSet1 of compareset.go (this check reads entries written {LevelClaim, "TAG", ...})')
    }
    const pageList = LEFT_OUT.map((e) => `${e.level}:${e.tag}`)
    for (const e of twice(codeList)) fail(`compareset.go lists ${e} twice`)
    for (const e of twice(pageList)) fail(`the page's leave-out list holds ${e} twice`)
    for (const e of codeList) if (!pageList.includes(e)) fail(`compareset.go leaves out ${e}, and the page's list does not`)
    for (const e of pageList) if (!codeList.includes(e)) fail(`the page's list leaves out ${e}, and compareset.go does not`)
    for (const e of LEFT_OUT) {
      const name = `${e.level}:${e.tag}`
      if (e.level !== 'claim' && e.level !== 'line') fail(`${name} has the level “${e.level}”, not claim or line`)
      for (const field of ['tag', 'publicName', 'reason']) {
        if (typeof e[field] !== 'string' || e[field].trim() === '' || e[field].includes('\n')) fail(`${name} has no one-line ${field}`)
      }
      if (e.keys.length === 0 || e.keys.some((k) => !k.startsWith(`${e.level}.`))) fail(`${name} has keys that are not at its level: ${e.keys.join(', ')}`)
    }

    // ---- the other state data
    if (JSON.stringify(LEAF_STATES) !== JSON.stringify(['compared', 'compared, watched', 'left out', 'outside the claim', 'unlisted'])) fail(`the states are ${LEAF_STATES.join(', ')}`)
    if (JSON.stringify(WATCHED.map((r) => r.label)) !== JSON.stringify(['DTM*050', 'AMT*I', 'AMT*NL', 'NM1*TT'])) fail(`the watched list is ${WATCHED.map((r) => r.label).join(', ')}`)
    for (const w of WATCHED) if (pageList.includes(`claim:${w.label}`)) fail(`${w.label} is both watched and left out`)
    for (const entry of [...WATCHED, ...COMPARED, NAME_ELEMENTS]) if (entry.keys.length === 0) fail(`“${entry.label}” has no keys`)
    if (NAME_ELEMENTS.keys.some((k) => !k.startsWith('claim.'))) fail(`the name elements have keys outside a claim: ${NAME_ELEMENTS.keys.join(', ')}`)

    // ---- stateOf on keys the examples do not send
    const cases = [
      ['claim.DTM.DTM*999', 'unlisted'], // a qualifier no rule mentions
      ['claim.DTM.DTM*472', 'unlisted'], // DTM*472 is a service-line element
      ['line.DTM.DTM*472', 'compared'],
      ['claim.REF.REF*TJ', 'unlisted'], // REF*TJ is a line element and a payee element
      ['claim.REF.REF*6R', 'unlisted'], // REF*6R is left out in a service line only
      ['line.REF.REF*TJ', 'compared'],
      ['line.REF.REF*6R', 'left out'],
      ['payee.REF.REF*TJ', 'outside the claim'],
      ['payee.REF.REF*XX', 'outside the claim'],
      ['claim.NM1*XX.NM103', 'left out'], // names, every role
      ['claim.NM1*XX.NM109', 'unlisted'], // a role no rule mentions
      ['claim.NM1*82.NM109', 'compared'],
      ['claim.NM1*82.NM104', 'left out'],
      ['claim.NM1*QC.NM109', 'left out'],
      ['claim.NM1*QC.NM103', 'left out'],
      ['claim.NM1*TT.NM109', 'compared, watched'],
      ['claim.NM1*TT.NM105', 'left out'],
      ['claim.CLP.CLP15', 'unlisted'],
      ['claim.CLP.CLP14', 'left out'],
      ['claim.PER.PER*CX', 'left out'],
      ['claim.PER.PER*XX', 'left out'], // a claim PER, any qualifier
      ['claim.CAS.CAS04', 'compared'],
      ['line.CAS.CAS04', 'compared'],
      ['claim.QTY.QTY*CA', 'compared'],
      ['claim.QTY.QTY*ZO', 'compared'],
      ['line.QTY.QTY*CA', 'left out'], // a service-line QTY, any qualifier
      ['claim.AMT.AMT*ZK', 'compared'],
      ['line.AMT.AMT*ZK', 'unlisted'], // the line rule reads ZL–ZO
      ['line.AMT.AMT*ZL', 'compared'],
      ['claim.AMT.AMT*I', 'compared, watched'],
      ['claim.AMT.AMT*NL', 'compared, watched'],
      ['claim.DTM.DTM*050', 'compared, watched'],
      ['claim.ZZZ.ZZZ01', 'compared'], // an unplaced segment in a claim
      ['line.ZZZ.ZZZ01', 'compared'], // and in a service line
      ['claim.MIA.MIA24', 'compared'],
      ['claim.MIA.MIA05', 'compared'],
      ['adjustment.PLB.PLB05', 'outside the claim'], // outside the claim, so not left out: BPR01, PLB, and the payer and payee N104 and REF*TJ
      ['payment.BPR.BPR01', 'outside the claim'],
      ['payment.BPR.BPR02', 'outside the claim'],
      ['payment.TRN.TRN02', 'outside the claim'],
      ['payee.N1.N104', 'outside the claim'],
      ['payer.N1.N104', 'outside the claim'],
      ['payer.N1.N102', 'outside the claim'],
      ['summary.TS3.TS301', 'outside the claim'],
    ]
    for (const [key, want] of cases) {
      const got = stateOf(key).state
      if (got !== want) fail(`stateOf(${key}) is “${got}”, not “${want}”`)
    }
    // the hover text: a list entry gives its reason, a name element says it is not read
    const reasonOf = (level, tag) => LEFT_OUT.find((e) => e.level === level && e.tag === tag)?.reason
    if (stateOf('claim.CLP.CLP03').why !== reasonOf('claim', 'CLP03')) fail('the hover text of claim.CLP.CLP03 is not its entry’s reason')
    if (stateOf('claim.NM1*QC.NM103').why !== reasonOf('claim', 'NM1*QC')) fail('the hover text of claim.NM1*QC.NM103 is not the reason of its entry, NM1*QC')
    if (!stateOf('claim.NM1*82.NM104').why?.startsWith('Not read:')) fail('a name element has no such hover text')
    if (stateOf('adjustment.PLB.PLB05').why !== 'Not inside a claim.') fail('PLB has no outside-the-claim hover text')
    if (stateOf('claim.CLP.CLP04').why !== undefined) fail('a compared element has hover text')
    if (!stateOf('claim.DTM.DTM*050').why) fail('a watched element has no hover text')
    if (!stateOf('payment.BPR.BPR02').why || !stateOf('claim.DTM.DTM*999').why) fail('an element outside the claim or unlisted has no hover text')
    for (const [states, want] of [
      [[], ''],
      [['compared', 'compared', 'compared, watched', 'left out'], '3 compared (1 watched) · 1 left out'],
      [['left out', 'outside the claim', 'outside the claim', 'unlisted'], '1 left out · 2 outside the claim · 1 unlisted'],
      [['compared, watched'], '1 compared (1 watched)'],
    ]) {
      if (tallyText(states) !== want) fail(`tallyText(${JSON.stringify(states)}) is “${tallyText(states)}”, not “${want}”`)
    }

    // ---- the page, rendered on the server
    const render = () => renderToStaticMarkup(createElement(StaticRouter, { location: '/compared-fields' }, createElement(ComparedFields)))
    const html = render()

    // the sections: two, in order, titled as the navigation lists them
    const drawnSections = [...html.matchAll(/<section id="([^"]+)"[^>]*><h2[^>]*><a href="#[^"]+"[^>]*>([^<]*)<\/a>/g)].map((m) => [m[1], m[2]])
    if (JSON.stringify(drawnSections) !== JSON.stringify([['tree', 'In file order'], ['left-out', 'Left out']])) fail(`the page's sections are ${JSON.stringify(drawnSections)}, not tree “In file order” then left-out “Left out”`)
    if (count(html, '<section') !== 2) fail(`the page holds ${count(html, '<section')} sections, not 2`)
    const navBlock = /path: '\/compared-fields',[\s\S]*?sections: \[([\s\S]*?)\]/.exec(await read('docs/App.tsx'))?.[1] ?? ''
    const navSections = [...navBlock.matchAll(/\{ id: '([^']+)', label: '([^']+)' \}/g)].map((m) => [m[1], m[2]])
    if (JSON.stringify(navSections) !== JSON.stringify(drawnSections)) fail(`the navigation lists ${JSON.stringify(navSections)} for this page, the page draws ${JSON.stringify(drawnSections)}`)
    const at = ['tree', 'left-out'].map((id) => html.indexOf(`id="${id}"`))
    if (at.some((i) => i < 0)) fail('a section is missing')
    const titleAt = html.indexOf('Compared fields</h1>')
    if (titleAt < 0 || titleAt > at[0]) fail('the title is missing or comes after the tree')

    // the lede, under the title, with the five-cent rule as compare.go writes it
    const compareSource = await readFile(new URL('../../backend/lib/core/service/compare.go', import.meta.url), 'utf8')
    const goCents = Number(/\bconst AmountMatchCents = (\d+)\b/.exec(compareSource)?.[1])
    if (!Number.isInteger(goCents)) fail('compare.go has no "const AmountMatchCents = N"')
    else if (goCents !== data.AMOUNT_MATCH_CENTS) fail(`compare.go counts two amounts within ${goCents} cents as the same, the page says ${data.AMOUNT_MATCH_CENTS}`)
    const header = /<header[^>]*>(.*?)<\/header>/s.exec(html)?.[1] ?? ''
    const LEDE =
      'The 835 elements the comparison looks at, as compare set 1 in the code: everything inside a claim except the elements left out below. ' +
      `Two amounts within $${(goCents / 100).toFixed(2)} of each other are the same amount; everything else compares exactly. ` +
      'The record of how the set was decided is docs/design/compared-fields.md.'
    if (textOf(/<p[^>]*>(.*?)<\/p>/s.exec(header)?.[1] ?? '') !== LEDE) fail(`the lede reads “${textOf(header)}”`)
    if (count(header, '<p') !== 1) fail('the header holds more than the one lede')

    // the tree section: the tree first, one note under it
    const treeSection = html.slice(at[0], at[1])
    if (count(treeSection, '<p') !== 1 || treeSection.indexOf('<p') < treeSection.lastIndexOf('data-leaf=')) fail('the tree section is not the tree followed by one note')
    if (!treeSection.includes('href="/835-compare"') || !treeSection.includes('href="/835-viewer"')) fail('the note under the tree does not link to the 835 viewer and the 835 compare page')
    if (/role="checkbox"/.test(treeSection)) fail('the tree still has a show/hide box')

    // every leaf once, with exactly one state tag
    const chunks = [...html.matchAll(/<li data-leaf="([^"]+)">(.*?)<\/li>/gs)]
    const drawn = chunks.map((m) => m[1])
    if (drawn.length !== leaves.length) fail(`the page draws ${drawn.length} leaves, the tree has ${leaves.length}`)
    if (new Set(drawn).size !== drawn.length) fail('the page draws a leaf twice')
    for (const leaf of leaves) if (!drawn.includes(leaf.id)) fail(`the page does not draw ${leaf.id}`)
    if (JSON.stringify(drawn) !== JSON.stringify(leaves.map((l) => l.id))) fail('the page draws the leaves in another order than the tree')
    const TAG = { blue: 'text-link', gray: 'bg-surface text-ink-2', yellow: 'bg-warn/20' }
    const COLOR = { compared: 'blue', 'compared, watched': 'blue', 'left out': 'gray', 'outside the claim': 'gray', unlisted: 'yellow' }
    const tally = Object.fromEntries(LEAF_STATES.map((s) => [s, []]))
    for (const [, id, inner] of chunks) {
      const tags = [...inner.matchAll(/data-state="([^"]+)"[^>]*><span class="([^"]*)">([^<]*)<\/span>/g)]
      if (tags.length !== 1 || count(inner, 'data-state=') !== 1) {
        fail(`${id} has ${count(inner, 'data-state=')} state attributes and ${tags.length} tags, not 1`)
        continue
      }
      const [, state, cls, text] = tags[0]
      const want = stateOf(id).state
      if (state !== want) fail(`${id} is drawn “${state}”, stateOf says “${want}”`)
      if (text !== state) fail(`${id}: the tag reads “${text}”, its state is “${state}”`)
      if (!LEAF_STATES.includes(state)) fail(`${id} has the state “${state}”, which is not one of the five`)
      else {
        if (!cls.includes(TAG[COLOR[state]])) fail(`${id}: the “${state}” tag is not ${COLOR[state]}`)
        tally[state].push(id)
      }
      if (/text-danger|text-ok|bg-ok|bg-danger/.test(inner)) fail(`${id} is drawn in red or green`)
    }
    // no element is unlisted: each inside a claim is covered by the compared rule or the left-out rule
    if (tally.unlisted.length > 0) fail(`${tally.unlisted.length} elements are unlisted, covered by neither the compared rule nor the left-out rule: ${tally.unlisted.join(', ')}`)
    // left out is only ever inside a claim: the elements outside it are "outside the claim"
    const leftOutside = tally['left out'].filter((id) => !/^(claim|line)\./.test(id))
    if (leftOutside.length > 0) fail(`${leftOutside.length} elements outside a claim are tagged left out: ${leftOutside.join(', ')}`)
    // every part and group shows the tally of its leaves
    for (const node of nodes.filter((n) => n.children.length > 0)) {
      const m = new RegExp(`data-group="${reEscape(node.id)}">.*?data-tally="([^"]*)"`, 's').exec(html)
      const want = tallyText(leavesOf(node).map((l) => stateOf(l.id).state))
      if (!m) fail(`the page draws no part or group ${node.id}`)
      else if (m[1] !== want) fail(`${node.id} shows the tally “${m[1]}”, its leaves give “${want}”`)
      if (m && !html.includes(`>${want}</span>`)) fail(`${node.id}: the tally is not shown as text`)
    }
    const groupCount = nodes.filter((n) => n.children.length > 0).length
    const partCount = tree.length
    if (count(html, 'data-group="') !== groupCount) fail('the page draws a part or group the tree does not hold')
    // every group opens: nothing is collapsed at the start; every part and
    // group has its triangle, and a part's title opens and closes it too
    if (count(html, 'aria-expanded="false"') !== 0) fail('a part or group starts closed')
    if (count(html, 'aria-expanded="true"') !== groupCount + partCount) fail('a part or group is missing its triangle, or a part its title button')
    if (!html.includes('>Expand all</button>') || !html.includes('>Collapse all</button>')) fail('the tree has no Expand all or Collapse all')
    // every group's and leaf's code is a link into the lookup, searching the
    // last part of its key, with the public name beneath it where it has one
    // (the markup of a node's own row: a leaf's whole item, a group's item up
    // to the list of its children)
    const rowOf = (node) => {
      const group = node.children.length > 0
      const start = html.indexOf(group ? `data-group="${node.id}">` : `data-leaf="${node.id}">`)
      if (start < 0) return ''
      const end = html.indexOf(group ? '<ul' : '</li>', start)
      return html.slice(start, end < 0 ? undefined : end)
    }
    for (const node of nodes.filter((n) => !tree.includes(n))) {
      const own = rowOf(node)
      const search = node.id.slice(node.id.lastIndexOf('.') + 1)
      if (!own.includes(`data-code="${escapeHtml(search)}"`)) fail(`${node.id}: its code is not a link into the lookup searching ${search}`)
      if (node.names !== node.codes && !own.includes(`>${escapeHtml(node.names)}</span>`)) fail(`${node.id}: its public name is not shown beneath its code`)
    }
    // the controls: the tree's triangles and part titles, its code links, Expand
    // all and Collapse all, the leave-out list's tag links, and the Lookup button
    const links = count(html, 'data-code="')
    const wantLinks = nodes.length - partCount + LEFT_OUT.length
    if (links !== wantLinks) fail(`the page has ${links} code links, the tree and the list have ${wantLinks} codes`)
    if (count(html, '<button') !== groupCount + partCount + links + 3) fail(`the page has ${count(html, '<button')} buttons, not ${groupCount + partCount + links + 3}`)
    if (!html.includes('aria-pressed="false"') || !/>Lookup<\/button>/.test(html)) fail('the page has no Lookup button')
    if (html.includes('role="complementary"')) fail('the lookup panel is open with no lookup=1')
    if (/<select|<input|<textarea/.test(html)) fail('the page has a form control')

    // the Left out section: the list, one row per entry of compare set 1
    const leftSection = html.slice(at[1])
    const rows = [...leftSection.matchAll(/<tr data-left-out="([^"]+)">(.*?)<\/tr>/gs)]
    const drawnList = rows.map((m) => m[1])
    if (JSON.stringify(drawnList) !== JSON.stringify(pageList)) fail(`the page draws the leave-out list ${drawnList.join(', ')}, its data holds ${pageList.join(', ')}`)
    for (const e of codeList) if (!drawnList.includes(e)) fail(`the page does not draw ${e}, which compareset.go leaves out`)
    for (const e of drawnList) if (!codeList.includes(e)) fail(`the page draws ${e}, which compareset.go does not leave out`)
    if (count(leftSection, 'data-left-out=') !== LEFT_OUT.length) fail('the page draws a leave-out row twice or not at all')
    const headingBefore = (index) => [...leftSection.slice(0, index).matchAll(/<th [^>]*scope="colgroup"[^>]*>([^<]*)<\/th>/g)].pop()?.[1]
    for (const m of rows) {
      const entry = LEFT_OUT.find((e) => `${e.level}:${e.tag}` === m[1])
      const cells = [...m[2].matchAll(/<td class="([^"]*)">(.*?)<\/td>/gs)]
      if (!entry || cells.length !== 3) {
        fail(`${m[1]} is not drawn as three cells`)
        continue
      }
      const [tag, name, reason] = cells.map((c) => textOf(c[2]))
      if (!cells[0][1].includes('font-mono')) fail(`${m[1]}: the tag is not in mono`)
      if (!cells[0][2].includes(`data-code="${escapeHtml(entry.tag)}"`)) fail(`${m[1]}: the tag is not a link into the lookup`)
      if (tag !== entry.tag || name !== entry.publicName || reason !== entry.reason) fail(`${m[1]} is drawn “${tag} | ${name} | ${reason}”, its data holds “${entry.tag} | ${entry.publicName} | ${entry.reason}”`)
      const want = entry.level === 'claim' ? 'Claim level' : 'Service line level'
      if (headingBefore(m.index) !== want) fail(`${m[1]} sits under the heading “${headingBefore(m.index)}”, not “${want}”`)
    }
    // the line about the watched elements, below the list
    const watchedLine = /<p><strong>Compared and watched:<\/strong>(.*?)<\/p>/s.exec(leftSection)?.[1]
    if (!watchedLine) fail('the Left out section has no “Compared and watched:” line')
    else {
      const named = [...watchedLine.matchAll(/<code>([^<]*)<\/code>/g)].map((m) => m[1])
      if (JSON.stringify(named) !== JSON.stringify(['DTM*050', 'AMT*I', 'AMT*NL', 'NM1*TT'])) fail(`the watched line names ${named.join(', ')}`)
      for (const phrase of ['compared and recorded like every other element', 'settled from the real files', 'named on every run']) {
        if (!textOf(watchedLine).includes(phrase)) fail(`the watched line does not say “${phrase}”`)
      }
      if (leftSection.indexOf('Compared and watched:') < leftSection.lastIndexOf('data-left-out=')) fail('the watched line is not below the list')
    }
    // the one sentence under the watched line: the comparison reads an NM1's identifier, not its name elements
    const leftText = textOf(leftSection)
    const NAME_NOTE =
      'The comparison reads the identifier of an NM1 (NM109) under its role, not its name elements (NM103 to NM107), so those are left out without an entry in the list above.'
    if (!leftText.includes(NAME_NOTE)) fail(`the Left out section does not say “${NAME_NOTE}”`)
    else if (leftText.indexOf(NAME_NOTE) < leftText.indexOf('Compared and watched:')) fail('the sentence about the name elements is above the watched line')
    if (count(leftSection, '<p>') !== 2) fail(`the Left out section holds ${count(leftSection, '<p>')} notes under the list, not 2: the watched line and the one sentence about the name elements`)

    // ---- a gap shows: take the compared rule for one element out of the data, in
    // memory, and the element turns yellow, "unlisted", and is counted
    {
      const gapAt = COMPARED.findIndex((r) => r.label === 'DTM*036')
      if (gapAt < 0) fail('the compared rule for DTM*036 is not in the data, so the gap test cannot run')
      else {
        const [removed] = COMPARED.splice(gapAt, 1)
        try {
          const gapHtml = render()
          const chunk = /<li data-leaf="claim\.DTM\.DTM\*036">(.*?)<\/li>/s.exec(gapHtml)?.[1] ?? ''
          if (!chunk.includes('data-state="unlisted"') || !chunk.includes('bg-warn/20')) fail('an element no rule covers is not drawn "unlisted" in yellow')
          if (!gapHtml.includes('data-tally="2 compared (1 watched) · 1 left out · 1 unlisted"')) fail('the tally of its group does not count the gap')
        } finally {
          COMPARED.splice(gapAt, 0, removed)
        }
        if (render() !== html) fail('the page is not the same again after the gap test')
      }
    }

    // ---- the lookup panel, opened from the address: beside the page, its search
    // holding the code, that element's row shown
    {
      const open = renderToStaticMarkup(createElement(StaticRouter, { location: '/compared-fields?lookup=1&q=CLP03' }, createElement(ComparedFields)))
      const at = open.indexOf('role="complementary"')
      if (at < 0 || !open.includes('aria-label="835 field lookup"')) fail('lookup=1 does not open the 835 lookup panel')
      else {
        const panel = open.slice(at)
        if (!panel.includes('value="CLP03"')) fail('the panel opened with q=CLP03 does not hold CLP03 in its search')
        if (!panel.includes('data-row="claim:CLP03"')) fail('the panel searched for CLP03 does not show that element')
      }
      if (!open.includes('aria-pressed="true"')) fail('the Lookup button does not show the panel open')
      if (!open.includes('class="@container min-w-0 flex-1"') || !/class="sticky top-0 h-svh w-\[26rem\][^"]*"/.test(open)) fail('the panel is not a column beside the page')
    }

    // ---- the tree's code, the data, and the page keep clear of storage and the
    // network; the address holds only the lookup's lookup= and q=, written by
    // the lookup's own code
    for (const path of ['docs/compare835/StateTree.tsx', 'docs/compare835/unionTree.ts']) {
      const found = (await read(path)).match(/localStorage|sessionStorage|indexedDB|useSearchParams|setSearchParams|location\b|history\b|document\.cookie|fetch\(|XMLHttpRequest|sendBeacon/g)
      if (found) fail(`${path} mentions ${[...new Set(found)].join(', ')}`)
    }
    // the data file holds prose that may say "history" and "location", so only the browser's own names count there
    const dataSource = await read('data/comparedFields.ts')
    const dataFound = dataSource.match(/localStorage|sessionStorage|indexedDB|useSearchParams|setSearchParams|document\.cookie|fetch\(|XMLHttpRequest|sendBeacon/g)
    if (dataFound) fail(`data/comparedFields.ts mentions ${[...new Set(dataFound)].join(', ')}`)
    const pageSource = await read('docs/pages/ComparedFields.tsx')
    const pageFound = pageSource.match(/localStorage|sessionStorage|indexedDB|useSearchParams|setSearchParams|history\b|document\.cookie|fetch\(|XMLHttpRequest|sendBeacon/g)
    if (pageFound) fail(`ComparedFields.tsx mentions ${[...new Set(pageFound)].join(', ')}`)
    if (!pageSource.includes('useLookup()')) fail('ComparedFields.tsx does not keep its lookup panel with useLookup')
    // the lookup's own code writes lookup= and q= and no other key, and no storage
    const lookupSource = await read('docs/viewer/useLookup.ts')
    const keys = [...lookupSource.matchAll(/next\.(?:set|delete)\('([^']*)'/g)].map((m) => m[1])
    if (keys.length === 0 || keys.some((k) => k !== 'lookup' && k !== 'q')) fail(`useLookup.ts writes the address keys ${[...new Set(keys)].join(', ')}, not only lookup and q`)
    if (/localStorage|sessionStorage|indexedDB|document\.cookie|fetch\(|XMLHttpRequest|sendBeacon/.test(lookupSource)) fail('useLookup.ts touches storage or the network')

    if (failures.length > 0) {
      console.error(`compared fields SSR check FAILED (${failures.length} problems):`)
      for (const f of failures.slice(0, 40)) console.error(`  - ${f}`)
      if (failures.length > 40) console.error(`  … and ${failures.length - 40} more`)
      process.exitCode = 1
      return
    }

    // What the tree shows, for the record.
    const groups = nodes.filter((n) => n.children.length > 0 && !tree.includes(n))
    console.log(`compared fields SSR check passed: the tree holds ${tree.length} parts, ${groups.length} groups, and ${leaves.length} elements, the union of ${examples.length} examples (${sectionCount} sections, ${rowCount} rows); the page draws each element once with one state tag, and every part and group with its tally.`)
    console.log(`leave-out list: ${pageList.length} entries on the page, ${codeList.length} in compareset.go (set 1), none on one side only`)
    console.log('states:')
    for (const state of LEAF_STATES) console.log(`  ${state}: ${tally[state].length}`)
    console.log(`  total: ${leaves.length}`)
    console.log('by part:')
    for (const part of tree) {
      console.log(`  ${part.id}: ${tallyText(leavesOf(part).map((l) => stateOf(l.id).state))}`)
    }
    console.log(tally.unlisted.length === 0 ? 'unlisted: none' : `unlisted (${tally.unlisted.length}): ${tally.unlisted.join(', ')}`)
    const names = tally['left out'].filter((id) => stateOf(id).why?.startsWith('Not read:'))
    const listed = tally['left out'].filter((id) => !names.includes(id))
    console.log(`left out (${tally['left out'].length}), all inside a claim: ${listed.length} under the list's entries, ${names.length} name elements with no entry`)
    console.log(`  under the list's entries (${listed.length}): ${listed.join(', ')}`)
    console.log(`  name elements, no entry (${names.length}): ${names.join(', ')}`)
    console.log(`compared, watched (${tally['compared, watched'].length}): ${tally['compared, watched'].join(', ')}`)
  } finally {
    await server.close()
  }
}

main().catch((err) => {
  console.error(err)
  process.exitCode = 1
})
