package similarity_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/similarity"
)

func TestNewCosineSimilarizer(t *testing.T) {
	t.Run("SuccessDefaults", func(t *testing.T) {
		sim, err := similarity.NewCosineSimilarizer([]string{"this", "is", "a", "test"})
		require.NoError(t, err)
		require.NotNil(t, sim)
	})
}

func TestCosineSimilarity(t *testing.T) {
	// TODO test CosineSimilarizer.Similarity() with several different text chunks
}
