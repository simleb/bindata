# bindata

[![GoDoc](https://godoc.org/github.com/simleb/bindata?status.svg)](http://godoc.org/github.com/simleb/bindata)
[![Coverage Status](https://img.shields.io/coveralls/simleb/bindata.svg)](https://coveralls.io/r/simleb/bindata)
[![Build Status](https://drone.io/github.com/simleb/bindata/status.png)](https://drone.io/github.com/simleb/bindata/latest)

The `bindata` command embeds binary files as byte arrays into a Go source file.

It is designed with go generate in mind, but can be used on its own as well.

The data is stored as a map of byte slices or strings indexed by the file paths as specified on the command line. The default name of the map is `bindata` but a custom name can be specified on the command line (`-m`).

Multiple files and directories can be provided on the command line. Directories are treated recursively. The keys of the map are the paths of the files relative to the current directory. A different root for the paths can be specified on the command line (`-r`).

By default, the data are saved as byte slices. It is also possible to save them a strings (`-s`).

By default, the package name of the file containing the generate directive is used as the package name of the generated file, or `main` otherwise. A custom package name can also be specified on the command line (`-p`).

The output file can be specified on the command line (`-o`). If a file already exists at this location, it will be overwritten. The file produced is properly formatted and commented. If no output file is specified, the contents are printed on the standard output.

To see the full list of flags, run:

	bindata -h

## Example

Given a file `hello.go` containing:

	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, 世界")
	}

Running `bindata hello.go` will produce:

	package main

	// This file is generated. Do not edit directly.

	// bindata stores binary files as byte slices indexed by filepaths.
	var bindata = map[string][]byte{
		"hello.go": []byte{
			0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x20, 0x6d, 0x61, 0x69, 0x6e,
			0x0a, 0x0a, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x22, 0x66, 0x6d,
			0x74, 0x22, 0x0a, 0x0a, 0x66, 0x75, 0x6e, 0x63, 0x20, 0x6d, 0x61, 0x69,
			0x6e, 0x28, 0x29, 0x20, 0x7b, 0x0a, 0x09, 0x66, 0x6d, 0x74, 0x2e, 0x50,
			0x72, 0x69, 0x6e, 0x74, 0x6c, 0x6e, 0x28, 0x22, 0x48, 0x65, 0x6c, 0x6c,
			0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, 0x22, 0x29, 0x0a,
			0x7d, 0x0a,
		},
	}

## Example using Go generate

Add a command like this one anywhere in a source file:

	//go:generate bindata -o jpegs.go pic1.jpg pic2.jpg pic3.jpg

Then simply run `go generate` and the file `jpegs.go` will be created.

## Todo (maybe)

- [ ] add option to compress data (but then need accessor)


## License

The MIT License (MIT). See [LICENSE.txt](LICENSE.txt).
