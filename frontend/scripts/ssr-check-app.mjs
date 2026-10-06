// Loads the validator app's pages through Vite's SSR module graph, the
// way the docs-site checks do, and checks with no browser and no server:
//   - the overview's fixtures (src/app/fixtures/front-*.json,
//     agreement.json) add up: each side's outcome rows sum to its totals,
//     claims and MAC claims; MAC claims are every outcome but extra; its
//     status matrix sums to its pairs (equal plus different); every one of
//     the twelve dimensions is present and, with no selection, sums to the
//     totals, claims and MAC claims; the payment bands are the five, match
//     through over 5%, and sum to the pairs; the MAC side (the totals, the
//     set-aside row, every value of a dimension the MAC claim carries) is
//     the same for slots 1 and 2, as it is for two vendors compared with one
//     MAC submission; the legacy block is the same in both responses, lists
//     the nine dimensions the MAC claim carries, each adding up to its
//     totals, and counts at least the MAC claims the vendor rows hold; and
//     the agreement grid's rows and columns, named v1 and v2, add up to slot
//     1's and slot 2's outcomes;
//   - the overview (src/app/pages/Overview.tsx) is titled Overview, and so
//     are its navigation entry and its help;
//   - the overview's selection header (src/app/SelectionHeader.tsx), rendered
//     on the server from those fixtures and the /api/macs fixture: it is
//     sticky at the top of the window, at most 40% of the window tall with
//     its own scrolling, and sits directly in the view's root, so it holds
//     for the whole page and the content starts where it ends; it offers
//     the nine boxes in order, each with All and None at its top and then
//     its values as tabs with no count in them, the plain, picked, and
//     dimmed tabs each in one set of classes; the jurisdiction box sets the
//     sixteen jurisdictions four to a row in CMS's order (the order checked
//     against the /api/macs fixture), and the other boxes as many to a row
//     as their values or their cap; a box lists exactly the values with
//     claims in either slot's facet (the hidden-zero rule): a facet narrowed
//     the way /api/front narrows it under a picked jurisdiction or claim
//     type hides the values it no longer lists, so picking JE leaves only
//     JE's workloads, while a picked value with no claims stays, pressed and
//     dimmed; the workload box starts closed, one line with its count (all
//     of them, the number picked, or none), and open lists every workload
//     once, under its jurisdiction in CMS's order, a roll-up as a group
//     headed by its number, the numbers the MAC map does not hold last, as
//     chips smaller than the other boxes' tabs; no checkbox panel is left;
//   - the overview's view (src/app/OverviewView.tsx) under the header, from
//     the same fixtures: every table sits on the one grid of
//     four columns (overviewColumns.tsx GRID; the status table's MATRIX_GRID,
//     whose slot columns widen only for a matrix too wide for them), its
//     heads Legacy, 1, 2 in that order, the Legacy column empty in exactly
//     the tables where the legacy side has no figure of its own (payment
//     agreement, the bands chart, what differs, and the breakdowns by vendor
//     status, payment band, and outcome); the agreement grid spans the three
//     figure columns; there are twelve breakdown tables, one per dimension,
//     each with every value's row and slot 1's and slot 2's counts, the nine
//     by what the MAC claim carries with a legacy figure per row adding up
//     to the legacy totals; the outcome table has a legacy figure on its
//     Total line only; the status table has the legacy count beside each MAC
//     status; the bars, the agreement grid, the five payment bands, the
//     status matrices with their diagonals, and every tag are drawn, slot 1
//     before slot 2; a narrowed selection is drawn as such and leaves the
//     legacy figures alone when it narrows only what a vendor decides; one
//     empty slot leaves its column empty; an error and the updating mark
//     show;
//   - the selection round-trips through the address (src/app/
//     overviewFilters.ts): what writeSelection writes, readSelection reads
//     back, an empty list and the empty value included; a tab click picks a
//     value, a second unpicks it, and unpicking the last returns the box to
//     All; picking a workload writes workload= as /api/front takes it;
//   - the claim screen's view (src/app/ClaimCompareView.tsx) reads the
//     fixture claims with the compare page's claim reader (docs/compare835/
//     file.ts parseClaim), builds the claim and line rows, marks the cells
//     whose values differ from the MAC's, with the comparison's five-cent
//     rule for dollar amounts (slot 1's payment and adjustment two cents
//     off are not marked; moved six cents they are, with their signed
//     differences), tags the rows the compared set leaves out, and offers
//     the tree, the switch, and the presets;
//   - the code popovers (src/data/codeHints.ts, drawn by src/app/CodeHint.tsx)
//     for every tag in the fixtures and every claim and line element: a
//     title of a few words saying what the thing is, then one paragraph
//     that does not repeat it, the maintainer's text for an adjustment or
//     remark code (from backend/pkg/carc.json and rarc.json), and no
//     document citation anywhere;
//   - every app link names a root page: the whole app (src/app/App.tsx) is
//     rendered at each of its pages (the overview at /, the claim screen,
//     the four review pages, and a missing page), with the API answered
//     from the fixtures above and small ones written here, and every href
//     in the markup is read; each link to an app page names one of the
//     app's pages at the root. Each render shows the page's own links (a
//     claim, a tag, an outcome), so the pages are checked with their data,
//     not only in their loading state;
//   - nothing in the overview's or the claim screen's drawing code reads
//     the network at render time.
// Run from frontend/:
//   node scripts/ssr-check-app.mjs
import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { StaticRouter } from 'react-router'
import { readFile } from 'node:fs/promises'
import { createServer } from 'vite'

const count = (html, text) => html.split(text).length - 1
const escapeHtml = (text) => text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#x27;')
const sum = (list, read) => list.reduce((n, item) => n + read(item), 0)
const sameList = (x, y) => JSON.stringify(x) === JSON.stringify(y)
// An exact decimal string in cents, for adding amounts without floats.
const cents = (amount) => {
  const negative = amount.startsWith('-')
  const [whole, fraction = ''] = amount.replace('-', '').split('.')
  const n = Number(whole) * 100 + Number((fraction + '00').slice(0, 2))
  return negative ? -n : n
}

// The markup from `start` to the first of `ends` after it, or to the end.
function region(html, start, ends = []) {
  const from = html.indexOf(start)
  if (from < 0) return ''
  let to = html.length
  for (const end of ends) {
    const at = html.indexOf(end, from + start.length)
    if (at >= 0 && at < to) to = at
  }
  return html.slice(from, to)
}

// The legacy cells in some markup (overviewColumns.tsx LegacyCell): the MAC
// claims and whether the row is in the selection.
const legacyEntries = (html) => [...html.matchAll(/data-legacy-claims="(\d+)" data-in="([01])"/g)].map((m) => ({ claims: Number(m[1]), in: m[2] === '1' }))
const legacySum = (entries) => sum(entries.filter((e) => e.in), (e) => e.claims)

