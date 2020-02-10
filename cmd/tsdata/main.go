package main

import (
	"bufio"
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
var version string = "v0.2.2"

func main() {
	logger = log.New(os.Stderr, "", 0)
	app := cli.NewApp()
	app.Name = cmdname
	app.Version = version
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "validate",
			Usage:       "validates a TSDATA file",
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
				err := validate(c.Args().Get(0), c.Bool("stringent"))
				if err != nil {
					logger.Println(err)
				}
				return err
			},
		},
		cli.Command{
			Name:        "csv",
			Usage:       "converts a TSDATA file to CSV",
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
				err := csv(c.Args().Get(0), c.Args().Get(1))
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

func validate(infile string, stringent bool) error {
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
		_, err := ts.ValidateLine(scanner.Text())
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

func csv(infile string, outfile string) error {
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
	buffSize := 1 << 16 // 65536 byte buffer
	w := bufio.NewWriterSize(outf, buffSize)

	// Write CSV column headers
	_, err = w.WriteString(strings.Join(ts.Headers, ",") + "\n")
	if err != nil {
		return err
	}

	// Write CSV lines
	i := tsdata.HeaderSize
	for scanner.Scan() {
		i++
		data, err := ts.ValidateLine(scanner.Text())
		if err != nil {
			logger.Println(err)
			continue
		}
		_, err = w.WriteString(strings.Join(data.Fields, ",") + "\n")
		if err != nil {
			return err
		}
	}
	err = scanner.Err()
	if err != nil {
		return err
	}

	w.Flush()
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
