// Package help holds what every vdb command's help shares: the check for a
// help argument, the usage error, the section layout, and the colors.
// Help is printed straight to stdout, never through a pager. The sections
// run prose first (DESCRIPTION, HOW IT WORKS, DEFAULTS), then OPTIONS,
// USAGE, and EXAMPLES last, so the bottom of the help, nearest the prompt,
// is what the reader types from. Each command's own text lives beside its
// flags in the command's package.
package help

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrUsage is returned once a usage error and the help have been printed,
// so main exits non-zero without printing anything more.
var ErrUsage = errors.New("usage")

// Local is the description of the -local flag the pipeline commands share.
const Local = "use the local development database (localhost:5432/validator, TLS off) with S3 off"

// The colors for section headings: a few ANSI sequences, used only when
// stdout is a terminal.
const (
	ansiBold  = "\033[1m"
	ansiCyan  = "\033[36m"
	ansiReset = "\033[0m"
)

// IsHelp reports whether arg asks for help.
func IsHelp(arg string) bool {
	switch arg {
	case "help", "-help", "--help", "-h":
		return true
	}
	return false
}

// stdoutIsTerminal reports whether stdout is a terminal (a character
// device), so headings are colored only where the colors will be seen.
func stdoutIsTerminal() bool {
	info, err := os.Stdout.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// Section is one block of help: a heading and its body, already indented
// and ending in a newline.
type Section struct {
	Heading string
	Body    string
}

// Render joins a title line and its sections, coloring the headings when
// stdout is a terminal.
func Render(title string, sections []Section) string {
	color := stdoutIsTerminal()
	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n")
	for _, s := range sections {
		b.WriteString("\n")
		if color {
			b.WriteString(ansiBold + ansiCyan + s.Heading + ansiReset)
		} else {
			b.WriteString(s.Heading)
		}
		b.WriteString("\n")
		b.WriteString(s.Body)
	}
	return b.String()
}

// UsageError prints a usage error and the help it belongs with, and
// returns ErrUsage.
func UsageError(message, text string) error {
	fmt.Fprintf(os.Stderr, "vdb: %s\n\n", message)
	fmt.Print(text)
	return ErrUsage
}
