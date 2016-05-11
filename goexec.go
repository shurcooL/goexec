// goexec is a command line tool to execute Go code. Output is printed as goons to stdout.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shurcooL/go/trim"
	"golang.org/x/tools/imports"

	// We need go-goon to be available; this ensures getting goexec will get go-goon too.
	_ "github.com/shurcooL/go-goon"
)

var (
	quietFlag    = flag.Bool("quiet", false, "Do not dump the return values as a goon.")
	stdinFlag    = flag.Bool("stdin", false, "Read func parameters from stdin instead.")
	nFlag        = flag.Bool("n", false, "Print the generated source but do not run it.")
	compilerFlag = flag.String("compiler", "gc", `Compiler to use, one of: "gc", "gopherjs".`)
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: goexec [flags] [packages] [package.]function(parameters)
       echo parameters | goexec --stdin [flags] [packages] [package.]function`)
	flag.PrintDefaults()
}

func run(src string) error {
	// Create a temp folder.
	tempDir, err := ioutil.TempDir("", "goexec_")
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
	var cmd *exec.Cmd
	switch *compilerFlag {
	case "gc":
		cmd = exec.Command("go", "run", tempFile)
	case "gopherjs":
		cmd = exec.Command("gopherjs", "run", tempFile)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
		return
	}

	switch *compilerFlag {
	case "gc", "gopherjs":
	default:
		flag.Usage()
		os.Exit(2)
		return
	}

	args := flag.Args()
	importPaths := args[:len(args)-1] // All but last.
	cmd := args[len(args)-1]          // Last one.
	if *stdinFlag {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		cmd += "(" + trim.LastNewline(string(stdin)) + ")"
	}
	if !*quietFlag {
		cmd = "goon.Dump(" + cmd + ")"
	}

	// Generate source code.
	src := `package main

import (
`
	if !*quietFlag {
		src += `	"github.com/shurcooL/go-goon"
`
	}
	for _, importPath := range importPaths {
		src += `	. "` + importPath + `"
`
	}
	src += `)

func main() {
	` + cmd + `
}
`

	// Run `goimports` on the source code.
	{
		out, err := imports.Process("", []byte(src), nil)
		if err != nil {
			fmt.Print(src)
			fmt.Print("imports.Process: gen.go:", err, "\n") // No space after colon so the ouput is like "gen.go:8:18: expected ...".
			os.Exit(1)
		}
		src = string(out)
	}

	if *nFlag {
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
