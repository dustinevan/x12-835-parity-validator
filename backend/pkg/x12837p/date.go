package x12837p

// Date is a calendar date parsed from an X12 CCYYMMDD element (for example
// BHT04, GS04, or a DTP03 sent with format qualifier D8). Set is false when
// the element was not sent, in which case Year, Month, and Day are
// meaningless; Go's own zero date (0000-00-00) is not a valid X12 date, so
// it cannot stand in for absence. Shape copied from backend/pkg/x12835.
type Date struct {
	Year  int  `json:"year,omitempty"`
	Month int  `json:"month,omitempty"`
	Day   int  `json:"day,omitempty"`
	Set   bool `json:"set,omitempty"`
}

// Time is a clock time parsed from an X12 HHMM or HHMMSS element (for
// example BHT05, GS05, or ISA10). Second is zero when the element carried
// only hours and minutes, and HasSeconds says which form the file sent, so
// the element writes back with the same width. Set is false when the
// element was not sent. Shape copied from backend/pkg/x12835, plus
// HasSeconds.
type Time struct {
	Hour       int  `json:"hour,omitempty"`
	Minute     int  `json:"minute,omitempty"`
	Second     int  `json:"second,omitempty"`
	HasSeconds bool `json:"hasSeconds,omitempty"`
	Set        bool `json:"set,omitempty"`
}
