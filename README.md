goe
===

A tool to execute Go functions. The output is printed as comma separated goons to stdout.

Installation
------------
```bash
$ go get -u github.com/shurcooL/goe
```

Add `$GOPATH/bin` to your `$PATH`.

Usage
-----
```
goe [package ...] function(parameters)
```

Example
-----
```bash
$ goe strings 'Repeat("Go! ", 5)'
string("Go! Go! Go! Go! Go!")

$ goe strings 'Replace("Calling Go functions from the terminal is hard.", "hard", "easy", -1)'
string("Calling Go functions from the terminal is easy.")

$ goe gist.github.com/5504644.git 'GetDocPackage("fmt")'
# Dumps the doc.Package struct for "fmt".
```
