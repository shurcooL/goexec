package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	. "gist.github.com/5286084.git"
	. "gist.github.com/5498057.git"
	. "gist.github.com/5892738.git"
	_ "github.com/shurcooL/go-goon" // We need go-goon to be available; this ensures getting goe will get go-goon too
)

func run(src string) (output string, err error) {
	// Create a temp folder
	tempDir, err := ioutil.TempDir("", "goe_")
	CheckError(err)
	defer func() {
		err := os.RemoveAll(tempDir)
		CheckError(err)
	}()

	// Write the source code file
	tempFile := path.Join(tempDir, "gen.go")
	err = ioutil.WriteFile(tempFile, []byte(src), 0600)
	CheckError(err)

	// Compile and run the program
	cmd := exec.Command("go", "run", tempFile)
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	if nil == err {
		return string(out), nil
	} else {
		return "", errors.New(string(out))
	}
}

func usage() {
	fmt.Println("Usage: goe [--quiet] [package ...] [package.]function(parameters)")
	fmt.Println("       echo parameters | goe --stdin [--quiet] [package ...] [package.]function)")
	flag.PrintDefaults()
}

var quietFlag = flag.Bool("quiet", false, "Do not dump the return values as a goon.")
var stdinFlag = flag.Bool("stdin", false, "Read func parameters from stdin instead.")

func main() {
	flag.Parse()

	Args := flag.Args()
	if len(Args) < 1 {
		usage()
		return
	}

	imports := Args[:len(Args)-1] // All but last
	cmd := Args[len(Args)-1]      // Last one
	if *stdinFlag {
		cmd += "(" + TrimLastNewline(ReadAllStdin()) + ")"
	}
	if false == *quietFlag {
		cmd = "goon.Dump(" + cmd + ")"
	}

	// Generate source code
	src := "package main\n\nimport (\n"
	if *quietFlag == false {
		src += "\t\"github.com/shurcooL/go-goon\"\n"
	}
	for _, importPath := range imports {
		src += "\t. \"" + importPath + "\"\n"
	}
	src += ")\n\nfunc main() {\n\t" + cmd + "\n}"

	// Run `goimports` on the source code
	{
		cmd := exec.Command("goimports")
		cmd.Stdin = strings.NewReader(src)
		out, err := cmd.CombinedOutput()
		if err == nil {
			src = string(out)
		} else {
			panic(err)
		}
	}

	// Run the program and get its output
	output, err := run(src)

	if err == nil {
		fmt.Print(output)
	} else {
		fmt.Println("===== Error =====")
		fmt.Println(err.Error())
	}
}
