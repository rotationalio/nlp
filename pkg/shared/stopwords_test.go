package shared_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/shared"
)

func TestIsStopWordEnglish(t *testing.T) {
	// normal tests
	for _, word := range shared.StopWordsEnglish {
		require.True(t, shared.IsStopWord(word, enum.LanguageEnglish), "stop word not recognized")
	}

	// different case and spacing tests
	for _, word := range []string{"about", "ABOUT", "About", " ABOUT ", "AbOuT\n", "\n \t about \n \t"} {
		require.True(t, shared.IsStopWord(word, enum.LanguageEnglish), "stop word not recognized in uppercase")
	}

	// false tests
	for _, word := range []string{"careful", "zebra", "aardvark", "89u3dr4wuf7y8fy7t4h8", "RACECAR", "about1234", " !@#$%^&*() \n \t"} {
		require.False(t, shared.IsStopWord(word, enum.LanguageEnglish), "this is not a stop word")
	}
}
