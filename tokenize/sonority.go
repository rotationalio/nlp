package tokenize

import (
	"strings"
	"unicode"

	"go.rtnl.ai/nlp/errors"
	"go.rtnl.ai/nlp/language"
)

// ###############################################################################
// This Sonority Sequencing Principle (SSP) algorithm code has been ported to
// Go under the terms of the NLTK Project's Apache 2.0 license using the Python
// implementation described at https://github.com/nltk/nltk
// ###############################################################################

// Word syllable tokenizer that uses Sonority Sequencing Principle (SSP).
type SSPSyllableTokenizer struct {
	lang         language.Language
	runeScoreMap map[rune]int8
	vowels       string
}

// Ensure [SSPSyllableTokenizer] meets the [Tokenizer] interface requirements.
var _ Tokenizer = &SSPSyllableTokenizer{}

// Returns a new [SSPSyllableTokenizer] configured for the [language.Language]
// provided. If the language is unsupported, it will return
// [errors.ErrLanguageNotSupported].
func NewSSPSyllableTokenizer(lang language.Language) (*SSPSyllableTokenizer, error) {
	// NOTE: we only need the lower-case OR the upper-case runes to be added
	// to the `runeScoreMap` and `vowels` fields, and the implementation will
	// check for both upper and lower cases in the given runes.
	switch lang {
	case language.English:
		return &SSPSyllableTokenizer{
			lang: lang,
			runeScoreMap: mapRuneScores(map[int8][]rune{
				3: []rune("aeiouy"),      // vowels
				2: []rune("lmnrw"),       // nasals
				1: []rune("zvsf"),        // fricatives
				0: []rune("bcdgtkpqxhj"), // stops
			}),
			vowels: "aeiouy",
		}, nil
	}
	return nil, errors.ErrLanguageNotSupported
}

// Returns word syllables. ALWAYS returns nil for the error.
func (t *SSPSyllableTokenizer) Tokenize(word string) (syllables []string, alwaysNil error) {
	runeToken := []rune(word)

	// Deal with words that are two characters or less
	//
	// NOTE: This will result in any "word" that is just two punctuation chars
	// to be returned as a single syllable, which is a degenerative case for
	// punctuation that we probably are fine ignoring. There is a test to ensure
	// this behavior is recorded in a test.
	if len(runeToken) <= 2 {
		return []string{word}, nil
	}

	// Process runes by "trigrams", starting with the first rune being assigned
	// to the syllable so it isn't lost
	syllable := []rune{runeToken[0]}
	for i := 1; i <= len(runeToken)-2; i++ {
		focusRune := runeToken[i]

		// Treat certain characters as their own syllables:
		//  * punctuation
		//  * numbers
		//  * whitespace
		if unicode.IsPunct(focusRune) || unicode.IsNumber(focusRune) || unicode.IsSpace(focusRune) {
			syllables = append(syllables, string(syllable))
			syllables = append(syllables, string(focusRune))
			syllable = syllable[:0]
			continue
		}

		// Check for syllable breaks by score
		prevScore := t.runeScoreMap[runeToken[i-1]]
		focusScore := t.runeScoreMap[focusRune]
		nextScore := t.runeScoreMap[runeToken[i+1]]

		if prevScore >= focusScore && focusScore == nextScore {
			// Syllable breaks after the focus rune
			syllable = append(syllable, focusRune)
			syllables = append(syllables, string(syllable))
			syllable = syllable[:0]
			continue
		}

		if prevScore > focusScore && focusScore < nextScore {
			// Syllable breaks before the focus rune
			syllables = append(syllables, string(syllable))
			syllable = syllable[:0]
			syllable = append(syllable, focusRune)
			continue
		}

		// No syllable break
		syllable = append(syllable, focusRune)
	}

	// Append the last rune and syllable depending on what the last rune is
	if unicode.IsPunct(runeToken[len(runeToken)-1]) {
		syllables = append(syllables, string(syllable))
		syllables = append(syllables, string(runeToken[len(runeToken)-1]))
	} else {
		syllable = append(syllable, runeToken[len(runeToken)-1])
		syllables = append(syllables, string(syllable))
	}

	// Validate and return syllables
	return t.validateSyllables(syllables), nil
}

// Ensures all syllables have a vowel by appending syllables without vowels to
// the previous syllable. This function may return nil if there are no proper
// syllables.
func (t *SSPSyllableTokenizer) validateSyllables(syllables []string) (validatedSyllables []string) {
	// If there is 0 or 1 syllables return immediately
	if len(syllables) <= 1 {
		return syllables
	}

	// Process all of the syllables to ensure they all have vowels, if possible
	currentSyllable := ""
	for _, syllable := range syllables {
		// Skip empty string syllables
		if syllable == "" {
			continue
		}

		// Ignore certain "syllables"
		//  * punctuation
		//  * whitespace
		runeSyllable := []rune(syllable)
		if len(runeSyllable) == 1 {
			if unicode.IsPunct(runeSyllable[0]) || unicode.IsSpace(runeSyllable[0]) {
				if currentSyllable != "" {
					validatedSyllables = append(validatedSyllables, currentSyllable)
					currentSyllable = ""
				}
				continue
			}
		}

		// If there is no current syllable, set this one as the current one
		if currentSyllable == "" {
			currentSyllable = syllable
			continue
		}

		// If a syllable has a vowel (upper or lower cases), then add the
		// previous syllable and set this syllable as the start of the next
		// syllable
		if strings.ContainsAny(syllable, strings.ToLower(t.vowels)+strings.ToUpper(t.vowels)) {
			validatedSyllables = append(validatedSyllables, currentSyllable)
			currentSyllable = syllable
			continue
		}

		// The syllable does not have a vowel so add it to the previous syllable
		currentSyllable += syllable
	}

	// Append the last syllable if there is one
	if currentSyllable != "" {
		validatedSyllables = append(validatedSyllables, currentSyllable)
	}

	return validatedSyllables
}

// ############################################################################
// Helpers
// ############################################################################

// Returns a map of the runes to their hierarchy score.
func mapRuneScores(hierarchy map[int8][]rune) map[rune]int8 {
	revMap := make(map[rune]int8, 0)
	for score, runes := range hierarchy {
		for _, r := range runes {
			// Add the lower and uppercase runes
			revMap[unicode.ToLower(r)] = score
			revMap[unicode.ToUpper(r)] = score
		}
	}
	return revMap
}
