package bvtree

import (
    "fmt"
//    "strconv"
)

/**
 * BvTree is a struct that holds a bitvector representation
 * of a set of integeres between 0 and n.
 */
type BvTree struct {

    // The number of bits in the bitvector, size of the universe.
    numBits uint64

    // Bit vector holding the supporting binary tree.
    suptree []uint64

    // Bit vector holding the actual values in the tree.
    bitvector []uint64
}

func (bvTree *BvTree) Min() uint64 {
    return 0
}

func (bvTree *BvTree) Max() uint64 {
    return 0
}

/**
 * Returns the number above n in the tree.
 */
func (bvTree *BvTree) Successor(n uint64) uint64 {
    return 0
}

/**
 * returns true if the bvTree contains the given uint64.
 */
func (bvTree * BvTree) Contains(n uint64) bool {
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
func (bvTree *BvTree) Insert(n uint64) {
    // Add the bit to the data.
    idx, off := offsets(n)
    b := uint64(1 << (63 - off))
    bvTree.bitvector[idx] |= b

    // Update the supporting binary tree.
    sIdx := bvTree.supIndex(n)
    for sIdx > 0 {
        idx, off = offsets(sIdx)
        b = uint64(1 << (63 - off))
        fmt.Println("inserting into suptree", idx, off)
        bvTree.suptree[idx] |= b
        sIdx = parentIndex(sIdx)
    }
    bvTree.suptree[0] |= (1 << 63)
}

func (bvTree *BvTree) Remove(n uint64) {
    // Rmove from the bitvector
    idx, off := offsets(n)
    b := ^uint64(1 << (63 - off))
    bvTree.bitvector[idx] &= b

    cIdx := bvTree.supIndex(n)
    idx, off = offsets(cIdx)
    b = ^uint64(1 << (63 - off))

    rIdx, rOff := offsets((n / 2) * 2)
    lIdx, lOff := offsets((n / 2) * 2 + 1)

    rVal := bvTree.bitvector[rIdx] & uint64(1 << (63 - rOff))
    lVal := bvTree.bitvector[lIdx] & uint64(1 << (63 - lOff))
    if (rVal == 0 && lVal == 0) {
        bvTree.suptree[idx] &= b
    }

    for cIdx > 0 {

        // bit to clear.
        b = ^uint64(1 << (63 - off))

        // Values of left and right children.
        rIdx, rOff = offsets(rightIndex(cIdx))
        lIdx, lOff = offsets(leftIndex(cIdx))

        rVal := bvTree.bitvector[rIdx] & uint64(1 << (63 - rOff))
        lVal := bvTree.bitvector[lIdx] & uint64(1 << (63 - lOff))
        if (rVal == 0 && lVal == 0) {
            bvTree.suptree[idx] &= b
        }


        cIdx = parentIndex(cIdx)
    }

    // Clear out the root node if theres no data left.
    fb := bvTree.suptree[0]
    if fb == (uint64(1 << 63))  {
        bvTree.suptree[0] = 0
    }
}

func BuildBvTree(numBits uint64) BvTree {
    result := BvTree{}

    // The number of uints we need is size / 64
    numUints := numBits / 64

    // if size is too small, just allocate one.
    if numUints == 0 {
        numUints = 1
    }

    fmt.Printf("building a tree with %d uint64s.\n", numUints)
    result.suptree = make([]uint64, numUints)
    result.bitvector = make([]uint64, numUints)
    result.numBits = numBits
    return result
}

func (bvTree *BvTree) supIndex(n uint64) uint64 {
    return (bvTree.numBits / 2 - 1) + (n / 2)
}

func parentIndex(n uint64) uint64 {
    return (n - 1) / 2
}

func leftIndex(n uint64) uint64 {
    return (n * 1) + 1
}

func rightIndex(n uint64) uint64 {
    return (n * 1) + 2
}

func offsets(n uint64) (uint64, uint64) {
    return n / 64, n % 64
}

func (bvTree *BvTree) DbgPrint() {
    fmt.Println("DbgPrint: ")
    fmt.Println("suptree")
    for _, val := range(bvTree.suptree) {
        fmt.Printf("%b\n", val)
    }
    fmt.Println("\nbitvector")
    for _, val := range(bvTree.bitvector) {
        fmt.Printf("%b\n", val)
    }
    fmt.Println(" ")
}
