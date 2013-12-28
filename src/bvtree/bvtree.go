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



/**
 * Returns the number below n in the tree.
 * Assumes that the number passed in is greater than
 * the min value.
 */
func (bvTree *BvTree) Predecessor(n uint64) uint64 {
    treePos := bvTree.supIndex(n)
    goingUp := true

    //Now, the normal case, we have to keep searching up the tree.
    for treePos <= bvTree.maxLlIndex() {
        //fmt.Printf("treePos %d\n", treePos)
        // If we're at the lowest level, check whether the right child has a bit
        // and make sure that we're not just returning n.
        if bvTree.inLowestLevel(treePos) {
        lPos, rPos := bvTree.bvIndices(treePos)
            if rPos < n && bvTree.hasBvBit(rPos) {
                return rPos
            } else if lPos < n && bvTree.hasBvBit(lPos) {
                return lPos
            }
        }

        // If we didn't find the successor, we need to traverse the tree.

        if goingUp {
            nextLeftPos := leftIndex(parentIndex(treePos))

            if nextLeftPos != treePos  && bvTree.hasStBit(nextLeftPos) {
                treePos = nextLeftPos
                goingUp = false
            } else {
                treePos = parentIndex(treePos)
            }
        } else {
            // When going down the tree, just look right, if nothing, go left
            nextRightPos, nextLeftPos := bvTree.childrenIndices(treePos)
            if bvTree.hasStBit(nextLeftPos) {
                treePos = nextLeftPos
            } else {
                treePos = nextRightPos
            }
        }
    }

    panic("There was a problem with predecessor.")
}



/**
 * Returns the number above n in the tree.
 * Assumes that the number passed in is less than
 * the max value.
 */
func (bvTree *BvTree) Successor(n uint64) uint64 {
    treePos := bvTree.supIndex(n)
    goingUp := true

    //Now, the normal case, we have to keep searching up the tree.
    for treePos <= bvTree.maxLlIndex() {
        //fmt.Printf("treePos %d\n", treePos)
        // If we're at the lowest level, check whether the right child has a bit
        // and make sure that we're not just returning n.
        if bvTree.inLowestLevel(treePos) {
        lPos, rPos := bvTree.bvIndices(treePos)
            if lPos > n && bvTree.hasBvBit(lPos) {
                return lPos
            } else if rPos > n && bvTree.hasBvBit(rPos) {
                return rPos
            }
        }

        // If we didn't find the successor, we need to traverse the tree.

        if goingUp {
            nextRightPos := rightIndex(parentIndex(treePos))

            if nextRightPos != treePos  && bvTree.hasStBit(nextRightPos) {
                treePos = nextRightPos
                goingUp = false
            } else {
                treePos = parentIndex(treePos)
            }
        } else {
            // When going down the tree, just look right, if nothing, go left
            nextRightPos, nextLeftPos := bvTree.childrenIndices(treePos)
            if bvTree.hasStBit(nextRightPos) {
                treePos = nextRightPos
            } else {
                treePos = nextLeftPos
            }
        }
    }

    panic("There was a problem with successor.")
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

func getNumUints(numBits uint64) uint64 {
    result := uint64(2)
    for result < numBits {
        result *= 2
    }
    if result < uint64(64) {
        return 1
    }
    return (result / uint64(64))
}

func BuildBvTree(numBits uint64) BvTree {
    result := BvTree{}

    // The number of uints we need is size / 64
    numUints := getNumUints(numBits)

    // if size is too small, just allocate one.
    if numUints == 0 {
        numUints = 1
    }

    fmt.Printf("building a tree with %d uint64s.\n", numUints)
    result.suptree = make([]uint64, numUints)
    result.bitvector = make([]uint64, numUints)
    result.numBits = numUints * uint64(64)
    return result
}

/**
 * Returns true if the position is in the lowest level
 * of the tree (so that it's children will be in the bitvector)
 */
func (bvTree *BvTree) inLowestLevel(pos uint64) bool {
    return (pos > bvTree.llIndex() && pos <= bvTree.maxLlIndex())
}

// Return true if the supporting tree has the bit.
func (bvTree *BvTree) hasStBit(pos uint64) bool {
    idx, off := offsets(pos)
    return (bvTree.suptree[idx] & uint64(1 << (63 - off))) != 0
}

// Return true if the bitvector has the bit.
func (bvTree *BvTree) hasBvBit(pos uint64) bool {
    idx, off := offsets(pos)
    return (bvTree.bitvector[idx] & uint64(1 << (63 - off))) != 0
}
// Assuming the bitvector is of size 2^n, get the first index of the
// last level in the supporting tree.
// E.G. 
func (bvTree *BvTree) llIndex() uint64 {
    return bvTree.numBits / 2 - 1
}

/**
 * Returns the last index in the lower level of the supporting tree.
 */
func (bvTree *BvTree) maxLlIndex() uint64 {
    return bvTree.numBits - 2
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

/**
 * Assuming that n is inside the support tree, returns the indices
 * of its left and right children
 */
func (bvTree *BvTree) childrenIndices(n uint64) (uint64, uint64) {
    return leftIndex(n), rightIndex(n)
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

func (bvTree *BvTree) checkBit(n uint64) {
    if bvTree.hasStBit(n) {
        fmt.Printf("Has St bit:          %d\n", n)
    } else {
        fmt.Printf(" Didn't have St bit: %d\n", n)
    }
}

func (bvTree *BvTree) DbgPrint() {
    fmt.Println("DbgPrint: ")
    fmt.Println("suptree")
    //for _, val := range(bvTree.suptree) {
    //    dbgPrintBin(val)
    //}
    fmt.Println("\nbitvector")
    //for _, val := range(bvTree.bitvector) {
    //    dbgPrintBin(val)
    //}
    fmt.Println(" ")
}

func dbgPrintBin(n uint64) {
    for i := uint64(0); i < 64; i++ {
        b := uint64(1 << (63 - i))
        if (n & b) == 0 {
            fmt.Printf("0")
        } else {
            fmt.Printf("1")
        }
    }
    fmt.Printf("\n")

}
