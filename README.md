Darts
=====
This is a GO implementation of Double-ARray Trie System. It's a clone of the [C++ version](http://chasen.org/~taku/software/darts/) 

Darts can be used as simple hash dictionary. You can also do very fast Common Prefix Search which is essential for morphological analysis, such as word split for CJK text indexing/searching.

Reference
---------
[What is Trie](http://en.wikipedia.org/wiki/Trie)   
[An Implementation of Double-Array Trie](http://linux.thai.net/~thep/datrie/datrie.html)

NEWS
----------
* Support building Double-Array from [DAWG有向无环图](https://en.wikipedia.org/wiki/Deterministic_acyclic_finite_state_automaton), reduce the on-disk dict half as Trie. Lookup performance increases 25%.

TO DO list
----------
* Documentation/comments
* Benchmark

Switch from unicode to byte version
----------------------
```sh
gofmt -tabs=false -tabwidth=4 -r='rune /*Key_type*/ -> byte /*Key_type*/' -w darts.go
gofmt -tabs=false -tabwidth=4 -r='rune /*Key_type*/ -> byte /*Key_type*/' -w dawg.go
```

Usage
---------
# Input Dictionary Format
```sh
Key\tFreq
```
Each key occupies one line. The file should be utf-8 encoded

# Code example (unicode version)
```go
package main

import (
    "darts"
    "fmt"
)

func main() {
    d, err:= darts.Import("darts.txt", "darts.lib", true)
    if err == nil {
        if d.ExactMatchSearch([]rune("考察队员", 0)) {
            fmt.Println("考察队员 is in dictionary")
        }
    }
}
```
# Code example (byte version)
```go
package main

import (
    "darts"
    "fmt"
)

func main() {
    d, err := darts.Import("darts.txt", "darts.lib", true)
    if err == nil {
        key := []byte("考察队员")
        r := d.CommonPrefixSearch(key, 0)
        for i := 0; i < len(r); i++ {
            fmt.Println(string(key[:r[i].PrefixLen]))
        }
    }
}
```

Performance
-----------
Using a 100K item dictionary, a simple search on eath key takes go map 46 ms, takes byte\_key version of darts 14 ms, and for unicode\_key version of darts 9.5 ms.

LICENSE
-----------
Apache License 2.0
