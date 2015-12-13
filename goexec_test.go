package main

func Example_run() {
	const src = `package main

import "github.com/shurcooL/go-goon"

func main() {
	goon.Dump("goexec's run() is working.")
}
`
	err := run(src)
	if err != nil {
		panic(err)
	}

	// Output: (string)("goexec's run() is working.")
}
