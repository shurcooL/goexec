// A command line tool to execute Go functions. The output is printed as goons to stdout.
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
	"github.com/shurcooL/go/gists/gist5498057"
	"github.com/shurcooL/go/gists/gist5892738"

	// We need go-goon to be available; this ensures getting goe will get go-goon too.
	_ "github.com/shurcooL/go-goon"
)

func run(src string) (output string, err error) {
	// Create a temp folder.
	tempDir, err := ioutil.TempDir("", "goe_")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "warning: error removing temp dir:", err)
		}
	}()

	// Write the source code file.
	tempFile := filepath.Join(tempDir, "gen.go")
	err = ioutil.WriteFile(tempFile, []byte(src), 0600)
	if err != nil {
		panic(err)
	}

	// Compile and run the program.
	cmd := exec.Command("go", "run", "-a", tempFile)
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	if nil == err {
		return string(out), nil
	} else {
		return "", errors.New(string(out))
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: goe [flags] [packages] [package.]function(parameters)
       echo parameters | goe --stdin [flags] [packages] [package.]function`)
	flag.PrintDefaults()
	os.Exit(2)
}

var quietFlag = flag.Bool("quiet", false, "Do not dump the return values as a goon.")
var stdinFlag = flag.Bool("stdin", false, "Read func parameters from stdin instead.")
var nFlag = flag.Bool("n", false, "Print the generated source but do not run it.")

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		return
	}

	args := flag.Args()
	imports := args[:len(args)-1] // All but last.
	cmd := args[len(args)-1]      // Last one.
	if *stdinFlag {
		cmd += "(" + gist5892738.TrimLastNewline(gist5498057.ReadAllStdin()) + ")"
	}
	if false == *quietFlag {
		cmd = "goon.Dump(" + cmd + ")"
	}

	// Generate source code.
	src := "package main\n\nimport (\n"
	if *quietFlag == false {
		src += "\t\"github.com/shurcooL/go-goon\"\n"
	}
	for _, importPath := range imports {
		src += "\t. \"" + importPath + "\"\n"
	}
	src += ")\n\nfunc main() {\n\t" + cmd + "\n}"

	// Run `goimports` on the source code.
	{
		out, err := goimports.Process("", []byte(src), nil)
		if err == nil {
			src = string(out)
		} else {
			panic(err)
		}
	}

	if *nFlag == true {
		fmt.Print(src)
		return
	}

	// Run the program and get its output.
	output, err := run(src)

	if err == nil {
		fmt.Print(output)
	} else {
		fmt.Println("===== Error =====")
		fmt.Println(err.Error())
	}
}
