# tsdata

A Go module and command-line tool to process TSDATA files

See [https://github.com/armbrustlab/tsdataformat](https://github.com/armbrustlab/tsdataformat) for information on TSDATA files.

## Installation

To install the command-line tool,
from within a cloned repo directory run `go install ./...`.
This will install a new binary called tsdata in `$GOPATH/bin`.

## CLI

```
$ tsdata -help
NAME:
   tsdata - A new cli application

USAGE:
   tsdata [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
   validate  validates a TSDATA file
   csv       converts a TSDATA file to CSV
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

```
$ tsdata validate -help
NAME:
   tsdata validate - validates a TSDATA file

USAGE:
   tsdata validate INFILE

DESCRIPTION:
   Validates metadata and data in INFILE. Prints errors encountered to STDERR.

OPTIONS:
   --stringent, -s  Exit after the first data line validation error
   --quiet, -q      Suppress logging output
```

```
$ tsdata csv -help
NAME:
   tsdata csv - converts a TSDATA file to CSV

USAGE:
   tsdata csv INFILE OUTFILE

DESCRIPTION:
   Validates and converts a TSDATA file at INFILE to a CSV file at OUTFILE.

OPTIONS:
   --quiet, -q  Suppress logging output
```

The csv subcommand will report data line errors,
but otherwise will ignore those lines when writing output.

## Library

```golang
import "github.com/ctberthiaume/tsdata"
```

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
data, err := t.ValidateLine(line)
if err != nil {
    log.Fatalf("%v\n", err)
}
// data.Fields will be a []string of columns present in the line
// data.Time will be the parsed timestamp for this line
```
