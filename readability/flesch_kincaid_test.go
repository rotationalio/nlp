package readability_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/readability"
	"go.rtnl.ai/nlp/text"
)

func TestFleschKincaid(t *testing.T) {
	testcases := []struct {
		Name          string
		TextFilename  string
		ExpectedEase  float64
		ExpectedGrade float64
	}{
		{
			Name:          "DeclarationOfIndependence",
			TextFilename:  "testdata/declaration.txt",
			ExpectedEase:  17.247, // 29.60 at https://serpninja.io/tools/flesch-kincaid-calculator/
			ExpectedGrade: 20.253, // 18.28 at https://serpninja.io/tools/flesch-kincaid-calculator/
		},
		{
			Name:          "CatMat",
			TextFilename:  "testdata/cat_mat.txt",
			ExpectedEase:  119.19, // 116.15 at https://serpninja.io/tools/flesch-kincaid-calculator/
			ExpectedGrade: -2.620, // -1.45 at https://serpninja.io/tools/flesch-kincaid-calculator/
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			data, err := os.ReadFile(tc.TextFilename)
			require.NoError(t, err)
			require.NotNil(t, data)

			myText, err := text.New(string(data))
			require.NoError(t, err)
			require.NotNil(t, myText)

			actualEase := readability.FleschKincaidReadingEase(myText.WordCount(), myText.SentenceCount(), myText.SyllableCount())
			require.InDelta(t, tc.ExpectedEase, actualEase, 1e-3)

			actualGrade := readability.FleschKincaidGradeLevel(myText.WordCount(), myText.SentenceCount(), myText.SyllableCount())
			require.InDelta(t, tc.ExpectedGrade, actualGrade, 1e-3)
		})
	}
}

func TestFleschKincaidError(t *testing.T) {
	// Word count cannot be zero (div/0)
	ease := readability.FleschKincaidReadingEase(0, 1, 1)
	require.Equal(t, 0.0, ease)

	// Sentence count cannot be zero (div/0)
	ease = readability.FleschKincaidReadingEase(1, 0, 1)
	require.Equal(t, 0.0, ease)
}
