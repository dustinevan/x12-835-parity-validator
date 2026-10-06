package x12835

import "fmt"

// codeEntry is one code an element may carry: the code as the file carries
// it, its public name ("" when the reviewed sources give none), and that
// name with its citation. Each coded element's type keeps one list of
// these beside its constants (enums_*.go); the lists are small enough that
// a lookup reads them in order.
type codeEntry struct {
	code, short, description string
}

// entryOf finds code in list.
func entryOf(list []codeEntry, code string) (codeEntry, bool) {
	for _, e := range list {
		if e.code == code {
			return e, true
		}
	}
	return codeEntry{}, false
}

// known reports whether code is on list. The empty string, an element not
// sent, is never on a list.
func known(list []codeEntry, code string) bool {
	_, ok := entryOf(list, code)
	return ok
}

// shortOf is the public name of code, or "" when it has none or is not on
// list.
func shortOf(list []codeEntry, code string) string {
	e, _ := entryOf(list, code)
	return e.short
}

// descriptionOf is the public name of code with its citation, or "" when
// code is not on list.
func descriptionOf(list []codeEntry, code string) string {
	e, _ := entryOf(list, code)
	return e.description
}

// checkListed returns an error naming element when code is not empty and
// not on list. The caller keeps the value either way: a code outside the
// list is reported, never dropped.
func checkListed(element string, list []codeEntry, code string) error {
	if code == "" || known(list, code) {
		return nil
	}
	return fmt.Errorf("x12835: parsing %s: a value not on the list", element)
}
