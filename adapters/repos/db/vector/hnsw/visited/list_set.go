//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package visited

// ListSet is a reusable list with very efficient resets. Inspired by the C++
// implementation in hnswlib it can be reset with zero memrory writes in the
// array by moving the match target instead of altering the list. Only after a
// version overflow do we need to actually reset
//
// The new implemtation uses a slice where the first element is reserved for the marker.
// This allow us to use ListSet as a value (i.e. no pointer is required)
// The marker (i.e. set[0]) allows for reusing the same list without having to zero all elements on each list reset.
// Resetting the list takes place once the marker (i.e. set[0]) overflows
type ListSet struct {
	set []uint8 // set[0] is reserved for the marker (version)
}

//  Len returns the number of elements in the list.
func (l ListSet) Len() uint64 { return uint64(len(l.set)) - 1 }

// Free allocated slice. This list should not be resuable after this call.
func (l *ListSet) Free() { l.set = nil }

// NewList creates a new list. It allocates memory for elements and marker
func NewList(size int) ListSet {
	set := make([]uint8, size+1)
	set[0] = 1 // the marker starts always by 1 since on reset all element are set to 0
	return ListSet{set: set}
}

// Visit sets element at node to the marker value
func (l *ListSet) Visit(node uint64) {
	if node >= l.Len() { // resize
		newset := make([]uint8, node+1024)
		copy(newset, l.set)
		l.set = newset
	}
	l.set[node+1] = l.set[0]
}

// Visited checks if l contains the specified node
func (l *ListSet) Visited(node uint64) bool {
	return node < l.Len() && l.set[node+1] == l.set[0]
}

// Reset list only in case of an overflow.
func (l *ListSet) Reset() {
	l.set[0]++
	if l.set[0] == 0 { // if overflowed
		for i := range l.set {
			l.set[i] = 0
		}
		l.set[0] = 1 // restart counting
	}
}
