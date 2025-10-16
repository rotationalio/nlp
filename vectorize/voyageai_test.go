package vectorize_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/vectorize"
)

// Creates a [vectorize.VoyageAIEmbedder] and tests it's client functionality
// live. This requires the JSON file `testdata/voyageai_credentials.json` to be
// present and have a functional API key. If the file is missing, this test will
// be skipped instead of failing.
//
// NOTE: This test will use ~21 tokens of usage on VoyageAI. Use "go test -v" to
// print a log with the total tokens used.
func TestVoyageAIEmbedder(t *testing.T) {
	// Skip the test if '-short' is provided to 'go test'
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}

	// Load the test env variables, ignoring errors
	_ = godotenv.Load(filepath.Join(".env"))

	// Skip the test if any of the required configs are unset in the environment
	apiKey, ok := os.LookupEnv("VOYAGEAI_API_KEY")
	if !ok {
		t.Skip("could not load apiKey from environment")
	}

	endpoint, ok := os.LookupEnv("VOYAGEAI_EMBEDDING_ENDPOINT")
	if !ok {
		t.Skip("could not load endpoint from environment")
	}

	model, ok := os.LookupEnv("VOYAGEAI_EMBEDDING_MODEL")
	if !ok {
		t.Skip("could not load model from environment")
	}

	// Create a new embedder
	voyage, err := vectorize.NewVoyageAIEmbedder(
		vectorize.VoyageAIEmbedderWithAPIKey(apiKey),
		vectorize.VoyageAIEmbedderWithEndpoint(endpoint),
		vectorize.VoyageAIEmbedderWithModel(model),
	)
	require.NoError(t, err, "error creating VoyageAIEmbedder")

	// Run a test with a single embedding
	chunk1 := "A simple test."
	embedding, err := voyage.Vectorize(chunk1)
	require.NoError(t, err, "error making embedding request")
	require.Len(t, embedding, 1024, "embedding vector length should default to 1024")
	require.IsType(t, float64(0.0), embedding[0], "embedding vector element type should be float64")

	// Run a test with a few embeddings
	chunks := []string{chunk1, "A slightly more complex test, but only slightly.", "Number three!"}
	embeddings, err := voyage.VectorizeAll(chunks)
	require.NoError(t, err, "error making embedding request")
	require.Lenf(t, embeddings, len(chunks), "there should be %d embeddings", len(chunks))
	require.IsType(t, float64(0.0), embeddings[0][0], "embedding vector element type should be float64")
	require.InDeltaSlice(t, embedding, embeddings[0], 1e-2, "the elements of these two embeddings should be within a small delta")

	// How many tokens were used for the tests?
	t.Logf("used %d tokens for TestVoyageAIEmbedder", voyage.TotalTokensUsed())
}
