package vectorize_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/vectorize"
)

// Creates a [vectorize.VoyageAIEmbedder] and tests it's client functionality
// live. This requires the JSON file `testdata/voyageai_credentials.json` to be
// present and have a functional API key. If the file is missing, this test will
// be skipped instead of failing.
//
// NOTE: This test will use 25 tokens worth of usage on VoyageAI.
func TestVoyageAIEmbedder(t *testing.T) {
	// Track token usage between different clients
	tokensUsedForTests := 0

	// Load the API key credential from the JSON file
	credJson, err := os.ReadFile("testdata/voyageai_credentials.json")
	if errors.Is(err, os.ErrNotExist) {
		// If the file does not exist, then skip this test instead of failing.
		t.SkipNow()
	}
	require.NoError(t, err, "error reading file")

	var creds map[string]any
	err = json.Unmarshal(credJson, &creds)
	require.NoError(t, err, "error parsing json")

	apiKey, ok := creds["api_key"].(string)
	require.True(t, ok, "could not find key 'api_key'")

	// Create a new embedder
	vai, err := vectorize.NewVoyageAIEmbedder(apiKey)
	require.NoError(t, err, "error creating VoyageAIEmbedder")

	// Run a test with a single embedding
	chunk1 := "A simple test."
	embedding, err := vai.Vectorize(chunk1)
	require.NoError(t, err, "error making embedding request")
	require.Len(t, embedding, 1024, "embedding vector length should default to 1024")
	require.IsType(t, float64(0.0), embedding[0], "embedding vector element type should be float64")

	// Run a test with a few embeddings
	chunks := []string{chunk1, "A slightly more complex test, but only slightly.", "Number three!"}
	embeddings, err := vai.VectorizeAll(chunks)
	require.NoError(t, err, "error making embedding request")
	require.Lenf(t, embeddings, len(chunks), "there should be %d embeddings", len(chunks))
	require.IsType(t, float64(0.0), embeddings[0][0], "embedding vector element type should be float64")
	require.InDeltaSlice(t, embedding, embeddings[0], 1e-2, "the elements of these two embeddings should be within a small delta")
	tokensUsedForTests += vai.TotalTokensUsed()

	// Create a new embedder with a different model
	vai, err = vectorize.NewVoyageAIEmbedder(apiKey, vectorize.VoyageAIEmbedderWithModel("voyage-3.5"))
	require.NoError(t, err, "error creating VoyageAIEmbedder with different model")

	// Run a test with a single embedding
	embedding2, err := vai.Vectorize(chunk1)
	require.NoError(t, err, "error making embedding request")
	require.Len(t, embedding2, 1024, "embedding vector length should default to 1024")
	require.IsType(t, float64(0.0), embedding2[0], "embedding vector element type should be float64")
	tokensUsedForTests += vai.TotalTokensUsed()

	// How many tokens were used for the tests?
	require.Equal(t, 25, tokensUsedForTests, "expected to use 25 total tokens for all of the tests")
}
