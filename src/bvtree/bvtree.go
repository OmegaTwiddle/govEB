package bvtree

import (
    "fmt"
)

/**
 * BvTree is a struct that holds a bitvector representation
 * of a set of integeres between 0 and n.
 */
type BvTree struct {

    bitvector []uint64
}

/**
 * returns true if the bvTree contains the given uint64.
 */
func (bvTree * BvTree) Contains(n uint64) bool {
    idx, off := offsets(n)
    b := uint64(1 << off)
    if (bvTree.bitvector[idx] & b) == 0 {
        return false
    }
    return true
}

/**
 * Inserts the integer n into the bvTree.
 */
func (bvTree *BvTree) Insert(n uint64) {
    idx, off := offsets(n)
    b := uint64(1 << off)
    bvTree.bitvector[idx] |= b
}

func BuildBvTree(size uint64) BvTree {
    result := BvTree{}

    // The number of bits we need is size / 64
    numBits := size / 64

    // if size is too small, just allocate one.
    if numBits == 0 {
        numBits = 1
    }

    fmt.Printf("building a tree with %d uint64s.\n", numBits)
    result.bitvector = make([]uint64, numBits)
    return result
}


func offsets(n uint64) (uint64, uint64) {
    return n / 64, n % 64
}
