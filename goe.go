package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	goimports "code.google.com/p/go.tools/imports"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5498057.git"
	. "gist.github.com/5892738.git"

	// We need go-goon to be available; this ensures getting goe will get go-goon too
	_ "github.com/shurcooL/go-goon"
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
	tempFile := filepath.Join(tempDir, "gen.go")
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
	fmt.Fprint(os.Stderr, "Usage: goe [--quiet] [package ...] [package.]function(parameters)\n")
	fmt.Fprint(os.Stderr, "       echo parameters | goe --stdin [--quiet] [package ...] [package.]function\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var quietFlag = flag.Bool("quiet", false, "Do not dump the return values as a goon.")
var stdinFlag = flag.Bool("stdin", false, "Read func parameters from stdin instead.")

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		return
	}

	Args := flag.Args()
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
		out, err := goimports.Process("", []byte(src), nil)
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
