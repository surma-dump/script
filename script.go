package script

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Environment map[string]string

// Get emulates shell behavior where undefined variables evaluate to
// empty string.
func (env Environment) Get(key string) string {
	if v, ok := env[key]; ok {
		return v
	}
	return ""
}

const (
	_VARIABLE_NAME = `([A-Za-z_][A-Za-z0-9_]*)`
	_VARIABLE      = `\$(?:\{` + _VARIABLE_NAME + `\}|` + _VARIABLE_NAME + `)`
)

var (
	variableMatcher = regexp.MustCompile(_VARIABLE)
)

// Replaces occurences of "$varname" or "${varname}" with
// the context's environment value for that variable.
func (env Environment) Expand(line string) string {
	loc := variableMatcher.FindAllStringSubmatchIndex(line, -1)
	if loc == nil {
		return line
	}
	ret := ""
	lastidx := 0
	for _, submatch := range loc {
		if submatch[2] == -1 {
			submatch[2], submatch[3] = submatch[4], submatch[5]
		}
		if submatch[2] == -1 {
			continue
		}
		varname := line[submatch[2]:submatch[3]]
		ret += line[lastidx:submatch[0]]
		ret += env.Get(varname)
		lastidx = submatch[1]
	}
	ret += line[lastidx:]
	return ret
}

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
func NewContext(root string) (*Context, error) {
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
	_UNQUOTED_PARAM = `([^"\s\\]+)`
	_PARAM          = `^\s*(?:` + _UNQUOTED_PARAM + `|` + _QUOTED_PARAM + `)`
)

var (
	paramMatcher     = regexp.MustCompile(_PARAM)
	ErrMalformedLine = fmt.Errorf("Malformed line")
)

// Parses a command string (environment variable expansion and splitting)
// and returns the array of parameters which can be passed to os/exec.Command().
func (ctx *Context) Parse(line string) ([]string, error) {
	line = ctx.Environment.Expand(line)
	return ShellSplit(line)
}

// Splits a line into parameters according to shell behavior.
func ShellSplit(line string) ([]string, error) {
	r := make([]string, 0, 5)
	for len(strings.TrimSpace(line)) > 0 {
		loc := paramMatcher.FindStringSubmatchIndex(line)
		if loc == nil {
			return nil, ErrMalformedLine
		}
		// If it isn't a unquoted parameter, use the next
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
