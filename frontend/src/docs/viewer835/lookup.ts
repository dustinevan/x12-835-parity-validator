import { CODES_835, ELEMENTS_835 } from '../../data/lookupRows835'
import { SEGMENT_NAME } from '../../data/x12Names835'
import type { LookupData } from '../viewer/LookupFrame'

// The 835 lookup panel's rows: the one panel the 835 viewer, the 835 compare
// page, and the Compared fields page open beside their content.
export const LOOKUP_835: LookupData = {
  title: '835 field lookup',
  elements: ELEMENTS_835,
  codes: CODES_835,
  segmentNames: SEGMENT_NAME,
}
