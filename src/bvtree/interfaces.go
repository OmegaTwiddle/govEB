package bvtree

type DynamicSet interface {

    Contains(n uint64) bool
    Insert(n uint64)
    Remove(n uint64)

    Predecessor(n uint64) uint64
    Successor(n uint64) uint64

    Min() uint64
    Max() uint64

    // TODO:: Remove this...?
    DbgPrint()
}
