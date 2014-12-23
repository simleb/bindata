// The bindata command translates binary files into byte arrays in Go source.
//
// bindata is designed to work with go generate, but can be used on its own as well.
//
// The data is stored as a map of byte slices or strings indexed by the
// file paths as specified on the command line. The default name of the
// map is "bindata" but a custom name can be specified on the command line.
//
// By default, the package name of the file containing the generate directive
// is used as the package name of the generated file, or "main" otherwise.
// A custom package name can also be specified on the command line.
//
// The output file can be specified on the command line (with the -o flag).
// By default, the file produced is printed on the standard output.
// If a file already exists at this location, it will be overwritten.
// The file produced is properly formatted and commented.
//
// To see the full list of flags, run:
//  bindata -h
//
// Example using go generate:
//  //go:generate bindata -o jpegs.go pic1.jpg pic2.jpg pic3.jpg
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// tmpl is the template of the generated Go source file.
var tmpl = template.Must(template.New("bindata").Parse(`package {{.Pkg}}

// This file is generated. Do not edit directly.

// {{.Map}} stores binary files as {{if .AsString}}strings{{else}}byte slices{{end}} indexed by filepaths.
var {{.Map}} = map[string]{{if .AsString}}string{{else}}[]byte{{end}}{{"{"}}{{range $name, $data := .Files}}
	{{printf "%#v" $name}}: {{printf "%#v" $data}},{{end}}
}
`))

// vars contains the variables required by the template.
var vars struct {
	Pkg      string
	Map      string
	AsString bool
	Files    map[string]fmt.Formatter
}

func main() {
	if err := run(); err != nil {
		fmt.Println("bindata:", err)
		os.Exit(1)
	}
}

// run executes the program.
func run() error {
	// use GOPACKAGE (set by go generate) as default package name if available
	pkg := os.Getenv("GOPACKAGE")
	if pkg == "" {
		pkg = "main"
	}

	var out, prefix string
	fs := flag.NewFlagSet("bindata", flag.ExitOnError)
	fs.StringVar(&out, "o", "", "output file (default: stdout)")
	fs.StringVar(&vars.Pkg, "pkg", pkg, "name of the package")
	fs.StringVar(&vars.Map, "map", "bindata", "name of the map")
	fs.StringVar(&prefix, "prefix", "", "prefix to strip")
	fs.BoolVar(&vars.AsString, "string", false, "save data as strings?")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	vars.Files = make(map[string]fmt.Formatter)
	for _, path := range fs.Args() {
		if err := AddPath(path, prefix); err != nil {
			return err
		}
	}

	var file *os.File
	if out != "" {
		var err error
		if file, err = os.Create(out); err != nil {
			return err
		}
	} else {
		file = os.Stdout
	}

	return tmpl.Execute(file, vars)
}

// AddPath add files to the slice in vars recursively.
func AddPath(path, prefix string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		dir, err := os.Open(path)
		if err != nil {
			return err
		}
		files, err := dir.Readdirnames(0)
		if err != nil {
			return err
		}
		for _, file := range files {
			if err := AddPath(filepath.Join(path, file), prefix); err != nil {
				return err
			}
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		path, err := filepath.Rel(prefix, path)
		if err != nil {
			return err
		}
		if vars.AsString {
			vars.Files[path] = StringFormatter{file}
		} else {
			vars.Files[path] = ByteSliceFormatter{file}
		}
	}
	return nil
}

// A ByteSliceFormatter is a byte slice pretty printing io.Reader.
type ByteSliceFormatter struct {
	io.Reader
}

// Format pretty prints the bytes read from the ByteSliceFormatter.
func (f ByteSliceFormatter) Format(s fmt.State, c rune) {
	buf := bufio.NewReader(f)

	const cols = 12 // number of columns in the formatted byte slice.

	fmt.Fprintf(s, "[]byte{")
	b, err := buf.ReadByte()
	for i := 0; err == nil; i++ {
		if i%cols == 0 {
			fmt.Fprintf(s, "\n\t\t")
		} else {
			fmt.Fprintf(s, " ")
		}
		fmt.Fprintf(s, "%#02x,", b)
		b, err = buf.ReadByte()
	}
	fmt.Fprintf(s, "\n\t}")
}

// A StringFormatter is a string pretty printing io.Reader.
type StringFormatter struct {
	io.Reader
}

// Format pretty prints the bytes read from the StringFormatter.
func (f StringFormatter) Format(s fmt.State, c rune) {
	buf := bufio.NewReader(f)

	const cols = 16 // number of bytes per line in the formatted string.

	fmt.Fprintf(s, `"`)
	b, err := buf.ReadByte()
	for i := 0; err == nil; i++ {
		if i%cols == 0 {
			fmt.Fprintf(s, "\" +\n\t\t\"")
		}
		fmt.Fprintf(s, "\\x%02x", b)
		b, err = buf.ReadByte()
	}
	fmt.Fprintf(s, `"`)
}
