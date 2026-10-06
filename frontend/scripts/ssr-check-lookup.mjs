// Checks the 835 field lookup's entries (frontend/src/data/x12About835.ts,
// read through lookupRows835.ts and drawn by docs/viewer/Lookup.tsx) with
// no browser:
//   - every element row of the 835 lookup has an entry (what it is, with at
//     least one source) and a validator note, and every row's entry renders
//     on the server with its sentences in the output;
//   - every entry key in x12About835.ts is used by some row (an element's
//     key by that element, a segment's key by the rows of that segment at
//     that use), so a misspelled key cannot leave an entry orphaned;
//   - every coded-value row has a public name or an entry text, so no code
//     is listed with nothing said about it;
//   - every citation, on the rows, in the entries, and on the codes, names
//     a document under docs/reference/ and a PDF page inside it (the page
//     counts are listed here), a file in the repository, or a public web
//     page on the listed hosts;
//   - the validator notes agree with the Compared fields page: an element
//     the leave-out list names is "left out", a watched element is
//     "compared, watched", a qualifier element is "part of the tag", and
//     nothing inside a claim is "unlisted";
//   - the panel, rendered on the server with the search "CLP09", opens that
//     one row by itself and shows its entry: the element's own name, the
//     third position of the type of bill, its validator state, and its
//     flat-file page;
//   - a click on a code searches the panel for that code (docs/viewer/
//     Code.tsx), so every element row shows when the panel is searched for
//     its reference, its segment ID, or its qualified tag (DTM*232, NM1*QC);
//     it prints how many references leave one element, whose entry then
//     opens by itself.
// Run from frontend/:
//   node scripts/ssr-check-lookup.mjs
import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { existsSync } from 'node:fs'
import { resolve } from 'node:path'
import { createServer } from 'vite'

// The documents a citation may name, with the file under the repository
// root and its page count.
const DOCUMENTS = {
  'CMS 835 Flat File': { file: 'docs/reference/file-formats/cms-835-flatfile.pdf', pages: 15 },
  'CEDI 835 CG': { file: 'docs/reference/file-formats/cedi-835-companion-guide.pdf', pages: 33 },
  'CEDI 837P CG': { file: 'docs/reference/file-formats/cedi-837p-companion-guide.pdf', pages: 56 },
  'CMS 835 TI CG (2012)': { file: 'docs/reference/file-formats/cms-835-ti-companion-guide-2012.pdf', pages: 10 },
  'Transmittal 907 (2011)': { file: 'docs/reference/file-formats/cms-837-835-flatfiles-r907otn-2011.pdf', pages: 283 },
  'Transmittal 579 (2009)': { file: 'docs/reference/file-formats/cms-837-flatfiles-r579otn-2009.pdf', pages: 251 },
  'IOM 100-04 ch.22 (Rev. 11427)': { file: 'docs/reference/manuals/clm104c22.pdf', pages: 20 },
  'IOM 100-04 ch.22 (2003 ed.)': { file: 'docs/reference/manuals/clm104c22_remit_phase3_apr24-03_r2.pdf', pages: 82 },
  'MLN Remittance Advice Resources and FAQs (ICN 905367)': { file: 'docs/reference/manuals/mln-remittance-advice-resources-and-faqs-icn905367.pdf', pages: 15 },
  // Cited by URL, not stored: the guide's footer restricts reproduction.
  'MREP User Guide v4.4 (2016)': { url: 'https://www.cms.gov/', pages: 223 },
}
const HOSTS = ['x12.org', 'www.cms.gov', 'www.ngscedi.com']

const root = resolve(process.cwd(), '..')

