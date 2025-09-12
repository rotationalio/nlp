package tokenize

import (
	"slices"
	"strings"

	"go.rtnl.ai/nlp/language"
)

type SentenceSegmenter struct {
	lang                language.Language
	punctuation         string
	stopWords           []string
	whitespaceTokenizer *WhitespaceTokenizer
}

// Ensure [SentenceSegmenter] meets the [Tokenizer] interface requirements.
var _ Tokenizer = &SentenceSegmenter{}

// Returns a new [SentenceSegmenter]. Default language is [language.English].
func NewSentenceSegmenter(opts ...SentenceSegmenterOption) *SentenceSegmenter {
	// Set options
	segmenter := &SentenceSegmenter{}
	for _, fn := range opts {
		fn(segmenter)
	}

	// Set default language
	if segmenter.lang == language.Unknown {
		segmenter.lang = language.English
	}

	// Set language punctuation and stop words
	switch segmenter.lang {
	case language.English:
		segmenter.punctuation = language.SENTENCE_PUNCTUATION_ENGLISH
		segmenter.stopWords = language.SENTENCE_STOP_WORDS_ENGLISH

	}

	// Init WhitespaceTokenizer
	segmenter.whitespaceTokenizer = NewWhitespaceTokenizer()

	return segmenter
}

// Returns the sentences in the text chunk as a string slice. ALWAYS returns nil
// for the error.
func (s *SentenceSegmenter) Tokenize(chunk string) (sentences []string, alwaysNil error) {
	var (
		sentence string
		prevWord string
	)

	words, _ := s.whitespaceTokenizer.Tokenize(chunk) // error is ALWAYS nil
	for _, word := range words {
		// Check if it's a new sentence
		if prevWord == "" || endsSentence(prevWord, s.punctuation, s.stopWords) {
			// No space when starting a sentence
			sentence = word
		} else {
			// Add a space if it isn't the start
			sentence += " " + word
		}

		// Set prevWord for next iteration
		prevWord = word

		// Append sentence if this word ends it
		if endsSentence(word, s.punctuation, s.stopWords) {
			sentences = append(sentences, sentence)
		}
	}

	// Append the last sentence
	if sentence != "" {
		sentences = append(sentences, sentence)
	}

	return sentences, nil
}

// ############################################################################
// Helpers
// ############################################################################

// If the string ends in sentence punctuation and is not an abbreviation or
// initialism or in the set of stop words then it ends a sentence.
func endsSentence(token, punctuation string, stopWords []string) bool {
	if token == "" {
		return false
	}

	// True if the token ends in punctuation
	endsInPunct := strings.LastIndexAny(token, punctuation) == len(token)-1

	// True if the word is one like 'Dr.' which usually do not end a sentence
	// but end in sentence punctuation
	isStopWord := slices.Contains(stopWords, token)

	// True if there is more than one period such as 'Ph.D.' or 'F.B.I.'
	isInitialism := 1 < strings.Count(token, ".")

	return endsInPunct && !(isStopWord || isInitialism)
}

// ############################################################################
// RegexTokenizerOption
// ############################################################################

// SentenceSegmenterOption functions modify a [SentenceSegmenter].
type SentenceSegmenterOption func(t *SentenceSegmenter)

// Returns a function which sets the [language.Language] to use with the
// [SentenceSegmenter].
func SentenceSegmenterWithLanguage(lang language.Language) SentenceSegmenterOption {
	return func(t *SentenceSegmenter) {
		t.lang = lang
	}
}
