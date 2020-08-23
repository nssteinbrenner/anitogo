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
	err := tkns.insert(-1, tkn)
	if err == nil {
		t.Errorf("expected error %s, got nil", indexTooSmallErr)
	} else if err != nil && len(*tkns) > 0 {
		t.Errorf("expected insert to fail on error, but token was inserted")
	}
	err = tkns.insert(100, tkn)
	if err == nil {
		t.Errorf("expected error %s, got nil", indexTooLargeErr)
	} else if err != nil && len(*tkns) > 0 {
		t.Errorf("expected insert to fail on error, but token was inserted")
	}
	err = tkns.insert(0, tkn)
	if err != nil {
		t.Errorf("expected token to insert successfully, but got %s", err.Error())
	} else if err == nil && len(*tkns) == 0 {
		t.Errorf("expected insert to succeed without error returned, but token was not inserted")
	}
	oldUUID := (*tkns)[0].UUID
	err = tkns.insert(0, tkn)
	if err != nil {
		t.Errorf("expected token to insert successfully, but got %s", err.Error())
	} else if err == nil && len(*tkns) > 1 {
		t.Errorf("expected insert to fail for duplicate token, but token was inserted")
	} else if (*tkns)[0].UUID != oldUUID {
		t.Errorf("expected token at 0 index stay the same, but it was replaced")
	}
	tkn.Content = "test1"
	err = tkns.insert(0, tkn)
	if err != nil {
		t.Errorf("expected token to insert successfully, but got %s", err.Error())
	} else if err == nil && len(*tkns) > 1 {
		t.Errorf("expected insert to fail for duplicate token, but token was inserted")
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
	_, err := tkns.getList(1, tkn, tkn)
	if err == nil {
		t.Error("expected error, got nil")
	}
	tkn.UUID = "test"
	(*tkns)[0].UUID = "test"
	tkn1 := &token{
		Category: tokenCategoryUnknown,
		Content:  "test1",
		Enclosed: false,
		UUID:     "test1",
	}
	_, err = tkns.getList(-1, tkn, tkn1)
	if err == nil {
		t.Error("expected error, got nil")
	}
	*tkns = append(*tkns, tkn1)
	_, err = tkns.getList(-1, tkn, tkn1)
	if err != nil {
		t.Error("expected nil, got error")
	}
}

func TestTokensGetListFlag(t *testing.T) {
	tkns := &tokens{}
	retTkns := tkns.getListFlag(-1)
	if len(retTkns) != 0 {
		t.Error("expected empty token slice")
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
	i, err := tkns.getIndex(tkn, -1)
	if err == nil {
		t.Errorf("expected error %s, got nil", indexTooSmallErr)
	}
	i, err = tkns.getIndex(tkn, 100)
	if err == nil {
		t.Errorf("expected error %s, got nil", indexTooLargeErr)
	}
	tkn1 := (*tkns)[0]
	i, err = tkns.getIndex(*tkn1, 0)
	if err != nil {
		t.Errorf("expected success, got err %s", err.Error())
	}
	if i != 0 {
		t.Errorf("expected 0, got %d", i)
	}
	tkn2 := token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
		UUID:     "test1",
	}
	i, err = tkns.getIndex(tkn2, 0)
	if err != nil {
		t.Errorf("expected nil, got error %s", err.Error())
	}
	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}
}

func TestTokensFindPrevious(t *testing.T) {
	tkns := &tokens{}
	tkn := token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	_, _, err := tkns.findPrevious(tkn, -1)
	if err == nil {
		t.Error("expected error, got nil")
	}
	tkn1 := token{}
	*tkns = append(*tkns, &tkn1)
	_, _, err = tkns.findPrevious(tkn1, -1)
	if err != nil {
		t.Error("expected nil, got error")
	}
}

func TestTokensFindNext(t *testing.T) {
	tkns := &tokens{}
	tkn := token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	_, _, err := tkns.findNext(tkn, -1)
	if err == nil {
		t.Error("expected error, got nil")
	}
	tkn1 := token{}
	*tkns = append(*tkns, &tkn1)
	_, _, err = tkns.findPrevious(tkn1, -1)
	if err != nil {
		t.Error("expected nil, got error")
	}
}

func TestTokensIsTokenIsolated(t *testing.T) {
	tkns := &tokens{}
	tkn := token{}
	found, _ := tkns.isTokenIsolated(tkn)
	if found {
		t.Error("expected false, got true")
	}
	tkn = token{
		Category: tokenCategoryUnknown,
		Content:  "test",
		Enclosed: false,
	}
	_, err := tkns.isTokenIsolated(tkn)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTokensFindInTokens(t *testing.T) {
	tkns := &tokens{}
	_, found := tkns.findInTokens(1)
	if found != false {
		t.Errorf("expected false, got true")
	}
}
