package anitogo

import (
	"regexp"
	"testing"
)

func TestTokenizerTokenize(t *testing.T) {
	tkns := &tokens{}
	tkns.appendToken(token{
		Category: 2,
		Content:  "[",
		Enclosed: true,
	})
	tkns.appendToken(token{
		Category: 1,
		Content:  "HorribleSubs",
		Enclosed: true,
	})
	tkns.appendToken(token{
		Category: 2,
		Content:  "]",
		Enclosed: true,
	})
	tkns.appendToken(token{
		Category: 4,
		Content:  " ",
		Enclosed: false,
	})
	tkns.appendToken(token{
		Category: 1,
		Content:  "Gintama",
		Enclosed: false,
	})
	tkns.appendToken(token{
		Category: 4,
		Content:  " ",
		Enclosed: false,
	})
	tkns.appendToken(token{
		Category: 1,
		Content:  "-",
		Enclosed: false,
	})
	tkns.appendToken(token{
		Category: 4,
		Content:  " ",
		Enclosed: false,
	})
	tkns.appendToken(token{
		Category: 1,
		Content:  "111C",
		Enclosed: false,
	})
	tkns.appendToken(token{
		Category: 4,
		Content:  " ",
		Enclosed: false,
	})
	tkns.appendToken(token{
		Category: 2,
		Content:  "[",
		Enclosed: true,
	})
	tkns.appendToken(token{
		Category: 8,
		Content:  "1080p",
		Enclosed: true,
	})
	tkns.appendToken(token{
		Category: 2,
		Content:  "]",
		Enclosed: true,
	})
	psr := getTestParser("[HorribleSubs] Gintama - 111C [1080p].mkv")
	for i, v := range *psr.tokenizer.tokens {
		if v.Category != (*tkns)[i].Category {
			t.Errorf("expected %d, got %d", v.Category, (*tkns)[i].Category)
		}
		if v.Content != (*tkns)[i].Content {
			t.Errorf("expected \"%s\", got \"%s\"", v.Content, (*tkns)[i].Content)
		}
		if v.Enclosed != (*tkns)[i].Enclosed {
			t.Errorf("expected %t, got %t", v.Enclosed, (*tkns)[i].Enclosed)
		}
	}
}

func TestTokenizerSplitWith(t *testing.T) {
	re := regexp.MustCompile(" ")
	ret := splitWith(re, "", 0)
	if ret != nil {
		t.Error("expected nil, got slice")
	}
	ret = splitWith(re, "this is a test", 2)
	if len(ret) != 3 {
		t.Errorf("expected slice with len 3, got %d", len(ret))
	}

	ret = splitWith(re, "this is a test", -1)
	if len(ret) != 7 {
		t.Errorf("expected slice with len 7, got %d", len(ret))
	}
}
