// Loads the 837P viewer's own parsing and rendering code through Vite's SSR
// module graph — the same way the docs site resolves a `?raw` import — and
// exercises every one of the 24 worked examples under
// frontend/src/data/viewer837p/files/ with no browser involved: parses each
// raw file into its loop tree, confirms every segment the file sent landed
// somewhere in that tree (nothing silently dropped), and calls
// describeOpening/describeRow/describeClosing on every box and row to
// confirm nothing throws building the display cells. It also rewrites each
// example in other delimiters (one long line with | and >, LF and CRLF line
// ends, a BOM) and confirms the same segments and the same tree come back,
// since the parser reads its delimiters from the ISA. The examples are
// Medicare-shaped and send only HI0n-1 and HI0n-2, so it also draws a
// synthetic HI segment that carries HI0n-3 through HI0n-9 and confirms each
// shows as its own named cell, while a composite that stops at the code
// still draws as one cell. Run from frontend/:
//   node scripts/ssr-check-837p.mjs
import { createServer } from 'vite'
import { DELIMITER_VARIANTS, rewriteDelimiters, shape } from './x12-rewrite.mjs'

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
    const { VIEWER_837P_EXAMPLES } = await server.ssrLoadModule('/src/data/viewer837p/index.ts')
    const { parseLoopTree, splitSegments } = await server.ssrLoadModule('/src/docs/viewer837p/parse.ts')
    const { describeOpening, describeRow, describeClosing } = await server.ssrLoadModule('/src/docs/viewer837p/describe.ts')
    const { verifySample } = await server.ssrLoadModule('/src/docs/viewer837p/check.ts')

    if (!Array.isArray(VIEWER_837P_EXAMPLES) || VIEWER_837P_EXAMPLES.length === 0) {
      throw new Error('VIEWER_837P_EXAMPLES is empty — nothing to check')
    }

    const failures = []
    for (const example of VIEWER_837P_EXAMPLES) {
      try {
        const totalSegments = splitSegments(example.file).length
        const root = parseLoopTree(example.file)

        let visited = 0
        function walk(box) {
          visited++ // opening
          describeOpening(box)
          if (box.closing) {
            visited++
            describeClosing(box)
          }
          for (const node of box.children) {
            if (node.kind === 'box') {
              walk(node.box)
            } else {
              visited++
              describeRow(node.segment, box.boxType)
            }
          }
        }
        walk(root)

        if (visited !== totalSegments) {
          failures.push(`${example.id}: tree holds ${visited} segments but the file has ${totalSegments}`)
        }

        const balance = verifySample(example.file)
        if (!balance.ok) failures.push(`${example.id}: ${balance.details.join('; ')}`)

        for (const variant of DELIMITER_VARIANTS) {
          const rewritten = rewriteDelimiters(example.file, splitSegments, variant)
          const count = splitSegments(rewritten).length
          if (count !== totalSegments) {
            failures.push(`${example.id} (${variant.name}): ${count} segments read, the file has ${totalSegments}`)
            continue
          }
          const same = JSON.stringify(shape(parseLoopTree(rewritten))) === JSON.stringify(shape(root))
          if (!same) failures.push(`${example.id} (${variant.name}): the tree differs from the standard-delimiter tree`)
        }
      } catch (err) {
        failures.push(`${example.id}: ${err instanceof Error ? err.stack : String(err)}`)
      }
    }

    // HI0n-3 through HI0n-9 (Not Used by Medicare, defined by the standard).
    const names = await server.ssrLoadModule('/src/data/x12Names837p.ts')
    const hiRefs = (elements) => describeRow({ id: 'HI', elements, raw: '' }, 'CLAIM').cells.map((c) => c.ref)
    const hiCells = describeRow({ id: 'HI', elements: ['ABK:J069:D8:20260901:12.50:2.5:V1:IND:Y', '', 'ABF:R05'], raw: '' }, 'CLAIM').cells
    const wantRefs = ['HI01', 'HI01-3', 'HI01-4', 'HI01-5', 'HI01-6', 'HI01-7', 'HI01-8', 'HI01-9', 'HI02', 'HI03']
    if (hiCells.map((c) => c.ref).join(' ') !== wantRefs.join(' ')) {
      failures.push(`HI with components 3 to 9: cells ${hiCells.map((c) => c.ref).join(' ')}, want ${wantRefs.join(' ')}`)
    } else if (JSON.stringify(hiCells.slice(1, 8).map((c) => c.name)) !== JSON.stringify(names.HI_NOT_USED_EL_NAME)) {
      failures.push('HI with components 3 to 9: the cells HI01-3 to HI01-9 are not named from HI_NOT_USED_EL_NAME')
    }
    if (hiRefs(['ABK:J069', 'ABF:R05']).join(' ') !== 'HI01 HI02') {
      failures.push(`HI with two components: cells ${hiRefs(['ABK:J069', 'ABF:R05']).join(' ')}, want HI01 HI02`)
    }

    if (failures.length > 0) {
      console.error(`837P SSR check FAILED (${failures.length} problems; ${VIEWER_837P_EXAMPLES.length} examples checked):`)
      for (const f of failures) console.error(`  - ${f}`)
      process.exitCode = 1
      return
    }

    console.log(`837P SSR check passed: ${VIEWER_837P_EXAMPLES.length} examples parsed, fully covered (every segment landed in the tree), and rendered with no error; each also reads the same under ${DELIMITER_VARIANTS.length} other delimiter spellings.`)
  } finally {
    await server.close()
  }
}

main().catch((err) => {
  console.error(err)
  process.exitCode = 1
})
