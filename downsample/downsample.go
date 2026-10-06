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
	out := []Bucket{}
	for _, s := range samples {
		c.Scanned++
		idx := -1
		for i := range out {
			c.Scanned++
			if out[i].Start == s.TS {
				idx = i
				break
			}
		}
		if idx < 0 {
			out = append(out, Bucket{Start: s.TS})
			idx = len(out) - 1
		}
		out[idx].Sum += s.Value
		out[idx].Count++
	}
	return out
}
