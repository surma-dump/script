package script

import (
	"fmt"
	"os"
	"os/exec"
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

var ErrNoDir = fmt.Errorf("contexts can only be created on directries")

func (ctx *Context) buildNewRoot(rootf string, v ...interface{}) (string, error) {
	root := fmt.Sprintf(rootf, v...)
	root = filepath.Clean(root)
	if !filepath.IsAbs(root) {
		root = filepath.Clean(filepath.Join(ctx.Root, root))
	}
	fi, e := os.Stat(root)
	if e != nil {
		return "", e
	}
	if !fi.IsDir() {
		return "", ErrNoDir
	}
	return root, nil
}

// NewContext creates a new context bound to `path` and an empty set of
// environment variables.
func NewContext(rootf string, v ...interface{}) (*Context, error) {
	ctx := &Context{
		Root:        "",
		Environment: Environment{},
	}
	root, e := ctx.buildNewRoot(rootf, v...)
	if e != nil {
		return nil, e
	}
	ctx.Root = root
	return ctx, nil
}

const (
	_QUOTED_PARAM   = `"((?:\\.|[^"\\])+)"`
	_UNQUOTED_PARAM = `((?:\\.|[^"\s\\])+)`
	_PARAM          = `^\s*(?:` + _UNQUOTED_PARAM + `|` + _QUOTED_PARAM + `)`
)

var (
	paramMatcher     = regexp.MustCompile(_PARAM)
	ErrMalformedLine = fmt.Errorf("malformed line")
)

// Parses a command string (environment variable expansion, parameter splitting,
// and globbing) and returns the array of parameters which can be passed
// to os/exec.Command().
func (ctx *Context) Parse(linef string, v ...interface{}) ([]string, error) {
	line := fmt.Sprintf(linef, v...)
	line = ctx.Environment.Expand(line)
	params, e := ShellSplit(line)
	if e != nil {
		return nil, e
	}
	r := make([]string, 0, len(params))
	for _, param := range params {
		globs, e := ctx.Glob(param)
		if e != nil {
			return nil, e
		}
		if len(globs) == 0 {
			// Wasn't globbable. Copy unchanged.
			r = append(r, param)
			continue
		}
		for _, glob := range globs {
			r = append(r, glob.Path())
		}
	}
	return r, nil
}

// Glob is a convenience wrapper arount path/filepath.Glob.
// The matching is done in the current context unless the given path
// is absolute.
func (ctx *Context) Glob(globf string, v ...interface{}) ([]File, error) {
	glob := fmt.Sprintf(globf, v...)
	if !filepath.IsAbs(glob) {
		// If Path is not absolute, temporarily switch
		// to the context's root. Rhyme.
		wd, e := os.Getwd()
		if e != nil {
			return nil, e
		}
		defer os.Chdir(wd)
		e = os.Chdir(ctx.Root)
		if e != nil {
			return nil, e
		}
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

// Changes the root of the context to the given path. If the path is relative,
// it will interpreted relative to the old root.
func (ctx *Context) Cd(pathf string, v ...interface{}) error {
	root, e := ctx.buildNewRoot(pathf, v...)
	if e != nil {
		return e
	}
	ctx.Root = root
	return nil
}

var ErrNotFound = fmt.Errorf("executable file not found in $PATH")

func findExecutable(file string) error {
	d, e := os.Stat(file)
	if e != nil {
		return e
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}

// LookPath is identical to os/exec.LookPath but uses the context's
// $PATH insted of the main process'.
func (ctx *Context) LookPath(file string) (string, error) {
	if strings.Contains(file, "/") {
		e := findExecutable(file)
		if e != nil {
			return file, nil
		}
		return "", e
	}

	pathenv := ctx.Environment.Get("PATH")
	for _, dir := range strings.Split(pathenv, ":") {
		if dir == "" {
			dir = "."
		}
		path := filepath.Join(dir, file)
		if e := findExecutable(path); e == nil {
			return path, nil
		}
	}
	return "", ErrNotFound
}

// Run runs a command
func (ctx *Context) Run(runf string, v ...interface{}) (string, error) {
	run := fmt.Sprintf(runf, v...)
	cmd, e := ctx.Parse(run)
	if e != nil {
		return "", e
	}
	cmdpath, e := ctx.LookPath(cmd[0])
	if e != nil {
		return "", e
	}
	p := &exec.Cmd{
		Path: cmdpath,
		Args: cmd,
		Env:  ctx.Environment.Array(),
		Dir:  ctx.Root,
	}
	data, e := p.CombinedOutput()
	return string(data), e
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
