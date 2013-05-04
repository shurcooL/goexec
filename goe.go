package main

import (
	"fmt"
	. "gist.github.com/5498057.git"
	"github.com/shurcooL/go-goon"
	"github.com/sriram-srinivasan/gore/eval"
	"os"
	"strings"
)

var _ = goon.Dump

func main() {
	//os.Args = []string{`./goe`, `strings`, `Repeat("Go! ", 5)`}
	//os.Args = []string{`./goe`, `strings`, `Replace("Calling Go functions from the terminal is hard.", "hard", "easy", -1)`}
	//os.Args = []string{`./goe`, `strings`, `Repeat`}
	//os.Args = []string{`./goe`, `gist.github.com/4727543.git`, `GetForcedUse`}
	//os.Args = []string{`./goe`, `regexp`, `Compile`}

	if len(os.Args) <= 1 {
		fmt.Println("Usage: ...")
		return
	}

	x := os.Args[1:]
	imports := x[:len(x)-1]
	//goon.Dump(imports)
	cmd := x[len(x)-1]
	//goon.Dump(cmd)

	src := "import \"github.com/shurcooL/go-goon\"\n" // TODO: Check that it hasn't been imported manually
	for _, importPath := range imports {
		src += "import . \"" + importPath + "\"\n"
	}
	if -1 == strings.Index(cmd, "(") {
		cmd += "(" + ReadAllStdin() + ")"
	}
	src += "goon.Dump(" + cmd + ")" // TODO: DumpComma

	//print(src); return
	out, err := eval.Eval(src)

	if err == "" {
		fmt.Print(out)
	} else {
		fmt.Println("===== Error =====")
		fmt.Println(err)
	}
}