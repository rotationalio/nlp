package tokenize

// ############################################################################
// Tokenizer interface
// ############################################################################

type Tokenizer interface {
	Tokenize(chunk string) (tokens []string, err error)
}
