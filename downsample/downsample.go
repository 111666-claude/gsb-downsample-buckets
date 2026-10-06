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

// Downsample 把样本降采样成桶。
func Downsample(samples []Sample, width, maxLateness int, c *Counter) []Bucket {
	seen := make(map[string]struct{}, len(samples))
	agg := make(map[int]Bucket, len(samples))
	minStart, maxStart := 0, 0
	have := false
	watermark := 0
	hasWatermark := false
	for _, s := range samples {
		c.Scanned++
		if _, ok := seen[s.ID]; ok {
			continue
		}
		seen[s.ID] = struct{}{}
		if hasWatermark && watermark-s.TS > maxLateness {
			continue
		}
		if !hasWatermark || s.TS > watermark {
			watermark = s.TS
			hasWatermark = true
		}
		start := s.TS - floorMod(s.TS, width)
		b := agg[start]
		b.Start = start
		b.Sum += s.Value
		b.Count++
		agg[start] = b
		if !have || start < minStart {
			minStart = start
		}
		if !have || start > maxStart {
			maxStart = start
		}
		have = true
	}
	if !have {
		return []Bucket{}
	}
	out := make([]Bucket, 0, (maxStart-minStart)/width+1)
	for start := minStart; start <= maxStart; start += width {
		if b, ok := agg[start]; ok {
			out = append(out, b)
		} else {
			out = append(out, Bucket{Start: start})
		}
	}
	return out
}

// floorMod 返回 ts 对 width 的向下取整余数（width 为正）。
func floorMod(ts, width int) int {
	r := ts % width
	if r < 0 {
		r += width
	}
	return r
}
