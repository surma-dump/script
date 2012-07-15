package script

import (
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
