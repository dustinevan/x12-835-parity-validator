// Help for codegen, printed straight to stdout, never through a pager.
// Prose first (DESCRIPTION, HOW IT WORKS, DEFAULTS), then OPTIONS, USAGE,
// and EXAMPLES last, so the bottom of the help, nearest the prompt, is what
// the reader types from.
package main

import (
	"os"
	"strings"
)

// The colors for section headings: a few ANSI sequences, used only when
// stdout is a terminal.
const (
	ansiBold  = "\033[1m"
	ansiCyan  = "\033[36m"
	ansiReset = "\033[0m"
)

// isHelp reports whether arg asks for help.
func isHelp(arg string) bool {
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

// section is one block of help: a heading and its body, already indented
// and ending in a newline.
type section struct {
	heading string
	body    string
}

// renderHelp joins a title line and its sections, coloring the headings
// when stdout is a terminal.
func renderHelp(title string, sections []section) string {
	color := stdoutIsTerminal()
	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n")
	for _, s := range sections {
		b.WriteString("\n")
		if color {
			b.WriteString(ansiBold + ansiCyan + s.heading + ansiReset)
		} else {
			b.WriteString(s.heading)
		}
		b.WriteString("\n")
		b.WriteString(s.body)
	}
	return b.String()
}

func helpText() string {
	return renderHelp("codegen: update the CARC and RARC entries tables from their JSON lists", []section{
		{"DESCRIPTION", `  backend/pkg/carc/carc.go and backend/pkg/rarc/rarc.go are written by hand,
  all but their entries tables: the rows between the two lines
  // codegen: entries begin and // codegen: entries end. codegen reads
  backend/pkg/carc.json and backend/pkg/rarc.json and rewrites those rows,
  one per code in the JSON's order, and changes nothing else. Run it after
  refreshing either JSON list from the x12.org page its _meta names.
`},
		{"HOW IT WORKS", `  1. Each JSON list is checked: codes of letters and digits only, no code
     twice, and as many codes as its _meta count says.
  2. Each Go file must hold each marker on exactly one line, begin first;
     otherwise codegen names the file, stops, and writes nothing.
  3. Rows are written in one format every time: Code, Short, Description,
     and Status, then Start, Stop, LastModified, and Notes when the list has
     them. Short is the description up to its first ". " or ";".
  4. A file is written only when its rows change, so a refresh changes the
     rows of the codes that changed and a rerun writes nothing. Each run
     reports the codes added, changed, and removed.
`},
		{"DEFAULTS", `  -root . : the repository root is the working directory. make codes runs
  codegen from tools/codegen with -root ../..
`},
		{"OPTIONS", `  -root DIR   the repository root, where backend/pkg lives
  -dry-run    say what would change and write nothing
`},
		{"USAGE", `  codegen -root DIR [-dry-run]
  codegen help
`},
		{"EXAMPLES", `  codegen -root .                      # at the repository root, after refreshing the JSON lists
  codegen -root . -dry-run             # the codes a refresh would add, change, or remove; writes nothing
  codegen -root /path/to/repository    # from anywhere else
`},
	})
}
