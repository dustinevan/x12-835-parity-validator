package x12837i

// Date is a calendar date parsed from an X12 CCYYMMDD element (for example
// BHT04, DTP03, or a DMG02 birth date). Set is false when the element was
// not sent, in which case Year, Month, and Day are meaningless; Go's own
// zero date (0000-00-00) is not a valid X12 date, so it cannot stand in for
// absence. Copied in shape from backend/pkg/x12835/date.go.
type Date struct {
	Year  int  `json:"year,omitempty"`
	Month int  `json:"month,omitempty"`
	Day   int  `json:"day,omitempty"`
	Set   bool `json:"set,omitempty"`
}

// Time is a clock time parsed from an X12 HHMM or HHMMSS element (for
// example BHT05 or ISA10). Second is zero when the element carried only
// hours and minutes, and WithSeconds says which of the two forms the file
// sent, so the time is written back in the same form. Set is false when the
// element was not sent. Copied in shape from backend/pkg/x12835/date.go.
type Time struct {
	Hour        int  `json:"hour,omitempty"`
	Minute      int  `json:"minute,omitempty"`
	Second      int  `json:"second,omitempty"`
	WithSeconds bool `json:"withSeconds,omitempty"` // the element was HHMMSS, not HHMM
	Set         bool `json:"set,omitempty"`
}
