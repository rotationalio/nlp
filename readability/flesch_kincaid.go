package readability

import "go.rtnl.ai/nlp/errors"

// Returns the Flesch-Kincaid Reading Ease score. Returns an [errors.ErrUndefinedValue]
// when the sentence and/or word count is zero.
func FleschKincaidReadingEase(wordCount, sentenceCount, syllableCount int) (score float64, err error) {
	if sentenceCount == 0 || wordCount == 0 {
		return 0.0, errors.ErrUndefinedValue
	}

	return 206.835 - 1.015*(float64(wordCount)/float64(sentenceCount)) - 84.6*(float64(syllableCount)/float64(wordCount)), nil
}

// Returns the Flesch-Kincaid grade level. Returns an [errors.ErrUndefinedValue]
// when the sentence and/or word count is zero.
func FleschKincaidGradeLevel(wordCount, sentenceCount, syllableCount int) (score float64, err error) {
	if sentenceCount == 0 || wordCount == 0 {
		return 0.0, errors.ErrUndefinedValue
	}

	return 0.39*(float64(wordCount)/float64(sentenceCount)) + 11.8*(float64(syllableCount)/float64(wordCount)) - 15.59, nil
}

// 1362, 36, 2697
