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
    bvTree.Insert(33); bvTree.Remove(33)
    bvTree.Insert(34); bvTree.Remove(34)
    bvTree.Insert(35); bvTree.Remove(35)
    bvTree.Insert(37)
    bvTree.Insert(38); bvTree.Remove(38)
    bvTree.Insert(63)
    bvTree.Remove(63)
    bvTree.Remove(0)
    bvTree.Remove(1)
    bvTree.Remove(2)
    //bvTree.Remove(2)
    fmt.Println("Alright.", bvTree)
    bvTree.DbgPrint()

    vals := []uint64{0,1,2,3,4,5,32,33,34,35,36,37,38,39,62,63}
    for _, val := range vals {
        if bvTree.Contains(val) {
            fmt.Printf("it had %d!\n", val)
        } else {
            fmt.Printf("nope, didn't have %d...\n", val)
        }
    }

}
