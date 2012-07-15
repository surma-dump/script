package script

import (
	"fmt"
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
