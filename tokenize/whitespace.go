package tokenize

import (
	"fmt"
	"slices"
	"strings"
)

// Tokenize text by whitespace  Splits a string like [strings.Fields] would do.
//
// Example:
//
//	"The quick brown fox jumped over the quicker-- 105.4% quicker, in fact--
//	\n
//	brown fox because it owed the quicker fox $3.14!"
//
// Tokens:
//
//	{"The", "quick", "brown", "fox", "jumped", "over", "the", "quicker--",
//	"105.4%", "quicker,", "in", "fact--", "brown", "fox", "because", "it",
//	"owed", "the", "quicker", "fox", "$3.14!"}
const REGEX_WHITESPACE = `\s*\S+\s*`

// Tokenize the chunk of text using whitespace. Splits the string using the
// [strings.Fields] function family. Cached for faster subsequent calls.
//
// Example:
//
//	"The quick brown fox jumped over the quicker-- 105.4% quicker, in fact--
//	\n
//	brown fox because it owed the quicker fox $3.14!"
//
// Tokens:
//
//	{"The", "quick", "brown", "fox", "jumped", "over", "the", "quicker--",
//	"105.4%", "quicker,", "in", "fact--", "brown", "fox", "because", "it",
//	"owed", "the", "quicker", "fox", "$3.14!"}
func WhitespaceTokenize(chunk string) (tokens []Token) {
	// TODO: cache these word tokens
	// TODO: make a WhitespaceTokenizer and make this its Tokenize() function
	for word := range strings.FieldsSeq(chunk) {
		tokens = append(tokens, Token{Token: word})
	}
	return tokens
}

// TODO convert to a test for REGEX_WHITESPACE and WhitespaceTokenize
func TestWhitespaceTokenize() {
	text := Text{"The quick brown fox, Mr. Fox, jumped over the quicker-- 105.4% quicker, in fact-- \n brown fox, the Hon. Judge Fox, because it owed the quicker fox $3.14! Isn't that amazing!?\n I think so!\t Crazy times, indeed."}
	fmt.Println(text)
	for _, s := range text.Sentences() {
		fmt.Printf("Sentence: [%s]\n", s)
	}
}

// ############################################################################
// TODO MOVE TO TEXT PACKAGE
// ############################################################################

// Returns the sentences in the [Text] as a string slice. Cached for faster
// subsequent calls.
func (t *Text) Sentences() []string {
	words := WhitespaceTokenize(t.Text)
	var (
		sentence  string
		sentences []string
		prevWord  Token
	)
	for _, word := range words {
		if prevWord.Token == "" || prevWord.EndsSentence() {
			// No space when starting a sentence
			sentence = word.Token
		} else {
			// Add a space if it isn't the start
			sentence += " " + word.Token
		}
		prevWord = word

		// Append sentence if this word ends it
		if word.EndsSentence() {
			sentences = append(sentences, sentence)
		}
	}
	// TODO cache the sentences
	return sentences
}

// Returns the count of the words in the [Text].
func (t *Text) WordCount() int {
	return len(WhitespaceTokenize(t.Text))
}

// Returns the count of the sentences in the [Text].
func (t *Text) SentenceCount() int {
	words := WhitespaceTokenize(t.Text)
	var count int
	for _, word := range words {
		if word.EndsSentence() {
			count += 1
		}
	}
	return count
}

// ############################################################################
// TODO MOVE TO TOKEN PACKAGE
// ############################################################################

// If the token ends in sentence punctuation and is not an abbreviation or
// initialism then it ends a sentence.
// TODO: probably do something else and not have this as a token method
func (t *Token) EndsSentence() bool {
	if t.Token == "" {
		return false
	}

	//TODO only do this specific one for English tokens, if it is possible to do so
	// True if the token ends in punctuation
	endsInPunct := strings.LastIndexAny(t.Token, SENTENCE_PUNCTUATION_ENGLISH) == len(t.Token)-1
	// True if the word is one like 'Dr.' which usually do not end a sentence
	// but end in sentence punctuation
	isStopWord := slices.Contains(SENTENCE_STOP_WORDS, t.Token)
	// True if there is more than one period such as 'Ph.D.' or 'F.B.I.'
	isInitialism := 1 < strings.Count(t.Token, ".")
	return endsInPunct && !(isStopWord || isInitialism)
}

// ############################################################################
// TODO MOVE TO LANGUAGE PACKAGE?
// ############################################################################

// Punctuation marks which end an English sentence.
var SENTENCE_PUNCTUATION_ENGLISH = ".!?"

// Words that look like they end a sentence but usually do not. Non-exhaustive.
var SENTENCE_STOP_WORDS = []string{"Mr.", "Mrs.", "Ms.", "Dr.", "Hon."}

// ############################################################################
//XXX temporary!
// ############################################################################

type Text struct {
	Text string
}

type Token struct {
	Token string
}
