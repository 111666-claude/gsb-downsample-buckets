package downsample

import "fmt"

// Scenario 是一个固定场景，用来复现题面里的锚点。
type Scenario struct {
	Name string
	Run  func() map[string]any
}

func starts(bs []Bucket) []int {
	out := []int{}
	for _, b := range bs {
		out = append(out, b.Start)
	}
	return out
}

func sumAt(bs []Bucket, start int) int {
	for _, b := range bs {
		if b.Start == start {
			return b.Sum
		}
	}
	return 0
}

// Samples 返回全部场景。
func Samples() []Scenario {
	return []Scenario{
		{Name: "align", Run: func() map[string]any {
			ss := []Sample{{"a", 3, 1}, {"b", 15, 1}}
			return map[string]any{"buckets": starts(Downsample(ss, 10, 100, &Counter{}))}
		}},
		{Name: "gaps", Run: func() map[string]any {
			ss := []Sample{{"a", 5, 1}, {"b", 25, 1}}
			return map[string]any{"buckets": starts(Downsample(ss, 10, 100, &Counter{}))}
		}},
		{Name: "late", Run: func() map[string]any {
			ss := []Sample{{"a", 0, 1}, {"b", 20, 2}, {"c", 0, 100}}
			return map[string]any{"sum0": sumAt(Downsample(ss, 10, 5, &Counter{}), 0)}
		}},
		{Name: "dupid", Run: func() map[string]any {
			ss := []Sample{{"x", 0, 1}, {"x", 0, 1}}
			return map[string]any{"sum0": sumAt(Downsample(ss, 10, 100, &Counter{}), 0)}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	ss := make([]Sample, 0, n)
	for i := 0; i < n; i++ {
		ss = append(ss, Sample{ID: fmt.Sprintf("s%04d", i), TS: i, Value: 1})
	}
	c := &Counter{}
	Downsample(ss, 100000, 1000000, c)
	return map[string]any{"samples": n, "scanned": c.Scanned}
}
