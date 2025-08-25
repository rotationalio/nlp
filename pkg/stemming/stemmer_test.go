package stemming_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/stemming"
)

// ############################################################################
// Tests
// ############################################################################

// Tests all of the stemmers on large word sets.
func TestStemmers(t *testing.T) {
	testcases := []struct {
		TestName     string
		Stemmer      stemming.Stemmer
		InputPath    string
		ExpectedPath string
	}{
		{
			TestName:     "NoOpStemmer",
			Stemmer:      &stemming.NoOpStemmer{},
			InputPath:    "testdata/Porter2Stemmer/voc.txt", // uses Porter2 inputs
			ExpectedPath: "testdata/Porter2Stemmer/voc.txt", // no-op -> same as input
		},
		{
			TestName:     "Porter2Stemmer [English]",
			Stemmer:      mustNewPorter2Stemmer(enum.LanguageEnglish),
			InputPath:    "testdata/Porter2Stemmer/voc.txt",
			ExpectedPath: "testdata/Porter2Stemmer/output.txt",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.TestName, func(t *testing.T) {
			// Load 'input' data
			input, err := os.Open(tc.InputPath)
			require.NotNil(t, input, "unexpected nil input file")
			require.Nil(t, err, "error opening 'input' file")
			defer input.Close()

			// Load 'expected' data
			expected, err := os.Open(tc.ExpectedPath)
			require.NotNil(t, expected, "unexpected nil 'expected' file")
			require.Nil(t, err, "error opening 'expected' file")
			defer expected.Close()

			// Scan each line of the input and compare to the output of the stemmer
			inputScanner := bufio.NewScanner(input)
			expectedScanner := bufio.NewScanner(expected)
			for inputScanner.Scan() && expectedScanner.Scan() {
				in := inputScanner.Text()
				// NOTE: Uncomment below to see the 'in' word to debug panics
				// fmt.Printf("IN: %s\n", in)
				exp := expectedScanner.Text()
				act := tc.Stemmer.Stem(in)
				require.Equal(t, exp, act, fmt.Sprintf("wrong stem for |%s|: expected |%s|, got |%s|", in, exp, act))
			}
			// Ensure there were no scanning errors
			require.Nil(t, inputScanner.Err(), "error scanning 'input'")
			require.Nil(t, expectedScanner.Err(), "error scanning 'expected'")
		})
	}
}

// ############################################################################
// Benchmarking
// ############################################################################

//TODO: benchmark Porter2Stemmer

// ############################################################################
// Helpers
// ############################################################################

// Returns a new [quant.Porter2Stemmer] which supports the [quant.Language]
// given or panics on an error.
func mustNewPorter2Stemmer(lang enum.Language) (stemmer *stemming.Porter2Stemmer) {
	var err error
	if stemmer, err = stemming.NewPorter2Stemmer(lang); err != nil {
		panic(err)
	}
	return stemmer
}
