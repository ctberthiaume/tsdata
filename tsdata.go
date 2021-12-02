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
	lastTime        time.Time
	FileType        string
	Project         string
	FileDescription string
	Comments        []string
	Types           []string
	Units           []string
	Headers         []string
}

// Data holds validated information for one TSDATA file line, with the original
// column strings in Fields and time in Time.
type Data struct {
	Fields []string
	Time   time.Time
}

// ValidateLine checks values in a data line and returns all fields as a slice of
// strings. It returns an error for the first field that fails validation. It
// also returns an error if the timestamp in this line is earlier than the
// timestamp in the last line validated by this struct.
func (t *Tsdata) ValidateLine(line string) (Data, error) {
	fields := strings.Split(line, Delim)
	if len(fields) < 2 {
		// Need at least time column plus one data column
		return Data{}, fmt.Errorf("found %v columns, expected >= 2", len(fields))
	}
	if len(fields) != len(t.Headers) {
		return Data{}, fmt.Errorf("found %v columns, expected %v", len(fields), len(t.Headers))
	}
	// Validate first time column separately here to avoid parsing timestamp
	// twice and to make sure not NA
	fields[0] = strings.TrimSpace(fields[0]) // remove leading/trailing whitespace
	tline, err := time.Parse(time.RFC3339, fields[0])
	if err != nil {
		return Data{}, fmt.Errorf("first time column, bad value '%v'", fields[0])
	}
	// Turn off time order check for now, it's sometimes too stringent.
	//if tline.Sub(t.lastTime) < 0 {
	//	return Data{}, fmt.Errorf("timestamp less than previous line, %v < %v", tline, t.lastTime)
	//}
	for i := 1; i < len(fields); i++ { // skip first time column
		// Remove leading/trailing whitespace from each data field
		fields[i] = strings.TrimSpace(fields[i])
		if !t.checkers[i](fields[i]) {
			return Data{}, fmt.Errorf("column %v, bad value '%v'", i+1, fields[i])
		}
	}
	t.lastTime = tline
	return Data{Fields: fields, Time: tline}, nil
}

// ParseHeader parses and validates header metadata. Input should a string of
// all lines in the file's header section.
func (t *Tsdata) ParseHeader(header string) error {
	header = strings.TrimSuffix(header, "\n")
	headerLines := strings.Split(header, "\n")
	if len(headerLines) != HeaderSize {
		return fmt.Errorf("expected %v lines in header, found %v", HeaderSize, len(headerLines))
	}
	// Remove trailing whitespace from each line
	for i := 0; i < len(headerLines); i++ {
		headerLines[i] = strings.TrimRight(headerLines[i], " \t\r")
	}

	t.FileType = strings.Split(headerLines[0], Delim)[0]
	t.Project = strings.Split(headerLines[1], Delim)[0]
	t.FileDescription = strings.Split(headerLines[2], Delim)[0]
	if headerLines[3] != "" {
		t.Comments = strings.Split(headerLines[3], Delim)
		// Remove leading/trailing whitespace from each field
		for i := 0; i < len(t.Comments); i++ {
			t.Comments[i] = strings.TrimSpace(t.Comments[i])
		}
	}
	if headerLines[4] != "" {
		t.Types = strings.Split(headerLines[4], Delim)
		// Remove leading/trailing whitespace from each field
		for i := 0; i < len(t.Types); i++ {
			t.Types[i] = strings.TrimSpace(t.Types[i])
		}
	}
	if headerLines[5] != "" {
		t.Units = strings.Split(headerLines[5], Delim)
		// Remove leading/trailing whitespace from each field
		for i := 0; i < len(t.Units); i++ {
			t.Units[i] = strings.TrimSpace(t.Units[i])
		}
	}
	if headerLines[6] != "" {
		t.Headers = strings.Split(headerLines[6], Delim)
		// Remove leading/trailing whitespace from each field
		for i := 0; i < len(t.Headers); i++ {
			t.Headers[i] = strings.TrimSpace(t.Headers[i])
		}
	}

	t.checkers = make([]func(string) bool, len(t.Types))
	for i, ty := range t.Types {
		t.checkers[i] = typecheckers[ty]
	}
	return t.ValidateMetadata()
}

// ValidateMetadata checks for errors and inconsistencies in metadata values.
func (t *Tsdata) ValidateMetadata() error {
	// FileType
	if t.FileType == "" {
		return fmt.Errorf("missing or empty FileType")
	}

	// Project
	if t.Project == "" {
		return fmt.Errorf("missing or empty Project")
	}

	// Comments
	colCount := 0
	// Column comments may be a blank line so allow 0 columns
	if len(t.Comments) > 0 {
		colCount = len(t.Comments)
		for i, com := range t.Comments {
			if com == "" {
				return fmt.Errorf("empty comment in column %v", i+1)
			}
		}
	}

	// Types
	if len(t.Types) == 0 {
		return fmt.Errorf("missing or empty Types")
	}
	if colCount > 0 && len(t.Types) != colCount {
		return fmt.Errorf("inconsistent Types column count")
	}
	for i, t := range t.Types {
		_, ok := typecheckers[t]
		if !ok {
			return fmt.Errorf("bad Types value '%v' in column %v", t, i+1)
		}
	}
	colCount = len(t.Types)

	// Units
	if len(t.Units) == 0 {
		return fmt.Errorf("missing or empty Units")
	}
	if len(t.Units) != colCount {
		return fmt.Errorf("inconsistent Units column count")
	}
	for i, u := range t.Units {
		if u == "" {
			return fmt.Errorf("empty Units value in column %v", i+1)
		}
	}

	// Headers
	if len(t.Headers) == 0 {
		return fmt.Errorf("missing or empty Headers")
	}
	if len(t.Headers) != colCount {
		return fmt.Errorf("inconsistent Headers column count")
	}
	if t.Headers[0] != "time" {
		return fmt.Errorf("first Headers column should be 'time'")
	}
	for i, h := range t.Headers {
		if h == "" {
			return fmt.Errorf("empty Headers value in column %v", i+1)
		}
	}

	// Finally column count should be > 1, meaning at least one data column
	// after the first time column
	if colCount < 2 {
		return fmt.Errorf("no data columns after time")
	}

	return nil
}

// Header creates a TSData file metadata header paragraph.
func (t *Tsdata) Header() string {
	// TODO: should this ever produce a non-conforming TSData header?
	cols := len(t.Headers)
	text := t.FileType + "\n"
	text = text + t.Project + "\n"
	text = text + t.FileDescription + "\n"
	if len(t.Comments) == 0 {
		text = text + nas(cols) + "\n"
	} else {
		text = text + strings.Join(t.Comments, Delim) + "\n"
	}
	text = text + strings.Join(t.Types, Delim) + "\n"
	text = text + strings.Join(t.Units, Delim) + "\n"
	text = text + strings.Join(t.Headers, Delim) // note, doesn't end with blank line
	return text
}

func checkTime(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s == NA
	}
	return true
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

func checkCategory(s string) bool {
	return s != ""
}

func checkBoolean(s string) bool {
	return (s == "TRUE" || s == "FALSE" || s == NA)
}

var typecheckers = map[string]func(string) bool{
	"time":     checkTime,
	"float":    checkFloat,
	"integer":  checkInteger,
	"text":     checkText,
	"category": checkCategory,
	"boolean":  checkBoolean,
}

func nas(size int) string {
	s := make([]string, size)
	for i := range s {
		s[i] = NA
	}
	return strings.Join(s, Delim)
}
