package vectorize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.rtnl.ai/nlp/errors"
	"go.rtnl.ai/nlp/vector"
)

// ############################################################################
// VoyageEmbedder
// ############################################################################

/*
VoyageAIEmbedder can be used to vectorize text using the VoyageAI embeddings
API (https://docs.voyageai.com/reference/embeddings-api).

Usage Example:

	// Panics when the error is not nil
	func checkErr(err error) {
		if err != nil {
			panic(err)
		}
	}

	// We can also use the VoyageAI API to get embedding vectors; see the docs for
	// [NewVoyageAIEmbedder] for information on how to load it's configs
	// via environment variables or you can load them using the options functions
	// as shown below
	voyage, err := vectorize.NewVoyageAIEmbedder(
		vectorize.VoyageAIEmbedderWithAPIKey("your_voyageai_api_key_here"),
		vectorize.VoyageAIEmbedderWithEndpoint("https://api.voyageai.com/v1/embeddings"),
		vectorize.VoyageAIEmbedderWithModel("voyage-3.5-lite"),
		)
	checkErr(err)

	// Get a single embedding vector using the [Vectorizer.Vectorize] interface
	chunk1 := "A simple test."
	embedding, err := voyage.Vectorize(chunk1)
	checkErr(err)

	// Get several embeddings in a single VoyageAI API call
	chunks := []string{chunk1, "A slightly more complex test, but only slightly.", "Number three!"}
	embeddings, err := voyage.VectorizeAll(chunks)
	checkErr(err)

	// Print the number of VoyageAI usage tokens that were used
	fmt.Printf("used %d tokens\n", voyage.TotalTokensUsed())
*/
type VoyageAIEmbedder struct {
	apiKey          string
	endpoint        string
	model           string
	client          *http.Client
	totalTokensUsed int
}

// Ensure [VoyageAIEmbedder] meets the [Vectorizer] interface requirements.
var _ Vectorizer = &VoyageAIEmbedder{}

// ############################################################################
// VoyageAI Constructor and Options
// ############################################################################

// Create a new VoyageAI embedding vectorizer. Options will be loaded from
// "VOYAGEAI_*" environment variables if available (see `.env.template`). The
// environment values will be overridden by [VoyageAIEmbedderOption]s.If a
// necessary option is not provided or found in the environment, then
// [errors.ErrMissingConfig] will be returned alongside another more descriptive
// error, so ensure you use [errors.Is] to disambiguate the errors.
func NewVoyageAIEmbedder(opts ...VoyageAIEmbedderOption) (vectorizer *VoyageAIEmbedder, err error) {
	// Initialize with options from the environment; the user must load the env
	// vars somewhere else themselves.
	vectorizer = &VoyageAIEmbedder{
		apiKey:   os.Getenv("VOYAGEAI_API_KEY"),
		endpoint: os.Getenv("VOYAGEAI_EMBEDDING_ENDPOINT"),
		model:    os.Getenv("VOYAGEAI_EMBEDDING_MODEL"),
		client:   &http.Client{},
	}

	// Set user-provided options, overridding any environment variables
	for _, fn := range opts {
		fn(vectorizer)
	}

	// Ensure all required configs are set before returning the vectorizer
	if vectorizer.apiKey == "" {
		return nil, errors.Join(
			errors.ErrMissingConfig,
			errors.New("field 'apiKey' is required; use option 'VoyageAIEmbedderWithAPIKey()' or set the environment variable 'VOYAGEAI_API_KEY'."),
		)
	}
	if vectorizer.model == "" {
		return nil, errors.Join(
			errors.ErrMissingConfig,
			errors.New("field 'model' is required; use option 'VoyageAIEmbedderWithModel()' or set the environment variable 'VOYAGEAI_EMBEDDING_ENDPOINT'."),
		)
	}
	if vectorizer.endpoint == "" {
		return nil, errors.Join(
			errors.ErrMissingConfig,
			errors.New("field 'endpoint' is required; use option 'VoyageAIEmbedderWithEndpoint()' or set the environment variable 'VOYAGEAI_EMBEDDING_MODEL'."),
		)
	}

	return vectorizer, nil
}

// Returns the total number of VoyageAI tokens used by this client.
func (v *VoyageAIEmbedder) TotalTokensUsed() int {
	return v.totalTokensUsed
}

// ############################################################################
// VoyageAI Embedder Options
// ############################################################################

// VoyageAIEmbedderOption functions modify a [VoyageAIEmbedder].
type VoyageAIEmbedderOption func(c *VoyageAIEmbedder)

