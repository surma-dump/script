package script

import (
	"strings"
	"testing"
)

func TestContext_Glob(t *testing.T) {
	ctx, e := NewContext("./data_test")
	if e != nil {
		t.Fatalf("Creating context failed: %s", e)
	}
	g, e := ctx.Glob("a*")
	if e != nil {
		t.Fatalf("Globbing failed: %s", e)
	}
	if len(g) != 4 {
		t.Fatalf("Expected to find 4 files, found %d", len(g))
	}
}

func TestContext_Run(t *testing.T) {
	ctx, e := NewContext("./data_test/")
	if e != nil {
		t.Fatalf("Creating context failed: %s", e)
	}
	ctx.Environment.Add(DefaultEnvironment)
	ctx.Environment["PARAMS"] = "-l"
	data, e := ctx.Run("ls ${PARAMS} %s*", "a")
	if e != nil {
		t.Fatalf("Could not run command: %s", e)
	}
	if len(strings.Split(data, "\n")) != 4+1 {
		t.Fatalf("Unexpected number of result lines")
	}
}

func TestContext_Parse(t *testing.T) {
	ctx, e := NewContext(".")
	if e != nil {
		t.Errorf("Creating context failed: %s", e)
	}
	ctx.Environment["ROOT"] = ctx.Root
	ctx.Environment["HOME"] = "data_test"

	cmd, e := ctx.Parse("rm -rf ${ROOT}/${HOME}/b*")
	if e != nil {
		t.Fatalf("Could not parse line: %s", e)
	}
	if len(cmd) != 6 || cmd[0] != "rm" || cmd[1] != "-rf" {
		t.Fatalf("Parse failed")
	}
}
