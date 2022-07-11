goexec
======

[![Go Reference](https://pkg.go.dev/badge/github.com/shurcooL/goexec.svg)](https://pkg.go.dev/github.com/shurcooL/goexec)

goexec is a command line tool to execute Go code. Output is printed as goons to stdout.

Installation
------------

```sh
go install github.com/shurcooL/goexec@latest
```

Usage
-----

```
Usage: goexec [flags] [packages] [package.]function(parameters)
       echo parameters | goexec -stdin [flags] [packages] [package.]function
  -compiler string
    	Compiler to use, one of: "gc", "gopherjs". (default "gc")
  -n	Print the generated source but do not run it.
  -quiet
    	Do not dump the return values as a goon.
  -stdin
    	Read func parameters from stdin instead.
  -tags string
    	A comma-separated list of build tags to consider satisfied during the build.
```

Examples
--------

```sh
$ goexec 'strings.Repeat("Go! ", 5)'
(string)("Go! Go! Go! Go! Go! ")

$ goexec strings 'Replace("Calling Go functions from the terminal is hard.", "hard", "easy", -1)'
(string)("Calling Go functions from the terminal is easy.")

# Note that parser.ParseExpr returns 2 values (ast.Expr, error).
$ goexec 'parser.ParseExpr("5 + 7")'
(*ast.BinaryExpr)(&ast.BinaryExpr{
	X: (*ast.BasicLit)(&ast.BasicLit{
		ValuePos: (token.Pos)(1),
		Kind:     (token.Token)(5),
		Value:    (string)("5"),
	}),
	OpPos: (token.Pos)(3),
	Op:    (token.Token)(12),
	Y: (*ast.BasicLit)(&ast.BasicLit{
		ValuePos: (token.Pos)(5),
		Kind:     (token.Token)(5),
		Value:    (string)("7"),
	}),
})
(interface{})(nil)

# Run function RepoRootForImportPath from package "golang.org/x/tools/go/vcs".
$ goexec 'vcs.RepoRootForImportPath("rsc.io/pdf", false)'
(*vcs.RepoRoot)(...)
(interface{})(nil)

$ goexec -quiet 'fmt.Println("Use -quiet to disable output of goon; useful if you want to print to stdout.")'
Use -quiet to disable output of goon; useful if you want to print to stdout.

$ echo '"fmt"' | goexec -stdin 'gist4727543.GetForcedUse'
(string)("var _ = fmt.Errorf")
```

Alternatives
------------

-	[gommand](https://github.com/sno6/gommand) - Go one liner program, similar to python -c.
-	[gorram](https://github.com/natefinch/gorram) - Like go run for any Go function.
-	[goeval](https://github.com/dolmen-go/goeval) - Run Go snippets instantly from the command-line.

License
-------

-	[MIT License](LICENSE)
