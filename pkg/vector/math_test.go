package vector_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/errors"
	"go.rtnl.ai/nlp/pkg/vector"
)

func TestCosine(t *testing.T) {
	testcases := []struct {
		Name     string
		First    vector.Vector
		Second   vector.Vector
		Expected float64
		Error    error
	}{
		{
			Name:   "Success_SameVectors",
			First:  vector.Vector{1, 2, 3},
			Second: vector.Vector{1, 2, 3},
			// NOTE: this test will fail if the return value from [vector.Cosine]
			// is not clamped to [-1.0, 1.0]
			Expected: 1.0,
			Error:    nil,
		},
		{
			Name:     "Success_PerpendicularVectors",
			First:    vector.Vector{1, 0, 1},
			Second:   vector.Vector{0, 1, 0},
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name:   "Success_OppositeVectors",
			First:  vector.Vector{1, 1, 1},
			Second: vector.Vector{-1, -1, -1},
			// NOTE: this test will fail if the return value from [vector.Cosine]
			// is not clamped to [-1.0, 1.0]
			Expected: -1.0,
			Error:    nil,
		},
		{
			Name:     "Success_7056LengthVectors",
			First:    repeatingVector(7056/2, 1, 2),
			Second:   repeatingVector(7056/2, 2, 1),
			Expected: 0.8, // calculated in Python 3.13.4
			Error:    nil,
		},
		{
			Name:     "Error_ZeroVectors",
			First:    vector.Vector{0, 0, 0},
			Second:   vector.Vector{0, 0, 0},
			Expected: 0.0,
			Error:    errors.ErrUndefinedValue,
		},
		{
			Name:     "Error_OneZeroVector",
			First:    vector.Vector{0, 0, 0},
			Second:   vector.Vector{1, 2, 3},
			Expected: 0.0,
			Error:    errors.ErrUndefinedValue,
		},
		{
			Name:     "Error_ZeroLengthVectors",
			First:    vector.Vector{},
			Second:   vector.Vector{},
			Expected: 0.0,
			Error:    errors.ErrUndefinedValue,
		},
		{
			Name:     "Error_UnequalLengthVectors",
			First:    vector.Vector{1, 2, 3},
			Second:   vector.Vector{1, 2, 3, 4},
			Expected: 0.0,
			Error:    errors.ErrUnequalLengthVectors,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			val, err := vector.Cosine(tc.First, tc.Second)
			if tc.Error != nil {
				require.Error(t, tc.Error, err)
			} else {
				require.NoError(t, err)
			}
			require.InDeltaf(t, tc.Expected, val, 1e-12, "expected %f got %f a difference of %e", tc.Expected, val, math.Abs(tc.Expected-val))
		})
	}
}

func TestDotProduct(t *testing.T) {
	testcases := []struct {
		Name     string
		First    vector.Vector
		Second   vector.Vector
		Expected float64
		Error    error
	}{
		{
			Name:     "Success_SameVectors",
			First:    vector.Vector{1, 2, 3},
			Second:   vector.Vector{1, 2, 3},
			Expected: 14.0, // calculated in Python 3.13.4
			Error:    nil,
		},
		{
			Name:     "Success_PerpendicularVectors",
			First:    vector.Vector{1, 0, 1},
			Second:   vector.Vector{0, 1, 0},
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name:     "Success_OppositeVectors",
			First:    vector.Vector{1, 1, 1},
			Second:   vector.Vector{-1, -1, -1},
			Expected: -3.0, // calculated in Python 3.13.4
			Error:    nil,
		},
		{
			Name:     "Success_7056LengthVectors",
			First:    repeatingVector(7056/2, 1, 2),
			Second:   repeatingVector(7056/2, 2, 1),
			Expected: 14112.0, // calculated in Python 3.13.4
			Error:    nil,
		},
		{
			Name:     "Success_ZeroLengthVectors",
			First:    vector.Vector{},
			Second:   vector.Vector{},
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name:     "Success_ZeroVectors",
			First:    vector.Vector{0, 0, 0},
			Second:   vector.Vector{0, 0, 0},
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name:     "Sucess_OneZeroVector",
			First:    vector.Vector{0, 0, 0},
			Second:   vector.Vector{1, 2, 3},
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name:     "Error_UnequalLengthVectors",
			First:    vector.Vector{1, 2, 3},
			Second:   vector.Vector{1, 2, 3, 4},
			Expected: 0.0,
			Error:    errors.ErrUnequalLengthVectors,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			val, err := vector.DotProduct(tc.First, tc.Second)
			if tc.Error != nil {
				require.Error(t, tc.Error, err)
			} else {
				require.NoError(t, err)
			}
			require.InDeltaf(t, tc.Expected, val, 1e-12, "expected %f got %f a difference of %e", tc.Expected, val, math.Abs(tc.Expected-val))
		})
	}
}

func TestVectorLength(t *testing.T) {
	testcases := []struct {
		Name     string
		Vector   vector.Vector
		Expected float64
		Error    error
	}{
		{
			Name:     "OneTwoThreeVector",
			Vector:   vector.Vector{1, 2, 3},
			Expected: 3.7416573867739413, // calculated in Python 3.13.4
			Error:    nil,
		},
		{
			Name:     "OnesVector",
			Vector:   vector.Vector{1, 1, 1},
			Expected: 1.7320508075688772, // calculated in Python 3.13.4
			Error:    nil,
		},
		{
			Name:     "ZeroLengthVector",
			Vector:   vector.Vector{},
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name:     "ZeroVector",
			Vector:   vector.Vector{0, 0, 0},
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name: "ZeroTo9999_Vector",
			// Vector has 10_000 elements going from 0 to 9999
			Vector:   rangeVector(10_000),
			Expected: 577306.967739001, // calculated in Python 3.13.4
			Error:    nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			val := vector.Magnitude(tc.Vector)
			require.InDeltaf(t, tc.Expected, val, 1e-12, "expected %f got %f a difference of %e", tc.Expected, val, math.Abs(tc.Expected-val))
		})
	}
}

func TestLen(t *testing.T) {
	// NOTE: no actual testing is required while [vector.Vector] is implemented
	// as a wrapper over a `[]float64`; when this type changes then this test
	// will fail and alert the developer to write a test for [vector.Len] here.
	vec := vector.Vector{1, 2, 3}
	kind := reflect.TypeOf(vec).Kind()
	require.Equal(t, reflect.Slice, kind)
	kind2 := reflect.TypeOf(vec[0]).Kind()
	require.Equal(t, reflect.Float64, kind2)
}

// ############################################################################
// Helpers
// ############################################################################

// Returns a [vector.Vector] that goes from 0 to `upper`.
func rangeVector(upper int) (vec vector.Vector) {
	vec = vector.Vector{}
	for i := range upper {
		vec = append(vec, float64(i))
	}
	return vec
}

// Returns a [vector.Vector] that repeats the `values` list `n` times.
func repeatingVector(n int, values ...float64) (vec vector.Vector) {
	vec = vector.Vector{}
	for _ = range n {
		vec = append(vec, values...)
	}
	return vec
}
