package anitogo

import (
	"testing"
)

func TestTokenizerTokenize(t *testing.T) {
	tkns := &tokens{}
	elems := &Elements{}
	km := newKeywordManager()
	tkz := tokenizer{
		filename:       "",
		options:        DefaultOptions,
		tokens:         tkns,
		keywordManager: km,
		elements:       elems,
	}
	err := tkz.tokenize()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

/*func TestTokenizerValidateDelimiterTokens(t *testing.T) {
	tkns := &tokens{}
	elems := &Elements{}
	km := newKeywordManager()
	tkz := tokenizer{
		filename:       "",
		options:        DefaultOptions,
		tokens:         tkns,
		keywordManager: km,
		elements:       elems,
	}

    tkz.tokens.appendToken(token{
        Category: tokenCategoryDelimiter,
        Content: "test",
    })

    err := tkz.validateDelimiterTokens()

    (*tkz.tokens)[0].Category = tokenCategoryUnknown

    tkz.tokens.appendToken(token{
        Category: tokenCategoryDelimiter,
        Content: "test1",
    })
    err = tkz.validateDelimiterTokens()
}*/
