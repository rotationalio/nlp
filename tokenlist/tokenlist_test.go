package tokenlist_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/token"
	"go.rtnl.ai/nlp/pkg/tokenlist"
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
		require.Equal(t, tc.Size, len(tl.Tokens()))
		require.Equal(t, tc.ExpCap, cap(tl.Tokens()))
		// require all tokens to be the empty string
		for _, tok := range tl.Tokens() {
			require.Equal(t, token.New(""), tok)
		}
	}
}

func TestAppend(t *testing.T) {
	tl := tokenlist.NewEmpty(0, 0)
	require.Len(t, tl.Tokens(), 0)

	expected := token.New("one")
	tl.Append(expected)
	require.Len(t, tl.Tokens(), 1)

	token := tl.Tokens()[0]
	require.Equal(t, expected, token)

	tl.Append(expected)
	require.Len(t, tl.Tokens(), 2)

	token = tl.Tokens()[1]
	require.Equal(t, expected, token)
}

func TestReplace(t *testing.T) {
	tl := tokenlist.New([]string{"one", "two", "three"})
	require.Equal(t, 3, tl.Len())

	err := tl.Replace(1, token.New("replacement"))
	require.NoError(t, err)
	require.Equal(t, []string{"one", "replacement", "three"}, tl.Strings())
}

func TestLen(t *testing.T) {
	tl := tokenlist.NewEmpty(0, 0)
	require.Len(t, tl.Tokens(), 0)
	require.Equal(t, 0, tl.Len())

	expected := token.New("one")
	tl.Append(expected)
	require.Len(t, tl.Tokens(), 1)
	require.Equal(t, 1, tl.Len())
}
