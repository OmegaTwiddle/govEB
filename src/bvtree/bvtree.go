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

func (bvTree *BvTree) zeroRoot() bool {
    fb := bvTree.suptree[0] & uint64(1 << 63)
    return fb == 0
}

// Assuming the bitvector is of size 2^n, get the first index of the
// last level in the supporting tree.
// E.G. 
func (bvTree *BvTree) llIndex() uint64 {
    return bvTree.numBits / 2 - 1
}

func (bvTree *BvTree) Min() uint64 {
    cPos := uint64(0)
    if bvTree.zeroRoot() {
        return 0
    }

    for cPos < bvTree.llIndex() {
        rPos := rightIndex(cPos)
        lPos := leftIndex(cPos)
        lIdx, lOff := offsets(lPos)

        lVal := bvTree.suptree[lIdx] & uint64(1 << (63 - lOff))

        if lVal != 0 {
            cPos = lPos
        } else {
            cPos = rPos
        }

    }

    // Now that we're outside that loop, we need to 
    // reach into the bitvector.
    lPos, rPos := bvTree.bvIndices(cPos)
    lIdx, lOff := offsets(lPos)
    lVal := bvTree.bitvector[lIdx] & uint64(1 << (63 - lOff))

    if lVal != 0 {
        return lPos
    }

    return rPos
}

func (bvTree *BvTree) Max() uint64 {
    cPos := uint64(0)
    if bvTree.zeroRoot() {
        return 0
    }

    for cPos < bvTree.llIndex() {
        rPos := rightIndex(cPos)
        lPos := leftIndex(cPos)
        rIdx, rOff := offsets(rPos)

        rVal := bvTree.suptree[rIdx] & uint64(1 << (63 - rOff))

        if rVal != 0 {
            cPos = rPos
        } else {
            cPos = lPos
        }

    }

    // Now that we're outside that loop, we need to 
    // reach into the bitvector.
    lPos, rPos := bvTree.bvIndices(cPos)
    rIdx, rOff := offsets(rPos)
    rVal := bvTree.bitvector[rIdx] & uint64(1 << (63 - rOff))

    if rVal != 0 {
        return rPos
    }

    return lPos
}

func (bvTree *BvTree) hasBvBit(pos uint64) bool {
    idx, off := offsets(pos)
    return (bvTree.bitvector[idx] & uint64(1 << (63 - off))) != 0
}

/**
 * Returns the number above n in the tree.
 */
func (bvTree *BvTree) Successor(n uint64) uint64 {
    oldPos := n
    treePos := bvTree.supIndex(n)
    rPos, lPos := bvTree.bvIndices(treePos)
    if lPos > oldPos {
        if bvTree.hasBvBit(lPos) {
            return lPos
        } else {
            return rPos
        }
    }
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

    cIdx = parentIndex(cIdx)
    for cIdx > 0 {
        idx, off := offsets(cIdx)
        // bit to clear.
        b = ^uint64(1 << (63 - off))

        // Values of left and right children.
        rPos := rightIndex(cIdx)
        lPos := leftIndex(cIdx)
        rIdx, rOff = offsets(rPos)
        lIdx, lOff = offsets(lPos)

        rVal := bvTree.suptree[rIdx] & uint64(1 << (63 - rOff))
        lVal := bvTree.suptree[lIdx] & uint64(1 << (63 - lOff))
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

/**
 * Given an index at the lowest level of the support tree,
 * returns the left and right "children" indices inside
 * the bit vector.
 */
func (bvTree *BvTree) bvIndices(n uint64) (uint64, uint64) {
    k := (n - (bvTree.numBits / 2 - 1)) * 2
    return k, k + 1
}

func (bvTree *BvTree) supIndex(n uint64) uint64 {
    return (bvTree.numBits / 2 - 1) + (n / 2)
}

func parentIndex(n uint64) uint64 {
    return (n - 1) / 2
}

func leftIndex(n uint64) uint64 {
    return (n * 2) + 1
}

func rightIndex(n uint64) uint64 {
    return (n * 2) + 2
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
