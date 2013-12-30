package bvtree

import (
    "fmt"
)

/**
 * BvFhTree is a struct that holds a bitvector representation
 * of a set of integeres between 0 and n.
 */
type BvFhTree struct {

    // The number of bits in the bitvector, size of the universe.
    numBits uint64

    sqNumBits uint64

    empty bool

    // Bit vector holding a summary tree of fixed height
    summary []uint64

    // Bit vector holding the actual values in the tree.
    bitvector []uint64
}

func (bvTree *BvFhTree) Min() uint64 {
    if bvTree.empty {
        panic("No min on an empty tree...")
    }

    idx := uint64(0)
    for i := uint64(0); i < bvTree.sqNumBits; i++ {
        if bvTree.hasSumBit(i) {
            idx = i
            break
        }
    }

    min, max := bvTree.childrenRange(idx)

    for i := min; i <= max; i++ {
        if bvTree.hasBvBit(i) {
            return i
        }
    }

    panic("Min didn't find anything...")
}

func (bvTree *BvFhTree) Max() uint64 {
    if bvTree.empty {
        panic("No min on an empty tree...")
    }

    idx := bvTree.sqNumBits
    for i := bvTree.sqNumBits - 1; i >= uint64(0); i-- {
        if bvTree.hasSumBit(i) {
            idx = i
            break
        }
    }

    min, max := bvTree.childrenRange(idx)

    for i := max; i >= min; i-- {
        if bvTree.hasBvBit(i) {
            return i
        }
    }

    panic("Min didn't find anything...")
}



/**
 * Returns the number below n in the tree.
 * Assumes that the number passed in is greater than
 * the min value.
 */
func (bvTree *BvFhTree) Predecessor(n uint64) uint64 {

    // First check the sibling range.
    min, _ := bvTree.siblingRange(n)
    for i := n - 1; i >= min; i-- {
        if bvTree.hasBvBit(i) {
            return i
        }
    }

    // If we haven't found it yet, find the next bit in 
    // the summary vector
    sumToSearch := bvTree.sumIndex(n) - 1
    for i := sumToSearch; i >= 0; i-- {
        if bvTree.hasSumBit(i) {
            sumToSearch = i
            break
        }
    }

    min, max := bvTree.childrenRange(sumToSearch)
    for i := max; i >= min; i-- {
        if bvTree.hasBvBit(i) {
            return i
        }
    }

    panic("There was a problem with predecessor.")
}



/**
 * Returns the number above n in the tree.
 * Assumes that the number passed in is less than
 * the max value.
 */
func (bvTree *BvFhTree) Successor(n uint64) uint64 {

    // First check the sibling range.
    _, max := bvTree.siblingRange(n)
    for i := n + 1; i <= max; i++ {
        if bvTree.hasBvBit(i) {
            return i
        }
    }

    // If we haven't found it yet, find the next bit in 
    // the summary vector
    sumToSearch := bvTree.sumIndex(n) + 1
    for i := sumToSearch; i < bvTree.sqNumBits; i++ {
        if bvTree.hasSumBit(i) {
            sumToSearch = i
            break
        }
    }

    min, max := bvTree.childrenRange(sumToSearch)
    for i := min; i <= max; i++ {
        if bvTree.hasBvBit(i) {
            return i
        }
    }

    panic("There was a problem with successor.")
}



/**
 * returns true if the bvTree contains the given uint64.
 */
func (bvTree *BvFhTree) Contains(n uint64) bool {
    idx, off := offsets(n)
    b := uint64(1 << (63 - off))
    if (bvTree.bitvector[idx] & b) == 0 {
        return false
    }
    return true
}

/**
 * Inserts the integer n into the bvTree.
 */
func (bvTree *BvFhTree) Insert(n uint64) {
    // Add the bit to the data.
    idx, off := offsets(n)
    b := uint64(1 << (63 - off))
    bvTree.bitvector[idx] |= b

    // Update the supporting binary tree.
    sIdx := bvTree.sumIndex(n)
    idx, off = offsets(sIdx)
    b = uint64(1 << (63 - off))
    bvTree.summary[idx] |= b
    bvTree.empty = false
}

