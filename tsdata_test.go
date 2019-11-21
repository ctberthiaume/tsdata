package tsdata

import "testing"

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

	type args struct {
		header string
	}
	tests := []struct {
		name    string
		fields  tsdataFields
		header  string
		wantErr bool
	}{
		{
			name: "correct header fully populated",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "Some general comments about this file on a single line, not tab-delimited",
				Comments:        []string{"ISO8601 timestamp", "column2 notes", "NA", "column4 notes", "column5 notes", "column6 notes"},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			},
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes	column6 notes
time	float	integer	text	category	boolean
NA	m/s	km	NA	NA	NA
time	speed	distance	notes	color	hasTail`,
			wantErr: false,
		},
		{
			name: "correct header no FileDescription",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "",
				Comments:        []string{"ISO8601 timestamp", "column2 notes", "NA", "column4 notes", "column5 notes", "column6 notes"},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			},
			header: `fileType
project

ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes	column6 notes
time	float	integer	text	category	boolean
NA	m/s	km	NA	NA	NA
time	speed	distance	notes	color	hasTail`,
			wantErr: false,
		},
		{
			name: "correct header no Comments",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "Some general comments about this file on a single line, not tab-delimited",
				Comments:        []string{},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			},
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited

time	float	integer	text	category	boolean
NA	m/s	km	NA	NA	NA
time	speed	distance	notes	color	hasTail
`,
			wantErr: false,
		},
		{
			name: "wrong line count",
			header: `project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "no FileType",
			header: `
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "no Project",
			header: `fileType

Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "incomplete comments",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	
time	float	integer	text	category
NA	m/s	km	NA	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "empty types column",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	
NA	m/s	km	NA	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "empty units column",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km		NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "empty headers column",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA	NA
time	speed		notes	color
`,
			wantErr: true,
		},
		{
			name: "missing types column",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text
NA	m/s	km	NA	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "missing units column",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "missing headers column",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA	NA
time	speed	distance	notes
`,
			wantErr: true,
		},
		{
			name: "empty types line",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes

NA	m/s	km	NA	NA
time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "empty units line",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category

time	speed	distance	notes	color
`,
			wantErr: true,
		},
		{
			name: "empty headers line",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA	NA

`,
			wantErr: true,
		},
		{
			name: "first header column not time",
			header: `fileType
project
Some general comments about this file on a single line, not tab-delimited
ISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes
time	float	integer	text	category
NA	m/s	km	NA	NA
notTime	speed	distance	notes	color
`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Tsdata{}
			if err := d.ParseHeader(tt.header); (err != nil) != tt.wantErr {
				t.Errorf("Tsdata.ParseHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
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
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		fields  []string
		args    args
		wantErr bool
	}{
		{
			name:   "correct line TRUE boolean",
			fields: []string{"2017-05-06T19:52:57.601Z", "6.0", "100", "foo", "blue", "TRUE"},
			args: args{"2017-05-06T19:52:57.601Z	6.0	100	foo	blue	TRUE"},
			wantErr: false,
		},
		{
			name:   "correct line FALSE boolean",
			fields: []string{"2017-05-06T19:52:57.601Z", "6.0", "100", "foo", "blue", "FALSE"},
			args: args{"2017-05-06T19:52:57.601Z	6.0	100	foo	blue	FALSE"},
			wantErr: false,
		},
		{
			name: "NA timestamp",
			args: args{"NA	6.0	100	foo	blue	TRUE"},
			wantErr: true,
		},
		{
			name:   "NA boolean",
			fields: []string{"2017-05-06T19:52:57.601Z", "6.0", "100", "foo", "blue", "NA"},
			args: args{"2017-05-06T19:52:57.601Z	6.0	100	foo	blue	NA"},
			wantErr: false,
		},
		{
			name:   "NA float",
			fields: []string{"2017-05-06T19:52:57.601Z", "NA", "100", "foo", "blue", "TRUE"},
			args: args{"2017-05-06T19:52:57.601Z	NA	100	foo	blue	TRUE"},
			wantErr: false,
		},
		{
			name:   "NA integer",
			fields: []string{"2017-05-06T19:52:57.601Z", "6.0", "NA", "foo", "blue", "TRUE"},
			args: args{"2017-05-06T19:52:57.601Z	6.0	NA	foo	blue	TRUE"},
			wantErr: false,
		},
		{
			name: "empty timestamp",
			args: args{"	6.0	100	foo	blue	TRUE"},
			wantErr: true,
		},
		{
			name: "empty boolean",
			args: args{"2017-05-06T19:52:57.601Z	6.0	100	foo	blue	"},
			wantErr: true,
		},
		{
			name: "empty float",
			args: args{"2017-05-06T19:52:57.601Z		100	foo	blue	TRUE"},
			wantErr: true,
		},
		{
			name: "empty integer",
			args: args{"2017-05-06T19:52:57.601Z	6.0		foo	blue	TRUE"},
			wantErr: true,
		},
		{
			name: "bad timestamp",
			args: args{"2017-05-06aT19:52:57.601Z	6.0	100	foo	blue	TRUE"},
			wantErr: true,
		},
		{
			name: "bad boolean",
			args: args{"2017-05-06aT19:52:57.601Z	6.0	100	foo	blue	TaRUE"},
			wantErr: true,
		},
		{
			name: "bad float",
			args: args{"2017-05-06T19:52:57.601Z	6a.0	100	foo	blue	TRUE"},
			wantErr: true,
		},
		{
			name: "bad integer",
			args: args{"2017-05-06T19:52:57.601Z	6.0	100.3	foo	blue	TRUE"},
			wantErr: true,
		},
		{
			name:    "empty line",
			args:    args{""},
			wantErr: true,
		},
		{
			name: "missing column",
			args: args{"2017-05-06T19:52:57.601Z	6.0	100	foo"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Tsdata{
				checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "Some general comments about this file on a single line, not tab-delimited",
				Comments:        []string{"ISO8601 timestamp", "column2 notes", "NA", "column4 notes", "column5 notes", "column6 notes"},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			}
			fields, err := d.ValidateLine(tt.args.line)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Tsdata.ValidateLine() err %v, expected a non-nil error", err)
				}
			} else {
				if !stringSliceEqual(fields, tt.fields) {
					t.Errorf("Tsdata.ValidateLine() fields %v, expected %v", fields, tt.fields)
				}
			}
		})
	}
}

func TestTsdata_Header(t *testing.T) {
	type args struct {
		header string
	}
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
				FileDescription: "Some general comments about this file on a single line, not tab-delimited",
				Comments:        []string{"ISO8601 timestamp", "column2 notes", "NA", "column4 notes", "column5 notes", "column6 notes"},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			},
			header: "fileType\nproject\nSome general comments about this file on a single line, not tab-delimited\nISO8601 timestamp	column2 notes	NA	column4 notes	column5 notes	column6 notes\ntime	float	integer	text	category	boolean\nNA	m/s	km	NA	NA	NA\ntime	speed	distance	notes	color	hasTail",
		},
		{
			name: "header with no column comments",
			fields: tsdataFields{
				checkers:        []func(string) bool{checkTime, checkFloat, checkInteger, checkText, checkText, checkBoolean},
				FileType:        "fileType",
				Project:         "project",
				FileDescription: "Some general comments about this file on a single line, not tab-delimited",
				Comments:        []string{},
				Types:           []string{"time", "float", "integer", "text", "category", "boolean"},
				Units:           []string{"NA", "m/s", "km", "NA", "NA", "NA"},
				Headers:         []string{"time", "speed", "distance", "notes", "color", "hasTail"},
			},
			header: "fileType\nproject\nSome general comments about this file on a single line, not tab-delimited\nNA	NA	NA	NA	NA	NA\ntime	float	integer	text	category	boolean\nNA	m/s	km	NA	NA	NA\ntime	speed	distance	notes	color	hasTail",
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
