package main

import (
    "fmt"
    "./bvtree"
)

func main() {
    fmt.Println("Quick sample of the vEB tree")
    bvTree := bvtree.BuildBvTree(1024)
    bvTree.Insert(12)
    //bvTree.Insert(11)
    fmt.Println("Alright.", bvTree)

    vals := []uint64{10,11,12,13,14}
    for _, val := range vals {
        if bvTree.Contains(val) {
            fmt.Printf("it had %d!\n", val)
        } else {
            fmt.Printf("nope, didn't have %d...\n", val)
        }
    }
}
