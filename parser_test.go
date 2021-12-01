package gosqlparser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		name              string
		input             string
		expectedStatement Statement
		err               error
	}{
		{
			"broken statement",
			"table1",
			nil,
			fmt.Errorf("expected SELECT, INSERT, UPDATE, DELETE, CREATE or DROP, but got identifier: table1"),
		},
		{
			"unfinished SELECT statement",
			"SELECT table1",
			nil,
			fmt.Errorf("expected FROM, delimeter, but got end: \"\""),
		},
		{
			"unfinished SELECT FROM statement",
			"SELECT table1 FROM",
			nil,
			fmt.Errorf("expected identifier, but got end: \"\""),
		},
		{
			"full CREATE TABLE query",
			"CREATE TABLE table1 (col1 INTEGER, col2 STRING)",
			&CreateTable{"table1", []ColumnDefinition{{"col1", TypeInteger}, {"col2", TypeString}}},
			nil,
		},
		{
			"full DROP TABLE query",
			"DROP TABLE table1",
			&DropTable{"table1"},
			nil,
		},
		{
			"broken DROP TABLE",
			"DROP table1",
			nil,
			fmt.Errorf("expected TABLE, but got identifier: \"table1\""),
		},
		{
			"simple SELECT FROM",
			"SELECT col1, col2 FROM table1",
			&Select{"table1", []string{"col1", "col2"}, nil, ""},
			nil,
		},
		{
			"simple DELETE FROM",
			"DELETE FROM table1",
			&Delete{"table1", nil},
			nil,
		},
		{
			"simple INSERT INTO",
			"INSERT INTO table1 (col1, col2) VALUES (\"val1\", 25)",
			&Insert{"table1", []string{"col1", "col2"}, []string{"\"val1\"", "25"}},
			nil,
		},
		{
			"simple UPDATE",
			"UPDATE table1 SET col1 = \"val1\", col2 = 2",
			&Update{"table1", []string{"col1", "col2"}, []string{"\"val1\"", "2"}, nil},
			nil,
		},
		{
			"SELECT FROM with LIMIT",
			"SELECT col1, col2 FROM table1 LIMIT 10",
			&Select{"table1", []string{"col1", "col2"}, nil, "10"},
			nil,
		},
		{
			"SELECT FROM with simple WHERE",
			"SELECT col1, col2 FROM table1 WHERE col1 == col2",
			&Select{"table1", []string{"col1", "col2"}, &Where{ExprOperation{ExprIdentifier{"col1"}, OperatorEquals, ExprIdentifier{"col2"}}}, ""},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			statement, err := Parse(testCase.input)
			if testCase.err != nil {
				if err == nil {
					t.Errorf("expected error \"%s\", but got nil", testCase.err)
				} else if testCase.err.Error() != err.Error() {
					t.Errorf("expected error \"%s\", but got \"%s\"", testCase.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error \"%s\"", err)
				}
			}

			if testCase.expectedStatement != nil && !reflect.DeepEqual(testCase.expectedStatement, statement) {
				t.Errorf("expected %v, but got %v", testCase.expectedStatement, statement)
			}
		})
	}
}
