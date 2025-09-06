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
