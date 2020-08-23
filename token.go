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

func (t *tokens) empty() bool {
	return (len((*t)) == 0)
}

func (t *tokens) appendToken(tkn token) {
	tkn.UUID = uuid.New().String()
	(*t) = append(*t, &tkn)
}

func (t *tokens) insert(index int, tkn token) error {
	if index == 0 {
		if len(*t) == 0 {
			tkn.UUID = uuid.New().String()
			(*t) = append((*t), &tkn)
			return nil
		} else if len(*t) == 1 && tkn.Content != (*t)[index].Content {
			tkn.UUID = uuid.New().String()
			(*t)[index] = &tkn
			return nil
		}
	} else if index > len(*t)-1 {
		return traceError(indexTooLargeErr)
	} else if index < 0 {
		return traceError(indexTooSmallErr)
	}
	if (*t)[index].Content == tkn.Content {
		return nil
	}
	tkn.UUID = uuid.New().String()
	startList := append((*t)[:index], &tkn)
	(*t) = append(startList, (*t)[index:]...)

	return nil
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

func (t *tokens) getList(flag int, begin, end *token) (tokens, error) {
	var err error
	beginIndex := -1
	if begin.UUID != "" {
		beginIndex, err = t.getIndex(*begin, 0)
		if err != nil {
			return tokens{}, err
		}
	}
	if beginIndex > len(*t)-1 {
		return tokens{}, traceError(indexTooLargeErr)
	}
	if beginIndex < 0 {
		return tokens{}, traceError(indexTooSmallErr)
	}
	endIndex := len((*t)) - 1
	if end.UUID != "" {
		endIndex, err = t.getIndex(*end, beginIndex)
		if err != nil {
			return tokens{}, err
		}
	}
	if endIndex+1 > len(*t) {
		return tokens{}, traceError(indexTooLargeErr)
	}
	if endIndex < 0 {
		return tokens{}, traceError(indexTooSmallErr)
	}
	if endIndex < beginIndex {
		return tokens{}, traceError(endIndexTooSmallErr)
	}

	if flag == -1 {
		return (*t)[beginIndex : endIndex+1], nil
	}
	var retTkns tokens
	for _, tkn := range (*t)[beginIndex : endIndex+1] {
		if tkn.checkFlags(flag) {
			retTkns = append(retTkns, tkn)
		}
	}
	return retTkns, nil
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

func (t *tokens) getIndex(tkn token, index int) (int, error) {
	if index > len(*t)-1 {
		return -1, traceError(indexTooLargeErr)
	}
	if index < 0 {
		return -1, traceError(indexTooSmallErr)
	}
	for i, tk := range (*t)[index:] {
		if tk.UUID == tkn.UUID {
			return i + index, nil
		}
	}
	return -1, nil
}

func (t *tokens) distance(begin, end *token) (int, error) {
	beginIndex, err := t.getIndex(*begin, 0)
	if err != nil {
		return -1, err
	}
	endIndex, err := t.getIndex(*end, beginIndex)
	if err != nil {
		return -1, err
	}

	return endIndex - beginIndex, nil
}

func (t *tokens) find(flag int) (*token, bool) {
	tkn, found := t.findInTokens(flag)
	return tkn, found
}

func (t *tokens) findPrevious(tkn token, flag int) (*token, bool, error) {
	newTkns := tokens{}
	tokenIndex, err := t.getIndex(tkn, 0)
	if err != nil {
		return &token{}, false, err
	}
	if tokenIndex > len(*t)-1 {
		return &token{}, false, traceError(indexTooLargeErr)
	}
	if tokenIndex < 0 {
		return &token{}, false, traceError(indexTooSmallErr)
	}
	if tkn.empty() {
		newTkns = append(newTkns, *t...)
		newTkns = reverseOrder(newTkns)
	} else {
		newTkns = append(newTkns, (*t)[:tokenIndex]...)
		newTkns = reverseOrder(newTkns)
	}
	retTkn, found := findInTokens(newTkns, flag)
	return retTkn, found, nil
}

func (t *tokens) findNext(tkn token, flag int) (*token, bool, error) {
	var newTkns tokens

	if tkn.empty() {
		return &token{}, false, nil
	}
	tokenIndex, err := t.getIndex(tkn, 0)
	if err != nil {
		return &token{}, false, err
	}
	if tokenIndex > len(*t)-1 {
		return &token{}, false, traceError(indexTooLargeErr)
	}
	if tokenIndex < 0 {
		return &token{}, false, traceError(indexTooSmallErr)
	}
	if !tkn.empty() && tokenIndex+1 <= len(*t)-1 {
		newTkns = (*t)[tokenIndex+1:]
	} else if !tkn.empty() && tokenIndex <= len(*t)-1 {
		newTkns = (*t)[tokenIndex:]
	} else {
		return &token{}, false, nil
	}
	retTkn, found := findInTokens(newTkns, flag)
	return retTkn, found, nil
}

func (t *tokens) isTokenIsolated(tkn token) (bool, error) {
	if tkn.empty() {
		return false, nil
	}
	previousToken, found, err := t.findPrevious(tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	if previousToken.Category != tokenCategoryBracket {
		return false, nil
	}
	nextToken, found, err := t.findNext(tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	if nextToken.Category != tokenCategoryBracket {
		return false, nil
	}

	return true, nil
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
