package main

import (
    "fmt"
    "./bvtree"
)

func main() {
    fmt.Println("Quick sample of the vEB tree")
    bvTree := bvtree.BuildBvTree(64)
    bvTree.Insert(0)
    bvTree.Insert(1)
    bvTree.Insert(2)
    bvTree.Insert(63)
    //bvTree.Remove(2)
    fmt.Println("Alright.", bvTree)
    bvTree.DbgPrint()

    vals := []uint64{0,1,2,3,4,5}
    for _, val := range vals {
        if bvTree.Contains(val) {
            fmt.Printf("it had %d!\n", val)
        } else {
            fmt.Printf("nope, didn't have %d...\n", val)
        }
    }

}
