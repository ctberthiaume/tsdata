package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ctberthiaume/tsdata"
	"github.com/urfave/cli"
)

var logger *log.Logger
var cmdname string = "tsdata"
var version string = "v0.3.0"

func main() {
	logger = log.New(os.Stderr, "", 0)
	app := cli.NewApp()
	app.Name = cmdname
	app.Usage = "process time-series TSDATA files (https://github.com/armbrustlab/tsdataformat)"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:        "validate",
			Usage:       "Validates a TSDATA file",
			UsageText:   "tsdata validate INFILE",
			Description: "Validates metadata and data in INFILE. Prints errors encountered to STDERR. Use '-' for STDIN.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stringent, s",
					Usage: "Exit after the first data line validation error",
				},
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "Suppress logging output",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					err := fmt.Errorf("missing required INFILE argument")
					logger.Println(err)
					return err
				}
				if c.NArg() > 1 {
					err := fmt.Errorf("too many arguments")
					logger.Println(err)
					return err
				}
				if c.Bool("quiet") {
					logger.SetOutput(ioutil.Discard)
				}
				err := validateCmd(c.Args().Get(0), c.Bool("stringent"))
				if err != nil {
					logger.Println(err)
				}
				return err
			},
		},
		{
			Name:        "csv",
			Usage:       "Converts a TSDATA file to CSV",
			UsageText:   "tsdata csv INFILE OUTFILE",
			Description: "Validates and converts a TSDATA file at INFILE to a CSV file at OUTFILE. Use '-' for STDIN and STDOUT.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "Suppress logging output",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					err := fmt.Errorf("missing required INFILE and OUTIFLE arguments")
					logger.Println(err)
					return err
				}
				if c.NArg() < 2 {
					err := fmt.Errorf("missing required OUTFILE argument")
					logger.Println(err)
					return err
				}
				if c.Bool("quiet") {
					logger.SetOutput(ioutil.Discard)
				}
				err := csvCmd(c.Args().Get(0), c.Args().Get(1))
				if err != nil {
					logger.Println(err)
				}
				return err
			},
		},
		{
			Name:        "clean",
			Usage:       "Clean a TSDATA file",
			UsageText:   "tsdata clean INFILE OUTFILE",
			Description: "Fix common errors in a TSDATA file at INFILE, write to OUTFILE. Use '-' for STDIN and STDOUT.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "Suppress logging output",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					err := fmt.Errorf("missing required INFILE and OUTIFLE arguments")
					logger.Println(err)
					return err
				}
				if c.NArg() < 2 {
					err := fmt.Errorf("missing required OUTFILE argument")
					logger.Println(err)
					return err
				}
				if c.Bool("quiet") {
					logger.SetOutput(ioutil.Discard)
				}
				err := cleanCmd(c.Args().Get(0), c.Args().Get(1))
				if err != nil {
					logger.Println(err)
				}
				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func validateCmd(infile string, stringent bool) error {
	var r *os.File
	var err error
	if infile == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(infile)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	ts := tsdata.Tsdata{}
	scanner := bufio.NewScanner(r)
	header, err := readHeader(scanner)
	if err != nil {
		return err
	}
	err = ts.ParseHeader(header)
	if err != nil {
		return err
	}

	sawError := false
	i := tsdata.HeaderSize
	for scanner.Scan() {
		i++
		_, err := ts.ValidateLine(scanner.Text(), true)
		if err != nil {
			sawError = true
			logger.Printf("line %v, %v\n", i, err)
			if stringent {
				break
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return err
	}

	if sawError {
		return fmt.Errorf("%v failed validation", infile)
	}
	return nil
}

func csvCmd(infile string, outfile string) error {
	var r *os.File
	var err error
	if infile == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(infile)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	ts := tsdata.Tsdata{}
	scanner := bufio.NewScanner(r)
	header, err := readHeader(scanner)
	if err != nil {
		return err
	}
	err = ts.ParseHeader(header)
	if err != nil {
		return err
	}

	var outf *os.File
	if outfile == "-" {
		outf = os.Stdout
	} else {
		outf, err = os.Create(outfile)
		if err != nil {
			return err
		}
	}
	w := csv.NewWriter(outf)

	// Write CSV column headers
	err = w.Write(ts.Headers)
	if err != nil {
		return err
	}

	// Write CSV lines
	i := tsdata.HeaderSize
	for scanner.Scan() {
		i++
		data, err := ts.ValidateLine(scanner.Text(), false)
		if err != nil {
			logger.Printf("line %v, %v\n", i, err)
			continue
		}
		err = w.Write(data.Fields)
		if err != nil {
			return err
		}
	}
	err = scanner.Err()
	if err != nil {
		return err
	}

	w.Flush()
	err = w.Error()
	if err != nil {
		return err
	}
	if outfile == "-" {
		err = outf.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanCmd(infile string, outfile string) error {
	var r *os.File
	var err error
	if infile == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(infile)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	ts := tsdata.Tsdata{}
	scanner := bufio.NewScanner(r)
	header, err := readHeader(scanner)
	if err != nil {
		return err
	}
	err = ts.ParseHeader(header)
	if err != nil {
		return err
	}

	var outf *os.File
	if outfile == "-" {
		outf = os.Stdout
	} else {
		outf, err = os.Create(outfile)
		if err != nil {
			return err
		}
	}
	w := bufio.NewWriter(outf)

	// Write header section
	_, err = w.WriteString(ts.Header() + "\n")
	if err != nil {
		return err
	}

	// Write TSDATA lines
	i := tsdata.HeaderSize
	for scanner.Scan() {
		i++
		data, err := ts.ValidateLine(scanner.Text(), false)
		if err != nil {
			logger.Printf("line %v, %v\n", i, err)
			continue
		}
		_, err = w.WriteString(strings.Join(data.Fields, tsdata.Delim) + "\n")
		if err != nil {
			return err
		}
	}
	err = scanner.Err()
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}
	if outfile == "-" {
		err = outf.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func readHeader(scanner *bufio.Scanner) (header string, err error) {
	headerLines := make([]string, 7)
	var i int
	for i = 0; i < tsdata.HeaderSize; i++ {
		if !scanner.Scan() {
			err := scanner.Err()
			if err != nil {
				return "", err
			}
			break
		}
		headerLines[i] = scanner.Text()
	}
	header = strings.Join(headerLines, "\n")
	return header, nil
}
