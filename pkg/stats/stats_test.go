package stats

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadTestData() ([]float64, error) {
	data := make([]float64, 0, 1000000)

	f, err := os.Open("testdata/standardized.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewScanner(f)
	for buf.Scan() {
		var val float64
		if val, err = strconv.ParseFloat(buf.Text(), 64); err != nil {
			return nil, err
		}
		data = append(data, val)
	}

	if buf.Err() != nil {
		return nil, err
	}

	return data, nil
}

func ExampleStatistics() {
	stats := new(Statistics)
	samples, _ := loadTestData()

	for _, sample := range samples {
		stats.Update(sample)
	}

	data, _ := json.MarshalIndent(stats.Serialize(), "", "  ")
	fmt.Println(string(data))
	// Output:
	// {
	//   "maximum": 5.30507026071,
	//   "mean": 0.00041124313405184064,
	//   "minimum": -4.72206033824,
	//   "range": 10.02713059895,
	//   "samples": 1000000,
	//   "stddev": 0.9988808397330513,
	//   "total": 411.2431340518406,
	//   "variance": 0.9977629319858057
	// }
}

func TestStatistics(t *testing.T) {
	data, err := loadTestData()
	require.NoError(t, err)
	require.NotNil(t, data)

	stats := new(Statistics)

	for _, v := range data {
		stats.Update(v)
	}

	require.Equal(t, uint64(1000000), stats.N())
	require.Equal(t, 0.00041124313405184064, stats.Mean())
	require.Equal(t, 0.9988808397330513, stats.StdDev())
	require.Equal(t, 0.9977629319858057, stats.Variance())
	require.Equal(t, 5.30507026071, stats.Maximum())
	require.Equal(t, -4.7220603382400004, stats.Minimum())
	require.Equal(t, 10.02713059895, stats.Range())
}

func TestStatisticsBulk(t *testing.T) {
	data, err := loadTestData()
	require.NoError(t, err)
	require.NotNil(t, data)

	stats := new(Statistics)
	stats.Update(data...)

	require.Equal(t, uint64(1000000), stats.N())
	require.Equal(t, 0.00041124313405184064, stats.Mean())
	require.Equal(t, 0.9988808397330513, stats.StdDev())
	require.Equal(t, 0.9977629319858057, stats.Variance())
	require.Equal(t, 5.30507026071, stats.Maximum())
	require.Equal(t, -4.7220603382400004, stats.Minimum())
	require.Equal(t, 10.02713059895, stats.Range())
}

func TestStatisticsAppend(t *testing.T) {
	values := []float64{
		15.45832771, 11.11727874, 10.30855758, 14.63626755,
		5.85474266, 10.37473159, 11.02068524, 9.92171508,
		9.45442518, 11.84815447, 11.98722063, 11.54485569,
		8.49187437, 8.32798107, 9.85561918, 8.64735984,
		6.20092164, 7.33269192, 11.79721845, 7.57280214,
		7.32801938, 11.7176034, 10.27039045, 12.52726886,
		8.84401993, 6.79783127, 7.42687921, 7.53989174,
		9.29713199, 10.67506366, 6.63483678, 9.54300577,
		9.93653413, 13.92093238, 7.95542668, 12.00052091,
		11.82680248, 5.89729658, 8.54045647, 13.60981458,
		10.00865388, 7.92837157, 8.31076266, 9.18471422,
		7.84693233, 8.76741161, 10.87795873, 14.65658323,
		7.85521071, 9.04012243, 7.43535867, 10.15812301,
		12.46519105, 7.35042452, 9.95608467, 11.42583285,
		9.83193081, 9.67750682, 11.16649223, 8.94295236,
		10.01809469, 7.17197717, 7.55621033, 13.3999663,
		11.85703991, 9.20101557, 8.29058923, 7.20849446,
		8.86770357, 8.8384832, 8.79774152, 9.26089846,
		8.16864633, 10.87662162, 8.39197205, 7.41328472,
		13.22198834, 11.29517127, 12.1842384, 10.41771674,
		10.8701562, 10.02489038, 12.04101253, 10.32352415,
		10.77943047, 9.12459943, 11.04568103, 13.54620779,
		14.221192, 13.43122872, 8.32564618, 10.43884202,
		10.30555116, 7.36896287, 10.7156544, 10.96224612,
		5.70032716, 8.45044525, 5.51224787, 8.7881203,
	}

	t.Run("S_Empty", func(t *testing.T) {
		s := new(Statistics)
		o := new(Statistics)

		o.Update(values...)
		s.Append(o)

		require.Equal(t, 9.813435956500003, s.Mean())
		require.Equal(t, 2.1848902532568184, s.StdDev())
		require.Equal(t, 4.773745418776644, s.Variance())
		require.Equal(t, 15.45832771, s.Maximum())
		require.Equal(t, 5.51224787, s.Minimum())
		require.Equal(t, 9.94607984, s.Range())
	})

	t.Run("O_Empty", func(t *testing.T) {
		s := new(Statistics)
		o := new(Statistics)

		s.Update(values...)
		s.Append(o)

		require.Equal(t, 9.813435956500003, s.Mean())
		require.Equal(t, 2.1848902532568184, s.StdDev())
		require.Equal(t, 4.773745418776644, s.Variance())
		require.Equal(t, 15.45832771, s.Maximum())
		require.Equal(t, 5.51224787, s.Minimum())
		require.Equal(t, 9.94607984, s.Range())
	})

	t.Run("S_Range", func(t *testing.T) {
		s := new(Statistics)
		o := new(Statistics)

		for i, v := range values {
			if i%2 == 0 {
				s.Update(v)
			} else {
				o.Update(v)
			}
		}

		require.Equal(t, 15.45832771, s.Maximum())
		require.Equal(t, 5.51224787, s.Minimum())

		s.Append(o)

		require.Equal(t, 9.8134359565, s.Mean())
		require.Equal(t, 2.1848902532568206, s.StdDev())
		require.Equal(t, 4.773745418776653, s.Variance())
		require.Equal(t, 15.45832771, s.Maximum())
		require.Equal(t, 5.51224787, s.Minimum())
		require.Equal(t, 9.94607984, s.Range())
	})

	t.Run("O_Range", func(t *testing.T) {
		s := new(Statistics)
		o := new(Statistics)

		for i, v := range values {
			if i%2 == 0 {
				o.Update(v)
			} else {
				s.Update(v)
			}
		}

		require.Equal(t, 15.45832771, o.Maximum())
		require.Equal(t, 5.51224787, o.Minimum())

		s.Append(o)

		require.Equal(t, 9.8134359565, s.Mean())
		require.Equal(t, 2.1848902532568206, s.StdDev())
		require.Equal(t, 4.773745418776653, s.Variance())
		require.Equal(t, 15.45832771, s.Maximum())
		require.Equal(t, 5.51224787, s.Minimum())
		require.Equal(t, 9.94607984, s.Range())
	})

}

func BenchmarkStatistics_Update(b *testing.B) {
	source42 := rand.NewSource(42)
	rand42 := rand.New(source42)
	stats := new(Statistics)

	for i := 0; i < b.N; i++ {
		val := rand42.Float64()
		stats.Update(val)
	}
}

func BenchmarkStatistics_Sequential(b *testing.B) {
	data, _ := loadTestData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stats := new(Statistics)
		for _, val := range data {
			stats.Update(val)
		}
	}
}

func BenchmarkStatistics_BulkLoad(b *testing.B) {
	data, _ := loadTestData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stats := new(Statistics)
		stats.Update(data...)
	}
}
