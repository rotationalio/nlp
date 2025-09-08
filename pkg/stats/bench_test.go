package stats

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func loadBenchData() ([]time.Duration, error) {
	data := make([]time.Duration, 0, 1000000)

	f, err := os.Open("testdata/latencies.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewScanner(f)
	for buf.Scan() {
		var val time.Duration
		if val, err = time.ParseDuration(buf.Text()); err != nil {
			return nil, err
		}
		data = append(data, val)
	}

	if buf.Err() != nil {
		return nil, err
	}

	return data, nil
}

func ExampleBenchmark() {
	stats := new(Benchmark)
	samples, _ := loadBenchData()

	for _, sample := range samples {
		stats.Update(sample)
	}

	data, _ := json.MarshalIndent(stats.Serialize(), "", "  ")
	fmt.Println(string(data))
	// Output:
	// {
	//   "duration": "0s",
	//   "fastest": "41.219436ms",
	//   "mean": "120.993689ms",
	//   "range": "167.175236ms",
	//   "samples": 1000000,
	//   "slowest": "208.394672ms",
	//   "stddev": "17.283562ms",
	//   "throughput": 8.264893850648656,
	//   "timeouts": 0,
	//   "total": "33h36m33.689461785s",
	//   "variance": "298.721µs"
	// }
}

func TestBenchmark(t *testing.T) {
	data, err := loadBenchData()
	require.NoError(t, err)
	require.NotNil(t, data)

	stats := new(Benchmark)

	for _, v := range data {
		stats.Update(v)
	}

	require.Equal(t, time.Duration(120993689), stats.Mean())
	require.Equal(t, time.Duration(17283562), stats.StdDev())
	require.Equal(t, time.Duration(298721), stats.Variance())
	require.Equal(t, time.Duration(208394672), stats.Slowest())
	require.Equal(t, time.Duration(41219436), stats.Fastest())
	require.Equal(t, time.Duration(167175236), stats.Range())
}

func TestBenchmarkBulk(t *testing.T) {
	data, err := loadBenchData()
	require.NoError(t, err)
	require.NotNil(t, data)

	stats := new(Benchmark)
	stats.Update(data...)

	require.Equal(t, time.Duration(120993689), stats.Mean())
	require.Equal(t, time.Duration(17283562), stats.StdDev())
	require.Equal(t, time.Duration(298721), stats.Variance())
	require.Equal(t, time.Duration(208394672), stats.Slowest())
	require.Equal(t, time.Duration(41219436), stats.Fastest())
	require.Equal(t, time.Duration(167175236), stats.Range())
}

func TestThroughput(t *testing.T) {
	stats := new(Benchmark)

	// With no samples, throughput should be zero
	require.Equal(t, 0.0, stats.Throughput())

	latencies := []time.Duration{
		175 * time.Millisecond, 250 * time.Millisecond,
		225 * time.Millisecond, 150 * time.Millisecond,
		200 * time.Millisecond,
	}

	// With no specified duration, throughput should equal total
	stats.Update(latencies...)
	require.InDelta(t, 5.0, stats.Throughput(), 0.00001)

	// With a specified duration, throughput should use that duration
	stats.SetDuration(250 * time.Millisecond)
	require.InDelta(t, 20.0, stats.Throughput(), 0.00001)
}

func TestBenchmarkAppend(t *testing.T) {
	values := []time.Duration{
		15458327, 11117278, 10308557, 14636267, 5854742, 10374731,
		11020685, 9921715, 9454425, 11848154, 11987220, 11544855,
		8491874, 8327981, 9855619, 8647359, 6200921, 0,
		11797218, 7572802, 7328019, 11717603, 10270390, 12527268,
		8844019, 6797831, 5512247, 7539891, 9297131, 10675063,
		6634836, 0, 9936534, 13920932, 7955426, 12000520,
		11826802, 5897296, 8540456, 13609814, 10008653, 7928371,
		8310762, 9184714, 7846932, 8767411, 10877958, 14656583,
		7855210, 9040122, 7435358, 10158123, 12465191, 7350424,
		9956084, 11425832, 9831930, 9677506, 11166492, 8942952,
		10018094, 7171977, 7556210, 13399966, 11857039, 9201015,
		8290589, 7208494, 8867703, 8838483, 8797741, 9260898,
		8168646, 10876621, 8391972, 7413284, 13221988, 11295171,
		12184238, 10417716, 10870156, 10024890, 12041012, 10323524,
		10779430, 9124599, 0, 13546207, 14221192, 0,
		8325646, 10438842, 10305551, 7368962, 10715654, 10962246,
		5700327, 8450445, 8788120, 7426879,
	}

	t.Run("S_Empty", func(t *testing.T) {
		s := new(Benchmark)
		o := new(Benchmark)

		o.Update(values...)
		s.Append(o)

		require.Equal(t, time.Duration(9791572), s.Mean())
		require.Equal(t, time.Duration(2180586), s.StdDev())
		require.Equal(t, time.Duration(4754), s.Variance())
		require.Equal(t, time.Duration(15458327), s.Slowest())
		require.Equal(t, time.Duration(5512247), s.Fastest())
		require.Equal(t, time.Duration(9946080), s.Range())
		require.Equal(t, time.Duration(939990942), s.Total())
		require.InDelta(t, (102.12864359481392), s.Throughput(), DEFAULT_DELTA)
		require.Equal(t, uint64(4), s.timeouts)
	})

	t.Run("O_Empty", func(t *testing.T) {
		s := new(Benchmark)
		o := new(Benchmark)

		s.Update(values...)
		s.Append(o)

		require.Equal(t, s.Mean(), time.Duration(9791572))
		require.Equal(t, s.StdDev(), time.Duration(2180586))
		require.Equal(t, s.Variance(), time.Duration(4754))
		require.Equal(t, s.Slowest(), time.Duration(15458327))
		require.Equal(t, s.Fastest(), time.Duration(5512247))
		require.Equal(t, s.Range(), time.Duration(9946080))

		require.Equal(t, s.Total(), time.Duration(939990942))
		require.InDelta(t, s.Throughput(), 102.12864359481392, DEFAULT_DELTA)
		require.Equal(t, s.timeouts, uint64(4))
	})

	t.Run("S_Range", func(t *testing.T) {
		s := new(Benchmark)
		o := new(Benchmark)

		for i, v := range values {
			if i%2 == 0 {
				s.Update(v)
			} else {
				o.Update(v)
			}
		}

		require.Equal(t, s.Slowest(), time.Duration(15458327))
		require.Equal(t, s.Fastest(), time.Duration(5512247))

		s.Append(o)

		require.Equal(t, time.Duration(9791572), s.Mean())
		require.Equal(t, time.Duration(2180586), s.StdDev())
		require.Equal(t, time.Duration(4754), s.Variance())
		require.Equal(t, time.Duration(15458327), s.Slowest())
		require.Equal(t, time.Duration(5512247), s.Fastest())
		require.Equal(t, time.Duration(9946080), s.Range())

		require.Equal(t, time.Duration(939990943), s.Total())
		require.InDelta(t, 102.12864359481385, s.Throughput(), DEFAULT_DELTA)
		require.Equal(t, uint64(4), s.timeouts)
	})

	t.Run("O_Range", func(t *testing.T) {
		s := new(Benchmark)
		o := new(Benchmark)

		for i, v := range values {
			if i%2 == 0 {
				o.Update(v)
			} else {
				s.Update(v)
			}
		}

		require.Equal(t, time.Duration(15458327), o.Slowest())
		require.Equal(t, time.Duration(5512247), o.Fastest())

		s.Append(o)

		require.Equal(t, time.Duration(9791572), s.Mean())
		require.Equal(t, time.Duration(2180586), s.StdDev())
		require.Equal(t, time.Duration(4754), s.Variance())
		require.Equal(t, time.Duration(15458327), s.Slowest())
		require.Equal(t, time.Duration(5512247), s.Fastest())
		require.Equal(t, time.Duration(9946080), s.Range())

		require.Equal(t, time.Duration(939990943), s.Total())
		require.InDelta(t, 102.12864359481385, s.Throughput(), DEFAULT_DELTA)
		require.Equal(t, uint64(4), s.timeouts)
	})

}

func BenchmarkBenchmark_Update(b *testing.B) {
	source42 := rand.NewSource(42)
	rand42 := rand.New(source42)
	stats := new(Benchmark)

	for i := 0; i < b.N; i++ {
		val := time.Duration(rand42.Int31n(1000)) * time.Millisecond
		stats.Update(val)
	}
}

func BenchmarkBenchmark_Sequential(b *testing.B) {
	data, _ := loadBenchData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stats := new(Benchmark)
		for _, val := range data {
			stats.Update(val)
		}
	}
}

func BenchmarkBenchmark_BulkLoad(b *testing.B) {
	data, _ := loadBenchData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stats := new(Benchmark)
		stats.Update(data...)
	}
}
