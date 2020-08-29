package anitogo

import "github.com/google/uuid"

const (
	tokenCategoryUnknown = 1 << iota
	tokenCategoryBracket
	tokenCategoryDelimiter
	tokenCategoryIdentifier
	tokenCategoryInvalid
)

const (
	tokenFlagsNone    = 0
	tokenFlagsBracket = 1 << iota
	tokenFlagsNotBracket
	tokenFlagsDelimiter
	tokenFlagsNotDelimiter
	tokenFlagsIdentifier
	tokenFlagsNotIdentifier
	tokenFlagsUnknown
	tokenFlagsNotUnknown
	tokenFlagsValid
	tokenFlagsNotValid
	tokenFlagsEnclosed
	tokenFlagsNotEnclosed
	tokenFlagsMaskCategories = (tokenFlagsBracket | tokenFlagsNotBracket |
		tokenFlagsDelimiter | tokenFlagsNotDelimiter |
		tokenFlagsIdentifier | tokenFlagsNotIdentifier |
		tokenFlagsUnknown | tokenFlagsNotUnknown |
		tokenFlagsValid | tokenFlagsNotValid)
	tokenFlagsMaskEnclosed = tokenFlagsEnclosed | tokenFlagsNotEnclosed
)

type token struct {
	Category int
	Content  string
	Enclosed bool
	UUID     string
}

type tokens []*token

func (t *token) empty() bool {
	return (*t) == token{}
}

func (t *token) checkFlags(flag int) bool {
	if (flag & tokenFlagsMaskEnclosed) != 0 {
		var success bool
		if checkFlag(flag, tokenFlagsEnclosed) {
			success = t.Enclosed
		} else {
			success = !t.Enclosed
		}
		if !success {
			return false
		}
	}

	if (flag & tokenFlagsMaskCategories) != 0 {
		if t.checkCategory(tokenFlagsBracket, tokenFlagsNotBracket, flag, tokenCategoryBracket) {
			return true
		}
		if t.checkCategory(tokenFlagsDelimiter, tokenFlagsNotDelimiter, flag, tokenCategoryDelimiter) {
			return true
		}
		if t.checkCategory(tokenFlagsIdentifier, tokenFlagsNotIdentifier, flag, tokenCategoryIdentifier) {
			return true
		}
		if t.checkCategory(tokenFlagsUnknown, tokenFlagsNotUnknown, flag, tokenCategoryUnknown) {
			return true
		}
		if t.checkCategory(tokenFlagsNotValid, tokenFlagsValid, flag, tokenCategoryInvalid) {
			return true
		}
		return false
	}
	return true

}

func (t *token) checkCategory(fe, fn, sourceFlag, category int) bool {
	if checkFlag(sourceFlag, fe) {
		return t.Category == category
	} else if checkFlag(sourceFlag, fn) {
		return t.Category != category
	}
	return false
}

func (t *tokens) appendToken(tkn token) {
	tkn.UUID = uuid.New().String()
	(*t) = append(*t, &tkn)
}

func (t *tokens) insert(index int, tkn token) {
	if index == 0 {
		if len(*t) == 0 {
			tkn.UUID = uuid.New().String()
			(*t) = append((*t), &tkn)
			return
		} else if len(*t) == 1 && tkn.Content != (*t)[index].Content {
			tkn.UUID = uuid.New().String()
			(*t)[index] = &tkn
			return
		}
	} else if index > len(*t)-1 {
		return
	} else if index < 0 {
		return
	}
	if (*t)[index].Content == tkn.Content {
		return
	}
	tkn.UUID = uuid.New().String()
	startList := append((*t)[:index], &tkn)
	(*t) = append(startList, (*t)[index:]...)
}

func (t *tokens) update(tkns tokens) {
	(*t) = tkns
}

func (t *tokens) get(index int) (*token, bool) {
	if index <= len((*t))-1 && index >= 0 {
		return (*t)[index], true
	}
	return &token{}, false
}

