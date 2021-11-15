package gosqlparser

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedTokens []token
	}{
		{
			"select keyword only",
			"SELECT",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenEnd, ""},
			},
		},
		{
			"keywords and spaces",
			"SELECT FROM WHERE LIMIT",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenFrom, "FROM"},
				{tokenSpace, " "},
				{tokenWhere, "WHERE"},
				{tokenSpace, " "},
				{tokenLimit, "LIMIT"},
				{tokenEnd, ""},
			},
		},
		{
			"keywords and identifiers",
			"SELECT field FROM table WHERE LIMIT",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenIdentifier, "field"},
				{tokenSpace, " "},
				{tokenFrom, "FROM"},
				{tokenSpace, " "},
				{tokenIdentifier, "table"},
				{tokenSpace, " "},
				{tokenWhere, "WHERE"},
				{tokenSpace, " "},
				{tokenLimit, "LIMIT"},
				{tokenEnd, ""},
			},
		},
		{
			"equals only",
			"==",
			[]token{
				{tokenEquals, "=="},
				{tokenEnd, ""},
			},
		},
		{
			"keyword and equals",
			"SELECT ==",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenEquals, "=="},
				{tokenEnd, ""},
			},
		},
		{
			"keyword and broken equals at the end",
			"SELECT =",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenError, "expected ="},
			},
		},
		{
			"keyword and broken equals before keyword",
			"SELECT =FROM",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenError, "expected ="},
			},
		},
		{
			"full SELECT query",
			"SELECT c1, c2 FROM table1 WHERE c3 == c4 AND c5 == c6",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenIdentifier, "c1"},
				{tokenDelimeter, ","},
				{tokenSpace, " "},
				{tokenIdentifier, "c2"},
				{tokenSpace, " "},
				{tokenFrom, "FROM"},
				{tokenSpace, " "},
				{tokenIdentifier, "table1"},
				{tokenSpace, " "},
				{tokenWhere, "WHERE"},
				{tokenSpace, " "},
				{tokenIdentifier, "c3"},
				{tokenSpace, " "},
				{tokenEquals, "=="},
				{tokenSpace, " "},
				{tokenIdentifier, "c4"},
				{tokenSpace, " "},
				{tokenAnd, "AND"},
				{tokenSpace, " "},
				{tokenIdentifier, "c5"},
				{tokenSpace, " "},
				{tokenEquals, "=="},
				{tokenSpace, " "},
				{tokenIdentifier, "c6"},
				{tokenEnd, ""},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualTokens := tokenize(testCase.input)

			if !reflect.DeepEqual(testCase.expectedTokens, actualTokens) {
				t.Errorf("expected %v, but got %v", testCase.expectedTokens, actualTokens)
			}
		})
	}
}

func tokenize(input string) []token {
	tokens := lex(input)

	all := make([]token, 0)
	for {
		t, ok := <-tokens
		if !ok {
			break
		}

		all = append(all, t)
	}

	return all
}
