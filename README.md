# bindata

[![GoDoc](https://godoc.org/github.com/simleb/bindata?status.svg)](http://godoc.org/github.com/simleb/bindata)
[![Coverage Status](https://img.shields.io/coveralls/simleb/bindata.svg)](https://coveralls.io/r/simleb/bindata)
[![Build Status](https://drone.io/github.com/simleb/bindata/status.png)](https://drone.io/github.com/simleb/bindata/latest)

The `bindata` command translates binary files into byte arrays in Go source.

`bindata` is designed to work with go generate, but can be used on its own as well.

The data is stored as a map of byte slices or strings indexed by the
file paths as specified on the command line. The default name of the
map is `bindata` but a custom name can be specified on the command line.

By default, the package name of the file containing the generate directive
is used as the package name of the generated file, or `main` otherwise.
A custom package name can also be specified on the command line.

The output file can be specified on the command line (with the `-o` flag).
By default, the file produced is printed on the standard output.
If a file already exists at this location, it will be overwritten.
The file produced is properly formatted and commented.

To see the full list of flags, run:

	bindata -h

## Example

Example using go generate:

	//go:generate bindata -o jpegs.go pic1.jpg pic2.jpg pic3.jpg

## Todo (maybe)

- [ ] add option to compress data (but then need accessor)


## License

The MIT License (MIT)

	Copyright (c) 2014 Simon Leblanc
	
	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:
	
	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.
	
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
