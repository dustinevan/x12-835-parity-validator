// Loads the 835 compare page's own code (frontend/src/docs/compare835/) through
// Vite's SSR module graph, the way the docs site resolves a `?raw` import, and
// runs it over the 28 worked examples under frontend/src/data/viewer835/files/
// with no browser involved. It checks that
//   - each example's parse holds every segment of its file (nothing dropped),
//     and every claim in it has a match key;
//   - the list built for every ordered pair of examples (a reference and a
//     compared A, with a third example as B) throws nothing, and every
//     segment the file carries, outside the envelope, produced a row;
//   - the tree built from those rows holds every row, and its checkboxes hide
//     every row when a part is unchecked;
//   - an example set against itself has no row that differs (unless it holds
//     a repeated match key, which is never matched);
//   - the marks follow the comparison's five-cent rule (rows.ts sameValue,
//     the cents from data/comparedFields.ts AMOUNT_MATCH_CENTS): two dollar
//     amounts within five cents of each other, edges included, are the same
//     amount, so neither cell is marked and the row is not one of the Only
//     differences rows; one more cent is marked, with its signed difference
//     as before; a number that is not money (the DRG weight CLP12, the
//     discharge fraction CLP13) and an identifier still compare exactly. The
//     28 examples hold no two dollar amounts one to five cents apart, so the
//     marked-cell count over every list is the same as under exact
//     comparison (the run prints how many cells the rule left unmarked:
//     none); the rule is checked on copies of two examples with their
//     amounts moved by a few cents;
//   - the default three examples share one claim, and the address's
//     defaults resolve to distinct examples;
//   - every example, rewritten with other delimiters (element `|`, component
//     `>`, other segment terminators; one long line, CRLF, a newline as the
//     terminator, leading white space), parses to the same segments, rows,
//     and marks as the original, and the ISA keeps its own ISA16;
//   - text pasted into a box reads the way an example does, and text that
//     does not read (blank, no ISA, cut short, a bad segment ID, two
//     interchanges) comes back as an error or an empty column and never
//     throws, with no file kept for it;
//   - the presets (compare835/presets.ts): for the tree of every list above
//     and of every example alone, each built-in preset names only keys that
//     tree holds, in the shortest form, and shows exactly the rows its rule
//     says (checked row by row against the rules written out again here);
//     formatPreset then parsePreset gives back every built-in preset and
//     random hidden sets; malformed text comes back as an error and never
//     throws; the saved list reads whatever storage holds without throwing
//     and works with storage refused; the preset code never touches storage
//     or the address with the text of the box; and the page, rendered on the
//     server, has the presets select, Save as… and Delete, and the preset
//     text box, showing the preset the address (or storage) matches;
//   - every code on the page (a row's key and note, a tree group and element)
//     is a link into the 835 lookup and a CAS row's pair of coded values is
//     not, and lookup=1&q=NM1*QC opens the lookup panel in its own column
//     beside the page, its search holding NM1*QC and showing NM109 under QC.
// Run from frontend/:
//   node scripts/ssr-check-835-compare.mjs
import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { StaticRouter } from 'react-router'
import { readFile } from 'node:fs/promises'
import { createServer } from 'vite'

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
    const { splitSegments } = await server.ssrLoadModule('/src/docs/viewer835/parse.ts')
    const { parseFile, segmentsOf } = await server.ssrLoadModule('/src/docs/compare835/file.ts')
    const { COLS, listOf, rowDiffers, isMarked, sameValue } = await server.ssrLoadModule('/src/docs/compare835/rows.ts')
    const { AMOUNT_MATCH_CENTS } = await server.ssrLoadModule('/src/data/comparedFields.ts')
    const { subtract } = await server.ssrLoadModule('/src/format.ts')
    const { ListGrid } = await server.ssrLoadModule('/src/docs/compare835/ListGrid.tsx')
    const { treeOf, checkOf, toggled, rowPath } = await server.ssrLoadModule('/src/docs/compare835/tree.ts')
    const { DEFAULT_COLUMNS, keysOf, resolveColumns, readingOfText, exampleTextOf, listOfReadings } = await server.ssrLoadModule('/src/docs/compare835/columns.ts')
    const {
      BUILT_IN_NAMES, MAX_NAME, PRESETS_KEY, builtInPresets, formatPreset, inTreeOrder, isKey, loadSavedPresets, matchingOption, nameProblem, normalizeName, parsePreset,
      readSavedPresets, storeSavedPresets, viewOf, withPreset, withoutPreset,
    } = await server.ssrLoadModule('/src/docs/compare835/presets.ts')

    const examples = VIEWER_835_EXAMPLES
    if (!Array.isArray(examples) || examples.length === 0) throw new Error('VIEWER_835_EXAMPLES is empty, nothing to check')

    const failures = []
    const fail = (message) => failures.push(message)

    // ---- the built-in presets, stated again here from their rules and checked row by row
    // (presets.ts works them out from the tree's keys; this looks at the rows
    // each one shows and hides.)
    const segmentOf = (row) => row.group.split('*')[0]
    const atTheLevels = (row) => row.scope === 'claim' || row.scope === 'line'
    const SHOWN = {
      Everything: () => true,
      'Claim level': (row) => row.scope !== 'line',
      'Service lines': (row) => !(row.scope === 'claim' && segmentOf(row) !== 'CLP'),
      'Amounts and adjustments': (row) => atTheLevels(row) && ['CLP', 'CAS', 'SVC', 'AMT'].includes(segmentOf(row)),
      'Codes and dates': (row) => atTheLevels(row) && ['CLP', 'MOA', 'MIA', 'LQ', 'DTM'].includes(segmentOf(row)),
    }
    const idsInTree = (tree) => {
      const ids = new Set()
      const add = (node) => {
        ids.add(node.id)
        node.children.forEach(add)
      }
      tree.forEach(add)
      return ids
    }
    let presetTrees = 0
    const checkPresets = (label, sections, tree) => {
      presetTrees++
      const known = idsInTree(tree)
      const rows = sections.flatMap((s) => s.rows)
      const made = builtInPresets(tree)
      if (JSON.stringify(made.map((p) => p.name)) !== JSON.stringify(Object.keys(SHOWN))) fail(`${label}: the built-in presets are ${made.map((p) => p.name).join(', ')}`)
      for (const id of known) if (!isKey(id)) fail(`${label}: the tree's id ${id} is not a key the preset text accepts`)
      for (const preset of made) {
        const shows = SHOWN[preset.name]
        if (!shows) continue
        const hidden = new Set(preset.hide)
        if (hidden.size !== preset.hide.length) fail(`${label}: preset “${preset.name}” lists an id twice`)
        for (const id of preset.hide) {
          if (!known.has(id)) fail(`${label}: preset “${preset.name}” names ${id}, which the tree does not hold`)
          if (preset.hide.some((other) => other !== id && id.startsWith(`${other}.`))) fail(`${label}: preset “${preset.name}” lists ${id} under another listed id`)
        }
        for (const row of rows) {
          const visible = !rowPath(row).some((id) => hidden.has(id))
          if (visible !== shows(row)) {
            fail(`${label}: preset “${preset.name}” ${visible ? 'shows' : 'hides'} ${rowPath(row)[2]}, which it should ${shows(row) ? 'show' : 'hide'}`)
            break
          }
        }
        const back = parsePreset(formatPreset(preset.name, preset.hide))
        if (JSON.stringify(back) !== JSON.stringify({ name: preset.name, hide: preset.hide })) fail(`${label}: preset “${preset.name}” does not read back from its text`)
      }
    }

    // ---- each example's parse
    const parsed = new Map()
    for (const ex of examples) {
      try {
        const file = parseFile(ex.file)
        parsed.set(ex.id, file)
        const total = splitSegments(ex.file).length
        const held = segmentsOf(file).length
        if (held !== total) fail(`${ex.id}: the parse holds ${held} segments but the file has ${total}`)
        for (const claim of file.claims) {
          if (!claim.key) fail(`${ex.id}: claim ${claim.head.elements[0]} has no match key: ${claim.keyError}`)
        }
      } catch (err) {
        fail(`${ex.id}: ${err instanceof Error ? err.stack : String(err)}`)
      }
    }

    // ---- the list for every ordered pair of examples
    let lists = 0
    let rows = 0
    let marked = 0
    let withinRule = 0 // dollar amounts that differ as text, are within five cents, and so are not marked
    const n = examples.length
    for (let i = 0; i < n; i++) {
      for (let j = 0; j < n; j++) {
        const ids = { ref: examples[i].id, a: examples[j].id, b: examples[(i + j + 1) % n].id }
        const label = `${ids.ref} | ${ids.a} | ${ids.b}`
        try {
          const files = { ref: parsed.get(ids.ref), a: parsed.get(ids.a), b: parsed.get(ids.b) }
          const sections = listOf(files)
          lists++

          // every segment outside the envelope produced a row
          const covered = { ref: new Set(), a: new Set(), b: new Set() }
          for (const section of sections) {
            if (section.rows.length === 0) fail(`${label}: section “${section.title}” has no rows`)
            for (const row of section.rows) {
              rows++
              if (!row.key.codes || !row.key.names) fail(`${label}: a row of “${section.title}” has no key`)
              for (const col of COLS) {
                for (const seg of row.from[col] ?? []) covered[col].add(seg)
                if (row.values[col] !== undefined && !row.from[col]) fail(`${label}: ${row.leaf} has a ${col} value and no segment`)
                if (isMarked(section, row, col)) marked++
                else if (col !== 'ref' && row.money && section.held.ref && section.held[col] && !sameValue({ amount: true, money: false }, row.values[col], row.values.ref)) withinRule++
              }
            }
          }
          for (const col of COLS) {
            const envelope = new Set(files[col].envelope)
            for (const seg of segmentsOf(files[col])) {
              if (!envelope.has(seg) && !covered[col].has(seg)) fail(`${label}: ${col} segment ${seg.raw} produced no row`)
            }
          }

          // the tree holds every row, and unchecking a part hides all its rows
          const tree = treeOf(sections)
          const known = new Set()
          const collect = (node) => {
            known.add(node.id)
            node.children.forEach(collect)
          }
          tree.forEach(collect)
          for (const section of sections) {
            for (const row of section.rows) {
              for (const id of rowPath(row)) if (!known.has(id)) fail(`${label}: row path ${id} is not in the tree`)
            }
          }
          for (const root of tree) {
            if (checkOf(root, new Set()) !== 'on') fail(`${label}: ${root.id} is not checked with nothing hidden`)
            const hidden = toggled(tree, new Set(), root)
            if (checkOf(root, hidden) !== 'off') fail(`${label}: unchecking ${root.id} does not hide it`)
            for (const other of tree) if (other !== root && checkOf(other, hidden) !== 'on') fail(`${label}: unchecking ${root.id} changed ${other.id}`)
            const back = toggled(tree, hidden, root)
            if (back.size !== 0) fail(`${label}: checking ${root.id} again leaves ${[...back].join(',')} hidden`)
          }
          checkPresets(label, sections, tree)
        } catch (err) {
          fail(`${label}: ${err instanceof Error ? err.stack : String(err)}`)
        }
      }
    }

    // ---- an example against itself
    for (const ex of examples) {
      const file = parsed.get(ex.id)
      const keys = file.claims.map((c) => c.key?.id)
      const repeated = new Set(keys).size !== keys.length
      const sections = listOf({ ref: file, a: file, b: file })
      checkPresets(`${ex.id} alone`, sections, treeOf(sections))
      const differing = sections.flatMap((s) => s.rows.filter((r) => rowDiffers(r) || COLS.some((c) => isMarked(s, r, c))))
      if (!repeated && differing.length > 0) fail(`${ex.id} against itself: ${differing.length} rows differ`)
      if (repeated && differing.length === 0) fail(`${ex.id} holds a repeated key but its claims were matched`)
    }

    // ---- exact decimals, the five-cent rule for money, and identifiers kept as text
    const number = { amount: true, money: false } // a DRG weight, a discharge fraction
    const dollars = { amount: true, money: true }
    if (!sameValue(number, '10.0', '10.00')) fail('10.0 and 10.00 are not equal as numbers')
    if (sameValue(number, '10.0', '10.01')) fail('10.0 and 10.01 are equal as numbers that are not money')
    if (sameValue({ amount: false, money: false }, '0123', '123')) fail('0123 and 123 are equal as identifiers')
    if (sameValue(dollars, undefined, '0.00')) fail('an absent amount equals 0.00')
    // The edges, from the one constant: the limit is the same amount, one cent more is not.
    const limit = (AMOUNT_MATCH_CENTS / 100).toFixed(2)
    const past = ((AMOUNT_MATCH_CENTS + 1) / 100).toFixed(2)
    const fiveCentCases = [
      ['80.00', subtract('80.00', limit), true, 'the limit under'],
      ['80.00', subtract('80.00', `-${limit}`), true, 'the limit over'],
      ['80.00', subtract('80.00', past), false, 'one cent past the limit under'],
      ['80.00', subtract('80.00', `-${past}`), false, 'one cent past the limit over'],
      ['80', '80.03', true, 'three cents, one side written without decimals'],
      ['80.0', '80.00', true, 'the same amount spelled two ways'],
      ['-1.00', subtract('-1.00', limit), true, 'a negative amount, the limit lower'],
      ['0.02', '-0.02', true, 'four cents across zero'],
      ['0.03', '-0.03', false, 'six cents across zero'],
    ]
    for (const [x, y, same, what] of fiveCentCases) {
      if (sameValue(dollars, x, y) !== same) fail(`five-cent rule, ${what}: ${x} and ${y} are ${same ? 'not ' : ''}the same dollar amount`)
      if (sameValue(dollars, y, x) !== same) fail(`five-cent rule, ${what}, reversed: ${y} and ${x} are ${same ? 'not ' : ''}the same dollar amount`)
    }
    if (sameValue(number, '1.2345', '1.2346')) fail('a DRG weight one ten-thousandth apart is the same number')
    if (sameValue(number, '1.00', '1.03')) fail('a discharge fraction three hundredths apart is the same number')

    // Copies of two examples with amounts moved by a few cents, against the
    // original as the reference: within the limit no cell is marked and the
    // row does not differ; past it the cell is marked with its signed
    // difference, drawn as before. The DRG weight and discharge fraction are
    // marked at any difference.
    let movedCells = 0
    {
      // The first segment with the given ID (and, for CAS, the given group
      // and reason), with element `position` set to `value`.
      const setElement = (text, head, position, value) => {
        const at = text.indexOf(`\n${head}*`)
        if (at < 0) throw new Error(`no ${head} segment`)
        const end = text.indexOf('~', at)
        const parts = text.slice(at + 1, end).split('*')
        parts[position] = value
        return `${text.slice(0, at + 1)}${parts.join('*')}${text.slice(end)}`
      }
      const by = (amount, cents) => subtract(amount, (-cents / 100).toFixed(2))
      const original = examples.find((e) => e.id === 'paid-by-ach')
      const inpatient = examples.find((e) => e.id === 'inpatient')
      if (!original || !inpatient) fail('the examples paid-by-ach and inpatient, which the five-cent check edits, are missing')
      else {
        // paid-by-ach pays CLP04 80.00 and SVC03 80.00 after CAS*CO*45 of 20.00.
        const a = setElement(setElement(setElement(original.file, 'CLP', 4, by('80.00', -3)), 'SVC', 3, by('80.00', AMOUNT_MATCH_CENTS)), 'CAS*CO*45', 3, by('20.00', AMOUNT_MATCH_CENTS + 1))
        const b = setElement(setElement(setElement(original.file, 'CLP', 4, by('80.00', -(AMOUNT_MATCH_CENTS + 1))), 'SVC', 3, by('80.00', -AMOUNT_MATCH_CENTS)), 'CAS*CO*45', 3, by('20.00', -AMOUNT_MATCH_CENTS))
        const sections = listOf({ ref: parseFile(original.file), a: parseFile(a), b: parseFile(b) })
        const find = (leaf, codes) => {
          for (const s of sections) for (const r of s.rows) if (r.leaf === leaf && (codes === undefined || r.key.codes === codes)) return { s, r }
          return undefined
        }
        const want = [
          ['CLP04', undefined, { a: false, b: true }],
          ['SVC03', undefined, { a: false, b: false }],
          ['CAS03', 'CO-45', { a: true, b: false }],
        ]
        for (const [leaf, codes, marks] of want) {
          const found = find(leaf, codes)
          if (!found) {
            fail(`five-cent rule on paid-by-ach: no ${codes ?? leaf} row`)
            continue
          }
          if (!found.r.money) fail(`five-cent rule on paid-by-ach: ${codes ?? leaf} is not a dollar amount`)
          for (const col of ['a', 'b']) {
            movedCells++
            const is = isMarked(found.s, found.r, col)
            if (is !== marks[col]) fail(`five-cent rule on paid-by-ach: ${codes ?? leaf} in ${col} (${found.r.values[col]} against ${found.r.values.ref}) is ${is ? '' : 'not '}marked`)
          }
          if (rowDiffers(found.r) !== (marks.a || marks.b)) fail(`five-cent rule on paid-by-ach: Only differences ${rowDiffers(found.r) ? 'keeps' : 'drops'} the ${codes ?? leaf} row`)
        }
        // The signed differences beside the marked cells, and none beside the others.
        const grid = renderToStaticMarkup(createElement(ListGrid, { sections, head: { ref: 'Reference', a: 'A', b: 'B' } }))
        const deltas = [...grid.matchAll(/<span class="bg-warn px-1\.5[^"]*">([^<]*)<\/span>/g)].map((m) => m[1])
        const pastText = `$0.${String(AMOUNT_MATCH_CENTS + 1).padStart(2, '0')}`
        const wantDeltas = [`−${pastText}`, `+${pastText}`] // b's CLP04, then a's CO-45, in row order
        if (JSON.stringify(deltas) !== JSON.stringify(wantDeltas)) fail(`five-cent rule on paid-by-ach: the grid shows the differences ${JSON.stringify(deltas)}, not ${JSON.stringify(wantDeltas)}`)

        // inpatient carries a DRG weight (CLP12 1.2345) and a discharge fraction (CLP13 1.0000).
        const weights = setElement(setElement(inpatient.file, 'CLP', 12, '1.2346'), 'CLP', 13, '1.0300')
        const drg = listOf({ ref: parseFile(inpatient.file), a: parseFile(weights), b: parseFile(inpatient.file) })
        for (const leaf of ['CLP12', 'CLP13']) {
          const s = drg.find((x) => x.rows.some((r) => r.leaf === leaf))
          const r = s?.rows.find((x) => x.leaf === leaf)
          if (!r) {
            fail(`five-cent rule on inpatient: no ${leaf} row`)
            continue
          }
          movedCells++
          if (r.money) fail(`five-cent rule on inpatient: ${leaf} is counted as money`)
          if (!isMarked(s, r, 'a')) fail(`five-cent rule on inpatient: ${leaf} moved from ${r.values.ref} to ${r.values.a} is not marked`)
        }
      }
    }

    // ---- the same example under other delimiters
    // The examples use ~ between segments, * between elements, : in a
    // composite. Rewriting the examples' own segments with other separators,
    // and joining them in other ways, must change nothing the page shows.
    const rewrite = (text, { element, component, terminator, between = '', before = '', after = '' }) =>
      before +
      text
        .split('~')
        .map((s) => s.trim())
        .filter(Boolean)
        .map((s) => s.split('*').map((e) => e.replaceAll(':', component)).join(element) + terminator)
        .join(between) +
      after
    const variants = {
      'one long line, | > ~': { element: '|', component: '>', terminator: '~' },
      'CRLF between segments, | > ~': { element: '|', component: '>', terminator: '~', between: '\r\n' },
      'newline as the terminator, | >': { element: '|', component: '>', terminator: '\n' },
      'CRLF as the terminator, * :': { element: '*', component: ':', terminator: '\r\n' },
      'one long line, | > !': { element: '|', component: '>', terminator: '!' },
      'one long line, * :, leading white space': { element: '*', component: ':', terminator: '~', before: '﻿ \n  ', after: '\n\n' },
    }
    const judged = (sections) => sections.flatMap((s) => s.rows.filter((r) => rowDiffers(r) || COLS.some((c) => isMarked(s, r, c))))
    let rewritten = 0
    for (const ex of examples) {
      const original = parsed.get(ex.id)
      const base = listOf({ ref: original, a: original, b: original })
      const baseRows = base.reduce((n, s) => n + s.rows.length, 0)
      const baseDiffering = judged(base).length
      for (const [name, v] of Object.entries(variants)) {
        const label = `${ex.id} (${name})`
        try {
          const text = rewrite(ex.file, v)
          if (v.component !== ':' && !text.includes(`${v.element}${v.component}`)) fail(`${label}: the rewrite did not use the component separator`)
          const file = parseFile(text)
          rewritten++
          if (segmentsOf(file).length !== segmentsOf(original).length) fail(`${label}: ${segmentsOf(file).length} segments, the original has ${segmentsOf(original).length}`)
          if (file.claims.length !== original.claims.length) fail(`${label}: ${file.claims.length} claims, the original has ${original.claims.length}`)
          if (file.envelope[0].elements[15] !== v.component) fail(`${label}: ISA16 reads “${file.envelope[0].elements[15]}”, not “${v.component}”`)
          const alt = listOf({ ref: original, a: file, b: file })
          const altRows = alt.reduce((n, s) => n + s.rows.length, 0)
          if (altRows !== baseRows) fail(`${label}: ${altRows} rows, the original has ${baseRows}`)
          if (judged(alt).length !== baseDiffering) fail(`${label}: ${judged(alt).length} rows differ from the original's, ${baseDiffering} expected`)
          const same = listOf({ ref: file, a: original, b: file })
          if (judged(same).length !== baseDiffering) fail(`${label}: as the reference, ${judged(same).length} rows differ, ${baseDiffering} expected`)
        } catch (err) {
          fail(`${label}: ${err instanceof Error ? err.stack : String(err)}`)
        }
      }
    }

    // The comparison above must be able to see a difference: the payment amount (BPR02) changed in a rewrite is judged.
    {
      const ex = examples.find((e) => e.id === DEFAULT_COLUMNS.ref) ?? examples[0]
      const original = parsed.get(ex.id)
      const changed = parseFile(rewrite(ex.file, variants['one long line, | > ~']).replace(/(~BPR\|[^|]*\|)[^|]*\|/, (_, head) => `${head}9999.99|`))
      if (judged(listOf({ ref: original, a: changed, b: changed })).length === 0) fail(`${ex.id}: a changed amount in a rewritten file is not judged as a difference`)
    }

    // ---- text pasted into a box
    let pastedChecks = 0
    const isReading = (r) => r !== null && typeof r === 'object' && !(r.file && r.error) && (r.file === undefined || typeof r.file === 'object') && (r.error === undefined || typeof r.error === 'string')
    const expectError = (label, text) => {
      pastedChecks++
      let r
      try {
        r = readingOfText(text)
      } catch (err) {
        fail(`${label}: readingOfText threw ${err instanceof Error ? err.message : String(err)}`)
        return
      }
      if (!isReading(r)) fail(`${label}: not a reading: ${JSON.stringify(r)}`)
      else if (r.file) fail(`${label}: text that does not read has a file`)
      else if (!r.error || !r.error.startsWith('x12:')) fail(`${label}: no x12 error text (${JSON.stringify(r.error)})`)
    }
    const expectEmpty = (label, text) => {
      pastedChecks++
      const r = readingOfText(text)
      if (r.file || r.error) fail(`${label}: blank text is not an empty column: ${JSON.stringify(r.error)}`)
    }
    expectEmpty('empty text', '')
    expectEmpty('white space only', ' \r\n\t \n')
    expectError('prose', 'not an 835 at all')
    expectError('no ISA at the start', 'GS*HP*A*B*20261030*0905*1*X*005010X221A1~ST*835*0001~')
    expectError('ISA with fewer than 16 elements', 'ISA*00*          *00*          *ZZ*PAYER99002     ')
    expectError('ISA that ends at ISA16', 'ISA*00*          *00*          *ZZ*PAYER99002     *ZZ*VALLEYCLINIC   *261030*0905*^*00501*000000905*0*T*:')
    expectError('equal separators', 'ISA*00*          *00*          *ZZ*PAYER99002     *ZZ*VALLEYCLINIC   *261030*0905*^*00501*000000905*0*T***GS*HP~')
    const sample = examples.find((e) => e.id === DEFAULT_COLUMNS.ref) ?? examples[0]
    const lines = sample.file.split('~').map((s) => s.trim()).filter(Boolean)
    expectError('no IEA', lines.slice(0, -1).map((s) => `${s}~`).join('\n'))
    expectError('segments after the IEA', `${sample.file}${sample.file}`)
    expectError('a segment ID in lower case', sample.file.replace('\nCLP*', '\nclp*'))
    expectError('a segment ID that is too long', sample.file.replace('\nCLP*', '\nCLPX*'))
    // (a cut inside the last segment, the IEA, still holds the interchange's envelope, so the cuts stop before it)
    for (let cut = 1; cut < sample.file.lastIndexOf('IEA'); cut += 37) {
      pastedChecks++
      const r = readingOfText(sample.file.slice(0, cut))
      if (!isReading(r)) fail(`cut at ${cut}: not a reading: ${JSON.stringify(r)}`)
      else if (!r.error && !r.file && sample.file.slice(0, cut).trim() !== '') fail(`cut at ${cut}: neither a file nor an error`)
      else if (r.file) fail(`cut at ${cut}: a cut-short file read as whole`)
    }
    // A segment with nothing in it where the rows expect elements: whatever
    // the outcome, no throw, and a file that is returned builds its list.
    for (const [label, mend] of [
      ['an SVC with no elements', (t) => t.replace(/\nSVC\*[^~]*~/, '\nSVC~')],
      ['a CLP with no elements', (t) => t.replace(/\nCLP\*[^~]*~/, '\nCLP~')],
      ['a PLB with no elements', (t) => t.replace(/\nSE\*/, '\nPLB~\nSE*')],
      ['a BPR with no elements', (t) => t.replace(/\nBPR\*[^~]*~/, '\nBPR~')],
    ]) {
      pastedChecks++
      try {
        const r = readingOfText(mend(sample.file))
        if (!isReading(r)) fail(`${label}: not a reading: ${JSON.stringify(r)}`)
        else if (r.file) listOfReadings({ ref: { file: parsed.get(sample.id) }, a: r, b: {} })
      } catch (err) {
        fail(`${label}: threw ${err instanceof Error ? err.message : String(err)}`)
      }
    }
    // Pasted example text reads as the example, and sits against it with no difference.
    for (const ex of examples) {
      pastedChecks++
      if (exampleTextOf(ex.id) !== ex.file) fail(`${ex.id}: exampleTextOf is not the stored file`)
      const r = readingOfText(ex.file)
      if (!r.file) {
        fail(`${ex.id}: its own text does not read: ${r.error}`)
        continue
      }
      const original = parsed.get(ex.id)
      if (segmentsOf(r.file).length !== segmentsOf(original).length) fail(`${ex.id}: pasted text holds ${segmentsOf(r.file).length} segments, the example ${segmentsOf(original).length}`)
      const keys = original.claims.map((c) => c.key?.id)
      const repeated = new Set(keys).size !== keys.length
      const built = listOfReadings({ ref: { file: original }, a: r, b: r })
      if (built.error) fail(`${ex.id}: ${built.error}`)
      if (!repeated && judged(built.sections).length > 0) fail(`${ex.id}: pasted text differs from its own example in ${judged(built.sections).length} rows`)
    }
    // A column with no file (blank, or an error) has no cell, and nothing is marked in it.
    {
      const file = parsed.get(DEFAULT_COLUMNS.ref)
      const none = listOfReadings({ ref: { file }, a: {}, b: { error: 'x12: for the check' } })
      if (none.error) fail(`columns with no file: ${none.error}`)
      for (const s of none.sections) {
        if (s.held.a || s.held.b) fail(`columns with no file: “${s.title}” says a column holds it`)
        for (const r of s.rows) {
          if (r.values.a !== undefined || r.values.b !== undefined || r.from.a || r.from.b) fail(`columns with no file: a row of “${s.title}” has a value in a column with none`)
          for (const c of COLS) if (isMarked(s, r, c)) fail(`columns with no file: a row of “${s.title}” marks ${c}`)
        }
      }
      if (none.sections.length === 0) fail('a reference with no companions has no rows')
      const empty = listOfReadings({ ref: {}, a: {}, b: {} })
      if (empty.error || empty.sections.length !== 0) fail('three empty columns do not give an empty list')
      const alone = listOfReadings({ ref: {}, a: { file }, b: {} })
      if (alone.error || alone.sections.length === 0) fail('a compared column alone has no rows')
      for (const s of alone.sections) for (const r of s.rows) for (const c of COLS) if (isMarked(s, r, c)) fail(`a compared column alone is marked in “${s.title}”`)
    }

    // ---- presets: the text form, the reader's own, and what a hidden set shows
    // A page that has no storage, or storage that refuses or holds junk, still works.
    const withStorage = (storage, run) => {
      if (storage === undefined) return run()
      const had = Object.getOwnPropertyDescriptor(globalThis, 'localStorage')
      Object.defineProperty(globalThis, 'localStorage', { value: storage, configurable: true, writable: true })
      try {
        return run()
      } finally {
        if (had) Object.defineProperty(globalThis, 'localStorage', had)
        else delete globalThis.localStorage
      }
    }
    const fakeStorage = (initial = {}) => {
      const held = { ...initial }
      return { held, getItem: (k) => (k in held ? held[k] : null), setItem: (k, v) => void (held[k] = String(v)) }
    }
    let presetTexts = 0
    {
      const sameJSON = (x, y) => JSON.stringify(x) === JSON.stringify(y)
      const BOM = String.fromCharCode(0xfeff) // a byte order mark, written as its code so the source holds no invisible character
      // a generator with a fixed seed, so a failure repeats
      let seed = 20261002
      const rand = () => {
        seed = (Math.imul(seed, 1664525) + 1013904223) >>> 0
        return seed / 2 ** 32
      }
      const shuffled = (list) => list.map((x) => [rand(), x]).sort((p, q) => p[0] - q[0]).map(([, x]) => x)
      const names = ['', 'My view', '  padded   name ', 'a: b', 'preset: nested', 'hide:', 'Ünïcode · name', 'x'.repeat(MAX_NAME), 'tab\tseparated', 'two\nlines', 'Claim level']
      // keys of the right shape that a file may or may not hold
      const invented = ['claim.ZZZ', 'line.XYZ.XYZ01', 'payer.N1.N101', 'summary', 'payment.REF']

      // the built-in presets and random hidden sets (some with keys the tree lacks), read back as written
      for (const ex of examples) {
        const file = parsed.get(ex.id)
        const tree = treeOf(listOf({ ref: file, a: file, b: file }))
        const all = [...idsInTree(tree)]
        const sets = [...builtInPresets(tree).map((p) => p.hide), [], all]
        for (let n = 0; n < 12; n++) sets.push([...new Set(shuffled([...all.filter(() => rand() < 0.3), ...invented.filter(() => rand() < 0.4)]))])
        for (const hide of sets) {
          const name = names[Math.floor(rand() * names.length)]
          presetTexts++
          const back = parsePreset(formatPreset(name, hide))
          if (!sameJSON(back, { name: normalizeName(name), hide })) fail(`${ex.id}: ${JSON.stringify(name)} with ${hide.length} keys reads back as ${JSON.stringify(back).slice(0, 120)}`)
        }
      }

      // text that is not a preset: an error, short, and no throw
      const malformed = [
        ['empty', ''],
        ['white space only', ' \n\t \r\n'],
        ['prose', 'not a preset'],
        ['a name and no hide line', 'preset: Only a name'],
        ['a key on the hide line', 'hide: claim.CLP'],
        ['a key before the hide line', 'claim.CLP\nhide:\n  line'],
        ['a key with an empty part', 'hide:\n  claim..CLP'],
        ['an unknown part', 'hide:\n  foo.bar'],
        ['a part in capitals', 'hide:\n  Claim.CLP'],
        ['a key with four parts', 'hide:\n  claim.CLP.CLP03.x'],
        ['a space inside a key', 'hide:\n  claim.NM1*Q C'],
        ['a comma between keys', 'hide:\n  claim.CLP,line'],
        ['a bare dot', 'hide:\n  .'],
        ['a part and a dot', 'hide:\n  claim.'],
        ['two name lines', 'preset: a\npreset: b\nhide:\n  line'],
        ['a name line after hide', 'hide:\n  line\npreset: a'],
        ['a name over the limit', `preset: ${'x'.repeat(MAX_NAME + 1)}\nhide:`],
        ['an 835', examples[0].file],
        ['a very long line', `hide:\n${'q'.repeat(100000)}`],
        ['a null byte', 'hide:\n  claim\0'],
        ['null', null],
        ['undefined', undefined],
        ['a number', 42],
        ['an object', {}],
      ]
      for (const [label, text] of malformed) {
        presetTexts++
        let r
        try {
          r = parsePreset(text)
        } catch (err) {
          fail(`preset text (${label}): threw ${err instanceof Error ? err.message : String(err)}`)
          continue
        }
        if (r === null || typeof r !== 'object' || typeof r.error !== 'string' || r.error === '' || r.error.length > 200 || 'name' in r || 'hide' in r) {
          fail(`preset text (${label}): not one short error: ${JSON.stringify(r).slice(0, 160)}`)
        }
      }

      // text that is a preset, written loosely
      const accepted = [
        ['CRLF line ends', 'preset: A\r\nhide:\r\n  line\r\n', { name: 'A', hide: ['line'] }],
        ['no name line', 'hide:\n  line', { name: '', hide: ['line'] }],
        ['no keys', 'hide:', { name: '', hide: [] }],
        ['an empty name', 'preset:\nhide:\n  line', { name: '', hide: ['line'] }],
        ['blank lines, tabs, no indent', '\n\npreset:   My   view \n\nhide:\n\tline\nclaim.NM1*QC\n\n', { name: 'My view', hide: ['line', 'claim.NM1*QC'] }],
        ['a key written twice', 'hide:\n  line\n  line', { name: '', hide: ['line'] }],
        ['the words in capitals', 'PRESET: X\nHIDE:\n  line', { name: 'X', hide: ['line'] }],
        ['a byte order mark', `${BOM}hide:\n  line`, { name: '', hide: ['line'] }],
      ]
      for (const [label, text, want] of accepted) {
        presetTexts++
        const r = parsePreset(text)
        if (!sameJSON(r, want)) fail(`preset text (${label}) reads as ${JSON.stringify(r)}, not ${JSON.stringify(want)}`)
      }

      // random text from the pieces of a preset: no throw, and whatever reads is made of keys and reads back the same
      const pieces = ['preset:', 'hide:', 'claim', 'line', 'payment', '.', 'CLP', 'CLP03', '*', 'QC', '\n', '\n', '\n', '\r\n', ' ', '  ', '\t', ',', 'x', '', 'ÿ', BOM, '\0']
      for (let n = 0; n < 3000; n++) {
        const text = Array.from({ length: 1 + Math.floor(rand() * 24) }, () => pieces[Math.floor(rand() * pieces.length)]).join('')
        presetTexts++
        let r
        try {
          r = parsePreset(text)
        } catch (err) {
          fail(`random preset text ${JSON.stringify(text)}: threw ${err instanceof Error ? err.message : String(err)}`)
          continue
        }
        if ('error' in r) {
          if (typeof r.error !== 'string' || r.error === '') fail(`random preset text ${JSON.stringify(text)}: an empty error`)
          continue
        }
        if (!r.hide.every(isKey) || new Set(r.hide).size !== r.hide.length) fail(`random preset text ${JSON.stringify(text)}: keys that are not keys, or a key twice`)
        if (!sameJSON(parsePreset(formatPreset(r.name, r.hide)), r)) fail(`random preset text ${JSON.stringify(text)}: does not read back from its own text`)
      }

      // names
      if (nameProblem('') === undefined || nameProblem('   ') === undefined) fail('an empty name can be saved')
      if (nameProblem('x'.repeat(MAX_NAME + 1)) === undefined) fail('a name over the limit can be saved')
      for (const reserved of [...BUILT_IN_NAMES, 'Custom', 'everything', 'CLAIM LEVEL', ' custom ']) if (nameProblem(reserved) === undefined) fail(`the name “${reserved}” can be saved`)
      for (const fine of ['My view', 'Claim level 2', 'x'.repeat(MAX_NAME)]) if (nameProblem(fine) !== undefined) fail(`the name “${fine}” cannot be saved: ${nameProblem(fine)}`)

      // the reader's own list: adding, replacing (in any case), removing, and reading what storage held
      const list = [{ name: 'Mine', hide: ['claim.MOA', 'line'] }, { name: 'Other', hide: [] }]
      if (!sameJSON(withPreset(list, { name: 'Third', hide: ['payment'] }), [...list, { name: 'Third', hide: ['payment'] }])) fail('a new saved preset is not added at the end')
      if (!sameJSON(withPreset(list, { name: 'MINE', hide: ['payer'] }), [{ name: 'MINE', hide: ['payer'] }, list[1]])) fail('a saved preset of the same name is not replaced in place')
      if (!sameJSON(withoutPreset(list, 'mine'), [list[1]])) fail('a saved preset is not removed by its name')
      if (!sameJSON(readSavedPresets(JSON.stringify(list)), list)) fail('a saved list does not read back as stored')
      const stored = [
        ['no value', null, []],
        ['no value, undefined', undefined, []],
        ['an empty value', '', []],
        ['text that is not JSON', '{not json', []],
        ['an object', '{"name":"a","hide":[]}', []],
        ['a list of numbers', '[1,2,3]', []],
        ['a list holding null', '[null]', []],
        ['a name that is not text', '[{"name":1,"hide":[]}]', []],
        ['no hide', '[{"name":"a"}]', []],
        ['a hide that is a number', '[{"name":"a","hide":5}]', []],
        ['a key that is a number', '[{"name":"a","hide":[1]}]', []],
        ['a key that is not a key', '[{"name":"a","hide":["claim..CLP"]}]', []],
        ['a built-in name', '[{"name":"Everything","hide":[]}]', []],
        ['the name Custom', '[{"name":"custom","hide":[]}]', []],
        ['an empty name', '[{"name":" ","hide":[]}]', []],
        ['a name over the limit', `[{"name":"${'x'.repeat(MAX_NAME + 1)}","hide":[]}]`, []],
        ['hide as the text hide= holds', '[{"name":"a","hide":"claim.MOA,line"}]', [{ name: 'a', hide: ['claim.MOA', 'line'] }]],
        ['the same name twice, in two cases', '[{"name":"A","hide":["line"]},{"name":"a","hide":[]}]', [{ name: 'A', hide: ['line'] }]],
        ['a key twice', '[{"name":"a","hide":["line","line"]}]', [{ name: 'a', hide: ['line'] }]],
        ['good entries among bad ones', '[7,{"name":"a","hide":["line"]},{"name":2},{"name":"b","hide":[]}]', [{ name: 'a', hide: ['line'] }, { name: 'b', hide: [] }]],
      ]
      for (const [label, text, want] of stored) {
        presetTexts++
        let r
        try {
          r = readSavedPresets(text)
        } catch (err) {
          fail(`saved presets (${label}): threw ${err instanceof Error ? err.message : String(err)}`)
          continue
        }
        if (!sameJSON(r, want)) fail(`saved presets (${label}) read as ${JSON.stringify(r)}, not ${JSON.stringify(want)}`)
      }

      // storage: none, refusing, holding junk, and working
      if (typeof localStorage === 'undefined') {
        if (loadSavedPresets().length !== 0) fail('no storage: the saved list is not empty')
        if (storeSavedPresets(list) !== false) fail('no storage: saving the list says it was kept')
      }
      const refusing = { getItem: () => { throw new Error('refused') }, setItem: () => { throw new Error('refused') } }
      if (withStorage(refusing, loadSavedPresets).length !== 0) fail('storage that refuses a read gives a saved list')
      if (withStorage(refusing, () => storeSavedPresets(list)) !== false) fail('storage that refuses a write says the list was kept')
      if (withStorage({ getItem: () => '{not json', setItem: () => {} }, loadSavedPresets).length !== 0) fail('storage holding junk gives a saved list')
      {
        const storage = fakeStorage()
        if (withStorage(storage, () => storeSavedPresets(list)) !== true) fail('working storage says the list was not kept')
        const keys = Object.keys(storage.held)
        if (keys.length !== 1 || keys[0] !== PRESETS_KEY) fail(`saving the list wrote ${JSON.stringify(keys)}, not only ${PRESETS_KEY}`)
        if (!sameJSON(JSON.parse(storage.held[PRESETS_KEY]), list)) fail(`the stored list is ${storage.held[PRESETS_KEY]}`)
        if (!sameJSON(withStorage(storage, loadSavedPresets), list)) fail('the stored list does not load back')
      }

      // what a hidden set shows
      const file = parsed.get(DEFAULT_COLUMNS.ref)
      const tree = treeOf(listOf({ ref: file, a: file, b: file }))
      const first = tree[0].id
      const last = tree[tree.length - 1].id
      if (viewOf(tree, []) !== '') fail('an empty hidden set hides something')
      if (viewOf(tree, ['claim.ZZZ', 'line.XYZ', 'line.XYZ.XYZ01', 'payer.ZZZ.ZZZ01']) !== '') fail('ids the tree has no node for hide something')
      if (viewOf(tree, [first]) === '') fail(`hiding ${first} hides nothing`)
      const under = tree.flatMap(function ids(node) { return [node.id, ...node.children.flatMap(ids)] }).filter((id) => id.startsWith(`${first}.`))
      if (viewOf(tree, [first]) !== viewOf(tree, [first, ...under])) fail(`hiding ${first} and everything under it is another view than hiding ${first}`)
      if (!sameJSON(inTreeOrder(tree, [last, 'claim.ZZZ', first, last, 'summary.ZZZ']), [first, last, 'claim.ZZZ', 'summary.ZZZ'])) fail(`inTreeOrder gives ${JSON.stringify(inTreeOrder(tree, [last, 'claim.ZZZ', first, last, 'summary.ZZZ']))}`)
      // the option the select shows: the one whose rows match, the chosen one among equals, else none (Custom)
      {
        const options = [{ value: 'a', hide: [] }, { value: 'b', hide: ['claim.ZZZ'] }, { value: 'c', hide: ['line'] }]
        const pick = (hide, chosen) => matchingOption(tree, options, viewOf(tree, hide), chosen)?.value ?? null
        if (pick([], '') !== 'a') fail('the first option that shows the same rows is not the one matched')
        if (pick([], 'b') !== 'b') fail('the chosen option is not preferred among options that show the same rows')
        if (pick([], 'c') !== 'a') fail('a chosen option that shows other rows is matched')
        if (pick(['line'], 'a') !== 'c') fail('the option that shows the hidden set is not matched')
        if (pick(['line', 'line.SVC'], '') !== 'c') fail('a hidden set that lists a key under another does not match')
        if (pick([first], '') !== null) fail(`hiding ${first} matches an option that does not`)
      }
      const everything = builtInPresets(tree).find((p) => p.name === 'Everything')
      if (!everything || everything.hide.length !== 0) fail('Everything hides something')
      const others = builtInPresets(tree).filter((p) => p.name !== 'Everything')
      if (others.some((p) => viewOf(tree, p.hide) === '')) fail('a built-in preset other than Everything hides nothing in the default three files')
      if (new Set(builtInPresets(tree).map((p) => viewOf(tree, p.hide))).size !== BUILT_IN_NAMES.length) fail('two built-in presets show the same rows in the default three files')

      // the code of the preset controls keeps the box's text out of storage and the address, and storage is read and written only for the saved list
      for (const name of ['PresetControls.tsx', 'usePresets.ts']) {
        const source = await readFile(new URL(`../src/docs/compare835/${name}`, import.meta.url), 'utf8')
        const found = source.match(/localStorage|sessionStorage|indexedDB|useSearchParams|setSearchParams|location\b|history\b|document\.cookie|fetch\(|XMLHttpRequest|sendBeacon/g)
        if (found) fail(`${name} mentions ${[...new Set(found)].join(', ')}: the preset box's text must not reach storage, the address, or the network`)
      }
      const presetsSource = await readFile(new URL('../src/docs/compare835/presets.ts', import.meta.url), 'utf8')
      const uses = [...presetsSource.matchAll(/localStorage\.(\w+)\(([^,)]*)/g)]
      if (uses.length !== 2 || uses.some((m) => m[2].trim() !== 'PRESETS_KEY')) fail('presets.ts reads or writes storage other than the saved list')
    }

    // ---- the page and its panel, rendered on the server
    {
      const { Compare835 } = await server.ssrLoadModule('/src/docs/pages/Compare835.tsx')
      const { TextPanel } = await server.ssrLoadModule('/src/docs/compare835/TextPanel.tsx')
      const escape = (t) => t.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;').replaceAll('"', '&quot;').replaceAll("'", '&#x27;')
      const count = (html, text) => html.split(text).length - 1
      const ids = { ref: 'paid-by-ach', a: 'denied', b: 'two-claims' }
      const html = renderToStaticMarkup(
        createElement(StaticRouter, { location: `/835-compare?ref=${ids.ref}&a=${ids.a}&b=${ids.b}` }, createElement(Compare835)),
      )
      // the three raw boxes and the preset text box
      if (count(html, '<textarea') !== 4) fail(`the page has ${count(html, '<textarea')} text boxes, not 4`)
      for (const col of COLS) {
        if (!html.includes(`for="compare835-text-${col}"`)) fail(`no label for the ${col} box`)
        if (!html.includes(`>${escape(exampleTextOf(ids[col]))}</textarea>`)) fail(`the ${col} box does not hold the file of ${ids[col]}, as stored`)
      }
      if (html.includes('>Pasted text</option>')) fail('a column that shows an example has a "Pasted text" entry')
      if (html.includes('role="alert"')) fail('the page shows an error with nothing pasted')
      if (!html.includes('Pasted text stays on this page')) fail('the line saying pasted text stays on the page is missing')
      if (count(html, '>Copy</button>') !== 4 || count(html, '>Clear</button>') !== 3) fail('each raw box does not have a Copy and a Clear control, or the preset box has no Copy')

      // Every code is a link into the 835 lookup (docs/viewer/Code.tsx): a
      // row's key (CLP01), a qualified row's note (NM1*QC), a tree group (CLP)
      // and element (CLP01); a CAS row's key (a group code and a reason code,
      // both values) is not. The panel opens from the address beside the page.
      for (const code of ['CLP', 'CLP01', 'NM1*QC']) if (!html.includes(`data-code="${escape(code)}"`)) fail(`the page has no lookup link for ${code}`)
      if (/data-code="(CO|PR|OA)-/.test(html)) fail('a CAS row\'s key, a pair of coded values, is a lookup link')
      if (html.includes('role="complementary"')) fail('the lookup panel is open with no lookup=1')
      const open = renderToStaticMarkup(
        createElement(StaticRouter, { location: `/835-compare?ref=${ids.ref}&a=${ids.a}&b=${ids.b}&lookup=1&q=${encodeURIComponent('NM1*QC')}` }, createElement(Compare835)),
      )
      const at = open.indexOf('role="complementary"')
      if (at < 0) fail('lookup=1 does not open the 835 lookup panel')
      else {
        if (!open.slice(at).includes('value="NM1*QC"')) fail('the panel opened with q=NM1*QC does not hold it in its search')
        if (!open.slice(at).includes('data-row="claim:NM109@QC"')) fail('the panel searched for NM1*QC does not show NM109 under QC')
      }
      if (!open.includes('class="@container min-w-0 flex-1"') || !/class="sticky top-0 h-svh w-\[26rem\][^"]*"/.test(open)) fail('the lookup panel is not a column beside the page')

      // A box that holds text that does not read: its text as typed, the error under it, the other boxes quiet.
      const typed = 'ISA*not an interchange <&>'
      const bad = renderToStaticMarkup(
        createElement(TextPanel, { texts: { ref: typed, a: '', b: '' }, errors: { ref: 'x12: for the check' }, onText: () => {} }),
      )
      if (!bad.includes(`>${escape(typed)}</textarea>`)) fail('a box with an error does not keep its text as typed')
      if (count(bad, 'role="alert"') !== 1 || !bad.includes('x12: for the check')) fail('the error is not shown once, under its box')
      if (count(bad, 'aria-invalid="true"') !== 1) fail('only the box with the error is marked invalid')
      if (count(bad, 'disabled=""') !== 4) fail('the Copy and Clear of the two empty boxes are not disabled')
    }

    // ---- the presets on the page, rendered on the server
    {
      const { Compare835 } = await server.ssrLoadModule('/src/docs/pages/Compare835.tsx')
      const escape = (t) => t.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;').replaceAll('"', '&quot;').replaceAll("'", '&#x27;')
      const count = (html, text) => html.split(text).length - 1
      const ids = { ref: 'paid-by-ach', a: 'denied', b: 'two-claims' }
      const pageTree = treeOf(listOf({ ref: parsed.get(ids.ref), a: parsed.get(ids.a), b: parsed.get(ids.b) }))
      const made = Object.fromEntries(builtInPresets(pageTree).map((p) => [p.name, p.hide]))
      // one element of the CLP line: hiding it alone matches no built-in preset
      const leaf = pageTree.find((n) => n.id === 'claim')?.children.find((g) => g.id === 'claim.CLP')?.children[0]?.id
      if (!pageTree.some((n) => n.id === 'line') || !leaf) {
        fail('the three examples of the page check have no service line, or no CLP element')
      } else {
        const render = (hide, storage) =>
          withStorage(storage, () =>
            renderToStaticMarkup(
              createElement(
                StaticRouter,
                { location: `/835-compare?ref=${ids.ref}&a=${ids.a}&b=${ids.b}${hide === undefined ? '' : `&hide=${encodeURIComponent(hide)}`}` },
                createElement(Compare835),
              ),
            ),
          )
        const selectOf = (html) => {
          const inner = /<select id="compare835-preset"[^>]*>(.*?)<\/select>/s.exec(html)?.[1] ?? ''
          return {
            groups: [...inner.matchAll(/<optgroup label="([^"]*)"/g)].map((m) => m[1]),
            options: [...inner.matchAll(/<option([^>]*)>(.*?)<\/option>/g)].map((m) => ({ label: m[2], selected: /\bselected=""/.test(m[1]) })),
          }
        }
        const boxOf = (html) => /<textarea id="compare835-preset-text"[^>]*>(.*?)<\/textarea>/s.exec(html)?.[1]
        const deleteDisabled = (html) => /<button[^>]*\bdisabled=""[^>]*>Delete<\/button>/.test(html)
        const BUILT = ['Everything', 'Claim level', 'Service lines', 'Amounts and adjustments', 'Codes and dates']
        const expect = (what, html, { selected, saved = [], custom = false, name, hide, canDelete = false }) => {
          const { groups, options } = selectOf(html)
          const labels = options.map((o) => o.label)
          const chosen = options.filter((o) => o.selected).map((o) => o.label)
          if (JSON.stringify(labels) !== JSON.stringify([...BUILT, ...saved, ...(custom ? ['Custom'] : [])])) fail(`${what}: the preset select lists ${JSON.stringify(labels)}`)
          if (JSON.stringify(groups) !== JSON.stringify(['Built in', ...(saved.length ? ['Saved'] : [])])) fail(`${what}: the preset select has the groups ${JSON.stringify(groups)}`)
          if (chosen.length !== 1 || chosen[0] !== selected) fail(`${what}: the preset select shows ${JSON.stringify(chosen)}, not ${selected}`)
          if (boxOf(html) !== escape(formatPreset(name, hide))) fail(`${what}: the preset box holds ${JSON.stringify(boxOf(html))}`)
          if (deleteDisabled(html) === canDelete) fail(`${what}: Delete is ${canDelete ? 'disabled' : 'enabled'} for ${selected}`)
        }

        const plain = render(undefined)
        expect('no hide in the address', plain, { selected: 'Everything', name: 'Everything', hide: [] })
        for (const part of ['id="compare835-preset"', 'id="compare835-preset-name"', 'id="compare835-preset-text"', 'id="compare835-preset-text-name"', 'aria-label="Preset text"']) {
          if (!plain.includes(part)) fail(`the page has no ${part}`)
        }
        if (count(plain, '>Save as…</label>') !== 2 || count(plain, '>Save</button>') !== 2 || count(plain, '>Apply</button>') !== 1 || count(plain, '>Delete</button>') !== 1) {
          fail('the page does not have two Save as… controls with a Save button each, one Apply, and one Delete')
        }
        if (plain.includes('role="alert"')) fail('the page shows an error with nothing typed in the preset box')
        // each built-in preset, as the address would hold it; where two show the same rows the first is the one shown
        for (const name of BUILT_IN_NAMES) {
          const hide = made[name]
          const first = BUILT_IN_NAMES.find((other) => viewOf(pageTree, made[other]) === viewOf(pageTree, hide))
          expect(`hide= holding ${name}`, render(hide.join(',')), { selected: first, name: first, hide })
        }
        expect('hide=line', render('line'), { selected: 'Claim level', name: 'Claim level', hide: ['line'] })
        expect('hide= with a key twice, and a key under another', render(`line.SVC,${leaf},line,${leaf}`), { selected: 'Custom', custom: true, name: '', hide: [leaf, 'line', 'line.SVC'] })
        expect('hide= one element', render(leaf), { selected: 'Custom', custom: true, name: '', hide: [leaf] })
        expect('hide= empty', render(''), { selected: 'Everything', name: 'Everything', hide: [] })
        expect('hide= naming only keys the tree lacks', render('claim.ZZZ'), { selected: 'Everything', name: 'Everything', hide: ['claim.ZZZ'] })

        // saved presets from storage
        const mine = JSON.stringify([{ name: 'Mine', hide: [leaf] }, { name: 'Plain', hide: ['line'] }])
        const storage = fakeStorage({ [PRESETS_KEY]: mine })
        expect('a saved preset, chosen by hide=', render(leaf, storage), { selected: 'Mine', saved: ['Mine', 'Plain'], name: 'Mine', hide: [leaf], canDelete: true })
        expect('a saved preset that shows what a built-in does', render('line', storage), { selected: 'Claim level', saved: ['Mine', 'Plain'], name: 'Claim level', hide: ['line'] })
        expect('a hide= that no saved preset has', render('payment', storage), { selected: 'Custom', saved: ['Mine', 'Plain'], custom: true, name: '', hide: ['payment'] })
        expect('the saved hidden set, with no hide= in the address', render(undefined, fakeStorage({ [PRESETS_KEY]: mine, 'docs.835-compare.hide': leaf })), { selected: 'Mine', saved: ['Mine', 'Plain'], name: 'Mine', hide: [leaf], canDelete: true })

        // storage that holds junk, or refuses
        expect('storage holding junk', render(undefined, fakeStorage({ [PRESETS_KEY]: '{not json' })), { selected: 'Everything', name: 'Everything', hide: [] })
        expect('storage holding a built-in name', render(undefined, fakeStorage({ [PRESETS_KEY]: '[{"name":"Everything","hide":[]}]' })), { selected: 'Everything', name: 'Everything', hide: [] })
        expect('storage that refuses', render(undefined, { getItem: () => { throw new Error('refused') }, setItem: () => { throw new Error('refused') } }), { selected: 'Everything', name: 'Everything', hide: [] })
      }
    }

    // ---- the defaults
    const shared = [...keysOf(DEFAULT_COLUMNS.ref)].filter((k) => keysOf(DEFAULT_COLUMNS.a).has(k) && keysOf(DEFAULT_COLUMNS.b).has(k))
    if (shared.length === 0) fail('the default examples share no claim')
    const none = resolveColumns({ ref: null, a: null, b: null })
    if (JSON.stringify(none) !== JSON.stringify(DEFAULT_COLUMNS)) fail('an empty address does not resolve to the defaults')
    const fromViewer = resolveColumns({ ref: 'paid-by-ach', a: null, b: null })
    if (new Set(Object.values(fromViewer)).size !== 3) fail(`a link naming only the reference resolves to ${JSON.stringify(fromViewer)}`)
    const bad = resolveColumns({ ref: 'no-such-example', a: 'two-claims', b: null })
    if (new Set(Object.values(bad)).size !== 3 || bad.a !== 'two-claims') fail(`an unknown id resolves to ${JSON.stringify(bad)}`)

    if (failures.length > 0) {
      console.error(`835 compare SSR check FAILED (${failures.length} problems):`)
      for (const f of failures.slice(0, 40)) console.error(`  - ${f}`)
      if (failures.length > 40) console.error(`  … and ${failures.length - 40} more`)
      process.exitCode = 1
      return
    }

    // What shares a claim, for the record.
    const byKey = new Map()
    for (const ex of examples) for (const k of keysOf(ex.id)) byKey.set(k, [...(byKey.get(k) ?? []), ex.id])
    const shares = [...byKey.values()].filter((ids) => ids.length > 1)
    console.log(`835 compare SSR check passed: ${examples.length} examples parsed with every segment placed and every claim keyed; ${lists} lists (${rows} rows, ${marked} marked cells) built with every segment outside the envelope on a row and every row in the tree.`)
    console.log(`five-cent rule: ${withinRule} cells in those lists hold a dollar amount within ${AMOUNT_MATCH_CENTS} cents of the reference's and are left unmarked; ${fiveCentCases.length * 2} value pairs and ${movedCells} cells of two examples with amounts moved by a few cents are marked as the rule says.`)
    console.log(`${rewritten} rewrites of the examples (${Object.keys(variants).length} delimiter and line-break variants each) read as the originals; ${pastedChecks} pasted-text readings checked (blank, unreadable, cut short, odd segments, and the examples' own text).`)
    console.log(`presets: ${presetTrees} trees (every list above, and each example alone) gave the ${BUILT_IN_NAMES.length} built-in presets with only keys the tree holds, each showing the rows its rule says; ${presetTexts} preset texts, saved lists, and storage cases checked (read back as written, malformed, random, junk, refused).`)
    console.log(`defaults ${JSON.stringify(DEFAULT_COLUMNS)} share ${shared.length} claim key(s); the reference alone resolves to ${JSON.stringify(fromViewer)}.`)
    console.log(`${shares.length} claim keys are held by more than one example:`)
    for (const ids of shares) console.log(`  ${ids.join(', ')}`)
  } finally {
    await server.close()
  }
}

main().catch((err) => {
  console.error(err)
  process.exitCode = 1
})
