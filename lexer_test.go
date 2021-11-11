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
			"simple identifier",
			"SELECT",
			[]token{
				{tokenIdentifier, "SELECT"},
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
