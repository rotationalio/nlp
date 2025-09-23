package readability

// Returns the Flesch-Kincaid Reading Ease score. Returns the value 0.0 when the
// sentence and/or word count is zero.
func FleschKincaidReadingEase(wordCount, sentenceCount, syllableCount int) (score float64) {
	if sentenceCount == 0 || wordCount == 0 {
		return 0.0
	}

	return 206.835 - 1.015*(float64(wordCount)/float64(sentenceCount)) - 84.6*(float64(syllableCount)/float64(wordCount))
}

// Returns the Flesch-Kincaid grade level. Returns the value 0.0 when the
// sentence and/or word count is zero.
func FleschKincaidGradeLevel(wordCount, sentenceCount, syllableCount int) (score float64) {
	if sentenceCount == 0 || wordCount == 0 {
		return 0.0
	}

	return 0.39*(float64(wordCount)/float64(sentenceCount)) + 11.8*(float64(syllableCount)/float64(wordCount)) - 15.59
}

// 1362, 36, 2697
