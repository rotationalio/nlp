package tokenlist

import (
	"go.rtnl.ai/nlp/token"
)

// A list of [token.Token]s.
type TokenList []token.Token

/*
[TokenList] is a list of [token.Token]s with useful features.

Usage:

	// Create a new [TokenList] from a slice of strings
	items := []string{"apple", "bananna", "zebra"}
	myTokens := tokenlist.New(items)

	// You can get a slice of strings back
	stringTokens := myTokens.Strings() // []string

	// You can get a new [TokenList] as a copy of another list
	aCopy := tokenlist.NewCopy(myTokens) // TokenList (copy of myTokens)

	// You can get an empty [TokenList] with a specific size and capacity
	emptyTokens := tokenlist.NewEmpty(0, 100) // TokenList (with no entries)
	emptyTokens = tokenlist.NewEmpty(10, 100) // TokenList (with 10 null string tokens)
	emptyTokens = tokenlist.NewEmpty(10, 5) // TokenList (capacity will be set to 10)

	// You can also use regular slice functions and operations on a [TokenList]
	length := len(myTokens) // 3
	myTokens = append(myTokens, myTokens[0]) // "apple", "bananna", "zebra", "apple"
	myTokens[0] = myTokens[1] // "bananna", "bananna", "zebra"
*/
func New(tokens []string) TokenList {
	tl := make([]token.Token, 0, len(tokens))
	for _, tok := range tokens {
		tl = append(tl, token.New(tok))
	}
	return tl
}

// Returns a new [TokenList] by copying another [TokenList].
func NewCopy(other TokenList) TokenList {
	tl := make([]token.Token, 0, len(other))
	tl = append(tl, other...)
	return tl
}

// Returns a new [TokenList] that is empty with a specific size and capacity. If
// capacity is smaller than size, then size will be used as the capacity.
func NewEmpty(size, capacity int) TokenList {
	if capacity < size {
		capacity = size
	}
	return make([]token.Token, size, capacity)
}

// Returns the tokens from the [TokenList] as a slice of strings.
func (t TokenList) Strings() []string {
	s := make([]string, 0, len(t))
	for _, tok := range t {
		s = append(s, tok.String())
	}
	return s
}
