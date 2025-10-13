package vectorize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.rtnl.ai/nlp/vector"
)

// ############################################################################
// VoyageEmbedder
// ############################################################################

// VoyageAIEmbedder can be used to vectorize text using the VoyageAI embeddings
// API (https://docs.voyageai.com/reference/embeddings-api).
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

// Create a new VoyageAI embedding vectorizer. Providing an API key is required.
//
// Defaults:
//   - Endpoint: "https://api.voyageai.com/v1/embeddings"	(Set with [VoyageAIEmbedderWithModel])
//   - Model: 	 "voyage-3.5-lite" 							(Set with [VoyageAIEmbedderWithEndpoint])
func NewVoyageAIEmbedder(apiKey string, opts ...VoyageAIEmbedderOption) (vectorizer *VoyageAIEmbedder, err error) {
	// Initialize with defaults
	vectorizer = &VoyageAIEmbedder{
		apiKey:   apiKey,
		endpoint: "https://api.voyageai.com/v1/embeddings",
		model:    "voyage-3.5-lite",
		client:   &http.Client{},
	}

	// Set options
	for _, fn := range opts {
		fn(vectorizer)
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

	// Get JSON body
	if jsonData, err = json.Marshal(data); err != nil {
		return nil, err
	}

	// Create HTTP request
	if request, err = http.NewRequest("POST", v.endpoint, bytes.NewBuffer(jsonData)); err != nil {
		return nil, err
	}

	// Set headers
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", v.apiKey))
	request.Header.Set("Content-Type", "application/json")

	// Make request
	if response, err = v.client.Do(request); err != nil {
		return nil, err
	}

	// Unmarshal data
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
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
