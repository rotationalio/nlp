package stem_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/language"
	"go.rtnl.ai/nlp/stem"
)

// NOTE: Full stemming tests implemented in [TestStemmers] in 'stemmers_test.go'

// Ensure that if the user bypasses setting the [quant.Language] in the
// [quant.Porter2Stemmer] that [quant.Porter2Stemmer.StemEnglish] still works.
func TestPorter2EnglishBypass(t *testing.T) {
	stemmer := &stem.Porter2Stemmer{} // no 'lang' set
	in := "seaweed"
	exp := "seawe"
	act := stemmer.StemEnglish(in)
	require.Equal(t, exp, act, fmt.Sprintf("wrong stem for |%s|: expected |%s|, got |%s|", in, exp, act))
}

// Use the following test to test a single word stem, for debugging.
func TestPorter2Single(t *testing.T) {
	// NOTE: skipping this test unless we're debugging a word:
	t.SkipNow()

	// Debug a single, specific word
	in := "seaweed"
	exp := "seawe"
	act := mustNewPorter2Stemmer(language.English).Stem(in)
	require.Equal(t, exp, act, fmt.Sprintf("wrong stem for |%s|: expected |%s|, got |%s|", in, exp, act))
}