// A citation is valid when every part of it (parts separated by "; ") names
// a document and page, a repository file, or a public web page.
function checkCitation(text, fail, where) {
  for (const part of text.split('; ')) {
    const trimmed = part.trim()
    const doc = Object.keys(DOCUMENTS).find((name) => trimmed.startsWith(name))
    if (doc) {
      const { url, pages } = DOCUMENTS[doc]
      if (url && !trimmed.includes(url)) fail(`${where}: ${doc} is cited without its URL: ${trimmed}`)
      const m = trimmed.match(/PDF p\.(\d+)(?:-(\d+))?/)
      if (!m) {
        fail(`${where}: no PDF page in "${trimmed}"`)
        continue
      }
      const from = Number(m[1])
      const to = m[2] ? Number(m[2]) : from
      if (from < 1 || to > pages || to < from) fail(`${where}: PDF p.${m[1]}${m[2] ? `-${m[2]}` : ''} is outside ${doc} (${pages} pages)`)
      continue
    }
    const url = trimmed.match(/https?:\/\/[^\s)]+/)
    if (url) {
      const host = new URL(url[0]).host
      if (!HOSTS.includes(host)) fail(`${where}: web page on an unlisted host: ${url[0]}`)
      continue
    }
    if (/^[\w./-]+$/.test(trimmed) && trimmed.includes('/')) continue
    fail(`${where}: unrecognized citation "${trimmed}"`)
  }
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
    const { ELEMENTS_835, CODES_835 } = await server.ssrLoadModule('/src/data/lookupRows835.ts')
    const { ABOUT_835, CODE_ABOUT_835 } = await server.ssrLoadModule('/src/data/x12About835.ts')
    const { SEGMENT_NAME } = await server.ssrLoadModule('/src/data/x12Names835.ts')
    const { LEFT_OUT, WATCHED, NAME_ELEMENTS } = await server.ssrLoadModule('/src/data/comparedFields.ts')
    const { LookupPanel, LookupEntry, rowKey } = await server.ssrLoadModule('/src/docs/viewer/Lookup.tsx')

    const failures = []
    const fail = (message) => failures.push(message)
    const texts = []

    // ---- every element row has an entry, a validator note, and renders
    const usedKeys = new Set()
    let withEntry = 0
    for (const row of ELEMENTS_835) {
      const where = `${row.use}`
      if (!row.use) fail(`${row.ref} (${row.where}): no use key`)
      if (!row.about) fail(`${where}: no entry`)
      else {
        withEntry++
        usedKeys.add(row.use)
        if (!row.about.is || row.about.is.length < 12) fail(`${where}: the entry's "is" is empty or too short`)
        if (!Array.isArray(row.about.sources) || row.about.sources.length === 0) fail(`${where}: the entry has no source`)
        for (const s of row.about.sources) checkCitation(s, fail, where)
        for (const text of [row.about.is, row.about.purpose, row.about.values, row.about.medicare]) if (text) texts.push([where, text])
      }
      if (row.segmentAbout) {
        const key = `${row.scope}:${row.segment}${row.use.includes('@') ? row.use.slice(row.use.indexOf('@')) : ''}`
        usedKeys.add(key)
        for (const s of row.segmentAbout.sources) checkCitation(s, fail, key)
        for (const text of [row.segmentAbout.is, row.segmentAbout.purpose, row.segmentAbout.medicare]) if (text) texts.push([key, text])
      }
      if (!row.validator) fail(`${where}: no validator note`)
      checkCitation(row.source, fail, `${where} (row)`)
      if (row.echo) texts.push([where, row.echo])
      const html = renderToStaticMarkup(createElement(LookupEntry, { row }))
      if (row.about && !html.includes(escape(row.about.is))) fail(`${where}: the rendered entry lacks its first sentence`)
      if (row.validator && !html.includes(`data-validator="${row.validator.state}"`)) fail(`${where}: the rendered entry lacks its validator state`)
    }
    for (const key of Object.keys(ABOUT_835)) if (!usedKeys.has(key)) fail(`x12About835.ts has an entry no row uses: ${key}`)

    // ---- every code row is named or explained, with valid citations
    let codesWithText = 0
    const codeKeys = new Set()
    for (const c of CODES_835) {
      const where = `${c.segment} ${c.element} ${c.code}`
      if (!c.name && !c.about) fail(`${where}: neither a public name nor an entry text`)
      if (c.about) {
        codesWithText++
        codeKeys.add(`${c.segment}:${c.element}:${c.code}`)
        for (const s of c.about.sources) checkCitation(s, fail, where)
        texts.push([where, c.about.text])
      }
      checkCitation(c.source, fail, `${where} (row)`)
    }
    for (const key of Object.keys(CODE_ABOUT_835)) if (!codeKeys.has(key)) fail(`x12About835.ts has a code entry no code row uses: ${key}`)

    // ---- the validator notes agree with the Compared fields page
    const QUALIFIERS = new Set(['NM101', 'REF01', 'DTM01', 'AMT01', 'QTY01', 'LQ01', 'PER01', 'CAS01', 'CAS02', 'CAS05', 'CAS08', 'CAS11', 'CAS14', 'CAS17'])
    const inside = ELEMENTS_835.filter((r) => r.scope === 'claim' || r.scope === 'line')
    for (const r of inside) {
      const expectQualifier = QUALIFIERS.has(r.ref)
      if (expectQualifier && r.validator.state !== 'part of the tag') fail(`${r.use}: a qualifier element should be "part of the tag", got "${r.validator.state}"`)
      if (!expectQualifier && r.validator.state === 'part of the tag') fail(`${r.use}: not a qualifier element, yet "part of the tag"`)
      if (r.validator.state === 'unlisted') fail(`${r.use}: unlisted`)
    }
    const outside = ELEMENTS_835.filter((r) => r.scope !== 'claim' && r.scope !== 'line')
    for (const r of outside) if (r.validator.state !== 'outside the claim') fail(`${r.use}: outside a claim yet "${r.validator.state}"`)
    const stateOfUse = (use) => ELEMENTS_835.find((r) => r.use === use)?.validator.state
    const expectState = (use, state) => {
      const got = stateOfUse(use)
      if (got !== state) fail(`${use}: expected "${state}", got "${got}"`)
    }
    // The leave-out list's entries, by their value positions.
    const leftOutUses = {
      CLP03: 'claim:CLP03',
      'NM1*QC': 'claim:NM109@QC',
      'DTM*232': 'claim:DTM02@232',
      CLP07: 'claim:CLP07',
      'REF*F8': 'claim:REF02@F8',
      CLP10: 'claim:CLP10',
      CLP14: 'claim:CLP14',
      PER: 'claim:PER04@CX',
      'REF*6R': 'line:REF02@6R',
      'QTY*': 'line:QTY02',
    }
    for (const entry of LEFT_OUT) {
      const use = leftOutUses[entry.tag]
      if (!use) {
        fail(`the leave-out list has an entry this check does not map: ${entry.tag}`)
        continue
      }
      // The line-level QTY is marked not used by Medicare and has no row.
      if (!ELEMENTS_835.some((r) => r.use === use)) {
        if (use !== 'line:QTY02') fail(`the leave-out entry ${entry.tag} maps to ${use}, which no row has`)
        continue
      }
      expectState(use, 'left out')
    }
    for (const group of WATCHED) {
      const use = { 'DTM*050': 'claim:DTM02@050', 'NM1*TT': 'claim:NM109@TT' }[group.label]
      if (use) expectState(use, 'compared, watched')
    }
    if (NAME_ELEMENTS.keys.length > 0) for (const role of ['QC', '74', '82', 'TT', 'PR']) expectState(`claim:NM103@${role}`, 'left out')
    expectState('claim:CLP04', 'compared')
    expectState('line:SVC02', 'compared')
    expectState('claim:CAS03', 'compared')
    expectState('claim:CAS04', 'compared')
    expectState('claim:MIA23', 'compared')
    expectState('claim:MOA07', 'compared')
    const amt = ELEMENTS_835.find((r) => r.use === 'claim:AMT02')
    if (!amt?.validator.detail?.includes('compared, watched: AMT*I, AMT*NL')) fail(`claim:AMT02: the detail should list the watched kinds, got "${amt?.validator.detail}"`)

    // ---- the panel opens the one row a search leaves, with its entry
    const html = renderToStaticMarkup(
      createElement(LookupPanel, {
        title: '835 field lookup',
        elements: ELEMENTS_835,
        codes: CODES_835,
        segmentNames: SEGMENT_NAME,
        query: 'CLP09',
        onQueryChange: () => {},
        onClose: () => {},
      }),
    )
    const clp09 = ELEMENTS_835.find((r) => r.use === 'claim:CLP09')
    for (const needle of [`data-entry="claim:CLP09"`, `data-validator="${clp09?.validator.state}"`]) {
      if (!html.includes(needle)) fail(`the panel searched for CLP09 does not show: ${needle}`)
    }
    for (const needle of ['Claim Frequency Code', 'third position of the type of bill', 'CMS 835 Flat File, PDF p.6']) {
      if (!html.includes(escape(needle))) fail(`the panel searched for CLP09 does not show: ${needle}`)
    }
    if ((html.match(/data-entry="/g) ?? []).length !== 1) fail('the panel searched for CLP09 should open exactly one entry')
    const panelFor = (query) =>
      renderToStaticMarkup(createElement(LookupPanel, { title: 't', elements: ELEMENTS_835, codes: CODES_835, segmentNames: SEGMENT_NAME, query, onQueryChange: () => {}, onClose: () => {} }))
    const closed = panelFor('')
    if (closed.includes('data-entry="')) fail('with no search, every entry should start closed')
    if (typeof rowKey !== 'function') fail('rowKey is not exported')

    // ---- a click on a code opens the panel searched for that code, and the
    // panel shows that element: every element's reference (a viewer cell, a
    // compare-page key), its segment ID (a row's lead, a tree group), and its
    // qualified tag where it has one (DTM*232, NM1*QC: a compare-page note, a
    // tree group or leaf). Each search is rendered once.
    const searched = new Map()
    const shownFor = (query) => {
      if (!searched.has(query)) searched.set(query, panelFor(query))
      return searched.get(query)
    }
    let tagged = 0
    for (const row of ELEMENTS_835) {
      if (row.tag) tagged++
      for (const query of [row.ref, row.segment, row.tag]) {
        if (query && !shownFor(query).includes(`data-row="${rowKey(row)}"`)) fail(`the panel searched for ${query} does not show ${row.use}`)
      }
    }
    const refs = [...new Set(ELEMENTS_835.map((r) => r.ref))]
    const opensAlone = refs.filter((ref) => (shownFor(ref).match(/data-entry="/g) ?? []).length === 1).length

    console.log(`elements: ${ELEMENTS_835.length}, with an entry: ${withEntry}; codes: ${CODES_835.length}, named: ${CODES_835.filter((c) => c.name).length}, with an entry text: ${codesWithText}; entry keys: ${Object.keys(ABOUT_835).length}; citations checked on ${texts.length} texts`)
    console.log(`code searches: ${searched.size} distinct (every element's reference, segment ID, and the tag of the ${tagged} rows that have one), each showing its rows; ${opensAlone} of ${refs.length} element references leave one element, which opens its entry by itself`)
    if (failures.length > 0) {
      console.error(`FAIL (${failures.length}):`)
      for (const f of failures) console.error(`  - ${f}`)
      process.exitCode = 1
    } else {
      console.log('OK')
    }
  } finally {
    await server.close()
  }
}

// React escapes text the way renderToStaticMarkup writes it.
function escape(text) {
  return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#x27;')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
