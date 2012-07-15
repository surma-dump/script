package script

import (
	"path/filepath"
	"strings"
)

// File unambiguously describes a file and contains
// the absolute path to that file.
// Methods starting with `Shell` escape the resulting
// string for safe shell usage.
type File string

// NewFile create a new file. Relative paths will
// be converted to absolute paths.
// Error is only non-nil if that conversion fails.
func NewFile(path string) (File, error) {
	f, e := filepath.Abs(path)
	return File(f), e
}

// Path returns the absolute path to the file.
func (f File) Path() string {
	return string(f)
}

// Fullname returns the filename include extension
func (f File) Fullname() string {
	_, name := filepath.Split(string(f))
	return name
}

// Name returns the filename without extension
// What the extension is is defined in the path/filepath package.
func (f File) Name() string {
	ext := filepath.Ext(string(f))
	fullname := f.Fullname()
	return fullname[0 : len(fullname)-len(ext)]
}

// Extension returns the extension of the file.
// What the extension is is defined in the path/filepath package.
func (f File) Extension() string {
	ext := filepath.Ext(string(f))
	return ext
}

func (f File) ShellPath() string {
	return ShellEscape(f.Path())
}

func (f File) ShellFullname() string {
	return ShellEscape(f.Fullname())
}

func (f File) ShellName() string {
	return ShellEscape(f.Name())
}

func (f File) ShellExtension() string {
	return ShellEscape(f.Extension())
}

var (
	shellescapes = []string{
		" ", "\\ ",
	}
	shellescaper = strings.NewReplacer(shellescapes...)
)

func ShellEscape(s string) string {
	return shellescaper.Replace(s)
}
