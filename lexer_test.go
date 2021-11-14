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
			"SELECT FROM AS WHERE LIMIT",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenFrom, "FROM"},
				{tokenSpace, " "},
				{tokenAlias, "AS"},
				{tokenSpace, " "},
				{tokenWhere, "WHERE"},
				{tokenSpace, " "},
				{tokenLimit, "LIMIT"},
				{tokenEnd, ""},
			},
		},
		{
			"keywords and identifiers",
			"SELECT field FROM table AS t WHERE LIMIT",
			[]token{
				{tokenSelect, "SELECT"},
				{tokenSpace, " "},
				{tokenIdentifier, "field"},
				{tokenSpace, " "},
				{tokenFrom, "FROM"},
				{tokenSpace, " "},
				{tokenIdentifier, "table"},
				{tokenSpace, " "},
				{tokenAlias, "AS"},
				{tokenSpace, " "},
				{tokenIdentifier, "t"},
				{tokenSpace, " "},
				{tokenWhere, "WHERE"},
				{tokenSpace, " "},
				{tokenLimit, "LIMIT"},
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
