package printer

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResultObject struct {
	Result testCSVDataObject `json:"result"`
}

type testCSVDataObject struct {
	Schools  []testSchoolObject `json:"schools"`
	Pupils   []int              `json:"pupils"`
	Teachers []int              `json:"teachers"`
}

type testSchoolObject struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var (
	testColumnHeaders = []string{"SCHOOL", "PUPILS", "TEACHERS"}
	testRowExpression = "{result.schools.#.name,result.pupils,result.teachers}"
	testCSVObject     = &testResultObject{
		Result: testCSVDataObject{
			Schools: []testSchoolObject{
				{
					Name:    "School1",
					Address: "Address1",
				},
				{
					Name:    "School2",
					Address: "Address2",
				},
				{
					Name:    "School3",
					Address: "Address3",
				},
				{
					Name:    "School4",
					Address: "Address4",
				},
			},
			Pupils:   []int{10, 20, 30, 40},
			Teachers: []int{1, 2, 3, 4},
		},
	}
)

func TestCsvPrinter_Print(t *testing.T) {
	cases := map[string]struct {
		expectedOutput   string
		noHeaders        bool
		headerAsComments bool
	}{
		"no settings specified": {
			expectedOutput: `SCHOOL,PUPILS,TEACHERS
School1,10,1
School2,20,2
School3,30,3
School4,40,4
`,
		},
		"print without headers": {
			expectedOutput: `School1,10,1
School2,20,2
School3,30,3
School4,40,4
`,
			noHeaders: true,
		},
		"print with headers as comments": {
			expectedOutput: `; SCHOOL,PUPILS,TEACHERS
School1,10,1
School2,20,2
School3,30,3
School4,40,4
`,
			headerAsComments: true,
		},
		"when specifying print no headers and print headers as comments - no headers should be printed with precedence": {
			expectedOutput: `School1,10,1
School2,20,2
School3,30,3
School4,40,4
`,
			noHeaders:        true,
			headerAsComments: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			out := strings.Builder{}
			printer := newCSVPrinter(testColumnHeaders, testRowExpression, c.noHeaders, c.headerAsComments)
			err := printer.Print(&testCSVObject, &out)
			require.NoError(t, err)
			assert.Equal(t, c.expectedOutput, out.String())
		})
	}
}

func TestCsvPrinter_Print_Failures(t *testing.T) {
	cases := map[string]struct {
		out            *strings.Builder
		headers        []string
		rowExpression  string
		object         interface{}
		err            error
		expectedOutput string
	}{
		"invalid JSON should cause a failure": {
			out:            &strings.Builder{},
			object:         make(chan int),
			headers:        testColumnHeaders,
			rowExpression:  testRowExpression,
			err:            &json.UnsupportedTypeError{Type: reflect.TypeOf(make(chan int))},
			expectedOutput: "",
		},
		"jagged input": {
			out:           &strings.Builder{},
			headers:       testColumnHeaders,
			rowExpression: testRowExpression,
			object: &testResultObject{Result: testCSVDataObject{
				Schools: []testSchoolObject{
					{
						Name: "School1",
					},
					{
						Name: "School2",
					},
					{
						Name: "School3",
					},
				},
				Pupils:   []int{10, 20},
				Teachers: []int{1, 2, 3},
			}},
			expectedOutput: "",
			err:            JaggedArrayError{maxAmount: 3, violatedAmount: 2, rowIndex: 0},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			p := newCSVPrinter(c.headers, c.rowExpression, false, false)
			err := p.Print(c.object, c.out)
			require.Error(t, err)
			assert.Equal(t, err, c.err)
			assert.Equal(t, c.expectedOutput, c.out.String())
		})
	}
}

func TestCsvPrinter_Print_EmptyData(t *testing.T) {
	cases := map[string]struct {
		out             *strings.Builder
		noHeader        bool
		headerAsComment bool
		expectedOutput  string
	}{
		"empty data should only print headers": {
			out:            &strings.Builder{},
			expectedOutput: "SCHOOL,PUPILS,TEACHERS\n",
		},
		"empty data with header as comment should print commented headers": {
			out:             &strings.Builder{},
			expectedOutput:  "; SCHOOL,PUPILS,TEACHERS\n",
			headerAsComment: true,
		},
		"empty data with no header set should print nothing": {
			out:            &strings.Builder{},
			expectedOutput: "",
			noHeader:       true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			p := newCSVPrinter(testColumnHeaders, testRowExpression, c.noHeader, c.headerAsComment)
			err := p.Print(&testResultObject{}, c.out)
			assert.NoError(t, err)
			assert.Equal(t, c.expectedOutput, c.out.String())
		})
	}
}

func TestCsvPrinter_ReadCSVOutputWithCommentedHeaders_Success(t *testing.T) {
	p := newCSVPrinter(testColumnHeaders, testRowExpression, false, true)
	out := &bytes.Buffer{}
	require.NoError(t, p.Print(testCSVObject, out))
	r := csv.NewReader(out)
	// Since Comment is per default not set, need to set it explicitly
	r.Comment = ';'
	records, err := r.ReadAll()
	require.NoError(t, err)
	for _, record := range records {
		assert.False(t, strings.HasPrefix(strings.Join(record, ","), ";"))
	}
}

func TestCsvPrinter_ReadCSVOutputWithCommentedHeaders_Failure(t *testing.T) {
	p := newCSVPrinter(testColumnHeaders, testRowExpression, false, true)
	out := &bytes.Buffer{}
	require.NoError(t, p.Print(testCSVObject, out))
	r := csv.NewReader(out)
	// Since Comment is per default not set, need to set it explicitly
	r.Comment = '#'
	records, err := r.ReadAll()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(strings.Join(records[0], ","), ";"))
}