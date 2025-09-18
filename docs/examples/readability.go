package main

/*
##################################################
## USAGE: `go run docs/examples/readability.go` ##
##################################################
*/

import "fmt"

// An example for how to perform readability scoring comparisons.
func ComparisonExample() {
	fmt.Println("TODO: ComparisonExample")
	// TODO in sc-34725
	//
	// load 2 short passages (A and B)
	// prints FRE and FKGL for each
	// calls compare(A,B)  to show the winner + the delta
}

// An example for how to perform readability scoring to get statistics on an
// array of text "chunks".
func StatsExample() {
	fmt.Println("TODO: StatsExample")
	// TODO in sc-34725
	//
	// load a small array of "context chunks"
	// prints FRE and FKGL for each
	// calls batch(contextChunks, method="fre", agg="micro") to return mean
	// calls batch(contextChunks, method="fre", agg="macro")  to return mean (edited)
}

func main() {
	ComparisonExample()
}
