package script

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Environment map[string]string

var (
	DefaultEnvironment = Environment{}
)

func init() {
	for _, env := range os.Environ() {
		i := strings.Index(env, "=")
		if i == -1 {
			continue
		}
		DefaultEnvironment[env[0:i]] = env[i+1:]
	}
}

// Add adds (or overwrites) all variables from otherenv to this environment.
func (env Environment) Add(otherenv Environment) {
	for k, v := range otherenv {
		env[k] = v
	}
}

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

func (env Environment) Array() []string {
	r := make([]string, len(env))
	for k, v := range env {
		r = append(r, fmt.Sprintf("%s=%s", k, v))
	}
	return r
}
