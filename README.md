# nlp

Natural Language Processing in Golang.

GitHub: <https://github.com/rotationalio/nlp>

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

TODO (sc-34050): Rewrite this section when the "Text/Blob" refactoring is complete in sc-34050

Generally, if we have an `Operation` that we want to perform on text, there will be an associated `Operation[izer|er]` type.
The `Operation[izer|er]` may be an `interface` in the cases where we might want more than one implementation of the `Operation`, such as with the `Stemmer` interface where we have a `Porter2Stemmer` implementation.
The `Operation[izer|er]` may alternately be a `struct` in the cases where we only need one implementation of the `Operation`, such as with the `TypeCounter` struct which does not need additional implementations.
Each of the `Operation[izer|er]` types will have a `NewOperation[izer|er](opts ...Operation[izer|er]Option)` which allows the user to either take the default configuration of the new instance, or they can include a variable number of arguments which can modify the options for the `Operation[izer|er]` instance returned.
Please see the documentation comments within the code for more information on how to use this library, and if anything is unclear please contact us to clarify!

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

### Planned

* Readability Scores (ASAP)
* Part-of-Speech Distributions (Future)
* Named Entities & Keyphrase Counts (Future)
* Custom Classifiers (Distant Future)

## Developing in nlp

Different feature categories are separated into different packages, for example we might have similarity metrics in `pkg/similarity` and text classifiers in `pkg/classifiers`.
If you want to add a new feature, please ensure it is placed in a package which fits the category, or create a new package if none yet exist.
Tests should be located next to each feature, for example `similarity_tests.go` would hold the tests for `similarity.go`.
Test data should go into the `testdata/` folder within the package where the test is located.
Documentation should go into each function's and package's docstrings so the documentation is accessible to the user while using the library in their local IDE and also available using Go's documentation tools.
Any documentation or research that isn't immediately relevant to the user in the code context should go into the `docs/` folder in the root.

## Sources and References

To ensure the algorithms in this package are accurate, we pulled information from several references, which have been recorded in [`docs/sources.md`](./docs/sources.md) and in the documentation and comments for the individual functions in this library.

## Research Notes

Research on different topics will go into the folder [`docs/research/`](./docs/research/).

* [Go NLP](./docs/research/go_nlp.md): notes on different NLP packages/libraries for Go

## License

See: [LICENSE](./LICENSE)
