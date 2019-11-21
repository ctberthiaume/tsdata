// Package tsdata provides tools to manage TSData files. See
// https://github.com/armbrustlab/tsdataformat for a description of TSData
// files.
package tsdata

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Delim is the field separator string
const Delim = "\t"

// NA is the string used to represent missing data
const NA = "NA"

// HeaderSize is the number of lines in a header section
const HeaderSize = 7

// Tsdata defines a TSData file
type Tsdata struct {
	checkers        []func(string) bool
	FileType        string
	Project         string
	FileDescription string
	Comments        []string
	Types           []string
	Units           []string
	Headers         []string
}

// ValidateLine checks values in a data line and returns all fields as a slice of
// strings. It returns an error for the first field that fails validation.
func (d *Tsdata) ValidateLine(line string) (fields []string, err error) {
	fields = strings.Split(line, Delim)
	if len(fields) != len(d.Headers) {
		return nil, fmt.Errorf("found %v columns, expected %v", len(fields), len(d.Headers))
	}
	for i, f := range fields {
		if !d.checkers[i](f) {
			return nil, fmt.Errorf("column %v, bad value '%v'", i+1, f)
		}
	}
	return fields, nil
}

// ParseHeader parses and validates header metadata. Input should a string of
// all lines in the file's header section.
func (d *Tsdata) ParseHeader(header string) error {
	header = strings.TrimSuffix(header, "\n")
	headerLines := strings.Split(header, "\n")
	if len(headerLines) != HeaderSize {
		return fmt.Errorf("expected %v lines in header, found %v", HeaderSize, len(headerLines))
	}
	d.FileType = strings.Split(headerLines[0], Delim)[0]
	d.Project = strings.Split(headerLines[1], Delim)[0]
	d.FileDescription = strings.Split(headerLines[2], Delim)[0]
	if headerLines[3] != "" {
		d.Comments = strings.Split(headerLines[3], Delim)
	}
	if headerLines[4] != "" {
		d.Types = strings.Split(headerLines[4], Delim)
	}
	if headerLines[5] != "" {
		d.Units = strings.Split(headerLines[5], Delim)
	}
	if headerLines[6] != "" {
		d.Headers = strings.Split(headerLines[6], Delim)
	}

	d.checkers = make([]func(string) bool, len(d.Types))
	for i, t := range d.Types {
		d.checkers[i] = typecheckers[t]
	}
	return d.ValidateMetadata()
}

// ValidateMetadata checks for errors and inconsistencies in metadata values.
func (d *Tsdata) ValidateMetadata() error {
	// FileType
	if d.FileType == "" {
		return fmt.Errorf("missing or empty FileType")
	}

	// Project
	if d.Project == "" {
		return fmt.Errorf("missing or empty Project")
	}

	// Comments
	colCount := 0
	// Column comments may be a blank line so allow 0 columns
	if len(d.Comments) > 0 {
		colCount = len(d.Comments)
		for i, com := range d.Comments {
			if com == "" {
				return fmt.Errorf("empty comment in column %v", i+1)
			}
		}
	}

	// Types
	if len(d.Types) == 0 {
		return fmt.Errorf("missing or empty Types")
	}
	if colCount > 0 && len(d.Types) != colCount {
		return fmt.Errorf("inconsistent Types column count")
	}
	for i, t := range d.Types {
		_, ok := typecheckers[t]
		if !ok {
			return fmt.Errorf("bad Types value '%v' in column %v", t, i+1)
		}
	}
	colCount = len(d.Types)

	// Units
	if len(d.Units) == 0 {
		return fmt.Errorf("missing or empty Units")
	}
	if len(d.Units) != colCount {
		return fmt.Errorf("inconsistent Units column count")
	}
	for i, u := range d.Units {
		if u == "" {
			return fmt.Errorf("empty Units value in column %v", i+1)
		}
	}

	// Headers
	if len(d.Headers) == 0 {
		return fmt.Errorf("missing or empty Headers")
	}
	if len(d.Headers) != colCount {
		return fmt.Errorf("inconsistent Headers column count")
	}
	if d.Headers[0] != "time" {
		return fmt.Errorf("first Headers column should be 'time'")
	}
	for i, h := range d.Headers {
		if h == "" {
			return fmt.Errorf("empty Headers value in column %v", i+1)
		}
	}

	return nil
}

// Header creates a TSData file metadata header paragraph.
func (d *Tsdata) Header() string {
	// TODO: should this ever produce a non-conforming TSData header?
	cols := len(d.Headers)
	text := d.FileType + "\n"
	text = text + d.Project + "\n"
	text = text + d.FileDescription + "\n"
	if len(d.Comments) == 0 {
		text = text + nas(cols) + "\n"
	} else {
		text = text + strings.Join(d.Comments, Delim) + "\n"
	}
	text = text + strings.Join(d.Types, Delim) + "\n"
	text = text + strings.Join(d.Units, Delim) + "\n"
	text = text + strings.Join(d.Headers, Delim) // note, doesn't end with blank line
	return text
}

func checkTime(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}

func checkFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s == NA
	}
	return true
}

func checkInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return s == NA
	}
	return true
}

func checkText(s string) bool {
	return true
}

func checkBoolean(s string) bool {
	return (s == "TRUE" || s == "FALSE" || s == NA)
}

var typecheckers = map[string]func(string) bool{
	"time":     checkTime,
	"float":    checkFloat,
	"integer":  checkInteger,
	"text":     checkText,
	"category": checkText,
	"boolean":  checkBoolean,
}

func nas(size int) string {
	s := make([]string, size)
	for i := range s {
		s[i] = NA
	}
	return strings.Join(s, Delim)
}