func (bvTree *BvFhTree) Remove(n uint64) {
    // Rmove from the bitvector
    idx, off := offsets(n)
    b := ^uint64(1 << (63 - off))
    bvTree.bitvector[idx] &= b

    min, max := bvTree.siblingRange(n)

    if bvTree.emptyRange(min, max) {
        idx, off := offsets(bvTree.sumIndex(n))
        b = ^uint64(1 << (63 - off))
        bvTree.summary[idx] &= b
    }

}

func getFhNumUints(numBits uint64) (uint64, uint64) {
    result := uint64(2)
    for result * result < numBits {
        result *= 2
    }
    if result * result <= uint64(64) {
        return 1, 1
    }
    return (result / uint64(64)), (result * result / uint64(64))
}

func BuildBvFhTree(numBits uint64) *BvFhTree {
    result := BvFhTree{}

    // The number of uints we need is size / 64
    numSumUints, numBvUints := getFhNumUints(numBits)

    fmt.Printf("building a tree with %d/%d uint64s.\n", numSumUints, numBvUints)
    result.numBits = numBvUints * uint64(64)
    result.sqNumBits = getRoot(result.numBits)
    result.bitvector = make([]uint64, numBvUints)
    result.summary = make([]uint64, numSumUints)
    return &result
}

func (bvTree *BvFhTree) siblingRange(n uint64) (uint64, uint64) {
    min := (n / bvTree.sqNumBits) * bvTree.sqNumBits
    max := min + bvTree.sqNumBits
    return min, max
}

func (bvTree *BvFhTree) emptyRange(min uint64, max uint64) bool {
    for i := min; i <= max; i++ {
        idx, off := offsets(i)
        b := uint64(1 << (63 - off))
        val := bvTree.bitvector[idx] & b
        if val != 0 {
            return false
        }
    }
    return true
}

// Return true if the supporting tree has the bit.
func (bvTree *BvFhTree) hasSumBit(pos uint64) bool {
    idx, off := offsets(pos)
    return (bvTree.summary[idx] & uint64(1 << (63 - off))) != 0
}

// Return true if the bitvector has the bit.
func (bvTree *BvFhTree) hasBvBit(pos uint64) bool {
    idx, off := offsets(pos)
    return (bvTree.bitvector[idx] & uint64(1 << (63 - off))) != 0
}
// Assuming the bitvector is of size 2^n, get the first index of the
// last level in the supporting tree.
// E.G. 
func (bvTree *BvFhTree) llIndex() uint64 {
    return bvTree.numBits / 2 - 1
}

/**
 * Returns the last index in the lower level of the supporting tree.
 */
func (bvTree *BvFhTree) maxLlIndex() uint64 {
    return bvTree.numBits - 2
}


/**
 * Given an index at the lowest level of the support tree,
 * returns the left and right "children" indices inside
 * the bit vector.
 */
func (bvTree *BvFhTree) bvIndices(n uint64) (uint64, uint64) {
    k := (n - (bvTree.numBits / 2 - 1)) * 2
    return k, k + 1
}

/**
 * Returns the min/max of the range from the given index
 * into the summary bitvector
 */
func (bvTree *BvFhTree) childrenRange(n uint64) (uint64, uint64) {
    return n * bvTree.sqNumBits, (n + 1) * bvTree.sqNumBits - 1
}

func (bvTree *BvFhTree) sumIndex(n uint64) uint64 {
    return (n / bvTree.sqNumBits)
}

func (bvTree *BvFhTree) checkBit(n uint64) {
    if bvTree.hasSumBit(n) {
        fmt.Printf("Has Summary bit:     %d\n", n)
    } else {
        fmt.Printf(" Didn't have St bit: %d\n", n)
    }
}

func (bvTree *BvFhTree) DbgPrint() {
    fmt.Println("DbgPrint: ")
    fmt.Println("suptree")
    for _, val := range(bvTree.summary) {
        dbgPrintBin(val)
    }
    fmt.Println("\nbitvector")
    for _, val := range(bvTree.bitvector) {
        dbgPrintBin(val)
    }
    fmt.Println(" ")
}
