package script

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// A context mimicks a separate process environment for execution:
// bound to a directory path with its own set of environment variables.
type Context struct {
	Root string
	Environment
}

var (
	ErrNoDir = fmt.Errorf("Contexts can only be created on directries")
)

// NewContext creates a new context bound to `path` and an empty set of
// environment variables.
func NewContext(rootf string, v ...interface{}) (*Context, error) {
	root := fmt.Sprintf(rootf, v...)
	root = filepath.Clean(root)
	root, e := filepath.Abs(root)
	if e != nil {
		return nil, e
	}
	fi, e := os.Stat(root)
	if e != nil {
		return nil, e
	}
	if !fi.IsDir() {
		return nil, ErrNoDir
	}
	return &Context{
		Root:        root,
		Environment: Environment{},
	}, nil
}

const (
	_QUOTED_PARAM   = `"((?:\\.|[^"\\])+)"`
	_UNQUOTED_PARAM = `((?:\\.|[^"\s\\])+)`
	_PARAM          = `^\s*(?:` + _UNQUOTED_PARAM + `|` + _QUOTED_PARAM + `)`
)

var (
	paramMatcher     = regexp.MustCompile(_PARAM)
	ErrMalformedLine = fmt.Errorf("Malformed line")
)

// Parses a command string (environment variable expansion and splitting)
// and returns the array of parameters which can be passed to os/exec.Command().
func (ctx *Context) Parse(linef string, v ...interface{}) ([]string, error) {
	line := fmt.Sprintf(linef, v...)
	line = ctx.Environment.Expand(line)
	return ShellSplit(line)
}

// Glob is a convenience wrapper arount path/filepath.Glob.
// The matching is done in the current context unless the given path
// is absolute.
func (ctx *Context) Glob(globf string, v ...interface{}) ([]File, error) {
	glob := fmt.Sprintf(globf, v...)
	if !filepath.IsAbs(glob) {
		glob = filepath.Join(ctx.Root, glob)
	}
	matches, e := filepath.Glob(glob)
	if e != nil {
		return nil, e
	}
	ret := make([]File, len(matches))
	for i, match := range matches {
		ret[i] = File(match)
	}
	return ret, nil
}

// Splits a line into parameters similar to shell behavior.
func ShellSplit(line string) ([]string, error) {
	r := make([]string, 0, 5)
	for len(strings.TrimSpace(line)) > 0 {
		loc := paramMatcher.FindStringSubmatchIndex(line)
		if loc == nil {
			return nil, ErrMalformedLine
		}
		// If it isn't an unquoted parameter, use the next
		// match which is a quoted parameter without the quotes
		if loc[2] == -1 {
			loc[2], loc[3] = loc[4], loc[5]
		}
		param := line[loc[2]:loc[3]]
		param = strings.Replace(strings.TrimSpace(param), "\\", "", -1)
		r = append(r, param)
		line = line[loc[1]:]
	}
	return r, nil
}
