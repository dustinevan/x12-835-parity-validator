// Examples in the sidebar group "Claim" (docs/design/835-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import processedAsPrimary from './files/processed-as-primary.835?raw'
import denied from './files/denied.835?raw'
import reversalAndCorrection from './files/reversal-and-correction.835?raw'
import medicareSecondary from './files/medicare-secondary.835?raw'
import paidDownAPayerChain from './files/paid-down-a-payer-chain.835?raw'
import crossoverToASupplement from './files/crossover-to-a-supplement.835?raw'
import inpatient from './files/inpatient.835?raw'
import patientNameCorrected from './files/patient-name-corrected.835?raw'
import interestPaid from './files/interest-paid.835?raw'
import dmeRental from './files/dme-rental.835?raw'
import twoClaims from './files/two-claims.835?raw'

export const CLAIM: ViewerExample[] = [
  {
    id: 'processed-as-primary',
    group: 'Claim',
    title: 'Processed as primary',
    file: processedAsPrimary,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'denied',
    group: 'Claim',
    title: 'Denied',
    file: denied,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'reversal-and-correction',
    group: 'Claim',
    title: 'Reversal and correction',
    file: reversalAndCorrection,
    open: ['ISA.GS.ST.LX#0.CLP#1'],
  },
  {
    id: 'medicare-secondary',
    group: 'Claim',
    title: 'Medicare secondary',
    file: medicareSecondary,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'paid-down-a-payer-chain',
    group: 'Claim',
    title: 'Paid down a payer chain',
    file: paidDownAPayerChain,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'crossover-to-a-supplement',
    group: 'Claim',
    title: 'Crossover to a supplement',
    file: crossoverToASupplement,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'inpatient',
    group: 'Claim',
    title: 'Inpatient',
    file: inpatient,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'patient-name-corrected',
    group: 'Claim',
    title: 'Patient name corrected',
    file: patientNameCorrected,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'interest-paid',
    group: 'Claim',
    title: 'Interest paid',
    file: interestPaid,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'dme-rental',
    group: 'Claim',
    title: 'DME rental',
    file: dmeRental,
    open: ['ISA.GS.ST.LX#0.CLP#0'],
  },
  {
    id: 'two-claims',
    group: 'Claim',
    title: 'Two paid claims, one with three service lines',
    file: twoClaims,
    open: ['ISA.GS.ST.N1-PR', 'ISA.GS.ST.N1-PE', 'ISA.GS.ST.LX#0'],
  },
]
