govEB
=====

An attempted implementation of van Emde Boas trees in golang. The approach I'm taking in this repo is to follow the example from CLRS. First we start with a bit vector tree implementation, then the proto-vEB trees, then the full-blown vEB tree.


bvtree
===

A bit vector representation of the set, with a super imposed binary tree. Then another implementation with a superimposed tree of height 3 (the root, the summary bitvector, and the actual bitvector).


pvEBtree
===

An implementation of a proto-vEB tree.


vEBtree
===

An implementation of the full vEB tree.

