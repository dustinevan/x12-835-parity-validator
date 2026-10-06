package x12835

// Date is a calendar date parsed from an X12 CCYYMMDD element (for example
// DTM02, GS04, or BPR16). Set is false when the element was not sent, in
// which case Year, Month, and Day are meaningless; Go's own zero date
// (0000-00-00) is not a valid X12 date, so it cannot stand in for absence.
type Date struct {
	Year  int  `json:"year,omitempty"`
	Month int  `json:"month,omitempty"`
	Day   int  `json:"day,omitempty"`
	Set   bool `json:"set,omitempty"`
}

// Time is a clock time parsed from an X12 HHMM or HHMMSS element (for
// example GS05 or ISA10). Second is zero when the element carried only
// hours and minutes, and HasSeconds says which form the element took, so
// "0900" and "090000" write back out as they came. Set is false when the
// element was not sent.
type Time struct {
	Hour       int  `json:"hour,omitempty"`
	Minute     int  `json:"minute,omitempty"`
	Second     int  `json:"second,omitempty"`
	HasSeconds bool `json:"hasSeconds,omitempty"`
	Set        bool `json:"set,omitempty"`
}
