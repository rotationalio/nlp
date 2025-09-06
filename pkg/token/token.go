/*
[Token] is used for word tokens, word stem tokens, etc.

Usage example:

	// Create a new [Token]
	token := token.New("aardvarks") // Token

	// Get the token as other types
	stringToken := token.String() // string
	runeToken := token.Runes() // []rune
	byteToken := token.Bytes() // []byte

	// Get the number of runes and bytes in a token
	runeCount := token.Len() // 9 for "aadrdvarks"
	byteCount := token.ByteLen() // acts like the len() function
*/
package token

import (
	"unicode/utf8"
)

// A single word token.
type Token struct {
	token string
}

// Returns a new [Token] from a string.
func New(token string) Token {
	return Token{
		token: token,
	}
}

// Returns the number of UTF-8 runes in the token.
func (t *Token) Len() int {
	return utf8.RuneCountInString(t.token)
}

// Returns the number of bytes in the token.
func (t *Token) ByteLen() int {
	return len(t.token)
}

// Returns the token as a string.
func (t *Token) String() string {
	return t.token
}

// Returns the token as a slice of runes.
func (t *Token) Runes() []rune {
	return []rune(t.token)
}

// Returns the token as a slice of bytes.
func (t *Token) Bytes() []byte {
	return []byte(t.token)
}
