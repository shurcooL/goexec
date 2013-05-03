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

echo parameters | goe [package ...] function

goe [package ...] function1 | goe [package ...] function2
```

Example
-----
```bash
$ goe strings 'Repeat("Go! ", 5)'
string("Go! Go! Go! Go! Go!")

$ goe strings 'Replace("Calling Go functions from the terminal is hard.", "hard", "easy", -1)'
string("Calling Go functions from the terminal is easy.")

# Dumps the doc.Package struct for "fmt".
$ echo "fmt" | goe gist.github.com/5504644.git GetDocPackage
(*doc.Package)(...)

# Note that regexp.Compile returns 2 values (*regexp.Regexp, error)
$ echo "Go+" | goe regexp Compile
(*regexp.Regexp)(...),
(interface{})(nil)
```
