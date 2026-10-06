// Examples in the sidebar group "Service line" (docs/design/835-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import feeScheduleReduction from './files/fee-schedule-reduction.835?raw'
import deductibleAndCoinsurance from './files/deductible-and-coinsurance.835?raw'
import bundledIntoAnotherLine from './files/bundled-into-another-line.835?raw'
import lineRemark from './files/line-remark.835?raw'
import drugBilledByNdc from './files/drug-billed-by-ndc.835?raw'
import procedureCodeChanged from './files/procedure-code-changed.835?raw'

export const SERVICE_LINE: ViewerExample[] = [
  {
    id: 'fee-schedule-reduction',
    group: 'Service line',
    title: 'Fee-schedule reduction',
    file: feeScheduleReduction,
    open: ['ISA.GS.ST.LX#0.CLP#0.SVC#0'],
  },
  {
    id: 'deductible-and-coinsurance',
    group: 'Service line',
    title: 'Deductible and coinsurance',
    file: deductibleAndCoinsurance,
    open: ['ISA.GS.ST.LX#0.CLP#0.SVC#0'],
  },
  {
    id: 'bundled-into-another-line',
    group: 'Service line',
    title: 'Bundled into another line',
    file: bundledIntoAnotherLine,
    open: ['ISA.GS.ST.LX#0.CLP#0.SVC#1'],
  },
  {
    id: 'line-remark',
    group: 'Service line',
    title: 'Line remark',
    file: lineRemark,
    open: ['ISA.GS.ST.LX#0.CLP#0.SVC#0'],
  },
  {
    id: 'drug-billed-by-ndc',
    group: 'Service line',
    title: 'Drug billed by NDC',
    file: drugBilledByNdc,
    open: ['ISA.GS.ST.LX#0.CLP#0.SVC#0'],
  },
  {
    id: 'procedure-code-changed',
    group: 'Service line',
    title: 'Procedure code changed',
    file: procedureCodeChanged,
    open: ['ISA.GS.ST.LX#0.CLP#0.SVC#0'],
  },
]
