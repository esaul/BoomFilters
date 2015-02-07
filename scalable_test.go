package boom

import (
	"strconv"
	"testing"
)

// Ensures that NewDefaultScalableBloomFilter creates a Scalable Bloom Filter
// with hint = 10000 and r = 0.8.
func TestNewDefaultScalableBloomFilter(t *testing.T) {
	f := NewDefaultScalableBloomFilter(0.1)

	if f.fp != 0.1 {
		t.Errorf("Expected 0.1, got %f", f.fp)
	}

	if f.hint != 10000 {
		t.Errorf("Expected 10000, got %d", f.hint)
	}

	if f.r != 0.8 {
		t.Errorf("Expected 0.8, got %f", f.r)
	}
}

// Ensures that Capacity returns the sum of the capacities for the contained
// Bloom filters.
func TestScalableBloomCapacity(t *testing.T) {
	f := NewScalableBloomFilter(1, 0.1, 1)
	f.addFilter()
	f.addFilter()

	if capacity := f.Capacity(); capacity != 15 {
		t.Errorf("Expected 15, got %d", capacity)
	}
}

// Ensures that K returns the number of hash functions used in each Bloom
// filter.
func TestScalableBloomK(t *testing.T) {
	f := NewScalableBloomFilter(10, 0.1, 0.8)

	if k := f.K(); k != 4 {
		t.Errorf("Expected 4, got %d", k)
	}
}

// Ensures that FillRatio returns the average fill ratio of the contained
// filters.
func TestScalableFillRatio(t *testing.T) {
	f := NewScalableBloomFilter(1, 0.1, 0.8)
	for i := 0; i < 100; i++ {
		f.Add([]byte(strconv.Itoa(i)))
	}

	if ratio := f.FillRatio(); ratio != 0.5 {
		t.Errorf("Expected 0.5, got %f", ratio)
	}
}

// Ensures that Test, Add, and TestAndAdd behave correctly.
func TestScalableBloomTestAndAdd(t *testing.T) {
	f := NewScalableBloomFilter(100, 0.01, 0.8)

	// `a` isn't in the filter.
	if f.Test([]byte(`a`)) {
		t.Error("`a` should not be a member")
	}

	if f.Add([]byte(`a`)) != f {
		t.Error("Returned ScalableBloomFilter should be the same instance")
	}

	// `a` is now in the filter.
	if !f.Test([]byte(`a`)) {
		t.Error("`a` should be a member")
	}

	// `a` is still in the filter.
	if !f.TestAndAdd([]byte(`a`)) {
		t.Error("`a` should be a member")
	}

	// `b` is not in the filter.
	if f.TestAndAdd([]byte(`b`)) {
		t.Error("`b` should not be a member")
	}

	// `a` is still in the filter.
	if !f.Test([]byte(`a`)) {
		t.Error("`a` should be a member")
	}

	// `b` is now in the filter.
	if !f.Test([]byte(`b`)) {
		t.Error("`b` should be a member")
	}

	// `c` is not in the filter.
	if f.Test([]byte(`c`)) {
		t.Error("`c` should not be a member")
	}

	for i := 0; i < 10000; i++ {
		f.Add([]byte(strconv.Itoa(i)))
	}

	// `x` should not be a false positive.
	if f.Test([]byte(`x`)) {
		t.Error("`x` should not be a member")
	}
}

// Ensures that Reset removes all Bloom filters and resets the initial one.
func TestScalableBloomReset(t *testing.T) {
	f := NewScalableBloomFilter(10, 0.1, 0.8)
	for i := 0; i < 1000; i++ {
		f.Add([]byte(strconv.Itoa(i)))
	}

	if len(f.filters) < 2 {
		t.Errorf("Expected more than 1 filter, got %d", len(f.filters))
	}

	if f.Reset() != f {
		t.Error("Returned ScalableBloomFilter should be the same instance")
	}

	if len(f.filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(f.filters))
	}

	for _, partition := range f.filters[0].partitions {
		if partition.Any() {
			t.Error("Expected all bits to be unset")
		}
	}
}

func BenchmarkScalableBloomAdd(b *testing.B) {
	b.StopTimer()
	f := NewScalableBloomFilter(100000, 0.1, 0.8)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		f.Add(data[n])
	}
}

func BenchmarkScalableBloomTest(b *testing.B) {
	b.StopTimer()
	f := NewScalableBloomFilter(100000, 0.1, 0.8)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		f.Test(data[n])
	}
}

func BenchmarkScalableBloomTestAndAdd(b *testing.B) {
	b.StopTimer()
	f := NewScalableBloomFilter(100000, 0.1, 0.8)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		f.TestAndAdd(data[n])
	}
}
