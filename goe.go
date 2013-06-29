package main

import (
	"fmt"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5498057.git"
	. "gist.github.com/5892738.git"
	"github.com/shurcooL/go-goon"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var _ = goon.Dump

func run(src string) (output string, error string) {
	tempDir, err := ioutil.TempDir("", "goe_")
	CheckError(err)

	tempFile := path.Join(tempDir, "gen.go")
	err = ioutil.WriteFile(tempFile, []byte(src), 0600)
	CheckError(err)

	cmd := exec.Command("go", "run", tempFile)
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	err = os.RemoveAll(tempDir)
	CheckError(err)

	if nil == err {
		return string(out), ""
	} else {
		return "", string(out)
	}
}

func usage() {
	fmt.Println("Usage: ...")
}

func main() {
	//os.Args = []string{`./goe`, `strings`, `Repeat("Go! ", 5)`}
	//os.Args = []string{`./goe`, `strings`, `Replace("Calling Go functions from the terminal is hard.", "hard", "easy", -1)`}
	//os.Args = []string{`./goe`, `strings`, `Repeat`}
	//os.Args = []string{`./goe`, `gist.github.com/4727543.git`, `GetForcedUse`}
	//os.Args = []string{`./goe`, `regexp`, `Compile`}

	type OutputType int
	const (
		Quiet = iota
		Goon
	)
	Output := Goon

	Args := os.Args[1:]

	if len(Args) < 1 {
		usage()
		return
	}

	if "--quiet" == Args[0] {
		Output = Quiet
		Args = Args[1:]
	}

	if len(Args) < 1 {
		usage()
		return
	}

	imports := Args[:len(Args)-1] // All but last
	//goon.Dump(imports)
	cmd := Args[len(Args)-1] // Last one
	//goon.Dump(cmd)

	src := "package main\n\nimport (\n"
	if Goon == Output {
		// TODO: Check that it hasn't already been imported
		src += "\t\"github.com/shurcooL/go-goon\"\n"
	} else if Quiet == Output {
	}
	for _, importPath := range imports {
		src += "\t. \"" + importPath + "\"\n"
	}
	if -1 == strings.Index(cmd, "(") { // BUG: What if the bracket is a part of a comment or a string...
		cmd += "(" + TrimLastNewline(ReadAllStdin()) + ")"
	}
	src += ")\n\nfunc main() {\n\t"
	if Goon == Output {
		src += "goon.Dump(" + cmd + ")" // TODO: DumpComma
	} else if Quiet == Output {
		src += cmd
	}
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