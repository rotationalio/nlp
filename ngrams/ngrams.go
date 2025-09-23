package ngrams

// Returns the n-grams for the given sequence. If the sequence length is less
// than n items then the returned value will be nil.
func Ngrams[T any](sequence []T, n int) (ngrams [][]T) {
	if len(sequence) < n {
		return nil
	}

	ngrams = make([][]T, 0, len(sequence)-(n-1))
	for i := range sequence {
		if len(sequence) < i+n {
			break
		}
		ngrams = append(ngrams, sequence[i:i+n])
	}

	return ngrams
}

// Returns the trigrams for the given sequence. If the sequence length is less
// than 3 items then the returned value will be nil.
func Trigrams[T any](sequence []T) (trigrams [][]T) {
	return Ngrams(sequence, 3)
}
