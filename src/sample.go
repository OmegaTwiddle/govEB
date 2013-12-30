package main

import (
    "fmt"
    "./bvtree"
    "math/rand"
    "time"
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    numBits := uint64(14336)
    numToInsert := 20
    for i := 0; i < 10; i++ {
        bvTree := bvtree.BuildBvFhTree(numBits)
        vals := make(map[uint64] bool, numToInsert)
        myMin := uint64(numBits)
        myMax := uint64(0)
        for j := 0; j < numToInsert; j++ {
            n := uint64(rand.Int63n(int64(numBits)))
            fmt.Println(n)
            vals[n] = true
            bvTree.Insert(n)
            if n < myMin {
                myMin = n
            }
            if n > myMax {
                myMax = n
            }
        }
        checkTree(bvTree, myMin, myMax, vals, []uint64{})

    }
}


func buildBvFhTree(vals []uint64, ghosts []uint64, size uint64) *bvtree.BvFhTree {

    bvTree := bvtree.BuildBvFhTree(size)

    bvTree.DbgPrint()
    for _, val := range(vals) {
        bvTree.Insert(val)
    }

    for _, ghost := range(ghosts) {
        bvTree.Insert(ghost)
        bvTree.Remove(ghost)
    }
    return bvTree
}
func buildBvTree(vals []uint64, ghosts []uint64, size uint64) *bvtree.BvTree {

    bvTree := bvtree.BuildBvTree(size)
    for _, val := range(vals) {
        bvTree.Insert(val)
    }

    for _, ghost := range(ghosts) {
        bvTree.Insert(ghost)
        bvTree.Remove(ghost)
    }
    return bvTree
}
func bvmain() {
    fmt.Println("Quick sample of the vEB tree")
    vals := []uint64{418,2208,8086,751,2770,10610,8021,9497,4221,3506,5223,2424,13567,8030,10316,1602,11062,1094,12052,2852}
    ghosts := []uint64{10,20,30,40,50,60,70,80,90,100}
    bvTree := buildBvTree(vals, ghosts, 14336)
    //bvTree.DbgPrint()

    mapVals := make(map[uint64] bool, 20)
    for _, val := range(vals) {
        mapVals[val] = true
    }

    checkTree(bvTree, 418, 13567, mapVals, ghosts)

}

func bvfhmain() {

    fmt.Println("Quick sample of the vEB tree")
    vals := []uint64{3,4,5,12,13,14,33,34,35,54,55,56,60,61,62}
    ghosts := []uint64{}
    bvTree := buildBvFhTree(vals, ghosts, 64)
    bvTree.DbgPrint()

    mapVals := make(map[uint64] bool, 20)
    for _, val := range(vals) {
        mapVals[val] = true
    }

    checkTree(bvTree, 3, 62, mapVals, ghosts)

}

func checkTree(bvTree bvtree.DynamicSet, myMin uint64, myMax uint64, vals map[uint64] bool, ghosts []uint64) {
        fmt.Println(bvTree)
        fmt.Printf("min/max were %d/%d\n ", bvTree.Min(), bvTree.Max())
        //bvTree.DbgPrint()

        for _, ghost := range(ghosts) {
            if bvTree.Contains(ghost) {
                panic("Uh oh, ghost in the system.")
            }
        }

        if myMin != bvTree.Min() {
            panic(fmt.Sprintf("myMin: %d, bvTre.Min(): %d\n", myMin, bvTree.Min()))
        }

        if myMax != bvTree.Max() {
            panic("Incorrect max calculated.")
        }

        cur := myMin
        for cur < myMax {
            fmt.Printf("Found a val: %d\n", cur)
            if !bvTree.Contains(cur) {
                panic ("doesn't contain a certain successor!")
            }

            if !vals[cur] {
                panic("had a value i didn't put in!")
            }

            cur = bvTree.Successor(cur)
        }

        if cur != myMax {
            panic ("Successor didn't return max...!")
        }

        for cur > myMin {
            fmt.Printf("Found a val: %d\n", cur)
            if !bvTree.Contains(cur) {
                panic("doesn't contain a certain predecessor!")
            }

            if !vals[cur] {
                panic("had a value i didn't put in!")
            }

            cur = bvTree.Predecessor(cur)
        }

        if cur != myMin {
            panic("Predecessor didn't return min...")
        }

        for val, _ := range(vals) {
            if !bvTree.Contains(val) {
                panic ("didn't have a value I put in!")
            }
        }
}

