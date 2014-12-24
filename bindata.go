// The bindata command embeds binary files as byte arrays into a Go source file.
//
// It is designed with go generate in mind, but can be used on its own as well.
//
// The data is stored as a map of byte slices or strings indexed by the
// file paths as specified on the command line. The default name of the
// map is "bindata" but a custom name can be specified on the command line (-m).
//
// Multiple files and directories can be provided on the command line.
// Directories are treated recursively. The keys of the map are the paths
// of the files relative to the current directory. A different root for
// the paths can be specified on the command line (-r).
//
// By default, the data are saved as byte slices.
// It is also possible to save them a strings (-s).
//
// By default, the package name of the file containing the generate directive
// is used as the package name of the generated file, or "main" otherwise.
// A custom package name can also be specified on the command line (-p).
//
// The output file can be specified on the command line (-o).
// If a file already exists at this location, it will be overwritten.
// The file produced is properly formatted and commented.
// If no output file is specified, the contents are printed on the standard output.
//
// To see the full list of flags, run:
//  bindata -h
//
// Example
//
// Given a file hello.go containing:
//
//  package main
//
//  import "fmt"
//
//  func main() {
//  	fmt.Println("Hello, 世界")
//  }
//
// Running `bindata hello.go` will produce:
//
//  package main
//
//  // This file is generated. Do not edit directly.
//
//  // bindata stores binary files as byte slices indexed by filepaths.
//  var bindata = map[string][]byte{
//  	"hello.go": []byte{
//  		0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x20, 0x6d, 0x61, 0x69, 0x6e,
//  		0x0a, 0x0a, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x22, 0x66, 0x6d,
//  		0x74, 0x22, 0x0a, 0x0a, 0x66, 0x75, 0x6e, 0x63, 0x20, 0x6d, 0x61, 0x69,
//  		0x6e, 0x28, 0x29, 0x20, 0x7b, 0x0a, 0x09, 0x66, 0x6d, 0x74, 0x2e, 0x50,
//  		0x72, 0x69, 0x6e, 0x74, 0x6c, 0x6e, 0x28, 0x22, 0x48, 0x65, 0x6c, 0x6c,
//  		0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, 0x22, 0x29, 0x0a,
//  		0x7d, 0x0a,
//  	},
//  }
//
// Example using go generate
//
// Add a command like this one anywhere in a source file:
//  //go:generate bindata -o jpegs.go pic1.jpg pic2.jpg pic3.jpg
// Then simply run
//  go generate
// and the file jpegs.go will be created.
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

// {{.Map}} stores binary files as {{if .AsString}}strings{{else}}byte slices{{end}} indexed by file paths.
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
	fs.StringVar(&vars.Pkg, "p", pkg, "name of the package")
	fs.StringVar(&vars.Map, "m", "bindata", "name of the map variable")
	fs.StringVar(&prefix, "r", "", "root path for map keys")
	fs.BoolVar(&vars.AsString, "s", false, "save data as strings")
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
