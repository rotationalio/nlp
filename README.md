# nlp

Natural Language Processing in Golang.

GitHub: <https://github.com/rotationalio/nlp>
Go Docs: <https://go.rtnl.ai/nlp>

## The Vision

`nlp` will house high-performance natural language processing and quantitative metrics, especially ones that can be computed over text.
Think statistical or structural properties of a string or document.
The first use case is text analysis inside [Endeavor](https://github.com/rotationalio/endeavor).

`nlp` will enable NLP for AI engineers: numeric metrics that help us reason about what a model did, what the humans expected, and how to compare the two.

We want this package to be:

* Performant (written in Go, reused across tools)
* Composable (individual functions that do one thing well)
* Extensible (easy to add new metrics as we learn more)

## Design Goals

* Each metric should be self-contained and independently callable
* Avoid hard dependencies on LLMs or external services
* Define a common interface (e.g., all metrics take `string`, return `float64` or `map[string]float64`)
* Organize by category (similarity, counts, readability, lexical, etc.)
* Stub out room for future metrics, even weird ones

## End-User API Usage

There are two ways you can use this library:

1) Use the unified `text.Text` interface (see example below) to perform all of the possible NLP operations using a single object that is configured with the specific tools you wish to use when it is created via `text.New(chunk string) *text.Text`.
This also includes using the `token.Token` and `tokenlist.TokenList` types which have their own useful features.
2) Use the various tools in the lower level packages such as the `stem` or the `tokenize` packages on an as-needed basis.
These tools generally use basic Go types such as strings, ints, floats, and slices of the same.

### text.Text API Example

```Go
// Create a [Text] with the default settings
myText, err := text.New("apple aardvarks zebra bananna aardvark")

// Get all of the word tokens
myTokens, err := myText.Tokens() // TokenList

// Get all word stem tokens which use the same underlying types as the full
// word tokens above (ignoring errors in this example)
myStems, err := myText.Stems() // TokenList

// The stems are 1:1 count with the tokens
if len(myTokens) != len(myStems) { // 5 == 5
  panic("this should never occur")
}

// You can also get a type count, which returns the count of each unique
// word stem (ignoring errors) ("aardvark" has a 2 count for this example)
myCount, err := myText.TypeCount() // map[string]int

// These are a [tokenlist.TokenList], but if you need a slice of strings...
stringTokens := myTokens.Strings() // []string

// You can also use regular slice functions and operations on a [tokenlist.TokenList]
length := len(myTokens) // 5
myTokens = append(myTokens, myTokens[0]) // "apple", "aardvarks", "zebra", "bananna", "aardvark", "apple"
myTokens[0] = myTokens[1] // "aardvarks", "aardvarks", "zebra", "bananna", "aardvark", "apple"

// Get an individual token
firstToken := myTokens[0] // Token

// You can also get a token as another type
stringToken := firstToken.String() // string
runeToken := firstToken.Runes() // []rune
byteToken := firstToken.Bytes() // []byte

// For these examples, we need to re-create the [Text] with a vocabulary,
// so the [vectorize.CountVectorizer] will work without an error to get
// cosine similarity. You could also use a different vectorization method.
myText, err = text.New(
  "cars have engines like motorcycles have engines",
  text.WithVocabulary([]string{"car", "engine", "brakes", "transmission"}),
)
otherText, err := text.New(
  "engines are attached to transmissions",
  text.WithVocabulary([]string{"car", "engine", "brakes", "transmission"}),
)

// Cosine similarity with another string
similarity, err := myText.CosineSimilarity(otherText) // ~0.5

// We can also get a one-hot or frequency vectorization of our text
myOneHotVector, err := myText.VectorizeOneHot() // vector.Vector{1, 1, 0, 0}
myFrequencyVector, err := myText.VectorizeFrequency() // vector.Vector{1, 2, 0, 0}
```

See the [NLP Go docs](https://go.rtnl.ai/nlp) for this library for more details.

## Features, metrics, and tools

* Tokenization and type counting
  * Regex tokenization with custom expression support
* Stemming
  * Porter2/Snowball stemming algorithm
* Similarity metrics
  * Cosine similarity
* Vectors & vectorization
  * One-hot encoding
  * Frequency (count) encoding
* Descriptive statistics (minimum, maximum, mean, stddev, variance, etc.)
  * See the stats package [README.md](./stats/README.md) for more information

### Planned

* Readability Scores (ASAP)
* Part-of-Speech Distributions (Future)
* Named Entities & Keyphrase Counts (Future)
* Custom Classifiers (Distant Future)

## Developing in nlp

Different feature categories are separated into different packages, for example we might have similarity metrics in `similarity/` and text classifiers in `classifiers/`.
If you want to add a new feature, please ensure it is placed in a package which fits the category, or create a new package if none yet exist.
Tests should be located next to each feature, for example `similarity_tests.go` would hold the tests for `similarity.go`.
Test data should go into the `testdata/` folder within the package where the test is located.
Documentation should go into each function's and package's docstrings so the documentation is accessible to the user while using the library in their local IDE and also available using Go's documentation tools.
Documentation can also be included in separate Markdown files as-needed in the `docs/` folder or in this README, such as for the `text.Text` API examples.
Any documentation or research that isn't immediately relevant to the user in the code context should go into the `docs/` folder in the root.

## Sources and References

To ensure the algorithms in this package are accurate, we pulled information from several references, which have been recorded in [`docs/sources.md`](./docs/sources.md) and in the documentation and comments for the individual functions in this library.

## Research Notes

Research on different topics will go into the folder [`docs/research/`](./docs/research/).

* [Go NLP](./docs/research/go_nlp.md): notes on different NLP packages/libraries for Go

## License

See: [LICENSE](./LICENSE)
