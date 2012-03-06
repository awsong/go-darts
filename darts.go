package darts

import (
//    "os"
    )

type key_type byte

type node struct{
    code key_type
    depth, left, right int
}

type unit struct{
    base int
    check int
}

type Darts []unit

type dartsBuild struct{
    array	Darts
    used	[]bool
    size	int
    keySize	int
    key		[][]key_type
    value	[]int
    nextCheckPos  int
    err		int
}

func Build(key [][]key_type, value []int) Darts{
    var d = new(dartsBuild)

    d.key = key
    d.value = value
    d.resize(20)

    d.array[0].base = 1
    d.nextCheckPos = 0

    var rootNode node
    rootNode.depth = 0
    rootNode.left = 0
    rootNode.right = len(key)

    siblings := d.fetch(rootNode)
    d.insert(siblings)

    d.size += (1 << 8 * d.keySize) + 1
    if d.size > len(d.array) {
	d.resize(d.size)
    }

    return d.array
}

func (d *dartsBuild) resize(newSize int) {
    if newSize > cap (d.array) {
	d.array = make(Darts, newSize)
    }else{
	d.array = d.array[:newSize]
    }
}

func (d *dartsBuild) fetch(parent node) []node{
    var siblings = make([]node, 0, 2)
    if d.err < 0 {
	return siblings[0:0]
    }
    var prev key_type = 0

    for i := parent.left; i < parent.right; i++ {
	if len(d.key[i]) < parent.depth {
	    continue
	}

	tmp := d.key[i]

	var cur key_type = 0
	if len(d.key[i]) != parent.depth {
	    cur = tmp[i] + 1
	}

	if prev > cur {
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

	if d.array[pos].check > 0 {
	    nonZeroNum++
	    continue
	}else if !first {
	    d.nextCheckPos = pos
	    first = true
	}

	begin = pos - int(siblings[0].code)
	if len(d.array) <= (begin + int(siblings[len(siblings) - 1].code)){
	    d.resize(len(d.array) * 105/100) //overflow
	}

	if d.used[begin] {
	    continue
	}

	for i := 1; i < len(siblings); i++ {
	    if 0 != d.array[begin + int(siblings[i].code)].check {
		goto next
	    }
	}
	break
    }

    if 1.0 * float32(nonZeroNum)/float32(pos - d.nextCheckPos + 1) >= 0.95{
	d.nextCheckPos = pos
    }
    d.used[begin] = true
    d.size = max(d.size, begin + int(siblings[len(siblings) - 1].code) + 1)

    for i := 0; i < len(siblings); i++ {
	d.array[begin + int(siblings[i].code)].check = begin
    }

    for i := 0; i < len(siblings); i++ {
	newSiblings := d.fetch(siblings[i])
	if len(newSiblings) == 0{
	    d.array[begin + int(siblings[i].code)].base = -d.value[siblings[i].left] - 1
	    if -d.value[siblings[i].left]-1 >= 0 {
		d.err = -2
		return 0
	    }
	}else{
	    h := d.insert(newSiblings)
	    d.array[begin + int(siblings[i].code)].base = h
	}
    }

    return begin
}
