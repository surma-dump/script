package script

import (
	"fmt"
)

func ExampleFile() {
	f := File("/tmp/backups and stuff/dungeons and dragons.tar.gz")
	fmt.Printf("%#v\n", f.Path())
	fmt.Printf("%#v\n", f.Fullname())
	fmt.Printf("%#v\n", f.Name())
	fmt.Printf("%#v\n", f.Extension())

	fmt.Printf("%#v\n", f.ShellPath())
	fmt.Printf("%#v\n", f.ShellFullname())
	fmt.Printf("%#v\n", f.ShellName())
	fmt.Printf("%#v\n", f.ShellExtension())

	// Output:
	// "/tmp/backups and stuff/dungeons and dragons.tar.gz"
	// "dungeons and dragons.tar.gz"
	// "dungeons and dragons.tar"
	// ".gz"

	// "/tmp/backups\\ and\\ stuff/dungeons\\ and\\ dragons.tar.gz"
	// "dungeons\\ and\\ dragons.tar.gz"
	// "dungeons\\ and\\ dragons.tar"
	// ".gz"
}