// The texts of the count slots (overviewColumns.tsx Slots, left) in some markup, in order.
const countSlots = (html) => [...html.matchAll(/class="w-16 shrink-0 text-right font-mono[^"]*"[^>]*>([^<]*)</g)].map((m) => m[1])

// The figure columns' marks in some markup, in order: a head (legacy, v1,
// v2) or an empty column (overviewColumns.tsx Gap: gap:legacy, gap:v1, gap:v2).
const columnMarks = (html) => [...html.matchAll(/data-(column|gap)="(\w+)"/g)].map((m) => (m[1] === 'gap' ? `gap:${m[2]}` : m[2]))

// How many columns a grid template lays out: repeat(n, ...) counts n.
const trackCount = (template) => {
  let n = 0
  const rest = template.replace(/repeat\((\d+),\s*(?:[^()]|\([^()]*\))*\)/g, (_, k) => {
    n += Number(k)
    return ''
  })
  return n + rest.split(/\s+/).filter((x) => x !== '').length
}

// The nine boxes of the selection header, in its order, and how many tabs
// each sets in a row at most.
const BOXES = ['jurisdiction', 'claim_type', 'adjudication', 'mac_status', 'vendor_status', 'frequency', 'service_month', 'outcome', 'workload']
const COLUMNS = { jurisdiction: 4, claim_type: 3, adjudication: 3, mac_status: 5, vendor_status: 5, frequency: 4, service_month: 12, outcome: 6 }

// The pages read the API through useApi (src/api.ts), whose effect never
// runs in a server render, so a page rendered here shows only its loading
// state. For the link check, this check's own module graph has useApi first
// ask globalThis.answerApi(key), when the check sets it, so a page renders
// from fixtures; with nothing set, useApi is unchanged. src/api.ts itself is
// not edited.
const USE_API_LINE = 'const current = state.key === key'
const answerApi = {
  name: 'answer-api-from-fixtures',
  enforce: 'pre',
  transform(code, id) {
    if (!id.split('?')[0].endsWith('/src/api.ts')) return null
    if (!code.includes(USE_API_LINE)) throw new Error(`src/api.ts: useApi no longer has “${USE_API_LINE}”, where the link check answers it from fixtures`)
    return code.replace(USE_API_LINE, [
      'const answered = key === null ? undefined : (globalThis as any).answerApi?.(key)',
      '  if (answered !== undefined) return { data: answered, error: undefined, loading: false }',
      `  ${USE_API_LINE}`,
    ].join('\n'))
  },
}

// How deep in the markup the tag starting at `at` sits: the tags opened
// before it and not yet closed.
function depthAt(html, at) {
  let depth = 0
  for (const m of html.slice(0, at).matchAll(/<(\/?)[a-zA-Z][^>]*?(\/?)>/g)) {
    if (m[1] === '/') depth--
    else if (m[2] !== '/') depth++
  }
  return depth
}

async function main() {
  const server = await createServer({
    configFile: false,
    root: process.cwd(),
    plugins: [answerApi],
    server: { middlewareMode: true },
    appType: 'custom',
    logLevel: 'warn',
  })

  try {
    const { FRONT_FILTERS } = await server.ssrLoadModule('/src/api.ts')
    const filters = await server.ssrLoadModule('/src/app/overviewFilters.ts')
    const { GRID, MATRIX_GRID } = await server.ssrLoadModule('/src/app/overviewColumns.tsx')
    const { OverviewView } = await server.ssrLoadModule('/src/app/OverviewView.tsx')
    const { Overview } = await server.ssrLoadModule('/src/app/pages/Overview.tsx')
    const { HELP } = await server.ssrLoadModule('/src/app/help.ts')
    const { ClaimCompareView } = await server.ssrLoadModule('/src/app/ClaimCompareView.tsx')
    const { parseClaim } = await server.ssrLoadModule('/src/docs/compare835/file.ts')
    const { SelectionProvider } = await server.ssrLoadModule('/src/app/selection.tsx')
    const { count: countText, money, statusLabel } = await server.ssrLoadModule('/src/format.ts')
    const { WorkloadBox } = await server.ssrLoadModule('/src/app/SelectionHeader.tsx')
    const { workloadGroups } = await server.ssrLoadModule('/src/app/macLines.ts')
    const hints = await server.ssrLoadModule('/src/data/codeHints.ts')
    const { ELEMENTS_835 } = await server.ssrLoadModule('/src/data/lookupRows835.ts')
    const { default: App } = await server.ssrLoadModule('/src/app/App.tsx')

    const fixture = async (name) => JSON.parse(await readFile(new URL(`../src/app/fixtures/${name}`, import.meta.url), 'utf8'))
    const source = (path) => readFile(new URL(`../${path}`, import.meta.url), 'utf8')
    const v1 = await fixture('front-1.json')
    const v2 = await fixture('front-2.json')
    const agreement = await fixture('agreement.json')
    const claim = await fixture('claim.json')
    const macs = await fixture('macs.json')
    const legacy = v1.legacy

    // Each slot with its name in failure messages.
    const SLOTS = [['v1', 'slot 1', v1], ['v2', 'slot 2', v2]]
    // The five payment bands, in order, and the dimensions the legacy block carries.
    const BANDS = ['match', 'up to 1.00', 'up to 1%', 'up to 5%', 'over 5%']
    const LEGACY = filters.LEGACY_DIMENSIONS

    const failures = []
    const fail = (message) => failures.push(message)

    // ---- the fixtures add up
    for (const [, name, data] of SLOTS) {
      const total = data.totals.claims
      const macTotal = data.totals.mac_claims
      const outcomes = sum(data.outcomes, (o) => o.claims)
      if (outcomes !== total) fail(`${name}: the outcome rows sum to ${outcomes}, the totals say ${total}`)
      const macOutcomes = sum(data.outcomes, (o) => o.mac_claims)
      if (macOutcomes !== macTotal) fail(`${name}: the outcome rows hold ${macOutcomes} MAC claims, the totals say ${macTotal}`)
      for (const o of data.outcomes) {
        const want = o.outcome === 'extra' ? 0 : o.claims
        if (o.mac_claims !== want) fail(`${name}: outcome ${o.outcome} holds ${o.mac_claims} MAC claims, not ${want}`)
      }
      const extra = data.outcomes.find((o) => o.outcome === 'extra')?.claims ?? 0
      if (total - macTotal !== extra) fail(`${name}: claims less MAC claims is ${total - macTotal}, but ${extra} claims are extra`)
      const pairs = sum(data.outcomes.filter((o) => o.outcome === 'equal' || o.outcome === 'different'), (o) => o.claims)
      const matrix = sum(data.status_matrix, (c) => c.claims)
      if (matrix !== pairs) fail(`${name}: the status matrix sums to ${matrix}, the pairs are ${pairs}`)
      const same = sum(data.status_matrix.filter((c) => c.mac_status === c.vendor_status), (c) => c.claims)
      if (same !== data.totals.same_status) fail(`${name}: the diagonal sums to ${same}, same_status says ${data.totals.same_status}`)
      const names = data.dimensions.map((d) => d.name)
      if (!sameList(names, [...FRONT_FILTERS])) fail(`${name}: the dimensions are ${names.join(', ')}, not the twelve in order`)
      for (const d of data.dimensions) {
        const n = sum(d.values, (v) => v.claims)
        if (n !== total) fail(`${name}: dimension ${d.name} sums to ${n}, not ${total}`)
        const m = sum(d.values, (v) => v.mac_claims)
        if (m !== macTotal) fail(`${name}: dimension ${d.name} holds ${m} MAC claims, not ${macTotal}`)
      }
      const bands = data.dimensions.find((d) => d.name === 'payment_band')
      const bandValues = bands.values.map((v) => v.value).filter((v) => v !== '')
      if (!sameList([...bandValues].sort(), [...BANDS].sort())) fail(`${name}: the payment bands are ${bandValues.join(', ')}, not the five`)
      const bandPairs = sum(bands.values.filter((v) => v.value !== ''), (v) => v.claims)
      if (bandPairs !== pairs) fail(`${name}: the payment bands sum to ${bandPairs} pairs, not ${pairs}`)
    }
    if (!sameList(filters.PAYMENT_BANDS, [...BANDS, ''])) fail(`overviewFilters.ts lists the payment bands ${filters.PAYMENT_BANDS.join(', ')}, not the five and the empty value`)
    // The MAC side is the same for both vendors.
    if (v1.totals.mac_claims !== v2.totals.mac_claims || v1.totals.mac_paid !== v2.totals.mac_paid) fail('the totals’ MAC claims or MAC paid differ between slots 1 and 2')
    const setAside = (data) => data.outcomes.find((o) => o.outcome === 'set_aside')
    if (setAside(v1).mac_claims !== setAside(v2).mac_claims || setAside(v1).mac_paid !== setAside(v2).mac_paid) fail('the set-aside row differs between slots 1 and 2')
    for (const name of LEGACY) {
      const legacyOf = (data) =>
        JSON.stringify(data.dimensions.find((d) => d.name === name).values.filter((v) => v.mac_claims > 0).map((v) => [v.value, v.mac_claims, v.mac_paid]).sort())
      if (legacyOf(v1) !== legacyOf(v2)) fail(`the MAC claims by ${name} differ between slots 1 and 2`)
    }
    // The legacy block: the same in both responses, the nine dimensions in
    // order, each adding up to the totals, no empty MAC status (only extra
    // claims have none), and at least the MAC claims the vendor rows count
    // (it counts claims; a vendor counts keys).
    if (!sameList(v1.legacy, v2.legacy)) fail('the legacy block differs between slots 1 and 2')
    if (!sameList(legacy.dimensions.map((d) => d.name), [...LEGACY])) fail(`the legacy dimensions are ${legacy.dimensions.map((d) => d.name).join(', ')}, not the nine in order`)
    for (const d of legacy.dimensions) {
      const n = sum(d.values, (v) => v.claims)
      if (n !== legacy.totals.claims) fail(`legacy ${d.name} sums to ${n} claims, the legacy totals say ${legacy.totals.claims}`)
      const paidCents = sum(d.values, (v) => cents(v.paid))
      if (paidCents !== cents(legacy.totals.paid)) fail(`legacy ${d.name} adds to ${paidCents / 100} paid, the legacy totals say ${legacy.totals.paid}`)
      const billedCents = sum(d.values, (v) => cents(v.billed))
      if (billedCents !== cents(legacy.totals.billed)) fail(`legacy ${d.name} adds to ${billedCents / 100} billed, the legacy totals say ${legacy.totals.billed}`)
    }
    if (legacy.dimensions.find((d) => d.name === 'mac_status').values.some((v) => v.value === '')) fail('the legacy block lists an empty MAC status')
    if (legacy.totals.claims < v1.totals.mac_claims) fail(`the legacy totals count ${legacy.totals.claims} claims, fewer than the ${v1.totals.mac_claims} MAC claims the vendor rows hold`)
    // The agreement grid counts each MAC claim not set aside once. Its cells
    // name slot 1's outcome v1 and slot 2's v2 (/api/agreement).
    if (!agreement.v1 || !agreement.v2 || 'a' in agreement || 'b' in agreement) fail('the agreement response does not name its submissions v1 and v2')
    for (const [key, name, data] of SLOTS) {
      for (const outcome of ['equal', 'different', 'missing', 'duplicated']) {
        const n = sum(agreement.cells.filter((c) => c[key] === outcome), (c) => c.claims)
        const want = data.outcomes.find((o) => o.outcome === outcome)?.claims ?? 0
        if (n !== want) fail(`the agreement grid holds ${n} claims ${outcome} in ${name}, its outcomes say ${want}`)
      }
    }

    // ---- the name
    const render = (element, location) =>
      renderToStaticMarkup(createElement(StaticRouter, { location }, createElement(SelectionProvider, null, element)))
    if (HELP.overview.title !== 'Overview') fail(`the help is titled “${HELP.overview.title}”, not Overview`)
    if (!/\{ path: '\/', label: 'Overview' \}/.test(await source('src/app/Shell.tsx'))) fail('the navigation does not name the page at / Overview')
    if (!render(createElement(Overview), '/').includes('>Overview</h1>')) fail('the page is not titled Overview')

    // ---- the view from the fixtures
    const sides = [
      { side: 'v1', submission: v1.submission, data: v1, loading: false, error: undefined },
      { side: 'v2', submission: v2.submission, data: v2, loading: false, error: undefined },
    ]
    const grid = { cells: agreement.cells, error: undefined, loading: false }
    const view = (props) => render(createElement(OverviewView, { sides, selection: {}, onSelection: () => {}, agreement: grid, macs, ...props }), '/')
    const html = view({})
    const legacyFacet = (name) => legacy.dimensions.find((d) => d.name === name)?.values ?? []
    const valuesOf = (name) =>
      [...new Set([
        ...[v1, v2].flatMap((data) => data.dimensions.find((d) => d.name === name).values.map((v) => v.value)),
        ...(LEGACY.includes(name) ? legacyFacet(name).map((v) => v.value) : []),
      ])].sort((x, y) => filters.compareValues(name, x, y))
    const facet = (data, name, value) => data.dimensions.find((d) => d.name === name).values.find((v) => v.value === value)

    // ---- the selection header
    // Sticky at the top, at most 40% of the window with its own scrolling,
    // in the page's flow, directly in the view's root (which holds the whole
    // page), between the column names and the first table.
    const headerAt = html.indexOf('<section aria-label="Selection" data-tab-row=""')
    const headerClass = (/^<section aria-label="Selection" data-tab-row="" class="([^"]*)"/.exec(html.slice(Math.max(0, headerAt)))?.[1] ?? '').split(' ')
    for (const cls of ['sticky', 'top-0', 'max-h-[40svh]', 'overflow-y-auto', 'bg-bg']) if (!headerClass.includes(cls)) fail(`the selection header is not ${cls}`)
    if (headerClass.some((c) => c === 'fixed' || c === 'absolute')) fail('the selection header is taken out of the page’s flow')
    if (headerAt < 0 || !html.startsWith('<div class="space-y-10">') || depthAt(html, headerAt) !== 1) fail('the selection header is not directly in the view’s root, so it would not stick for the whole page')
    if (!(html.indexOf('Vendor 2') < headerAt && headerAt < html.indexOf('data-table="outcomes"'))) fail('the selection header is not between the column names and the first table')

    // The nine boxes in order; in each, All and None at the top, then the
    // values as tabs whose text is the tab label alone, no count.
    const headerOf = (markup) => region(markup, 'data-tab-row=""', ['data-table="outcomes"'])
    const boxesShown = [...headerOf(html).matchAll(/data-tab-group="(\w+)"/g)].map((m) => m[1])
    if (!sameList(boxesShown, BOXES)) fail(`the header offers the boxes ${boxesShown.join(', ')}, not the nine in order`)
    if (html.includes('type="checkbox"')) fail('a checkbox is still drawn')
    const buttonsOf = (markup) =>
      [...markup.matchAll(/<button([^>]*)>([^<]*)<\/button>/g)]
        .map((m) => ({
          tab: /data-tab="([^"]*)"/.exec(m[1])?.[1],
          className: /class="([^"]*)"/.exec(m[1])?.[1] ?? '',
          text: m[2],
          pressed: /aria-pressed="true"/.test(m[1]),
          zero: /data-zero=""/.test(m[1]),
        }))
        .filter((x) => x.tab !== undefined)
    const boxOf = (markup, name) => region(headerOf(markup), `data-tab-group="${name}"`, ['data-tab-group="', 'data-address-only'])
    const facetOf = (data, name) => data.dimensions.find((d) => d.name === name).values
    // The values a box lists: those with claims in either slot's facet, in
    // the box's order (CMS's for the jurisdictions), and every picked one.
    const rank = (v) => (filters.JURISDICTIONS.includes(v) ? filters.JURISDICTIONS.indexOf(v) : 99 + (v === '' ? 1 : 0))
    const listedValues = (datas, name, picked = []) => {
      const all = [...new Set([...datas.flatMap((d) => facetOf(d, name).filter((v) => v.claims > 0).map((v) => v.value)), ...picked])]
      return all.sort(name === 'jurisdiction' ? (x, y) => rank(x) - rank(y) || x.localeCompare(y) : (x, y) => filters.compareValues(name, x, y))
    }
    const hasClaims = (datas, name, value) => datas.some((d) => facetOf(d, name).some((v) => v.value === value && v.claims > 0))
    // Each kind of tab in one set of classes across the header (monospace aside).
    const tones = { plain: new Set(), picked: new Set(), dimmed: new Set(), end: new Set() }
    const classKey = (cls) => cls.split(' ').filter((c) => c !== 'font-mono').join(' ')
    let tabs = 0
    let boxChecks = 0
    const checkBox = (markup, name, datas, selection, what) => {
      boxChecks++
      const box = boxOf(markup, name)
      const buttons = buttonsOf(box)
      const picked = selection[name] ?? []
      const values = listedValues(datas, name, picked)
      const want = [
        { tab: 'all', text: 'All', pressed: selection[name] === undefined },
        { tab: 'none', text: 'None', pressed: selection[name] !== undefined && picked.length === 0 },
        ...values.map((v) => ({ tab: escapeHtml(`${name}:${v}`), text: escapeHtml(filters.tabLabel(name, v)), pressed: picked.includes(v) })),
      ]
      if (!sameList(buttons.map((x) => [x.tab, x.text, x.pressed]), want.map((x) => [x.tab, x.text, x.pressed]))) {
        fail(`${what}: the ${name} box reads ${JSON.stringify(buttons.map((x) => [x.text, x.pressed]))}, not ${JSON.stringify(want.map((x) => [x.text, x.pressed]))}`)
      }
      const grid = box.indexOf('class="grid')
      if (values.length > 0 && (grid < 0 || box.indexOf('data-tab="none"') > grid)) fail(`${what}: All and None are not at the top of the ${name} box`)
      const columns = Number(/data-columns="(\d+)"/.exec(box)?.[1])
      if (values.length > 0 && columns !== Math.min(COLUMNS[name], values.length)) fail(`${what}: the ${name} box sets ${columns} tabs to a row, not ${Math.min(COLUMNS[name], values.length)}`)
      for (const x of buttons) {
        if (x.tab === 'all' || x.tab === 'none') {
          tones.end.add(`${x.pressed}:${x.className}`)
          continue
        }
        const value = values.find((v) => escapeHtml(`${name}:${v}`) === x.tab)
        const empty = x.pressed && !hasClaims(datas, name, value)
        if (x.zero !== empty) fail(`${what}: the ${name} tab “${x.text}” is ${x.zero ? '' : 'not '}dimmed`)
        tones[empty ? 'dimmed' : x.pressed ? 'picked' : 'plain'].add(classKey(x.className))
      }
      return buttons.length - 2
    }
    const both = [v1, v2]
    for (const name of BOXES.filter((b) => b !== 'workload')) tabs += checkBox(html, name, both, {}, 'no selection')
    // The workload box starts closed, one line with how many it lists.
    const workloadBox = boxOf(html, 'workload')
    if (!workloadBox.includes('data-open="0"') || !workloadBox.includes('aria-expanded="false"')) fail('the workload box does not start closed')
    const countOf = (markup) => /data-workload-count="">([^<]*)</.exec(boxOf(markup, 'workload'))?.[1]
    const allWorkloads = listedValues(both, 'workload')
    if (countOf(html) !== `all ${allWorkloads.length}`) fail(`the closed workload box reads “Workload: ${countOf(html)}”, not “Workload: all ${allWorkloads.length}”`)
    if (buttonsOf(workloadBox).length !== 0) fail('the closed workload box shows tabs')
    if (countOf(view({ selection: { workload: ['01111', '02202'] } })) !== '2 selected') fail('the closed workload box does not count the picked workloads')
    if (countOf(view({ selection: { workload: [] } })) !== 'none') fail('the closed workload box does not say none under None')

    // The jurisdictions in CMS's order, four to a row: the order matches
    // /api/macs, and a facet holding all sixteen (sent in value order, as
    // the server sends it) draws them in that order in four rows of four.
    if (!sameList(filters.JURISDICTIONS, macs.jurisdictions.map((j) => j.jurisdiction))) fail(`overviewFilters.ts JURISDICTIONS is not /api/macs's order: ${macs.jurisdictions.map((j) => j.jurisdiction).join(', ')}`)
    const withFacet = (data, name, values) => ({ ...data, dimensions: data.dimensions.map((d) => (d.name === name ? { ...d, values } : d)) })
    const sixteen = [...filters.JURISDICTIONS].sort().map((value) => ({ ...facetOf(v1, 'jurisdiction')[0], value, claims: 1, mac_claims: 1 }))
    const fullData = both.map((d) => withFacet(d, 'jurisdiction', sixteen))
    const fullHtml = render(createElement(OverviewView, { sides: sides.map((s, i) => ({ ...s, data: fullData[i] })), selection: {}, onSelection: () => {}, macs }), '/')
    tabs += checkBox(fullHtml, 'jurisdiction', fullData, {}, 'all sixteen jurisdictions')
    const drawn = buttonsOf(boxOf(fullHtml, 'jurisdiction')).slice(2).map((x) => x.text)
    if (!sameList(drawn, filters.JURISDICTIONS) || !boxOf(fullHtml, 'jurisdiction').includes('data-columns="4"')) fail(`the sixteen jurisdictions read ${drawn.join(' ')}, not CMS's order four to a row`)

    // The hidden-zero rule. A value with no claims under the rest of the
    // selection is not listed unless it is picked; a picked one stays,
    // pressed and dimmed. J5 has no claims in the fixtures.
    tabs += checkBox(view({ selection: { jurisdiction: ['J5'] } }), 'jurisdiction', both, { jurisdiction: ['J5'] }, 'J5 picked')
    // Picking JE: /api/front narrows every other facet to JE's claims; here
    // the workload facets keep only JE's numbers, as the server sends them.
    const jeData = both.map((d) => withFacet(d, 'workload', facetOf(d, 'workload').filter((v) => v.value.startsWith('01'))))
    const jeSelection = { jurisdiction: ['JE'] }
    const jeHtml = render(createElement(OverviewView, { sides: sides.map((s, i) => ({ ...s, data: jeData[i] })), selection: jeSelection, onSelection: () => {}, macs }), '/')
    tabs += checkBox(jeHtml, 'jurisdiction', jeData, jeSelection, 'JE picked')
    const jeWorkloads = listedValues(jeData, 'workload')
    if (countOf(jeHtml) !== `all ${jeWorkloads.length}` || jeWorkloads.some((w) => !w.startsWith('01'))) fail(`with JE picked the closed workload box reads “${countOf(jeHtml)}”, not all ${jeWorkloads.length} of JE's`)
    // Under None on the jurisdictions every other facet is empty.
    const emptyData = both.map((d) => withFacet(d, 'workload', []))
    const emptyHtml = render(createElement(OverviewView, { sides: sides.map((s, i) => ({ ...s, data: emptyData[i] })), selection: { jurisdiction: [] }, onSelection: () => {}, macs }), '/')
    if (countOf(emptyHtml) !== 'none with claims') fail(`with no workload under the selection the closed workload box reads “${countOf(emptyHtml)}”`)
    // Picking a claim type: the status and frequency facets keep only the
    // values that claim type carries; a status picked earlier that it does
    // not carry stays, dimmed.
    const dmeData = both.map((d) => withFacet(withFacet(d, 'mac_status', facetOf(d, 'mac_status').filter((v) => v.value === '1' || v.value === '4')), 'frequency', facetOf(d, 'frequency').filter((v) => v.value === '1')))
    const dmeSelection = { claim_type: ['dme'], mac_status: ['22'] }
    const dmeHtml = render(createElement(OverviewView, { sides: sides.map((s, i) => ({ ...s, data: dmeData[i] })), selection: dmeSelection, onSelection: () => {}, macs }), '/')
    for (const name of ['claim_type', 'mac_status', 'frequency']) tabs += checkBox(dmeHtml, name, dmeData, dmeSelection, 'DME and status 22 picked')
    const statusTabs = buttonsOf(boxOf(dmeHtml, 'mac_status')).slice(2).map((x) => `${x.text}${x.zero ? ' (dimmed)' : ''}`)
    if (!sameList(statusTabs, ['1', '4', '22 (dimmed)'])) fail(`with DME picked the MAC status box reads ${statusTabs.join(', ')}, not 1, 4, and 22 dimmed`)
    for (const [kind, set] of Object.entries(tones)) {
      if (kind === 'end' ? set.size !== 2 : set.size !== 1) fail(`the header draws ${kind} tabs in ${set.size} sets of classes: ${[...set].join(' | ')}`)
    }

    // The workload box opened: every workload once, under its jurisdiction
    // in CMS's order, a roll-up as a group headed by its number, the numbers
    // the MAC map does not hold last; chips smaller than the boxes' tabs.
    const openBox = (values, selection) =>
      render(createElement(WorkloadBox, { values, selection, onSelection: () => {}, macs, open: true, onOpen: () => {} }), '/')
    const workloadValues = filters.boxValues(both.map((d) => facetOf(d, 'workload')), 'workload', {})
    const opened = openBox(workloadValues, {})
    const chips = buttonsOf(opened).filter((x) => x.tab.startsWith('workload:'))
    if (!opened.includes('data-open="1"') || !opened.includes('aria-expanded="true"') || !opened.includes('data-tab="all"') || !opened.includes('data-tab="none"')) fail('the opened workload box has no header with All and None')
    if (!sameList(chips.map((x) => x.tab.slice('workload:'.length)).sort(), [...allWorkloads].sort())) fail(`the opened workload box lists ${chips.map((x) => x.text).join(' ')}, not each of ${allWorkloads.join(' ')} once`)
    const placedIn = [...opened.matchAll(/data-workload-jurisdiction="([^"]*)"/g)].map((m) => m[1])
    if (!sameList(placedIn, ['JE', 'JF', ''])) fail(`the opened workload box groups by ${JSON.stringify(placedIn)}, not JE, JF, then the unplaced`)
    const rollUp = region(opened, 'data-rollup="01011"', ['data-tab="workload:01112"'])
    if (!rollUp.includes('>01011<') || !rollUp.includes('data-tab="workload:01111"') || !rollUp.includes('viewBox="0 0 16 16"')) fail('roll-up 01011 does not head 01111 with its folder and number')
    const unplaced = region(opened, 'data-workload-jurisdiction=""', ['</div></div>'])
    if (!unplaced.includes('>other<') || !unplaced.includes('data-tab="workload:01412"') || !unplaced.includes('data-tab="workload:01912"')) fail('the numbers the MAC map does not hold are not listed last, under other')
    const tabSize = /text-\[(\d+)px\]/.exec([...tones.plain][0] ?? '')?.[1]
    if (chips.some((x) => /text-\[(\d+)px\]/.exec(x.className)?.[1] === tabSize || !/text-\[10px\]/.test(x.className))) fail('the workload chips are not smaller than the other boxes’ tabs')
    const jeOpened = openBox(filters.boxValues(jeData.map((d) => facetOf(d, 'workload')), 'workload', {}), {})
    if (jeOpened.includes('data-workload-jurisdiction="JF"') || /data-tab="workload:0[23]/.test(jeOpened)) fail('with JE picked the opened workload box still lists JF’s workloads')
    const strayPick = openBox(filters.boxValues(jeData.map((d) => facetOf(d, 'workload')), 'workload', { workload: ['02202'] }), { workload: ['02202'] })
    if (!/data-tab="workload:02202" aria-pressed="true" data-zero=""/.test(strayPick) || !strayPick.includes('data-workload-jurisdiction="JF"')) fail('a picked workload of another jurisdiction is not kept, dimmed, under its jurisdiction')
    // The two tiers cover the whole MAC map: every number once, sixteen
    // jurisdictions in CMS's order, nothing unplaced.
    const everyNumber = [...new Set(macs.jurisdictions.flatMap((j) => j.parts.flatMap((p) => [...p.roll_ups.flatMap((r) => [r.id, ...r.workloads.map((w) => w.id)]), ...p.workloads, ...p.hhh, ...p.other].map((n) => (typeof n === 'string' ? n : n.id)))))]
    const tiers = workloadGroups(macs, everyNumber)
    const placedCount = sum(tiers.groups, (g) => sum(g.items, (item) => (item.kind === 'rollup' ? (item.own ? 1 : 0) + item.workloads.length : 1)))
    if (tiers.unplaced.length !== 0 || placedCount !== everyNumber.length || !sameList(tiers.groups.map((g) => g.jurisdiction), filters.JURISDICTIONS)) {
      fail(`the workload box's tiers place ${placedCount} of the MAC map's ${everyNumber.length} numbers (${tiers.unplaced.length} unplaced) in ${tiers.groups.length} jurisdictions`)
    }

    // Every table sits on the one grid, four columns, its heads Legacy, 1, 2
    // in that order, with the Legacy column empty where the legacy side has
    // no figure of its own.
    const templates = new Map([...html.matchAll(/style="grid-template-columns:([^"]*)" data-(?:table|chart|grid)="([\w-]+)"/g)].map((m) => [m[2], m[1]]))
    const tablesOf = (markup) => [
      ...[...markup.matchAll(/data-(table|chart)="([\w-]+)"/g)].map((m) => ({ key: `${m[1]}:${m[2]}`, start: m[0], name: m[2] })),
      ...[...markup.matchAll(/data-breakdown="(\w+)"/g)].map((m) => ({ key: `breakdown:${m[1]}`, start: m[0], name: m[1], breakdown: true })),
    ]
    const withLegacy = (t) => (t.breakdown ? LEGACY.includes(t.name) : t.name === 'outcomes' || t.name === 'status')
    const headsOf = (markup, t) => {
      const at = markup.indexOf(t.start)
      return columnMarks(markup.slice(at)).slice(0, 3)
    }
    const tables = tablesOf(html)
    for (const t of tables) {
      const want = [withLegacy(t) ? 'legacy' : 'gap:legacy', 'v1', 'v2']
      const heads = headsOf(html, t)
      if (!sameList(heads, want)) fail(`${t.key} heads read ${heads.join(' ')}, not ${want.join(' ')}`)
      const template = t.breakdown
        ? /style="grid-template-columns:([^"]*)"/.exec(region(html, t.start, ['data-breakdown="']))?.[1]
        : templates.get(t.name)
      const wantTemplate = t.name === 'status' ? MATRIX_GRID.gridTemplateColumns : GRID.gridTemplateColumns
      if (template !== wantTemplate) fail(`${t.key} lays out its columns as “${template}”, not the page grid “${wantTemplate}”`)
      if (template && trackCount(template) !== 4) fail(`${t.key} has ${trackCount(template)} columns, not 4`)
    }
    if (trackCount(GRID.gridTemplateColumns) !== 4 || trackCount(MATRIX_GRID.gridTemplateColumns) !== 4) fail('the page grid does not have four columns')
    for (const marker of ['data-table="outcomes"', 'data-table="payment"', 'data-chart="payment-bands"', 'data-table="status"', 'data-table="differs"']) {
      if (!html.includes(marker)) fail(`the page has no ${marker}`)
    }
    const legacyTables = tables.filter(withLegacy).length
    // The agreement grid sits on the page grid: the slot 1 outcomes in the
    // label column, the slot 2 outcomes across the three figure columns.
    if (templates.get('agreement') !== GRID.gridTemplateColumns) fail(`the agreement grid lays out its columns as “${templates.get('agreement')}”, not the page grid`)
    const across = count(region(html, 'data-grid="agreement"', ['data-table="payment"']), 'grid-column:2 / span 3')
    if (across !== 5) fail(`the agreement grid spans the figure columns in ${across} rows, not 5 (its head and four rows)`)

    // The twelve breakdown tables, in the response's order, each with every
    // value and slot 1's and slot 2's counts; the nine by what the MAC claim
    // carries with a legacy figure per row adding up to the legacy totals,
    // the three by what a vendor decides with an empty Legacy column.
    const breakdowns = [...html.matchAll(/data-breakdown="(\w+)"/g)].map((m) => m[1])
    if (!sameList(breakdowns, [...FRONT_FILTERS])) fail(`the breakdown tables are ${breakdowns.join(', ')}, not the twelve in order`)
    let rows = 0
    for (const name of FRONT_FILTERS) {
      const table = region(html, `data-breakdown="${name}"`, ['data-breakdown="'])
      const entries = legacyEntries(table)
      const values = valuesOf(name)
      if (LEGACY.includes(name)) {
        if (entries.length !== values.length) fail(`By ${name} has ${entries.length} legacy figures for ${values.length} rows`)
        const n = legacySum(entries)
        if (n !== legacy.totals.claims) fail(`By ${name}: the legacy column adds to ${n} MAC claims, the legacy totals say ${legacy.totals.claims}`)
      } else {
        if (entries.length > 0) fail(`By ${name} has legacy figures, but a vendor decides its rows`)
        if (count(table, 'data-gap="legacy"') !== values.length + 1) fail(`By ${name} leaves the Legacy column empty in ${count(table, 'data-gap="legacy"')} places, not its head and ${values.length} rows`)
      }
      for (const value of values) {
        const at = table.indexOf(`data-row="${escapeHtml(value)}"`)
        if (at < 0) {
          fail(`By ${name} has no row for “${value}”`)
          continue
        }
        rows += 1
        const row = region(table.slice(at), `data-row="${escapeHtml(value)}"`, ['data-row="', 'aria-expanded'])
        if (!row.includes(`title="${escapeHtml(filters.valueLabel(name, value))}"`)) fail(`By ${name} does not label “${value}” as “${filters.valueLabel(name, value)}”`)
        const slots = countSlots(row)
        const want = [v1, v2].map((data) => facet(data, name, value)).map((v) => (v ? countText(v.claims) : '—'))
        if (!sameList(slots.slice(-2), want)) fail(`By ${name}, “${value}”: slots 1 and 2 read ${slots.slice(-2).join(' and ')}, not ${want.join(' and ')}`)
        if (LEGACY.includes(name)) {
          const l = legacyFacet(name).find((v) => v.value === value)
          const wantLegacy = l ? countText(l.claims) : '—'
          if (slots.length !== 3 || slots[0] !== wantLegacy) fail(`By ${name}, “${value}”: Legacy reads ${slots.length === 3 ? slots[0] : 'nothing'}, not ${wantLegacy}`)
        } else if (slots.length !== 2) fail(`By ${name}, “${value}”: ${slots.length} figures, not slot 1’s and slot 2’s alone`)
      }
    }
    const workload = region(html, 'data-breakdown="workload"', ['data-breakdown="'])
    if (count(workload, 'class="hidden" data-row=') !== 4 || !workload.includes('Show all 10')) fail('By workload does not start closed on its first 6 of 10 rows')

    // The outcome table: one row per outcome with no legacy figure, and the
    // total with the legacy MAC claims and MAC paid.
    const outcomeTable = region(html, 'data-table="outcomes"', ['data-section="agreement"', 'data-table="payment"'])
    const outcomeRows = outcomeTable.slice(0, outcomeTable.indexOf('>Total<'))
    const outcomeTotal = outcomeTable.slice(outcomeTable.indexOf('>Total<'))
    for (const [, slotName, data] of SLOTS) {
      for (const o of data.outcomes) if (!outcomeTable.includes(`>${countText(o.claims)}<`)) fail(`the outcome table does not show ${slotName}'s ${o.outcome} = ${o.claims}`)
      if (!outcomeTable.includes(`>${countText(data.totals.claims)}<`)) fail(`the outcome table does not show ${slotName}'s total ${data.totals.claims}`)
    }
    if (legacyEntries(outcomeRows).length !== 0) fail('the outcome table has a legacy figure on an outcome’s row')
    if (count(outcomeRows, 'data-gap="legacy"') !== filters.OUTCOMES.length) fail(`the outcome table leaves Legacy empty on ${count(outcomeRows, 'data-gap="legacy"')} rows, not the ${filters.OUTCOMES.length} outcomes`)
    if (!sameList(legacyEntries(outcomeTotal).map((e) => e.claims), [legacy.totals.claims])) fail(`the outcome table's legacy total is not ${legacy.totals.claims}, shown once`)
    if (!outcomeTotal.includes(escapeHtml(money(legacy.totals.paid)))) fail(`the outcome table's legacy total does not show the MAC paid ${legacy.totals.paid}`)
    if (!html.includes('data-bar="legacy"') || !html.includes(`${countText(legacy.totals.claims)} MAC claims`) || !html.includes('data-bar="v1"') || !html.includes('data-bar="v2"')) fail('the outcome bars are not drawn for Legacy, 1, and 2')

    // Payment agreement and its bands have no legacy figure; the same payment
    // is the match band.
    const payment = region(html, 'data-table="payment"', ['data-chart="payment-bands"'])
    if (legacyEntries(payment).length !== 0) fail('the payment table has a legacy figure')
    const samePayment = region(payment, 'Same payment (CLP04)', ['Different payment'])
    for (const [, slotName, data] of SLOTS) {
      const match = countText(facet(data, 'payment_band', 'match')?.claims ?? 0)
      if (!sameList(countSlots(samePayment), [countText(facet(v1, 'payment_band', 'match')?.claims ?? 0), countText(facet(v2, 'payment_band', 'match')?.claims ?? 0)])) {
        fail(`the same-payment line does not read the match band, ${match} for ${slotName}`)
      }
    }
    const chart = region(html, 'data-chart="payment-bands"', ['data-table="status"'])
    if (legacyEntries(chart).length !== 0) fail('the bands chart has a legacy figure')
    const charted = [...chart.matchAll(/data-band="v1:([^"]*)"/g)].map((m) => m[1].replace(/&amp;/g, '&'))
    if (!sameList(charted, BANDS)) fail(`the bands chart draws ${charted.join(', ')}, not the five bands in order`)

    // The rest of the page.
    for (const data of [v1, v2]) {
      if (!html.includes(escapeHtml(money(data.totals.gross_difference)))) fail(`the payment tile does not show the gross difference ${data.totals.gross_difference}`)
      for (const t of data.tags) {
        if (!html.includes(escapeHtml(t.tag))) fail(`What differs does not list ${t.tag}`)
        if (!html.includes(escapeHtml(t.label))) fail(`What differs does not name ${t.tag} (“${t.label}”)`)
      }
    }
    if (legacyEntries(region(html, 'data-table="differs"', ['data-section="breakdowns"'])).length !== 0) fail('What differs has a legacy figure')
    for (const [slot, slotName, data] of SLOTS) {
      for (const band of BANDS) {
        const bar = new RegExp(`data-band="${slot}:${band.replace(/[.%]/g, (c) => `\\${c}`)}"><span[^>]*></span><span[^>]*>([^<]*)<`).exec(html)
        const want = countText(facet(data, 'payment_band', band)?.claims ?? 0)
        if (bar?.[1] !== want) fail(`the bands chart reads ${bar?.[1]} for ${slotName} “${band}”, not ${want}`)
      }
    }
    // Both matrices: diagonal cells green, the others amber; the legacy
    // count beside each MAC status.
    const diagonal = sum([v1, v2], (data) => data.status_matrix.filter((c) => c.mac_status === c.vendor_status).length)
    const off = sum([v1, v2], (data) => data.status_matrix.length) - diagonal
    if (count(html, 'data-status-cell="same"') !== diagonal || count(html, 'bg-ok/15 text-ink') !== diagonal) fail(`the status matrices mark ${count(html, 'data-status-cell="same"')} diagonal cells, not ${diagonal}`)
    if (count(html, 'data-status-cell="other"') !== off || count(html, 'bg-warn/20 text-ink') !== off) fail(`the status matrices mark ${count(html, 'data-status-cell="other"')} off-diagonal cells, not ${off}`)
    const macCodes = [...new Set([v1, v2].flatMap((data) => data.status_matrix.map((c) => c.mac_status)))]
    for (const m of macCodes) {
      const cell = new RegExp(`data-legacy-status="${m}"><span class="block" data-legacy-claims="(\\d+)"`).exec(html)
      const want = legacyFacet('mac_status').find((v) => v.value === m)?.claims ?? 0
      if (Number(cell?.[1]) !== want) fail(`the status table shows ${cell?.[1]} legacy claims beside MAC status ${m}, not ${want}`)
    }
    if (count(html, 'data-agreement="') !== agreement.cells.length) fail(`the agreement grid draws ${count(html, 'data-agreement="')} cells, not ${agreement.cells.length}`)
    if (!html.includes(`${countText(sum(agreement.cells, (c) => c.claims))} MAC claims, set-aside claims left out`)) fail('the agreement grid does not print its total')
    if (html.indexOf('Vendor 1') > html.indexOf('Vendor 2')) fail('slot 2 is named before slot 1')
    for (const title of ['Every claim, by outcome', 'How the two agree, per MAC claim', 'Payment agreement among pairs', 'Status agreement among pairs', 'What differs', 'Breakdowns']) {
      if (!html.includes(title)) fail(`the page has no “${title}”`)
    }
    if (!html.includes('>Reset selection<')) fail('the selection header has no Reset selection')
    // Slot 2's over-5% pairs add to a negative paid amount; it keeps its sign.
    if (!region(html, 'data-breakdown="payment_band"', ['data-breakdown="']).includes('>−$503.00<')) fail('a negative paid sum is shown without its sign')

    // ---- a narrowed selection is drawn as such
    const narrowedHtml = view({ selection: { payment_band: ['match', 'up to 1.00'] } })
    const narrowedChart = region(narrowedHtml, 'data-chart="payment-bands"', ['data-table="status"'])
    if (count(narrowedChart, 'data-left-out="payment_band"') !== 3) fail(`${count(narrowedChart, 'data-left-out="payment_band"')} bands are struck through in the chart, not 3`)
    const bandTable = region(narrowedHtml, 'data-breakdown="payment_band"', ['data-breakdown="'])
    if (count(bandTable, 'data-left-out="payment_band"') !== 4) fail(`${count(bandTable, 'data-left-out="payment_band"')} rows are struck through in By payment difference band, not 4`)
    if (!narrowedHtml.includes('data-address-only') || !narrowedHtml.includes('Payment difference band: match, up to $1.00 apart')) fail('the address-only narrowing is not named')
    if (/<button[^>]*disabled=""[^>]*>Reset selection/.test(narrowedHtml)) fail('Reset selection is disabled on a narrowed page')
    // A payment band is a vendor's result: the legacy figures stay as they are.
    const byJurisdiction = legacyEntries(region(narrowedHtml, 'data-breakdown="jurisdiction"', ['data-breakdown="']))
    if (byJurisdiction.some((e) => !e.in) || legacySum(byJurisdiction) !== legacy.totals.claims) fail('a payment band selection changes the legacy column of By jurisdiction')
    if (!narrowedHtml.includes('data-bar="legacy"') || !narrowedHtml.includes(`${countText(legacy.totals.claims)} MAC claims`)) fail('the legacy bar is not drawn under a payment band selection')
    // A legacy dimension left out is struck through in its own table, legacy figure included.
    const jfLeftOut = region(view({ selection: { jurisdiction: ['JE'] } }), 'data-breakdown="jurisdiction"', ['data-breakdown="'])
    if (!sameList(legacyEntries(jfLeftOut).map((e) => e.in), [true, false])) fail('with only JE kept, By jurisdiction does not mark JF’s legacy figure left out')

    // ---- one empty slot, an error, and the updating mark
    const only1 = render(createElement(OverviewView, { sides: [sides[0]], selection: {}, onSelection: () => {}, agreement: grid }), '/')
    for (const t of tablesOf(only1)) {
      const want = [withLegacy(t) ? 'legacy' : 'gap:legacy', 'v1', 'gap:v2']
      const heads = headsOf(only1, t)
      if (!sameList(heads, want)) fail(`with slot 2 empty, ${t.key} heads read ${heads.join(' ')}, not ${want.join(' ')}`)
    }
    if (only1.includes('data-section="agreement"')) fail('with slot 2 empty, the agreement grid is still drawn')
    const message = '500 Internal Server Error: unable to read the database (request 1790345126747061000)'
    const failing = render(createElement(OverviewView, { sides: [{ ...sides[0], loading: true, error: new Error(message) }, sides[1]], selection: {}, onSelection: () => {} }), '/')
    if (!failing.includes(`Submission 1: ${escapeHtml(message)}`)) fail('an error is not shown with its request ID')
    if (!failing.includes('updating…') || !view({ pending: true }).includes('updating…') || html.includes('updating…')) fail('the updating mark does not follow loading and pending')

    // ---- the selection round-trips through the address
    const sample = { jurisdiction: ['JE'], part: ['A', ''], payment_band: [], outcome: ['different', 'equal'] }
    const params = new URLSearchParams('v1=x&v2=y')
    filters.writeSelection(params, sample)
    if (params.get('v1') !== 'x' || params.get('v2') !== 'y') fail('writing the selection disturbed the submission parameters')
    if (params.get('payment_band') !== filters.NONE) fail(`an empty list is written as “${params.get('payment_band')}”, not “${filters.NONE}”`)
    const back = filters.readSelection(params)
    if (!sameList(back, sample)) fail(`the selection does not round-trip: ${JSON.stringify(back)}`)
    const request = filters.filtersOf(sample)
    if (request.part !== 'A,' || request.outcome !== 'different,equal' || request.payment_band !== filters.NONE || request.workload !== undefined) {
      fail(`the request filters are wrong: ${JSON.stringify(request)}`)
    }
    // A click picks, a second unpicks, unpicking the last returns to All;
    // from None a click picks that value alone.
    const one = filters.toggleValue({}, 'jurisdiction', 'JF')
    if (!sameList(one.jurisdiction, ['JF'])) fail(`picking JF with no list gave ${JSON.stringify(one)}`)
    const two = filters.toggleValue(one, 'jurisdiction', 'JE')
    if (!sameList(two.jurisdiction, ['JF', 'JE'])) fail(`picking JE as well gave ${JSON.stringify(two)}`)
    const back1 = filters.toggleValue(two, 'jurisdiction', 'JF')
    if (!sameList(back1.jurisdiction, ['JE'])) fail(`unpicking JF gave ${JSON.stringify(back1)}`)
    const toAll = filters.toggleValue(back1, 'jurisdiction', 'JE')
    if ('jurisdiction' in toAll) fail(`unpicking the last value did not return the box to All: ${JSON.stringify(toAll)}`)
    const fromNone = filters.toggleValue(filters.selectNone({}, 'jurisdiction'), 'jurisdiction', 'JE')
    if (!sameList(fromNone.jurisdiction, ['JE'])) fail(`picking JE after None gave ${JSON.stringify(fromNone)}`)
    if (!filters.isPicked(two, 'jurisdiction', 'JE') || filters.isPicked({}, 'jurisdiction', 'JE')) fail('isPicked does not follow the list')
    // Picking a workload writes workload= as /api/front takes it.
    const workloadParams = new URLSearchParams('v1=x')
    filters.writeSelection(workloadParams, filters.toggleValue(filters.toggleValue({}, 'workload', '01111'), 'workload', '02202'))
    if (workloadParams.get('workload') !== '01111,02202' || filters.filtersOf(filters.readSelection(workloadParams)).workload !== '01111,02202') fail(`picking two workloads writes workload=${workloadParams.get('workload')}`)
    const none = filters.selectNone({}, 'part')
    if (!sameList(none.part, [])) fail('None did not empty the dimension')
    if (filters.isSelected(none, 'part', 'A')) fail('a value is selected under None')
    if (!filters.isSelected({}, 'part', 'A')) fail('a value is not selected with no list')
    const key1 = filters.selectionKey(sample)
    const key2 = filters.selectionKey(filters.readSelection(params))
    if (key1 !== key2) fail('the selection key is not stable across a round trip')

    // ---- the claim screen from the fixture claims
    const claims = {}
    for (const col of ['mac', 'v1', 'v2']) claims[col] = parseClaim(claim[col], claim.delimiters)
    if (claims.mac.head.id !== 'CLP' || claims.mac.lines.length !== 1 || claims.mac.lines[0].line_control_number !== 'L000301') fail('the MAC claim did not read as one claim with one line')
    const head = { mac: 'MAC-processed', v1: 'Vendor 1', v2: 'Vendor 2' }
    const claimHtml = render(createElement(ClaimCompareView, { claims, head }), '/claim?match_key=abc')
    for (const text of ['CLP04', 'SVC03', 'Only differences', 'Preset', 'Elements shown', 'MAC-processed', 'Vendor 1', 'Vendor 2']) {
      if (!claimHtml.includes(text)) fail(`the claim screen does not show “${text}”`)
    }
    // Slot 1 pays $79.98 on the claim and the line with $20.02 of CO-45
    // where the MAC pays $80.00 with $20.00: two cents apart, within the
    // comparison's five cents, so none of the three is marked; its control
    // number (CLP07) differs. Slot 2's CLP07, DTM*050, and MOA list differ.
    const marked = count(claimHtml, 'border-warn bg-warn/15')
    if (marked !== 4) fail(`${marked} cells are marked, not 4: slot 1's CLP07, and slot 2's CLP07, DTM*050, and MA18`)
    if (claimHtml.includes('$0.02')) fail('slot 1’s payment and adjustment, two cents from the MAC’s, are marked with their differences')
    if (!claimHtml.includes('data-state="left out"')) fail('no row is tagged left out (CLP03 and CLP07 should be)')
    const leftOut = count(claimHtml, 'data-state="left out"')
    if (leftOut < 2) fail(`${leftOut} rows are tagged left out; at least CLP03, CLP07, and the names should be`)
    if (!claimHtml.includes('data-state="compared, watched"')) fail('DTM*050 is not tagged watched')
    const onlyDifferences = render(createElement(ClaimCompareView, { claims, head }), '/claim?match_key=abc&differs=1')
    if (count(onlyDifferences, '<tr class="border-b border-line') >= count(claimHtml, '<tr class="border-b border-line')) fail('Only differences does not reduce the rows')
    if (onlyDifferences.includes('>CLP01<')) fail('Only differences still shows CLP01, which is the same in every column')
    if (onlyDifferences.includes('>CLP04<') || onlyDifferences.includes('>SVC03<')) fail('Only differences keeps CLP04 or SVC03, two cents apart, which the five-cent rule counts as the same')
    if (!onlyDifferences.includes('>CLP07<')) fail('Only differences drops CLP07, which differs')
    // Moved six cents, past the rule, the same three cells are marked, each
    // with its signed difference.
    const moved = { ...claims, v1: parseClaim(claim.v1.replaceAll('79.98', '79.94').replace('20.02', '20.06'), claim.delimiters) }
    const movedHtml = render(createElement(ClaimCompareView, { claims: moved, head }), '/claim?match_key=abc')
    const movedMarks = count(movedHtml, 'border-warn bg-warn/15')
    if (movedMarks !== marked + 3) fail(`with slot 1 six cents off, ${movedMarks} cells are marked, not ${marked + 3}`)
    if (count(movedHtml, '−$0.06') !== 2 || count(movedHtml, '+$0.06') !== 1) fail('slot 1’s payment and line payment, six cents under, and its adjustment, six cents over, do not show −$0.06, −$0.06, and +$0.06')

    // ---- the code popovers (data/codeHints.ts): a title of a few words, then
    // one paragraph that does not repeat it, the maintainer's text for an
    // adjustment or remark code, and no citation.
    const codeList = async (file) => {
      const parsed = JSON.parse(await readFile(new URL(`../../backend/pkg/${file}`, import.meta.url), 'utf8'))
      return Object.fromEntries(parsed.codes.map((c) => [c.code, { description: c.description, status: c.status }]))
    }
    const codes = { carc: await codeList('carc.json'), rarc: await codeList('rarc.json') }
    const CITATION = /PDF p\.|https?:|www\.|Flat File|CEDI|\bIOM\b|MREP|Transmittal|MLN|\bp\.\d/
    const wordCount = (text) => text.trim().split(/\s+/).length
    let popovers = 0
    const checkPopover = (what, e) => {
      popovers++
      for (const text of [e.title, e.body ?? '']) if (CITATION.test(text)) fail(`the ${what} popover cites a document: “${text}”`)
      if (!e.title) fail(`the ${what} popover has no title`)
      if (e.body && e.body.toLowerCase().startsWith(e.title.toLowerCase())) fail(`the ${what} popover's body repeats its title: “${e.title}” / “${e.body}”`)
      if ('lines' in e || 'summary' in e) fail(`the ${what} popover still carries lines of codes and citations`)
    }
    // Every tag the fixtures carry: the title two to five words, a body.
    for (const t of [...v1.tags, ...v2.tags]) {
      const e = hints.describeElement(t.tag, t.label, codes)
      checkPopover(t.tag, e)
      if (wordCount(e.title) < 2 || wordCount(e.title) > 5) fail(`the ${t.tag} popover's title “${e.title}” is not two to five words`)
      if (!e.body) fail(`the ${t.tag} popover has no body`)
    }
    // Every claim and line element as a bare tag, every qualified segment
    // as the comparison tags it, and a code of each kind.
    for (const r of ELEMENTS_835.filter((x) => x.scope === 'claim' || x.scope === 'line')) {
      checkPopover(r.ref, hints.describeElement(r.ref, undefined, codes))
      if (r.tag) checkPopover(r.tag, hints.describeElement(r.tag, undefined, codes))
    }
    for (const code of Object.keys(codes.carc).slice(0, 50)) checkPopover(`CARC ${code}`, hints.describeCarc(code, codes))
    for (const code of Object.keys(codes.rarc).slice(0, 50)) checkPopover(`RARC ${code}`, hints.describeRarc(code, codes))
    for (const code of ['1', '2', '3', '4', '19', '20', '21', '22', '23']) checkPopover(`status ${code}`, hints.describeStatus(code, statusLabel(code)))
    // The example: CAS*CO*253 reads as its label, then the group, the
    // reason, and the reason's text as its maintainers write it.
    const sequestration = hints.describeElement('CAS*CO*253', 'Sequestration adjustment', codes)
    const wantBody = `The amount adjusted under group CO (Contractual Obligation) for reason 253: ${codes.carc['253'].description}`
    if (sequestration.title !== 'Sequestration adjustment' || sequestration.body !== wantBody) fail(`CAS*CO*253 reads “${sequestration.title}” / “${sequestration.body}”`)
    const claimPayment = hints.describeElement('CLP04', 'Claim payment', codes)
    if (claimPayment.title !== 'Claim payment amount' || !claimPayment.body?.startsWith('What Medicare pays the provider for this claim')) fail(`CLP04 reads “${claimPayment.title}” / “${claimPayment.body}”`)
    const remark = hints.describeRarc('MA01', codes)
    if (remark.title !== 'Remark code MA01' || remark.body !== codes.rarc.MA01.description) fail(`MA01 reads “${remark.title}” / “${remark.body}”`)
    if (!hints.describeCarc('253', undefined).loading) fail('a CARC popover does not say its text is loading before /api/codes answers')

    // ---- every app link names a root page
    // The API's answers for the pages: the fixtures above where there are
    // some (the overview's, the agreement grid's, the MAC map's, the claim's
    // segments), and small ones here for the review pages. A claim is named
    // by its hex match key.
    const MATCH_KEY = 'a1b2c3d4'
    const macSubmission = { prefix: 'mac/2026-10-01', side: 'mac', label: 'MAC', submitted_on: '2026-10-01', status: 'summarized', claim_count: legacy.totals.claims }
    const fronts = new Map([[v1.submission.prefix, v1], [v2.submission.prefix, v2]])
    const measure = (claimType, workload, macClaims) => ({
      claim_type: claimType, workload, mac_claims: macClaims, equal: macClaims - 5, different: 2, missing: 1, extra: 1, set_aside: 1, duplicated: 1,
      outcome_match_rate: '0.9000', paid_claims: macClaims - 2, median_deviation: '0.0010', p95_deviation: '0.0090',
    })
    const kinds = {
      different_claims: 2,
      rows: v1.tags.slice(0, 2).map((t) => ({ tag: t.tag, label: t.label, claims: t.claims, dollars: t.dollars, previous_claims: null })),
    }
    // A claim as /api/claim sends it: its bytes as text and its file's delimiters.
    const claimOf = (segments) => ({ segments, delimiters: claim.delimiters })
    const unmatchedRow = (pcn) => ({ match_key: MATCH_KEY, key: claim.key, claim_type: 'professional', patient_control_number: pcn, paid: '80.00' })
    // Each endpoint's answer, from the inputs in its useApi key:
    // claim:<prefix>|<match key>, unmatched:<prefix>|<outcome>,
    // tops:<prefix>|<workloads>, and so on.
    const ANSWERS = {
      submissions: () => ({ submissions: [macSubmission, v1.submission, v2.submission] }),
      measures: (prefix) => ({
        submission: fronts.get(prefix).submission, previous: null, previous_measures: [],
        measures: [
          measure('all', 'all', 20), measure('institutional', 'all', 10), measure('professional', 'all', 10),
          measure('all', '01111', 10), measure('all', '01112', 10), measure('institutional', '01111', 10), measure('professional', '01112', 10),
        ],
      }),
      differences: () => kinds,
      tops: (_, workloads) => workloads.split(',').map((w) => [w, kinds.rows[0]]),
      claims: () => ({
        claims: 1, dollars: '0.02', histogram: [{ from: '-0.02', to: '-0.02', count: 1 }], page: 1, pages: 1,
        rows: [{ match_key: MATCH_KEY, key: claim.key, claim_type: 'professional', mac_paid: '80.00', vendor_paid: '79.98', delta: '-0.02' }],
      }),
      claim: (prefix) => {
        const vendor = claimOf(prefix === v2.submission.prefix ? claim.v2 : claim.v1)
        return {
          key: claim.key, outcome: 'different', mac: claimOf(claim.mac), vendor, vendor_claims: [vendor],
          differences: [{ element: 'CLP04', occurrence: 1, line_control_number: null, mac_value: '80.00', vendor_value: '79.98', amount_difference: '-0.02' }],
        }
      },
      unmatched: (_, outcome) => ({
        outcome, count: 2, share: '0.0020', by_claim_type: [{ claim_type: 'professional', count: 2 }], by_workload: [{ workload: '01112', count: 2 }],
        rows: [unmatchedRow('PCN3001'), unmatchedRow('PCN3002')],
      }),
      front: (prefix) => fronts.get(prefix),
      agreement: () => agreement,
      macs: () => macs,
    }
    const answer = (key) => {
      const at = key.indexOf(':')
      return at < 0 ? ANSWERS[key]?.() : ANSWERS[key.slice(0, at)]?.(...key.slice(at + 1).split('|'))
    }
    // The whole app at one address. ThemeToggle reads the page's theme while
    // rendering (document.documentElement); the server has no document, so
    // the app renders with a stand-in that says light.
    const renderApp = (location) => {
      const had = Object.getOwnPropertyDescriptor(globalThis, 'document')
      Object.defineProperty(globalThis, 'document', { value: { documentElement: { getAttribute: () => 'light' } }, configurable: true, writable: true })
      globalThis.answerApi = answer
      try {
        return renderToStaticMarkup(createElement(StaticRouter, { location }, createElement(App)))
      } finally {
        delete globalThis.answerApi
        if (had) Object.defineProperty(globalThis, 'document', had)
        else delete globalThis.document
      }
    }
    // Every href in some markup, its &amp; read back. A link to an app page
    // names no other site (a scheme, or //) and is not an anchor on the page.
    const hrefsOf = (markup) => [...markup.matchAll(/ href="([^"]*)"/g)].map((m) => m[1].replace(/&amp;/g, '&'))
    const isAppLink = (href) => !/^([a-z][a-z0-9+.-]*:|\/\/|#)/i.test(href)
    const pathOf = (href) => href.split(/[?#]/)[0]
    const paramsOf = (href) => new URLSearchParams(/\?([^#]*)/.exec(href)?.[1] ?? '')
    // The app's pages, every one at the root, and its renders: the address,
    // then the links the page itself must draw (a page and a query parameter
    // its link carries), or text it must show where its own links have the
    // same addresses as the navigation's.
    const PAGES = ['/', '/claims', '/differences', '/unmatched', '/workloads', '/claim']
    const RENDERS = [
      { at: '/', wants: ['/claims?v1_outcome', '/claims?v2_outcome', '/unmatched?outcome', '/claims?tag'] },
      { at: '/claims', wants: ['/claim?match_key'] },
      { at: '/differences', wants: ['/claims?tag'] },
      { at: '/unmatched', text: 'All workloads →' },
      { at: '/unmatched?outcome=set_aside', wants: ['/claim?match_key'] },
      { at: '/workloads', wants: ['/differences?workload', '/claims?tag'] },
      { at: '/claim', text: 'Open a claim from' },
      { at: `/claim?match_key=${MATCH_KEY}`, text: 'Only differences' },
      { at: '/no-such-page', text: 'No page at' },
    ]
    let links = 0
    for (const r of RENDERS) {
      const markup = renderApp(r.at)
      const appLinks = hrefsOf(markup).filter(isAppLink)
      links += appLinks.length
      if (appLinks.length === 0) fail(`${r.at} draws no links`)
      for (const href of appLinks) {
        const path = pathOf(href)
        if (!PAGES.includes(path)) fail(`${r.at}: a link names ${href}, which is not one of the app's pages at the root`)
      }
      for (const want of r.wants ?? []) {
        const [path, param] = want.split('?')
        if (!appLinks.some((href) => pathOf(href) === path && paramsOf(href).has(param))) fail(`${r.at}: no link to ${path} with ${param}=, so the page did not draw its own links`)
      }
      if (r.text !== undefined && !markup.includes(r.text)) fail(`${r.at} does not show “${r.text}”`)
    }

    // ---- nothing in the drawing code reads the network at render time
    const drawing = ['OverviewView.tsx', 'SelectionHeader.tsx', 'Breakdowns.tsx', 'overviewColumns.tsx', 'overviewFilters.ts', 'ClaimCompareView.tsx']
    for (const file of drawing) {
      if (/\bfetch\(/.test(await source(`src/app/${file}`))) fail(`src/app/${file} calls fetch`)
    }

    if (failures.length > 0) {
      console.error('App SSR check failed:')
      for (const f of failures) console.error(`  - ${f}`)
      process.exitCode = 1
      return
    }
    console.log(
      `App SSR check passed: the sticky selection header offers ${BOXES.length} boxes; ${boxChecks} box drawings with ${tabs} value tabs checked ` +
      `(no selection, the sixteen jurisdictions four to a row, and the hidden-zero rule under J5, JE, and DME picked), and the workload box, closed with its count and opened in two tiers over the MAC map's ${macs.jurisdictions.length} jurisdictions; ` +
      `${tables.length} tables sit on the four-column grid, ${legacyTables} of them with a Legacy column and ${tables.length - legacyTables} with it empty; ` +
      `${breakdowns.length} breakdown tables draw ${rows} rows whose legacy columns add to ${legacy.totals.claims} MAC claims; ` +
      `${v1.outcomes.length} outcomes, ${BANDS.length} payment bands, ${agreement.cells.length} agreement cells, ${v1.status_matrix.length}+${v2.status_matrix.length} matrix cells, and ${v1.tags.length}+${v2.tags.length} tags are drawn, slot 1 before slot 2; ` +
      `the selection round-trips; the claim screen reads 3 claims into ${claimHtml.split('<tbody').length - 1} sections with ${marked} marked cells (${movedMarks} with slot 1 six cents off) and ${leftOut} rows tagged left out; ` +
      `all ${links} app links in ${RENDERS.length} renders name the app's pages at the root; ` +
      `${popovers} code popovers carry a short title, a body that does not repeat it, and no citation.`,
    )
  } finally {
    await server.close()
  }
}

main().catch((err) => {
  console.error(err)
  process.exitCode = 1
})
