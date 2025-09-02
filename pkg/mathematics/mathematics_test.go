package mathematics_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/mathematics"
)

func TestBoundToRange(t *testing.T) {
	// x, upper, lower, expected
	testcases := [][4]float64{
		{
			0.0, -1.0, 1.0, 0.0,
		},
		{
			2.0, -1.0, 1.0, 1.0,
		},
		{
			-2.0, -1.0, 1.0, -1.0,
		},
		{
			-1.0, -1.0, 1.0, -1.0,
		},
		{
			1.0, -1.0, 1.0, 1.0,
		},
		{
			math.MaxFloat64, -1.0, 1.0, 1.0,
		},
		{
			-math.MaxFloat64, -1.0, 1.0, -1.0,
		},
	}

	for _, tc := range testcases {
		actual := mathematics.BoundToRange(tc[0], tc[1], tc[2])
		require.Equal(t, tc[3], actual)
	}
}
