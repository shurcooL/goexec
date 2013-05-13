package main

import (
	"fmt"
	. "gist.github.com/5498057.git"
	"github.com/shurcooL/go-goon"
	. "gist.github.com/5286084.git"
	"os"
	"strings"
	"os/exec"
	"path"
	"io/ioutil"
)

var _ = goon.Dump

func run(src string) (output string, error string) {
	tempDir, err := ioutil.TempDir("", "goe_")
	CheckError(err)

	tempFile := path.Join(tempDir, "gen.go")
	err = ioutil.WriteFile(tempFile, []byte(src), 0600)
	CheckError(err)

	cmd := exec.Command("go", "run", tempFile)
	out, err := cmd.CombinedOutput()

	err = os.RemoveAll(tempDir)
	CheckError(err)

	if nil == err {
		return string(out), ""
	} else {
		return "", string(out)
	}
}

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

	Args := os.Args[1:]
	imports := Args[:len(Args)-1]		// All but last
	//goon.Dump(imports)
	cmd := Args[len(Args)-1]			// Last one
	//goon.Dump(cmd)

	src := "package main\n\nimport (\n\t\"github.com/shurcooL/go-goon\"\n" // TODO: Check that it hasn't been imported manually
	for _, importPath := range imports {
		src += "\t. \"" + importPath + "\"\n"
	}
	if -1 == strings.Index(cmd, "(") {		// BUG: What if the bracket is a part of a comment or a string...
		cmd += "(" + ReadAllStdin() + ")"
	}
	src += ")\n\nfunc main() {\n\t"
	src += "goon.Dump(" + cmd + ")" // TODO: DumpComma
	src += "\n}"

	//print(src); return
	//out, err := eval.Eval(src)
	out, err := run(src)

	if err == "" {
		fmt.Print(out)
	} else {
		fmt.Println("===== Error =====")
		fmt.Println(err)
	}
}