// VoyageAIEmbedderWithAPIKey sets the API key string to use with the
// [VoyageAIEmbedder].
func VoyageAIEmbedderWithAPIKey(apiKey string) VoyageAIEmbedderOption {
	return func(c *VoyageAIEmbedder) {
		c.apiKey = apiKey
	}
}

// VoyageAIEmbedderWithModel sets the model to use with the [VoyageAIEmbedder].
func VoyageAIEmbedderWithModel(model string) VoyageAIEmbedderOption {
	return func(c *VoyageAIEmbedder) {
		c.model = model
	}
}

// VoyageAIEmbedderWithEndpoint sets the endpoint URL to use with the
// [VoyageAIEmbedder].
func VoyageAIEmbedderWithEndpoint(endpoint string) VoyageAIEmbedderOption {
	return func(c *VoyageAIEmbedder) {
		c.endpoint = endpoint
	}
}

// ############################################################################
// VoyageAI Text Embedding Functionality
// ############################################################################

// Performs embedding vectorization on a chunk of text using the VoyageAI API.
func (v *VoyageAIEmbedder) Vectorize(chunk string) (embedding vector.Vector, err error) {
	if embeddings, err := v.VectorizeAll([]string{chunk}); err != nil {
		return nil, err
	} else {
		return embeddings[0], nil
	}
}

// Performs embedding vectorization on several chunks of text in a single
// request using the VoyageAI API.
func (v *VoyageAIEmbedder) VectorizeAll(chunks []string) (embeddings []vector.Vector, err error) {
	data := &VoyageAIEmbeddingsRequest{
		Model: v.model,
		Input: chunks,
	}
	return v.makeTextEmbeddingsRequest(data)
}

// Makes a request to the VoyageAI
func (v *VoyageAIEmbedder) makeTextEmbeddingsRequest(data *VoyageAIEmbeddingsRequest) (embeddings []vector.Vector, err error) {
	var (
		jsonData []byte
		request  *http.Request
		response *http.Response
		body     *VoyageAIEmbeddingsResponse
	)

	// Marshal JSON body
	if jsonData, err = json.Marshal(data); err != nil {
		return nil, errors.Join(err, errors.New("error encoding VoyageAIEmbeddingsRequest"))
	}

	// Create HTTP request
	if request, err = http.NewRequest("POST", v.endpoint, bytes.NewBuffer(jsonData)); err != nil {
		return nil, errors.Join(err, errors.New("error creating embeddings http request"))
	}

	// Set headers
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", v.apiKey))
	request.Header.Set("Content-Type", "application/json")

	// Make request
	if response, err = v.client.Do(request); err != nil {
		return nil, errors.Join(err, errors.New("embeddings request failed"))
	}

	// Unmarshal data
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, errors.Join(err, errors.New("error decoding VoyageAIEmbeddingsResponse"))
	}

	// Count the tokens used so we have a running total for this client
	v.totalTokensUsed += body.Usage.TotalTokens

	// Gather and return the embeddings
	for _, embedding := range body.Data {
		embeddings = append(embeddings, embedding.Embedding)
	}

	return embeddings, nil
}

// ############################################################################
// VoyageAI Request Struct
// ############################################################################

// Model for VoyageAI request (non-exhaustive).
// See: https://docs.voyageai.com/reference/embeddings-api
type VoyageAIEmbeddingsRequest struct {
	// Required
	Model string   `json:"model"`
	Input []string `json:"input"`

	// Optional
	// TODO: if these are needed create a [VoyageAIEmbedderOption] for each
	Truncation      *bool   `json:"truncation,omitempty"`
	OutputDimension *int    `json:"output_dimension,omitempty"`
	OutputDType     *string `json:"output_dtype,omitempty"`
	EncodingFormat  *string `json:"encoding_format"`
}

// ############################################################################
// VoyageAI Response Structs
// ############################################################################

// Model for VoyageAI response.
// See: https://docs.voyageai.com/reference/embeddings-api
type VoyageAIEmbeddingsResponse struct {
	Object string               `json:"object"` // always "list"
	Data   []*VoyageAIEmbedding `json:"data"`
	Model  string               `json:"model"`
	Usage  VoyageAIUsage        `json:"usage"`
}

// Model for VoyageAI embeddings inside the [VoyageAIEmbeddingsResponse].
type VoyageAIEmbedding struct {
	Object    string    `json:"object"` // always "embedding"
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// Model for VoyageAI usage data inside the [VoyageAIEmbeddingsResponse].
type VoyageAIUsage struct {
	TotalTokens int `json:"total_tokens"`
}
