package main

import (
    "fmt"
    "./bvtree"
    "math/rand"
    "time"
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    for i := 0; i < 10; i++ {
        bvTree := bvtree.BuildBvTree(64)
        vals := make(map[uint64] bool, 5)
        myMin := uint64(64)
        myMax := uint64(0)
        for j := 0; j < 5; j++ {
            n := uint64(rand.Int63n(64))
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
        checkTree(bvTree, myMin, myMax, vals)

    }
}

func oldmain() {
    fmt.Println("Quick sample of the vEB tree")
    bvTree := bvtree.BuildBvTree(64)
    bvTree.Insert(0)
    bvTree.Insert(1)
    bvTree.Insert(2)
    bvTree.Insert(33); bvTree.Remove(33)
    bvTree.Insert(34); bvTree.Remove(34)
    bvTree.Insert(35); bvTree.Remove(35)
    //bvTree.Insert(37)
    bvTree.Insert(38); bvTree.Remove(38)
    //bvTree.Insert(63)
    //bvTree.Remove(0)
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

    fmt.Println("The min of the tree was: ", bvTree.Min())

}

func checkTree(bvTree bvtree.BvTree, myMin uint64, myMax uint64, vals map[uint64] bool) {
        fmt.Println(bvTree)
        fmt.Printf("min/max were %d/%d\n ", bvTree.Min(), bvTree.Max())
        bvTree.DbgPrint()
        if myMin != bvTree.Min() {
            panic("Incorrect min calculated.")
        }

        if myMax != bvTree.Max() {
            panic("Incorrect max calculated.")
        }

        cur := myMin
        for cur < myMax {
            if !bvTree.Contains(cur) {
                panic ("doesn't contain a certain successor!")
            }

            if !vals[cur] {
                panic("had a value i didn't put in!")
            }
            fmt.Printf("successor was %d\n", bvTree.Successor(cur))
            cur = myMax//bvTree.Successor(cur)
        }

        for val, _ := range(vals) {
            if !bvTree.Contains(val) {
                panic ("didn't have a value I put in!")
            }
        }
}

