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
	// Get trigrams for the token
	runeToken := []rune(word)
	syllable := []rune{runeToken[0]} // start with first rune
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

	// Append the last syllable
	syllable = append(syllable, runeToken[len(runeToken)-1])
	syllables = append(syllables, string(syllable))

	// Validate and return syllables
	return t.validateSyllables(syllables), nil
}

// Ensures all syllables have a vowel by appending syllables without vowels to
// the previous syllable.
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

		// If there is no current syllable, set this one as the current one
		if currentSyllable == "" {
			currentSyllable = syllable
			continue
		}

		// Treat certain characters as their own syllable:
		//  * punctuation
		//  * whitespace
		runeSyllable := []rune(syllable)
		if len(runeSyllable) == 1 {
			if unicode.IsPunct(runeSyllable[0]) || unicode.IsSpace(runeSyllable[0]) {
				validatedSyllables = append(validatedSyllables, currentSyllable)
				validatedSyllables = append(validatedSyllables, syllable)
				currentSyllable = ""
				continue
			}
		}

		// If a syllable has a vowel, then add the previous syllable and set
		// this syllable as the start of the next syllable
		if strings.ContainsAny(syllable, t.vowels) {
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
			revMap[r] = score
		}
	}
	return revMap
}
