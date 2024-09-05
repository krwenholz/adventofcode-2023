package dayTwentyThree

// CellHeap is a heap implementation for the Cell struct
type CellHeap []*Cell

// Len returns the length of the heap
func (h CellHeap) Len() int {
	return len(h)
}

// Less compares two cells based on their f values
func (h CellHeap) Less(i, j int) bool {
	if h[i].f == h[j].f {
		return h[i].g < h[j].g
		//if len(h[i].prevs) < len(h[j].prevs) {
		//	return true
		//}
	}
	return h[i].f < h[j].f
}

// Swap swaps two cells in the heap
func (h CellHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push adds a cell to the heap
func (h *CellHeap) Push(x interface{}) {
	*h = append(*h, x.(*Cell))
}

// Pop removes and returns the top cell from the heap
func (h *CellHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
