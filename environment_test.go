package script

import (
	"fmt"
	"testing"
)

func ExampleShellSplit() {
	var params []string
	params, _ = ShellSplit(`echo a b c`)
	fmt.Printf("%#v\n", params)
	// Excessive spacing, tabbing, quoting and escaping
	params, _ = ShellSplit(`    echo      "longer parameter 1"` +
		`		small\ param2          "another long \"weird\" parameter"    `)
	fmt.Printf("%#v\n", params)
	// Output:
	// []string{"echo", "a", "b", "c"}
	// []string{"echo", "longer parameter 1", "small param2", "another long \"weird\" parameter"}
}

func ExampleEnvironment_Expand() {
	env := Environment{
		"A":          "Variable A",
		"variable_b": "Variable B",
	}
	fmt.Printf("%#v\n", env.Expand("Noise $A Noise"))
	fmt.Printf("%#v\n", env.Expand("Noise ${VARIABLE_B} Noise ${variable_b}"))
	// Output:
	// "Noise Variable A Noise"
	// "Noise  Noise Variable B"
}

func TestEnvironment_Parse(t *testing.T) {
	ctx, e := NewContext(".")
	if e != nil {
		t.Errorf("Creating context failed: %s", e)
	}
	ctx.Environment["ROOT"] = ctx.Root
	ctx.Environment["HOME"] = "tree_test"

	cmd, e := ctx.Parse("rm -rf ${ROOT}/${HOME}")
	if e != nil {
		t.Fatalf("Could not parse line: %s", e)
	}
	got := fmt.Sprintf("%#v", cmd)
	expected := `[]string{"rm", "-rf", "` + ctx.Root + `/tree_test"}`
	if got != expected {
		t.Fatalf("Expected \"%s\", got \"%s\"", expected, got)
	}
}
