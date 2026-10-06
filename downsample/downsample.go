// Package downsample 做时间序列降采样。
package downsample

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Sample 是一个样本。
type Sample struct {
	ID    string
	TS    int
	Value int
}

// Bucket 是一个聚合桶。
type Bucket struct {
	Start int
	Sum   int
	Count int
}

// bucketStart 把时间戳向下取整到宽度的整数倍（对负数同样成立）。
func bucketStart(ts, width int) int {
	q := ts / width
	if r := ts % width; r < 0 {
		q--
	}
	return q * width
}

// Downsample 把样本降采样成桶。
func Downsample(samples []Sample, width, maxLateness int, c *Counter) []Bucket {
	if len(samples) == 0 || width <= 0 {
		return []Bucket{}
	}

	seen := make(map[string]bool)
	sums := make(map[int]int)
	counts := make(map[int]int)

	haveWatermark := false
	watermark := 0
	haveRange := false
	minStart, maxStart := 0, 0

	for _, s := range samples {
		if c != nil {
			c.Scanned++
		}

		// 去重：同一样本 ID 只计第一次出现。
		if seen[s.ID] {
			continue
		}
		seen[s.ID] = true

		// 迟到：比已见到的最大时间戳早超过迟到窗口的样本丢弃。
		if haveWatermark && s.TS < watermark-maxLateness {
			continue
		}
		if !haveWatermark || s.TS > watermark {
			watermark = s.TS
		}
		haveWatermark = true

		start := bucketStart(s.TS, width)
		if !haveRange {
			minStart, maxStart = start, start
			haveRange = true
		} else if start < minStart {
			minStart = start
		} else if start > maxStart {
			maxStart = start
		}
		sums[start] += s.Value
		counts[start]++
	}

	if !haveRange {
		return []Bucket{}
	}

	// 补缺口：按起点升序输出，空桶求和与条数为零。
	out := []Bucket{}
	for start := minStart; start <= maxStart; start += width {
		out = append(out, Bucket{
			Start: start,
			Sum:   sums[start],
			Count: counts[start],
		})
	}
	return out
}
