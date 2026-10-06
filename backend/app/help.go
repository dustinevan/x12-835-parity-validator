// The only arguments validator takes: help. validator alone starts the
// service; validator help prints what the program is and the APP_SERVER_*
// settings it reads, straight to stdout, before the settings are parsed.
package main

import (
	"errors"
	"strings"
)

// errUsage is returned once an unknown argument and the help have been
// printed, so main exits non-zero without logging anything more.
var errUsage = errors.New("usage")

// isHelp reports whether arg asks for help.
func isHelp(arg string) bool {
	switch arg {
	case "help", "-help", "--help", "-h":
		return true
	}
	return false
}

// helpText is the two-line description and the settings. usage is what
// conf formats from the configuration struct; only its ENVIRONMENT section
// is shown, since the --flags it also lists are never read: the service's
// only arguments are help.
func helpText(usage string) string {
	_, env, found := strings.Cut(usage, "ENVIRONMENT\n")
	if !found {
		env = usage
	}
	return "validator: the X12 835 Parity Validator service\n" +
		"\n" +
		"DESCRIPTION\n" +
		"  The web app and its API on one port (APP_SERVER_WEB_API_HOST, default 127.0.0.1:8089); validator alone starts it.\n" +
		"  Its settings come from APP_SERVER_* environment variables, listed below with their defaults; it takes no flags.\n" +
		"\n" +
		"SETTINGS\n" +
		env
}
