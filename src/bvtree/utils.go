package bvtree

import (
    "fmt"
)

func getRoot(n uint64) uint64 {
    result := uint64(2)
    for result * result < n {
        result = result * 2
    }
    return result
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
