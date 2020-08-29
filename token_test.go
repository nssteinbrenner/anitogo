package anitogo

import (
	"testing"
)

func TestTokensCheckFlags(t *testing.T) {
	tkn := &token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	res := tkn.checkFlags(8721)
	if res != true {
		t.Error("expected true, got false")
	}
	targSlice := []bool{true, true, false, false, true, true, false, false, false, false}
	for i := 0; i < 10; i++ {
		if tkn.checkFlags(i) != targSlice[i] {
			t.Errorf("expected %t, got %t", targSlice[i], tkn.checkFlags(i))
		}
	}
	tkn.Category = tokenCategoryBracket
	targSlice = []bool{true, true, true, true, false, false, true, true, false, false}
	for i := 0; i < 10; i++ {
		if tkn.checkFlags(i) != targSlice[i] {
			t.Errorf("expected %t, got %t", targSlice[i], tkn.checkFlags(i))
		}
	}
	tkn.Category = tokenCategoryDelimiter
	targSlice = []bool{true, true, false, false, true, true, false, false, true, true}
	for i := 0; i < 10; i++ {
		if tkn.checkFlags(i) != targSlice[i] {
			t.Errorf("expected %t, got %t", targSlice[i], tkn.checkFlags(i))
		}
	}
	tkn.Category = tokenCategoryIdentifier
	targSlice = []bool{true, true, false, false, true, true, false, false, false, false}
	for i := 0; i < 10; i++ {
		if tkn.checkFlags(i) != targSlice[i] {
			t.Errorf("expected %t, got %t", targSlice[i], tkn.checkFlags(i))
		}
	}
	tkn.Category = tokenCategoryInvalid
	targSlice = []bool{true, true, false, false, true, true, false, false, false, false}
	for i := 0; i < 10; i++ {
		if tkn.checkFlags(i) != targSlice[i] {
			t.Errorf("expected %t, got %t", targSlice[i], tkn.checkFlags(i))
		}
	}
}

func TestTokensInsert(t *testing.T) {
	tkns := &tokens{}
	tkn := token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	tkns.insert(-1, tkn)
	if len(*tkns) > 0 {
		t.Errorf("expected insert to fail on negative index, but token was inserted")
	}
	tkns.insert(100, tkn)
	if len(*tkns) > 0 {
		t.Errorf("expected insert to fail on too large index, but token was inserted")
	}
	tkns.insert(0, tkn)
	if len(*tkns) == 0 {
		t.Errorf("expected insert to succeed, but token was not inserted")
	}
	oldUUID := (*tkns)[0].UUID
	tkns.insert(0, tkn)
	if len(*tkns) > 1 {
		t.Errorf("expected insert to not do anything on duplicate token, but token was inserted")
	} else if (*tkns)[0].UUID != oldUUID {
		t.Errorf("expected token at 0 index stay the same, but it was replaced")
	}
	tkn.Content = "test1"
	tkns.insert(0, tkn)
	if len(*tkns) > 1 {
		t.Errorf("expected insert to not do anything on duplicate token, but token was inserted")
	} else if (*tkns)[0].UUID == oldUUID {
		t.Errorf("expected token at 0 index be replaced, but it stayed the same")
	}
}

func TestTokensAppendToken(t *testing.T) {
	tkns := &tokens{}
	tkns.appendToken(token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	})
	if len(*tkns) == 0 {
		t.Error("token was not appended")
	}
	if (*tkns)[0].UUID == "" {
		t.Error("token did not have a UUID added")
	}
}

func TestTokensGet(t *testing.T) {
	tkns := &tokens{}
	_, found := tkns.get(-1)
	if found != false {
		t.Errorf("expected false, got true")
	}
	tkns.appendToken(token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	})
	_, found = tkns.get(0)
	if found != true {
		t.Errorf("expected true, got false")
	}
}

func TestTokensGetList(t *testing.T) {
	tkns := &tokens{}
	tkn := &token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	*tkns = append(*tkns, tkn)
	retTkns := tkns.getList(1, tkn, tkn)
	if len(retTkns) > 0 {
		t.Error("expected empty tokens")
	}
	tkn.UUID = "test"
	(*tkns)[0].UUID = "test"
	tkn1 := &token{
		Category: tokenCategoryUnknown,
		Content:  "test1",
		Enclosed: false,
		UUID:     "test1",
	}
	retTkns = tkns.getList(-1, tkn, tkn1)
	if len(retTkns) > 0 {
		t.Error("expected empty tokens")
	}
	*tkns = append(*tkns, tkn1)
	retTkns = tkns.getList(tokenFlagsUnknown, tkn, tkn1)
	if len(retTkns) == 0 {
		t.Error("expected tokens with length greater than 0")
	}
}

