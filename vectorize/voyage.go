package vectorize

import "go.rtnl.ai/nlp/vector"

// ############################################################################
// VoyageEmbedder
// ############################################################################

// VoyageEmbedder can be used to vectorize text using the frequency or one-hot
// text vectorization algorithms.
type VoyageEmbedder struct {
	//TODO
}

// Ensure [VoyageEmbedder] meets the [Vectorizer] interface requirements.
var _ Vectorizer = &VoyageEmbedder{}

// Performs embedding vectorization on a chunk of text using the VoyageAI API.
func (v *VoyageEmbedder) Vectorize(chunk string) (embedding vector.Vector, err error) {
	return nil, nil //TODO
}
