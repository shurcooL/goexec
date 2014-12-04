// A command line tool to execute Go functions. The output is printed as goons to stdout.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shurcooL/go/gists/gist5892738"
	goimports "golang.org/x/tools/imports"

	// We need go-goon to be available; this ensures getting goe will get go-goon too.
	_ "github.com/shurcooL/go-goon"
)

func run(src string) error {
	// Create a temp folder.
	tempDir, err := ioutil.TempDir("", "goe_")
	if err != nil {
		return err
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
		return err
	}

	// Compile and run the program.
	cmd := exec.Command("go", "run", "-a", tempFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		cmd += "(" + gist5892738.TrimLastNewline(string(stdin)) + ")"
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
		if err != nil {
			fmt.Print("gen.go:", err, "\n") // No space after colon so the ouput is like "gen.go:8:18: expected ...".
			os.Exit(1)
		}
		src = string(out)
	}

	if *nFlag == true {
		fmt.Print(src)
		return
	}

	// Run the program.
	err := run(src)
	if err != nil {
		fmt.Println("### Error ###")
		fmt.Println(err)
	}
}
