/*
[TokenList] is a list of [token.Token]s with useful features.

Usage:

	// Create a new [TokenList] from a slice of strings
	myTokens := tokenlist.New([]string{"apple", "bananna", "zebra"})

	// Get the tokens as a slice of [token.Token]
	tokens := myTokens.Tokens() // []Token

	// You can get a slice of strings instead of tokens
	stringTokens := myTokens.Strings() // []string

	// You can get a new [TokenList] as a copy of another list
	aCopy := tokenlist.NewCopy(myTokens) // TokenList (copy of myTokens)

	// You can get an empty [TokenList] with a specific size and capacity
	emptyTokens := tokenlist.Empty(0, 100) // TokenList (with no entries)
	emptyTokens := tokenlist.Empty(10, 100) // TokenList (with 10 null string tokens)
	emptyTokens := tokenlist.Empty(10, 5) // TokenList (capacity will be set to 10)
*/
package tokenlist

import (
	"go.rtnl.ai/nlp/pkg/errors"
	"go.rtnl.ai/nlp/pkg/token"
)

// A list of [token.Token].
type TokenList struct {
	tokens []token.Token
}

// Returns a new [TokenList] from a slice of strings.
func New(tokens []string) *TokenList {
	tl := &TokenList{
		tokens: make([]token.Token, 0, len(tokens)),
	}
	for _, tok := range tokens {
		tl.tokens = append(tl.tokens, token.New(tok))
	}
	return tl
}

// Returns a new [TokenList] by copying another [TokenList].
func NewCopy(other *TokenList) *TokenList {
	tl := &TokenList{
		tokens: make([]token.Token, 0, len(other.tokens)),
	}
	tl.tokens = append(tl.tokens, other.tokens...)
	return tl
}

// Returns a new [TokenList] that is empty with a specific size and capacity. If
// capacity is smaller than size, then size will be used as the capacity.
func NewEmpty(size, capacity int) *TokenList {
	if capacity < size {
		capacity = size
	}
	return &TokenList{
		tokens: make([]token.Token, size, capacity),
	}
}

// Append a [Token] to the [TokenList].
func (t *TokenList) Append(tok token.Token) {
	if t.tokens == nil {
		t.tokens = make([]token.Token, 0)
	}
	t.tokens = append(t.tokens, tok)
}

// Replaces the [token.Token] at the [TokenList] index with another token.
func (t *TokenList) Replace(idx int, replacement token.Token) error {
	// No tokens to replace
	if t.tokens == nil {
		return errors.ErrInvalidIndex
	}

	// Index is out of bounds
	if len(t.tokens) <= idx {
		return errors.ErrInvalidIndex
	}

	t.tokens[idx] = replacement
	return nil
}

// Returns the number of tokens in the [TokenList].
func (t *TokenList) Len() int {
	return len(t.tokens)
}

// Returns the tokens from the [TokenList] as a slice of strings.
func (t *TokenList) Strings() []string {
	s := make([]string, 0, len(t.tokens))
	for _, tok := range t.tokens {
		s = append(s, tok.String())
	}
	return s
}

// Returns the tokens from the [TokenList] as a slice of [token.Token].
func (t *TokenList) Tokens() []token.Token {
	return t.tokens
}
