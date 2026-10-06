package downsample

import "testing"

func TestEmptySamples(t *testing.T) {
	if got := Downsample(nil, 10, 5, &Counter{}); len(got) != 0 {
		t.Fatalf("空输入应无桶：%+v", got)
	}
}

func TestSingleSample(t *testing.T) {
	got := Downsample([]Sample{{"a", 0, 5}}, 10, 5, &Counter{})
	if len(got) != 1 || got[0].Sum != 5 {
		t.Fatalf("应有一个和为 5 的桶：%+v", got)
	}
}

func TestTwoAdjacentBuckets(t *testing.T) {
	got := Downsample([]Sample{{"a", 0, 1}, {"b", 10, 1}}, 10, 100, &Counter{})
	if len(got) != 2 {
		t.Fatalf("应有两个桶：%+v", got)
	}
}

func TestSumAccumulates(t *testing.T) {
	got := Downsample([]Sample{{"a", 0, 1}, {"b", 0, 2}}, 10, 100, &Counter{})
	if got[0].Sum != 3 {
		t.Fatalf("同桶应累加为 3：%+v", got)
	}
}

func TestBucketStartField(t *testing.T) {
	got := Downsample([]Sample{{"a", 0, 1}}, 10, 100, &Counter{})
	if got[0].Start != 0 {
		t.Fatalf("起点应是 0：%+v", got)
	}
}

func TestCounterCounts(t *testing.T) {
	c := &Counter{}
	Downsample([]Sample{{"a", 0, 1}}, 10, 100, c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}