func (t *tokens) getList(flag int, begin, end *token) tokens {
	beginIndex := -1
	if begin.UUID != "" {
		beginIndex = t.getIndex(*begin, 0)
	}
	if beginIndex < 0 {
		return tokens{}
	}
	endIndex := len(*t) - 1
	if end.UUID != "" {
		endIndex = t.getIndex(*end, beginIndex)
	}
	if endIndex < 0 {
		return tokens{}
	}

	if flag == -1 {
		return (*t)[beginIndex : endIndex+1]
	}
	var retTkns tokens
	for _, tkn := range (*t)[beginIndex : endIndex+1] {
		if tkn.checkFlags(flag) {
			retTkns = append(retTkns, tkn)
		}
	}
	return retTkns
}

func (t *tokens) getListFlag(flag int) tokens {
	if flag == -1 {
		return (*t)
	}
	var retTkns tokens
	for _, tkn := range *t {
		if tkn.checkFlags(flag) {
			retTkns = append(retTkns, tkn)
		}
	}
	return retTkns
}

func (t *tokens) getIndex(tkn token, index int) int {
	if index > len(*t)-1 {
		return -1
	}
	if index < 0 {
		return -1
	}
	for i, tk := range (*t)[index:] {
		if tk.UUID == tkn.UUID {
			return i + index
		}
	}
	return -1
}

func (t *tokens) distance(begin, end *token) int {
	beginIndex := t.getIndex(*begin, 0)
	endIndex := t.getIndex(*end, beginIndex)

	return endIndex - beginIndex
}

func (t *tokens) find(flag int) (*token, bool) {
	tkn, found := t.findInTokens(flag)
	return tkn, found
}

func (t *tokens) findPrevious(tkn token, flag int) (*token, bool) {
	newTkns := tokens{}
	tokenIndex := t.getIndex(tkn, 0)
	if tokenIndex < 0 {
		return &token{}, false
	}
	if tkn.empty() {
		newTkns = append(newTkns, *t...)
		newTkns = reverseOrder(newTkns)
	} else {
		newTkns = append(newTkns, (*t)[:tokenIndex]...)
		newTkns = reverseOrder(newTkns)
	}
	retTkn, found := findInTokens(newTkns, flag)
	return retTkn, found
}

func (t *tokens) findNext(tkn token, flag int) (*token, bool) {
	var newTkns tokens

	if tkn.empty() {
		return &token{}, false
	}
	tokenIndex := t.getIndex(tkn, 0)
	if tokenIndex < 0 {
		return &token{}, false
	}
	if !tkn.empty() && tokenIndex+1 <= len(*t)-1 {
		newTkns = (*t)[tokenIndex+1:]
	}
	retTkn, found := findInTokens(newTkns, flag)
	return retTkn, found
}

func (t *tokens) isTokenIsolated(tkn token) bool {
	if tkn.empty() {
		return false
	}
	previousToken, found := t.findPrevious(tkn, tokenFlagsNotDelimiter)
	if !found {
		return false
	}
	if previousToken.Category != tokenCategoryBracket {
		return false
	}
	nextToken, found := t.findNext(tkn, tokenFlagsNotDelimiter)
	if !found {
		return false
	}
	if nextToken.Category != tokenCategoryBracket {
		return false
	}

	return true
}

func (t *tokens) findInTokens(flag int) (*token, bool) {
	if len(*t) == 0 {
		return &token{}, false
	}
	for _, tkn := range *t {
		if tkn.checkFlags(flag) {
			return tkn, true
		}
	}
	return &token{}, false
}

func checkFlag(sourceFlag, targetFlag int) bool {
	return (sourceFlag & targetFlag) == targetFlag
}

func reverseOrder(t tokens) tokens {
	if len(t) == 0 {
		return tokens{}
	}
	lastIndex := len(t) - 1
	var tkns = t
	for i := 0; i < len(t)/2; i++ {
		tkns[i], tkns[lastIndex-i] = tkns[lastIndex-i], tkns[i]
	}
	return tkns
}

func findInTokens(t tokens, flag int) (*token, bool) {
	if len(t) == 0 {
		return &token{}, false
	}
	for _, tkn := range t {
		if tkn.checkFlags(flag) {
			return tkn, true
		}
	}
	return &token{}, false
}