func TestTokensGetListFlag(t *testing.T) {
	tkn := &token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	tkns := &tokens{}
	retTkns := tkns.getListFlag(-1)
	if len(retTkns) != 0 {
		t.Error("expected empty token slice")
	}
	tkns.appendToken(*tkn)
	tkns.appendToken(*tkn)
	tkns.appendToken(*tkn)
	retTkns = tkns.getListFlag(tokenFlagsUnknown)
	if len(retTkns) == 0 {
		t.Error("expected tokens with length greater than 0")
	}
}

func TestTokensGetIndex(t *testing.T) {
	tkns := &tokens{}
	tkn := token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	tkns.appendToken(tkn)
	i := tkns.getIndex(tkn, -1)
	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}
	i = tkns.getIndex(tkn, 100)
	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}
	tkn1 := (*tkns)[0]
	i = tkns.getIndex(*tkn1, 0)
	if i != 0 {
		t.Errorf("expected 0, got %d", i)
	}
	tkn2 := token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
		UUID:     "test1",
	}
	i = tkns.getIndex(tkn2, 0)
	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}
}

func TestTokensFindPrevious(t *testing.T) {
	tkns := &tokens{}
	tkn := token{
		Category: tokenCategoryBracket,
		Content:  "test",
		Enclosed: false,
	}
	retTkn, found := tkns.findPrevious(tkn, -1)
	if found {
		t.Error("expected false, got true")
	}
	if !retTkn.empty() {
		t.Error("expected empty token")
	}
	tkn1 := token{}
	*tkns = append(*tkns, &tkn1)
	retTkn, found = tkns.findPrevious(tkn1, -1)
	if found {
		t.Error("expected false, got true")
	}
	if !retTkn.empty() {
		t.Error("expected empty token")
	}
	psr := getTestParser("")
	retTkn, found = psr.tokenizer.tokens.findPrevious(*(*psr.tokenizer.tokens)[len(*psr.tokenizer.tokens)-1], tokenFlagsUnknown)
	if !found {
		t.Error("expected true, got false")
	}
	if retTkn.empty() {
		t.Error("expected non-empty token")
	}
}

func TestTokensFindNext(t *testing.T) {
	tkns := &tokens{}
	tkn := token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	retTkn, found := tkns.findNext(token{}, -1)
	if found {
		t.Error("expected false, got true")
	}
	if !retTkn.empty() {
		t.Error("expected empty token")
	}
	retTkn, found = tkns.findNext(tkn, -1)
	if found {
		t.Error("expected false, got true")
	}
	if !retTkn.empty() {
		t.Error("expected empty token")
	}
	tkn1 := token{}
	*tkns = append(*tkns, &tkn1)
	retTkn, found = tkns.findNext(tkn1, -1)
	if found {
		t.Error("expected false, got true")
	}
	if !retTkn.empty() {
		t.Error("expected empty token")
	}
	psr := getTestParser("")
	retTkn, found = psr.tokenizer.tokens.findNext(*(*psr.tokenizer.tokens)[0], tokenFlagsUnknown)
	if !found {
		t.Error("expected true, got false")
	}
	if retTkn.empty() {
		t.Error("expected non-empty token")
	}
}

func TestTokensIsTokenIsolated(t *testing.T) {
	tkns := &tokens{}
	tkn := token{}
	if tkns.isTokenIsolated(tkn) {
		t.Error("expected false, got true")
	}
	tkn = token{
		Category: tokenCategoryIdentifier,
		Content:  "test",
		Enclosed: false,
	}
	if tkns.isTokenIsolated(tkn) {
		t.Error("expected false, got true")
	}

	psr := getTestParser("")
	(*psr.tokenizer.tokens)[len(*psr.tokenizer.tokens)-2].Category = tokenCategoryBracket
	if psr.tokenizer.tokens.isTokenIsolated(*(*psr.tokenizer.tokens)[len(*psr.tokenizer.tokens)-1]) {
		t.Error("expected false, got true")
	}

	var testTkn *token
	psr = getTestParser("")
	for _, v := range *psr.tokenizer.tokens {
		if v.Content == "Toradora!" {
			testTkn = v
			break
		}
	}
	if !psr.tokenizer.tokens.isTokenIsolated(*testTkn) {
		t.Error("expected true, got false")
	}
}

func TestTokensFindInTokens(t *testing.T) {
	tkns := &tokens{}
	_, found := tkns.findInTokens(1)
	if found != false {
		t.Error("expected false, got true")
	}
	psr := getTestParser("")
	_, found = psr.tokenizer.tokens.findInTokens(tokenFlagsUnknown)
	if !found {
		t.Error("expected true, got false")
	}
}
