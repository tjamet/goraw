[![GoDoc](https://godoc.org/github.com/tjamet/goraw?status.svg)](https://godoc.org/github.com/tjamet/goraw)
[![Build Status](https://travis-ci.org/tjamet/goraw.svg?branch=master)](https://travis-ci.org/tjamet/goraw)

# goraw
A go library to read raw files

This library currently exposes simple parsers to access Exif metadata inside the raw files.

Are currently supported:

- canon `import _ "github.com/tjamet/goraw/canon"`
- nikon `import _ "github.com/tjamet/goraw/nikon"`
- fujifilm `import _ "github.com/tjamet/goraw/fuji"`

To use the library, please see  example_test.go
