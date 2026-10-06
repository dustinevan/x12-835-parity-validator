import { useMemo, type ReactNode } from 'react'
import { Link } from 'react-router'
import { DocSection } from '../../components/DocSection'
import { PageTitle } from '../../components/PageTitle'
import type { TagColor } from '../../components/Tag'
import { AMOUNT_MATCH_CENTS, LEFT_OUT, WATCHED, stateOf, tallyText, type LeafState, type Level } from '../../data/comparedFields'
import { VIEWER_835_EXAMPLES } from '../../data/viewer835Examples'
import { parsedOf } from '../compare835/columns'
import { StateTree } from '../compare835/StateTree'
import { unionTreeOf } from '../compare835/unionTree'
import { CodeLink } from '../viewer/Code'
import { LookupButton, LookupFrame } from '../viewer/LookupFrame'
import { useLookup } from '../viewer/useLookup'
import { LOOKUP_835 } from '../viewer835/lookup'

// The 835 elements the comparison looks at, drawn two ways. First, the elements
// of the worked examples in file order, the tree of the 835 compare page, each
// with a state (the tree is compare835/unionTree.ts; the states are stateOf in
// data/comparedFields.ts). Then the leave-out list of compare set 1, the
// elements the comparison does not look at, each with its reason. The set is
// backend/lib/core/service/compareset.go; scripts/ssr-check-compared-fields.mjs
// reads that file and fails when this page's list differs from it. Every code
// in the tree and the list is a link into the 835 lookup panel, which opens
// beside the page (viewer/LookupFrame.tsx). The panel's lookup= and q= are the
// only things the page puts in the address; which groups of the tree are
// closed is page state only, and nothing is stored.

// ---------------------------------------------------------------- small parts

function Th({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <th className={`border-b border-line pb-2 pr-4 text-left font-mono text-[11px] uppercase tracking-wide text-ink-2 ${className ?? ''}`}>
      {children}
    </th>
  )
}

function Td({ children, strong }: { children: ReactNode; strong?: boolean }) {
  return <td className={`border-b border-line py-2 pr-4 align-top ${strong ? 'font-mono text-[13px] text-ink' : ''}`}>{children}</td>
}

// ---------------------------------------------------------------- the file in order

// Compared is blue, left out and outside the claim are gray, the gap is
// yellow. Never red or green.
const STATE_COLOR: Record<LeafState, TagColor> = {
  compared: 'blue',
  'compared, watched': 'blue',
  'left out': 'gray',
  'outside the claim': 'gray',
  unlisted: 'yellow',
}

function FileTree() {
  const tree = useMemo(() => unionTreeOf(VIEWER_835_EXAMPLES.map((example) => parsedOf(example.id))), [])
  return (
    <StateTree
      tree={tree}
      markOf={(leaf) => {
        const { state, why } = stateOf(leaf.id)
        return { state, color: STATE_COLOR[state], title: why }
      }}
      tallyOf={tallyText}
    />
  )
}

// ---------------------------------------------------------------- the leave-out list

const LEVELS: { level: Level; title: string }[] = [
  { level: 'claim', title: 'Claim level' },
  { level: 'line', title: 'Service line level' },
]

function LeftOutTable() {
  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse text-[14px]">
        <thead>
          <tr>
            <Th className="min-w-[8rem]">tag</Th>
            <Th className="min-w-[14rem]">public name</Th>
            <Th className="min-w-[20rem]">reason</Th>
          </tr>
        </thead>
        {LEVELS.map(({ level, title }) => {
          const entries = LEFT_OUT.filter((entry) => entry.level === level)
          if (entries.length === 0) return null
          return (
            <tbody key={level}>
              <tr>
                <th colSpan={3} scope="colgroup" className="border-b border-line pb-1.5 pt-8 text-left font-semibold text-ink">
                  {title}
                </th>
              </tr>
              {entries.map((entry) => (
                <tr key={entry.tag} data-left-out={`${entry.level}:${entry.tag}`}>
                  <Td strong>
                    <CodeLink code={entry.tag} />
                  </Td>
                  <Td>{entry.publicName}</Td>
                  <Td>{entry.reason}</Td>
                </tr>
              ))}
            </tbody>
          )
        })}
      </table>
    </div>
  )
}

// ---------------------------------------------------------------- page

export function ComparedFields() {
  const lookup = useLookup()
  return (
    <LookupFrame lookup={lookup} data={LOOKUP_835}>
      <div className="mx-auto max-w-[96rem] px-10 py-12">
        <PageTitle
          lede={
            <>
              The 835 elements the comparison looks at, as compare set 1 in the code: everything inside a claim except the
              elements left out below. Two amounts within ${(AMOUNT_MATCH_CENTS / 100).toFixed(2)} of each other are the same
              amount; everything else compares exactly. The record of how the set was decided is{' '}
              <code className="bg-surface px-1.5 py-0.5 font-mono text-[13px] text-ink">docs/design/compared-fields.md</code>.
            </>
          }
          actions={<LookupButton lookup={lookup} />}
        >
          Compared fields
        </PageTitle>

        <DocSection id="tree" title="In file order">
          <FileTree />
          <p className="text-[13px] leading-6">
            Each part and group shows a tally of its elements. The tree is the union of the {VIEWER_835_EXAMPLES.length} worked
            examples in the{' '}
            <Link to="/835-viewer" className="text-link hover:underline">
              835 viewer
            </Link>
            , in the order of the{' '}
            <Link to="/835-compare" className="text-link hover:underline">
              835 compare
            </Link>{' '}
            page, so an element none of them sends has no entry (for example <code>CAS04</code>, a CAS quantity, and{' '}
            <code>AMT*NL</code>). Where a segment carries a qualifier (DTM, REF, AMT, QTY, LQ, PER), each entry is one
            qualified tag, as the comparison labels it; the compare page lists <code>DTM01</code> and <code>DTM02</code> once
            for every qualifier, which cannot hold one state per qualifier. The envelope (ISA, GS, ST, SE, GE, IEA) is not
            listed there or here; it is outside the claim.
          </p>
        </DocSection>

        <DocSection id="left-out" title="Left out">
          <LeftOutTable />
          <p>
            <strong>Compared and watched:</strong>{' '}
            {WATCHED.map((entry, i) => (
              <span key={entry.label}>
                {i > 0 && ', '}
                <code>{entry.label}</code>
              </span>
            ))}
            . They are compared and recorded like every other element. Whether an element belongs in the compared set is
            settled from the real files, and the set is named on every run, so a change is visible.
          </p>
          <p>
            The comparison reads the identifier of an NM1 (<code>NM109</code>) under its role, not its name elements (
            <code>NM103</code> to <code>NM107</code>), so those are left out without an entry in the list above.
          </p>
        </DocSection>
      </div>
    </LookupFrame>
  )
}
