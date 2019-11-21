# tsdata

A Go module and command-line tool to process TSDATA files

See [https://github.com/armbrustlab/tsdataformat](https://github.com/armbrustlab/tsdataformat) for information on TSDATA files.

## Installation

To install the command-line tool of the same name

```sh
env GO111MODULE=on go get github.com/ctberthiaume/tsdata/cmd/tsdata
```

As of Go version 1.13 the default value for `GO111MODULE` is `auto`,
which may cause dependency compatiblity problems.
Prepending with `env GO111MODULE=on` ensures the tool builds with the correct dependency versions.
This may be unnecessary in future versions of Go.
See [https://golang.org/cmd/go/#hdr-Module_support](https://golang.org/cmd/go/#hdr-Module_support).

## Library

To define the metadata for a TSDATA file, either construct a `Tsdata` struct directly, e.g.

```golang
t := tsdata.Tsdata{
    checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
    FileType:        "fileType",
    Project:         "project",
    FileDescription: "Some general comments about this file on a single line, not tab-delimited",
    Comments:        []string{"ISO8601 timestamp", "column2 notes", "NA", "column4 notes", "column5 notes", "column6 notes"},
    Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
    Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
    Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
}
err := t.ValidateMetaData()
if err != nil {
    log.Fatalf("%v\n", err)
}
```

or parse the header lines from a TSDATA file, e.g.

```golang
// ... get header paragraph string
t := tsdata.Tsdata{}
// ValidateMetaData is called at the end of header parsing and the same errors
// will be returned by ParseHeader
err := t.ParseHeader(header)
if err != nil {
    log.Fatalf("%v\n", err)
}
```

Once a `Tsdata` struct has been created with validated header metadata,
data lines can be validated with `ValidateLine`.

```golang
fields, err := t.ValidateLine(line)
if err != nil {
    log.Fatalf("%v\n", err)
}
// fields will be a []string of columns present in the line
```
