package tsdata

import (
	"testing"
	"time"
)

type tsdataFields struct {
	checkers        []func(string) bool
	FileType        string
	Project         string
	FileDescription string
	Comments        []string
	Types           []string
	Units           []string
	Headers         []string
}

func TestTsdata_ParseHeader(t *testing.T) {
	tests := []struct {
		name    string
		fields  tsdataFields
		header  string
		wantErr bool
	}{
		{
			name: "correct header fully populated",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "file description",
				Comments:        []string{"ISO8601 timestamp", "NA"},
				Types:           []string{"time", "float"},
				Units:           []string{"NA", "NA"},
				Headers:         []string{"time", "col1"},
			},
			header: `fileType
project
file description
ISO8601 timestamp	NA
time	float
NA	NA
time	col1`,
			wantErr: false,
		},
		{
			name: "correct header no FileDescription",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "",
				Comments:        []string{"ISO8601 timestamp", "NA"},
				Types:           []string{"time", "float"},
				Units:           []string{"NA", "NA"},
				Headers:         []string{"time", "col1"},
			},
			header: `fileType
project

ISO8601 timestamp	NA
time	float
NA	NA
time	col1`,
			wantErr: false,
		},
		{
			name: "correct header no Comments",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "file description",
				Comments:        []string{},
				Types:           []string{"time", "float"},
				Units:           []string{"NA", "NA"},
				Headers:         []string{"time", "col1"},
			},
			header: `fileType
project
file description

time	float
NA	NA
time	col1`,
			wantErr: false,
		},
		{
			name: "correct header with leading/trailing whitespace",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "file description",
				Comments:        []string{},
				Types:           []string{"time", "float"},
				Units:           []string{"NA", "NA"},
				Headers:         []string{"time", "col1"},
			},
			header: `fileType  
project   
file description   

time  	  float  
NA	NA
time  	  col1  `,
			wantErr: false,
		},
		{
			name: "no data columns",
			header: `project
file description
ISO8601 timestamp
time
NA
time
`,
			wantErr: true,
		},
		{
			name: "wrong line count",
			header: `project
file description
ISO8601 timestamp	NA
time	float
NA	NA
time	col1
`,
			wantErr: true,
		},
		{
			name: "no FileType",
			header: `
project
file description
ISO8601 timestamp	NA
time	float
NA	NA
time	col1
`,
			wantErr: true,
		},
		{
			name: "no Project",
			header: `fileType

file description
ISO8601 timestamp	NA
time	float
NA	NA
time	col1
`,
			wantErr: true,
		},
		{
			name: "incomplete comments",
			header: `fileType
project
file description
ISO8601 timestamp	NA	
time	float	integer
NA	NA	NA
time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "bad type",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	notfloat	integer
NA	NA	NA
time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "empty types column",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time		integer
NA	NA	NA
time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "empty units column",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	float	integer
NA		NA
time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "empty headers column",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	float	integer
NA	NA	NA
time		col2
`,
			wantErr: true,
		},
		{
			name: "missing types column",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	integer
NA	NA	NA
time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "missing units column",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	float	integer
NA	NA
time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "missing headers column",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	float	integer
NA	NA	NA
time	col1
`,
			wantErr: true,
		},
		{
			name: "empty types line",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA

NA	NA	NA
time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "empty units line",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	float	integer

time	col1	col2
`,
			wantErr: true,
		},
		{
			name: "empty headers line",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	float	integer
NA	NA	NA

`,
			wantErr: true,
		},
		{
			name: "first header column not time",
			header: `fileType
project
file description
ISO8601 timestamp	NA	NA
time	float	integer
NA	NA	NA
nottime	col1	col2
`,
			wantErr: true,
		},
		{
			name: "no data columns",
			header: `fileType
project
file description
ISO8601 timestamp
time
NA
time
`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Tsdata{}
			err := d.ParseHeader(tt.header)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Tsdata.ParseHeader() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if len(d.checkers) != len(tt.fields.checkers) {
					t.Errorf("Tsdata.ParseHeader() len(checkers) = %v, expected %v", len(d.checkers), len(tt.fields.checkers))
				}
				if d.FileType != tt.fields.FileType {
					t.Errorf("Tsdata.ParseHeader() FileType = %v, expected %v", d.FileType, tt.fields.FileType)
				}
				if d.Project != tt.fields.Project {
					t.Errorf("Tsdata.ParseHeader() Project = %v, expected %v", d.Project, tt.fields.Project)
				}
				if d.FileDescription != tt.fields.FileDescription {
					t.Errorf("Tsdata.ParseHeader() FileDescription = %v, expected %v", d.FileDescription, tt.fields.FileDescription)
				}
				if !stringSliceEqual(d.Comments, tt.fields.Comments) {
					t.Errorf("Tsdata.ParseHeader() Comments = %v, expected %v", d.Comments, tt.fields.Comments)
				}
				if !stringSliceEqual(d.Types, tt.fields.Types) {
					t.Errorf("Tsdata.ParseHeader() Types = %v, expected %v", d.Types, tt.fields.Types)
				}
				if !stringSliceEqual(d.Units, tt.fields.Units) {
					t.Errorf("Tsdata.ParseHeader() Units = %v, expected %v", d.Units, tt.fields.Units)
				}
				if !stringSliceEqual(d.Headers, tt.fields.Headers) {
					t.Errorf("Tsdata.ParseHeader() Headers = %v, expected %v", d.Headers, tt.fields.Headers)
				}
			}
		})
	}
}

func TestTsdata_ValidateLine(t *testing.T) {
	tline, _ := time.Parse(time.RFC3339, "2017-05-06T19:52:57.601Z")
	floatFields := tsdataFields{
		checkers:        []func(string) bool{checkTime, checkFloat},
		FileType:        "fileType",
		Project:         "project",
		FileDescription: "file description",
		Comments:        []string{"ISO8601 timestamp", "NA"},
		Types:           []string{"time", "float"},
		Units:           []string{"NA", "NA"},
		Headers:         []string{"time", "col1"},
	}
	intFields := tsdataFields{
		checkers:        []func(string) bool{checkTime, checkInteger},
		FileType:        "fileType",
		Project:         "project",
		FileDescription: "file description",
		Comments:        []string{"ISO8601 timestamp", "NA"},
		Types:           []string{"time", "integer"},
		Units:           []string{"NA", "NA"},
		Headers:         []string{"time", "col1"},
	}
	textFields := tsdataFields{
		checkers:        []func(string) bool{checkTime, checkText},
		FileType:        "fileType",
		Project:         "project",
		FileDescription: "file description",
		Comments:        []string{"ISO8601 timestamp", "NA"},
		Types:           []string{"time", "text"},
		Units:           []string{"NA", "NA"},
		Headers:         []string{"time", "col1"},
	}
	categoryFields := tsdataFields{
		checkers:        []func(string) bool{checkTime, checkCategory},
		FileType:        "fileType",
		Project:         "project",
		FileDescription: "file description",
		Comments:        []string{"ISO8601 timestamp", "NA"},
		Types:           []string{"time", "category"},
		Units:           []string{"NA", "NA"},
		Headers:         []string{"time", "col1"},
	}
	boolFields := tsdataFields{
		checkers:        []func(string) bool{checkTime, checkBoolean},
		FileType:        "fileType",
		Project:         "project",
		FileDescription: "file description",
		Comments:        []string{"ISO8601 timestamp", "NA"},
		Types:           []string{"time", "boolean"},
		Units:           []string{"NA", "NA"},
		Headers:         []string{"time", "col1"},
	}
	timeFields := tsdataFields{
		checkers:        []func(string) bool{checkTime, checkTime},
		FileType:        "fileType",
		Project:         "project",
		FileDescription: "file description",
		Comments:        []string{"ISO8601 timestamp", "NA"},
		Types:           []string{"time", "time"},
		Units:           []string{"NA", "NA"},
		Headers:         []string{"time", "col1"},
	}
	multiFields := tsdataFields{
		checkers:        []func(string) bool{checkTime, checkFloat, checkInteger},
		FileType:        "fileType",
		Project:         "project",
		FileDescription: "file description",
		Comments:        []string{"ISO8601 timestamp", "NA", "NA"},
		Types:           []string{"time", "float", "integer"},
		Units:           []string{"NA", "NA", "NA"},
		Headers:         []string{"time", "col1", "col2"},
	}
	tests := []struct {
		name       string
		time       time.Time
		dataFields []string
		fields     tsdataFields
		line       string
		notStrict  bool
		wantErr    bool
	}{
		{
			name:       "correct line TRUE boolean",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "TRUE"},
			fields:     boolFields,
			line:       "2017-05-06T19:52:57.601Z	TRUE",
			wantErr:    false,
		},
		{
			name:       "correct line FALSE boolean",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "FALSE"},
			fields:     boolFields,
			line:       "2017-05-06T19:52:57.601Z	FALSE",
			wantErr:    false,
		},
		{
			name:       "correct line float",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "6.0"},
			fields:     floatFields,
			line:       "2017-05-06T19:52:57.601Z	6.0",
			wantErr:    false,
		},
		{
			name:       "correct line integer",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "100"},
			fields:     intFields,
			line:       "2017-05-06T19:52:57.601Z	100",
			wantErr:    false,
		},
		{
			name:       "correct line time",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "2017-05-07T00:00:00Z"},
			fields:     timeFields,
			line:       "2017-05-06T19:52:57.601Z	2017-05-07T00:00:00.000Z",
			wantErr:    false,
		},
		{
			name:       "accept RFC3339 without 'T', emit with 'T'",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "2017-05-07T00:00:00Z"},
			fields:     timeFields,
			line:       "2017-05-06 19:52:57.601Z	2017-05-07 00:00:00.000Z",
			wantErr:    false,
		},
		{
			name:       "correct line text",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "foo"},
			fields:     textFields,
			line:       "2017-05-06T19:52:57.601Z	foo",
			wantErr:    false,
		},
		{
			name:       "correct line category",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "foo"},
			fields:     categoryFields,
			line:       "2017-05-06T19:52:57.601Z	foo",
			wantErr:    false,
		},
		{
			name:       "empty text",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", ""},
			fields:     textFields,
			line:       "2017-05-06T19:52:57.601Z	",
			wantErr:    false,
		},
		{
			name:    "empty category",
			fields:  categoryFields,
			line:    "2017-05-06T19:52:57.601Z	",
			wantErr: true,
		},
		{
			name:    "NA first timestamp",
			fields:  floatFields,
			line:    "NA	6.0",
			wantErr: true,
		},
		{
			name:       "NA timestamp",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "NA"},
			fields:     timeFields,
			line:       "2017-05-06T19:52:57.601Z	NA",
			wantErr:    false,
		},
		{
			name:       "NA boolean",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "NA"},
			fields:     boolFields,
			line:       "2017-05-06T19:52:57.601Z	NA",
			wantErr:    false,
		},
		{
			name:       "NA float",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "NA"},
			fields:     floatFields,
			line:       "2017-05-06T19:52:57.601Z	NA",
			wantErr:    false,
		},
		{
			name:       "NA integer",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "NA"},
			fields:     intFields,
			line:       "2017-05-06T19:52:57.601Z	NA",
			wantErr:    false,
		},
		{
			name:       "leading whitespace in first time column",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "6.0"},
			fields:     floatFields,
			line:       "  2017-05-06T19:52:57.601Z	6.0",
			wantErr:    false,
		},
		{
			name:       "trailing whitespace in first time column",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "6.0"},
			fields:     floatFields,
			line:       "2017-05-06T19:52:57.601Z    	6.0",
			wantErr:    false,
		},
		{
			name:       "leading whitespace in data column",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "6.0"},
			fields:     floatFields,
			line:       "2017-05-06T19:52:57.601Z	  6.0",
			wantErr:    false,
		},
		{
			name:       "trailing whitespace in data column",
			time:       tline,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "6.0"},
			fields:     floatFields,
			line:       "2017-05-06T19:52:57.601Z	6.0  \r",
			wantErr:    false,
		},
		{
			name:    "empty first timestamp",
			fields:  floatFields,
			line:    "	6.0",
			wantErr: true,
		},
		{
			name:    "empty timestamp",
			fields:  timeFields,
			line:    "2017-05-06T19:52:57.601Z	",
			wantErr: true,
		},
		{
			name:    "empty boolean",
			fields:  boolFields,
			line:    "2017-05-06T19:52:57.601Z	",
			wantErr: true,
		},
		{
			name:    "empty float",
			fields:  floatFields,
			line:    "2017-05-06T19:52:57.601Z	",
			wantErr: true,
		},
		{
			name:    "empty integer",
			fields:  intFields,
			line:    "2017-05-06T19:52:57.601Z	",
			wantErr: true,
		},
		{
			name:    "bad first timestamp",
			fields:  floatFields,
			line:    "2017-05-06aT19:52:57.601Z	6.0",
			wantErr: true,
		},
		{
			name:    "bad timestamp",
			fields:  timeFields,
			line:    "2017-05-06T19:52:57.601Z	201a7-05-07T00:00:00.000Z",
			wantErr: true,
		},
		{
			name:    "bad boolean",
			fields:  boolFields,
			line:    "2017-05-06aT19:52:57.601Z	TaRUE",
			wantErr: true,
		},
		{
			name:    "bad float",
			fields:  floatFields,
			line:    "2017-05-06T19:52:57.601Z	6a.0",
			wantErr: true,
		},
		{
			name:    "bad integer",
			fields:  intFields,
			line:    "2017-05-06T19:52:57.601Z	100.3",
			wantErr: true,
		},
		{
			name:    "empty line",
			fields:  floatFields,
			line:    "",
			wantErr: true,
		},
		{
			name:    "no data column",
			fields:  multiFields,
			line:    "2017-05-06T19:52:57.601Z",
			wantErr: true,
		},
		{
			name:    "missing data column",
			fields:  multiFields,
			line:    "2017-05-06T19:52:57.601Z	100.3",
			wantErr: true,
		},
		{
			name:      "bad first timestamp, not strict",
			fields:    floatFields,
			line:      "2017-05-06aT19:52:57.601Z	6.0",
			wantErr:   true,
			notStrict: true,
		},
		{
			name:       "bad float, not strict",
			fields:     floatFields,
			dataFields: []string{"2017-05-06T19:52:57.601Z", "NA"},
			line:       "2017-05-06T19:52:57.601Z	6a.0",
			wantErr:    false,
			notStrict:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Tsdata{
				checkers:        tt.fields.checkers,
				FileType:        tt.fields.FileType,
				Project:         tt.fields.Project,
				FileDescription: tt.fields.FileDescription,
				Comments:        tt.fields.Comments,
				Types:           tt.fields.Types,
				Units:           tt.fields.Units,
				Headers:         tt.fields.Headers,
			}
			data, err := d.ValidateLine(tt.line, !tt.notStrict)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Tsdata.ValidateLine() err %v, expected a non-nil error", err)
				}
			} else {
				if err != nil {
					t.Errorf("Tsdata.ValidateLine() err %v, expected nil", err)
				}
				if !tline.Equal(data.Time) {
					t.Errorf("Tsdata.ValidateLine() Data.Time %v, expected %v", data.Time, tline)
				}
				if !stringSliceEqual(data.Fields, tt.dataFields) {
					t.Errorf("Tsdata.ValidateLine() fields %v, expected %v", data.Fields, tt.dataFields)
				}
			}
		})
	}
}

// This test has been changed so that it confirms that both in-order and
// out-of-order lines don't generate errors. In the future line order validation
// may be reinstated so this test and accompanying code will not be deleted yet.
func TestTsdata_ValidateLine_order(t *testing.T) {
	t.Run("validate line order", func(t *testing.T) {
		d := &Tsdata{
			checkers:        []func(string) bool{checkTime, checkFloat},
			FileType:        "fileType",
			Project:         "project",
			FileDescription: "",
			Comments:        []string{"ISO8601 timestamp", "column2 notes"},
			Types:           []string{"time", "float"},
			Units:           []string{"NA", "m/s"},
			Headers:         []string{"time", "speed"},
		}
		line0 := "2017-05-06T19:00:00.000Z	6.0"
		line1 := "2017-05-06T19:52:57.601Z	6.0"
		line2 := "2017-05-06T00:00:00.000Z	6.0"

		_, _ = d.ValidateLine(line0, true)
		_, err := d.ValidateLine(line1, true)
		if err != nil {
			t.Errorf("Tsdata.ValidateLine() expected nil error for in-order lines, saw %v", err)
		}
		_, err = d.ValidateLine(line2, true)
		if err != nil {
			t.Errorf("Tsdata.ValidateLine() expected nil error for out-of-order lines")
		}
	})
}

func TestTsdata_Header(t *testing.T) {
	tests := []struct {
		name   string
		fields tsdataFields
		header string
	}{
		{
			name: "full header",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "file description",
				Comments:        []string{"ISO8601 timestamp", "column2 notes", "NA", "column4 notes", "column5 notes", "column6 notes"},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			},
			header: "fileType\nproject\nfile description\nISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes	column6 notes\ntime	float	integer	text	category	boolean\nNA	m/s	km	NA	NA	NA\ntime	speed	distance	notes	color	hasTail",
		},
		{
			name: "header with no column comments",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "file description",
				Comments:        []string{},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			},
			header: "fileType\nproject\nfile description\nNA	NA	NA	NA	NA	NA\ntime	float	integer	text	category	boolean\nNA	m/s	km	NA	NA	NA\ntime	speed	distance	notes	color	hasTail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Tsdata{
				checkers:        tt.fields.checkers,
				FileType:        tt.fields.FileType,
				Project:         tt.fields.Project,
				FileDescription: tt.fields.FileDescription,
				Comments:        tt.fields.Comments,
				Types:           tt.fields.Types,
				Units:           tt.fields.Units,
				Headers:         tt.fields.Headers,
			}
			h := d.Header()
			if h != tt.header {
				t.Errorf("Tsdata.Header() Header = %v, expected %v", h, tt.header)
			}
		})
	}
}

func stringSliceEqual(a []string, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
