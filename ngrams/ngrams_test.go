package ngrams_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/ngrams"
)

func TestRune1grams(t *testing.T) {
	expected := [][]rune{[]rune("a"), []rune("b"), []rune("c"), []rune("d"), []rune("e"), []rune("f")}
	actual := ngrams.Ngrams([]rune("abcdef"), 1)
	require.Equal(t, expected, actual)
}

func TestRune3grams(t *testing.T) {
	expected := [][]rune{[]rune("abc"), []rune("bcd"), []rune("cde"), []rune("def")}
	actual := ngrams.Ngrams([]rune("abcdef"), 3)
	require.Equal(t, expected, actual)
}

func TestRuneTrigrams(t *testing.T) {
	expected := [][]rune{[]rune("abc"), []rune("bcd"), []rune("cde"), []rune("def")}
	actual := ngrams.Trigrams([]rune("abcdef"))
	require.Equal(t, expected, actual)
}

func TestRune5grams(t *testing.T) {
	expected := [][]rune{[]rune("abcde"), []rune("bcdef")}
	actual := ngrams.Ngrams([]rune("abcdef"), 5)
	require.Equal(t, expected, actual)
}

func TestRune6grams(t *testing.T) {
	expected := [][]rune{[]rune("abcdef")}
	actual := ngrams.Ngrams([]rune("abcdef"), 6)
	require.Equal(t, expected, actual)
}

func TestNgramsTooShortSequence(t *testing.T) {
	actual := ngrams.Ngrams([]rune("ab"), 3)
	require.Nil(t, actual)
}
