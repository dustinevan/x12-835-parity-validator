package x12837p

// ContactPER is the PER contact segment, shared by every contact-carrying
// loop in the 837P: submitter (1000A, SUBMITTER_PER_EL_NAME), billing
// provider (2010AA, BILLING_PROVIDER_PER_NAME), subscriber property and
// casualty contact (2010BA), ordering provider (2420E,
// ORDERING_PROVIDER_PER_NAME), and service facility location (2310C,
// SERVICE_FACILITY_PER_NAME). PER01 and the three communication-qualifier
// positions carry the same element name at every instance; PER02's own
// label (whose name it is) varies by loop - see
// frontend/src/data/x12Names837p.ts for the per-loop text (source:
// Transmittal 907). InquiryReferencePER09 is populated only in the
// submitter loop (1000A); every other loop marks PER09 Not Used.
//
// ContactFunctionCodePER01 and CommunicationNumberQualifierPER03 are
// defined in enums_provider.go.
type ContactPER struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	FunctionPER01 ContactFunctionCodePER01 `json:"functionPER01,omitempty"` // Contact Function Code

	NamePER02 string `json:"namePER02,omitempty"` // contact name

	QualifierPER03 CommunicationNumberQualifierPER03 `json:"qualifierPER03,omitempty"` // Communication Number Qualifier
	CommPER04      string                            `json:"commPER04,omitempty"`      // Communication Number

	QualifierPER05 CommunicationNumberQualifierPER03 `json:"qualifierPER05,omitempty"` // Communication Number Qualifier
	CommPER06      string                            `json:"commPER06,omitempty"`      // Communication Number

	QualifierPER07 CommunicationNumberQualifierPER03 `json:"qualifierPER07,omitempty"` // Communication Number Qualifier
	CommPER08      string                            `json:"commPER08,omitempty"`      // Communication Number

	InquiryReferencePER09 string `json:"inquiryReferencePER09,omitempty"` // Contact Inquiry Reference; submitter loop (1000A) only

	Extra []string `json:"extra,omitempty"` // elements after PER09, kept as sent
}
