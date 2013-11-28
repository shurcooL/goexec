goe
===

A command line tool to execute Go functions. The output is printed as goons to stdout.

Installation
------------
```bash
$ go get -u github.com/shurcooL/goe
```

Add `$GOPATH/bin` to your `$PATH` or copy `$GOPATH/bin/goe` to your `$PATH`.

Usage
-----
```bash
$ goe [--quiet] [package ...] package.function(parameters)
```

Examples
--------
```bash
$ goe 'strings.Repeat("Go! ", 5)'
(string)("Go! Go! Go! Go! Go! ")

$ goe 'strings.Replace("Calling Go functions from the terminal is hard.", "hard", "easy", -1)'
(string)("Calling Go functions from the terminal is easy.")

# Dumps the doc.Package struct for "fmt"
$ goe gist.github.com/5504644.git 'GetDocPackage(BuildPackageFromImportPath("fmt"))'
(*doc.Package)(...)

$ goe 'gist4727543.GetForcedUse("fmt")'
(string)("var _ = fmt.Errorf")

# Note that parser.ParseExpr returns 2 values (ast.Expr, error)
$ goe 'parser.ParseExpr("5 + 7")'
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

$ goe --quiet 'fmt.Println("Use --quiet to disable output of goon; useful if you want to print to stdout.")'
Use --quiet to disable output of goon; useful if you want to print to stdout.
```
