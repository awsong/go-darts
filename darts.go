package darts

import (
    "fmt"
    )

//type Key_type rune

type node struct{
    code rune /*Key_type*/
    depth, left, right int
}

type Unit struct{
    Base int
    Check int
}

type Darts []Unit

type dartsBuild struct{
    array	Darts
    used	[]bool
    size	int
    keySize	int
    key		[][]rune /*Key_type*/
    value	[]int
    nextCheckPos  int
    err		int
}

func Build(key [][]rune /*Key_type*/, value []int) Darts{
    var d = new(dartsBuild)

    d.key = key
    d.value = value
    d.resize(512)

    d.array[0].Base = 1
    d.nextCheckPos = 0

    var rootNode node
    rootNode.depth = 0
    rootNode.left = 0
    rootNode.right = len(key)

    siblings := d.fetch(rootNode)
    d.insert(siblings)

/*
    d.size += (1 << 8 * d.keySize) + 1
    if d.size > len(d.array) {
	d.resize(d.size)
    }
    */

    if d.err < 0 {
	panic("Build error")
    }
    return d.array
}

func (d *dartsBuild) resize(newSize int) {
    if newSize > cap(d.array) {
	d.array = append(d.array, make(Darts, (newSize - len(d.array)))...)
	d.used = append(d.used, make([]bool, (newSize - len(d.used)))...)
    }else{
	d.array = d.array[:newSize]
	d.used = d.used[:newSize]
    }
}

func (d *dartsBuild) fetch(parent node) []node{
    var siblings = make([]node, 0, 2)
    if d.err < 0 {
	return siblings[0:0]
    }
    var prev rune /*Key_type*/ = 0

    for i := parent.left; i < parent.right; i++ {
	if len(d.key[i]) < parent.depth {
	    continue
	}

	tmp := d.key[i]

	var cur rune /*Key_type*/ = 0
	if len(d.key[i]) != parent.depth {
	    cur = tmp[parent.depth] + 1
	}

	if prev > cur {
	    fmt.Println(prev, cur,i, parent.depth, d.key[i])
	    fmt.Println(d.key[i])
	    panic("fetch error 1")
	    d.err = -3
	    return siblings[0:0]
	}

	if cur != prev || len(siblings) == 0 {
	    var tmpNode node
	    tmpNode.depth = parent.depth + 1
	    tmpNode.code = cur
	    tmpNode.left = i
	    if len(siblings) != 0 {
		siblings[len(siblings) - 1].right = i
	    }

	    siblings = append(siblings, tmpNode)
	}

	prev = cur
    }

    if len(siblings) != 0 {
	siblings[len(siblings) - 1].right = parent.right
    }

    return siblings
}

func max(a, b int) int{
    if a > b {
	return a
    }
    return b
}
func (d *dartsBuild) insert(siblings []node) int{
    if d.err < 0 {
	panic("insert error")
	return 0
    }

    var begin int = 0
    var pos int = max(int(siblings[0].code) + 1, d.nextCheckPos) - 1
    var nonZeroNum int = 0
    first := false
    if len(d.array) <= pos {
	d.resize(pos + 1)
    }

    for {
next:
	pos++

	if len(d.array) <= pos {
	    d.resize(pos + 1)
	}

	if d.array[pos].Check > 0 {
	    nonZeroNum++
	    continue
	}else if !first {
	    d.nextCheckPos = pos
	    first = true
	}

	begin = pos - int(siblings[0].code)
	if len(d.array) <= (begin + int(siblings[len(siblings) - 1].code)){
	    d.resize(begin + int(siblings[len(siblings) - 1].code) + 400)
	}

	if d.used[begin] {
	    continue
	}

	for i := 1; i < len(siblings); i++ {
	    if begin + int(siblings[i].code) >= len(d.array){
		fmt.Println(len(d.array), begin + int(siblings[i].code), begin + int(siblings[len(siblings) - 1].code))
	    }
	    if 0 != d.array[begin + int(siblings[i].code)].Check {
		goto next
	    }
	}
	break
    }

    if float32(nonZeroNum)/float32(pos - d.nextCheckPos + 1) >= 0.95{
	d.nextCheckPos = pos
    }
    d.used[begin] = true
    d.size = max(d.size, begin + int(siblings[len(siblings) - 1].code) + 1)

    for i := 0; i < len(siblings); i++ {
	d.array[begin + int(siblings[i].code)].Check = begin
    }

    for i := 0; i < len(siblings); i++ {
	newSiblings := d.fetch(siblings[i])
	if len(newSiblings) == 0{
	    d.array[begin + int(siblings[i].code)].Base = -d.value[siblings[i].left] - 1
	    if -d.value[siblings[i].left]-1 >= 0 {
		d.err = -2
		panic("insert error 1")
		return 0
	    }
	}else{
	    h := d.insert(newSiblings)
	    d.array[begin + int(siblings[i].code)].Base = h
	}
    }

    return begin
}
func (d Darts) ExactMatchSearch(key []rune /*Key_type*/, nodePos int) bool{
    b := d[nodePos].Base
    var p int

    for i := 0; i < len(key); i++ {
	p = b + int(key[i]) + 1
	if b == d[p].Check {
	    b = d[p].Base
	}else{
	    return false
	}
    }

    p = b
    n := d[p].Base
    if b == d[p].Check && n < 0 {
	return true
    }

    return false
}
