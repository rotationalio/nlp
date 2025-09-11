package tokenlist_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/token"
	"go.rtnl.ai/nlp/tokenlist"
)

func TestNew(t *testing.T) {
	tokens := []string{"one", "two", "three"}
	tl := tokenlist.New(tokens)
	require.Equal(t, tl.Strings(), tokens)
}

func TestNewCopy(t *testing.T) {
	tokens := []string{"one", "two", "three"}
	tl1 := tokenlist.New(tokens)
	require.Equal(t, tl1.Strings(), tokens)
	tl2 := tokenlist.NewCopy(tl1)
	require.Equal(t, tl2.Strings(), tokens)
}

func TestNewEmpty(t *testing.T) {
	testcases := []struct {
		Size   int
		Cap    int
		ExpCap int // expected capacity can be different from the input
	}{
		{10, 10, 10},
		{0, 10, 10},
		{10, 15, 15},
		{10, 5, 10},
		{10, 0, 10},
	}
	for _, tc := range testcases {
		tl := tokenlist.NewEmpty(tc.Size, tc.Cap)
		require.Equal(t, tc.Size, len(tl))
		require.Equal(t, tc.ExpCap, cap(tl))
		// require all tokens to be the empty string
		for _, tok := range tl {
			require.Equal(t, token.New(""), tok)
		}
	}
}

// Tests that the docs for [tokenlist.New] work properly
func TestNewDocs(t *testing.T) {
	// Create a new [TokenList] from a slice of strings
	items := []string{"apple", "bananna", "zebra"}
	myTokens := tokenlist.New(items)

	// You can get a slice of strings back
	stringTokens := myTokens.Strings() // []string
	require.Equal(t, items, stringTokens)

	// You can get a new [TokenList] as a copy of another list
	aCopy := tokenlist.NewCopy(myTokens) // TokenList (copy of myTokens)
	require.Equal(t, myTokens, aCopy)

	// You can get an empty [TokenList] with a specific size and capacity
	emptyTokens := tokenlist.NewEmpty(0, 100) // TokenList (with no entries)
	require.Len(t, emptyTokens, 0)
	require.Equal(t, 100, cap(emptyTokens))
	emptyTokens = tokenlist.NewEmpty(10, 100) // TokenList (with 10 null string tokens)
	require.Len(t, emptyTokens, 10)
	require.Equal(t, 100, cap(emptyTokens))
	emptyTokens = tokenlist.NewEmpty(10, 5) // TokenList (capacity will be set to 10)
	require.Len(t, emptyTokens, 10)
	require.Equal(t, 10, cap(emptyTokens))

	// You can also use regular slice functions and operations on a [TokenList]
	length := len(myTokens) // 3
	require.Equal(t, 3, length)
	myTokens = append(myTokens, myTokens[0]) // "apple", "bananna", "zebra", "apple"
	require.Equal(t, []string{"apple", "bananna", "zebra", "apple"}, myTokens.Strings())
	myTokens[0] = myTokens[1] // "bananna", "bananna", "zebra"
	require.Equal(t, []string{"bananna", "bananna", "zebra", "apple"}, myTokens.Strings())
}
