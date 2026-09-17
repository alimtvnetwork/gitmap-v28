package cmdpipeline

import (
	"runtime"
	"sync"
)

// SliceLineMarker manages non-mutating two-pass line marking and extraction.
type SliceLineMarker struct {
	mu         sync.Mutex
	keepMask   []bool
	keptCounts []int
}

// NewSliceLineMarker initializes a marker structure with pre-allocated mask.
func NewSliceLineMarker(size int, workers int) *SliceLineMarker {
	return &SliceLineMarker{
		keepMask:   make([]bool, size),
		keptCounts: make([]int, workers),
	}
}

// MarkWorkerPartition evaluates lines in a chunk and records keep state in mask.
func (m *SliceLineMarker) MarkWorkerPartition(
	lines []string,
	startIdx int,
	endIdx int,
	workerId int,
	isKeepPredicate func(string) bool,
) {
	localKept := 0
	for i := startIdx; i < endIdx; i++ {
		isKeep := isKeepPredicate(lines[i])
		m.keepMask[i] = isKeep
		if isKeep {
			localKept++
		}
	}
	m.mu.Lock()
	m.keptCounts[workerId] = localKept
	m.mu.Unlock()
}

// TotalKeptCount calculates the exact destination array capacity across all workers.
func (m *SliceLineMarker) TotalKeptCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	total := 0
	for _, count := range m.keptCounts {
		total += count
	}

	return total
}

// MaterializeKeptLines creates the pre-allocated slice and copies kept lines.
func (m *SliceLineMarker) MaterializeKeptLines(lines []string) []string {
	totalKept := m.TotalKeptCount()
	if totalKept == 0 {
		return []string{}
	}

	output := make([]string, totalKept)
	writeIdx := 0
	for i, isKeep := range m.keepMask {
		if isKeep {
			output[writeIdx] = lines[i]
			writeIdx++
		}
	}

	return output
}

// ParallelFilterLines filters a line slice using a non-mutating two-pass algorithm.
func ParallelFilterLines(lines []string, isKeepPredicate func(string) bool) []string {
	total := len(lines)
	if total == 0 {
		return []string{}
	}
	if total <= 128 {
		return filterLinesSequential(lines, isKeepPredicate)
	}

	return filterLinesParallelTwoPass(lines, isKeepPredicate)
}

func filterLinesSequential(lines []string, isKeepPredicate func(string) bool) []string {
	mask := make([]bool, len(lines))
	keptCount := 0
	for i, line := range lines {
		isKeep := isKeepPredicate(line)
		mask[i] = isKeep
		if isKeep {
			keptCount++
		}
	}

	return buildMaterializedSlice(lines, mask, keptCount)
}

func buildMaterializedSlice(lines []string, mask []bool, keptCount int) []string {
	result := make([]string, keptCount)
	slot := 0
	for i, isKeep := range mask {
		if isKeep {
			result[slot] = lines[i]
			slot++
		}
	}

	return result
}

func filterLinesParallelTwoPass(lines []string, isKeepPredicate func(string) bool) []string {
	workerCount := resolveFilterWorkerCount(len(lines))
	marker := NewSliceLineMarker(len(lines), workerCount)
	chunkSize := (len(lines) + workerCount - 1) / workerCount

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		startIdx := w * chunkSize
		endIdx := resolvePartitionEnd(startIdx+chunkSize, len(lines))
		if startIdx >= endIdx {
			continue
		}
		wg.Add(1)
		go dispatchFilterWorker(&wg, marker, lines, startIdx, endIdx, w, isKeepPredicate)
	}
	wg.Wait()

	return marker.MaterializeKeptLines(lines)
}

func dispatchFilterWorker(
	wg *sync.WaitGroup,
	marker *SliceLineMarker,
	lines []string,
	startIdx int,
	endIdx int,
	workerId int,
	isKeepPredicate func(string) bool,
) {
	defer wg.Done()
	marker.MarkWorkerPartition(lines, startIdx, endIdx, workerId, isKeepPredicate)
}

func resolveFilterWorkerCount(totalLines int) int {
	cpuCores := runtime.NumCPU()
	if cpuCores < 1 {
		cpuCores = 1
	}
	maxWorkers := cpuCores
	if maxWorkers > 8 {
		maxWorkers = 8
	}
	neededWorkers := (totalLines + 255) / 256
	if neededWorkers < maxWorkers {
		return neededWorkers
	}

	return maxWorkers
}

func resolvePartitionEnd(targetEnd, totalLines int) int {
	if targetEnd > totalLines {
		return totalLines
	}

	return targetEnd
}
