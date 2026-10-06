# downsample-buckets

把时间序列降采样成固定宽度的桶：对齐、聚合、补缺口，并按迟到窗口丢弃过期样本。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/downsample --sample=align
go run ./cmd/downsample --sample=gaps
go run ./cmd/downsample --sample=late
go run ./cmd/downsample --sample=dupid
go run ./cmd/downsample --sample=work
```

## 口径

- **对齐**：桶起点是时间戳向下取整到宽度的整数倍。
- **聚合**：每个桶累加样本值与条数。
- **补缺口**：结果里相邻桶之间没有样本的桶也要补出来，求和与条数为零。
- **迟到**：时间戳比最大时间戳早超过迟到窗口的样本被丢弃。
- **去重**：同一样本 ID 只计一次。

## 不变量

- 输出桶按起点升序，起点都是宽度的整数倍，相邻起点相差一个宽度。
- 每个样本至多进入一个桶一次。
- `scanned` 不随样本数平方级放大：两千个样本的 `scanned` 不超过 6000。

## 输出契约

场景打印一行 JSON，键随场景而定：`buckets` 是桶起点列表，`sum0` 是起点为 0 的桶的和。
`--sample=work` 打印 `{"samples": N, "scanned": N}`。
