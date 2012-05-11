package darts

import(
    "sort"
    "fmt"
    )

type dawgNode struct{
    parents, children map[rune]*dawgNode
    lastChar rune
    acceptable bool
    index int
    freq int
}

func addSubTree(node *dawgNode, key []rune){
    current := node
    for _, char := range key{
	current.lastChar = char
	current.children = make(map[rune]*dawgNode)
	newNode := new(dawgNode)
	current.children[char] = newNode
	newNode.parents = make(map[rune]*dawgNode)
	if current.acceptable {
	    newNode.parents[-char] = current
	}else{
	    newNode.parents[char] = current
	}
	current = newNode
    }
    current.acceptable = true
}
func merge(current, start, end *dawgNode){
    // go to the tail of the not-yet-merged path
    c := current
    for c.children != nil{
	c = c.children[c.lastChar]
    }
    if c == end {
	return
    }

    r := end
    var char rune
    var pc *dawgNode
    // each node on the not-yet-merged path has only one parent
    for char, pc = range c.parents{
    }
    pr, found := r.parents[char]
    for found == true{
	c, r = pc, pr
	for char, pc = range c.parents{
	}
	pr, found = r.parents[char]
    }
    pc.children[char] = r
    r.parents[char] = pc
}
func BuildDAWG(keys [][]rune /*Key_type*/, freq []int) *dawgNode{
    start := new(dawgNode)
    /*
    end := new(dawgNode)
    end.acceptable = true
    end.parents = make(map[rune]*dawgNode)
    */
    var end *dawgNode
    first := true
    //end.index = -1

f0: for _, key := range keys{
	current := start
	for j, alphabet := range key{
	    if current.children == nil{
		// if we are here, means key is a super string of the previous one
		// like previous: abcd, key: abcdef
		addSubTree(current, key[j:])
		continue f0
	    }
	    if alphabet > current.lastChar{
		if first {
		    first = false
		    end := current
		    for end.children != nil{
			end = end.children[end.lastChar]
		    }
		}
		// order is important, merge() must be called before addSubTree()
		merge(current, start, end)
		addSubTree(current, key[j:])
		continue f0
	    }
	    current = current.children[current.lastChar]
	}
    }
    if first {
	first = false
	end = start
	for end.children != nil{
	    end = end.children[end.lastChar]
	}
    }
    merge(start, start, end)
    return start
}

func printDAWG(d *dawgNode){
    fmt.Printf("This: %p, acceptable:%t\n", d, d.acceptable)
    for i, p := range d.parents{
	fmt.Printf("Parent %d: %p\n", i, p)
    }
    for i, p := range d.children{
	fmt.Printf("Child %d: %p\n", i, p)
    }
    for _, p := range d.children{
	printDAWG(p)
    }
}
func BuildFromDAWG(keys [][]rune /*Key_type*/, freq []int) Darts {
    dawg := BuildDAWG(keys, freq)

    //printDAWG(dawg)
    var d = new(dartsBuild)
    d.key = keys
    d.freq = freq
    d.resize(512)

    d.darts.Base[0] = 1
    d.nextCheckPos = 0

    siblings := d.fetchDAWG(dawg)
    d.insertDAWG(siblings)

    if d.err < 0 {
        panic("Build error")
    }
    return d.darts
}

type Pair struct {
  Char rune
  node *dawgNode
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Char < p[j].Char }

// A function to turn a map into a PairList, then sort and return it. 
func sortMapByValue(m map[rune]*dawgNode) PairList {
   p := make(PairList, len(m))
   i := 0
   for k, v := range m {
      p[i] = Pair{k+1, v}
      i++
   }
   sort.Sort(p)
   return p
}

func (d *dartsBuild) fetchDAWG(parent *dawgNode) PairList {
    if parent.acceptable {
	newNode := new(dawgNode)
	if nil == parent.children{
	    parent.children = make(map[rune]*dawgNode)
	}
	parent.children[-1] = newNode
    }
    return sortMapByValue(parent.children)
}

func (d *dartsBuild) insertDAWG(siblings PairList) int {
    if d.err < 0 {
        panic("insertDAWG error")
        return 0
    }

    var begin int = 0
    var pos int = max(int(siblings[0].Char)+1, d.nextCheckPos) - 1
    var nonZeroNum int = 0
    first := false
    if len(d.darts.Base) <= pos {
        d.resize(pos + 1)
    }

    for {
    next:
        pos++

        if len(d.darts.Base) <= pos {
            d.resize(pos + 1)
        }

        if d.darts.Check[pos] > 0 {
            nonZeroNum++
            continue
        } else if !first {
            d.nextCheckPos = pos
            first = true
        }

        begin = pos - int(siblings[0].Char)
        if len(d.darts.Base) <= (begin + int(siblings[len(siblings)-1].Char)) {
            d.resize(begin + int(siblings[len(siblings)-1].Char) + 400)
        }

        if d.used[begin] {
            continue
        }

        for i := 1; i < len(siblings); i++ {
            if begin+int(siblings[i].Char) >= len(d.darts.Base) {
                fmt.Println(len(d.darts.Base), begin+int(siblings[i].Char), begin+int(siblings[len(siblings)-1].Char))
            }
            if 0 != d.darts.Check[begin+int(siblings[i].Char)] {
                goto next
            }
        }
        break
    }

    if float32(nonZeroNum)/float32(pos-d.nextCheckPos+1) >= 0.95 {
        d.nextCheckPos = pos
    }
    d.used[begin] = true
    d.size = max(d.size, begin+int(siblings[len(siblings)-1].Char)+1)

    for i := 0; i < len(siblings); i++ {
        d.darts.Check[begin+int(siblings[i].Char)] = begin
    }

    for i := 0; i < len(siblings); i++ {
	if siblings[i].node.index > 0 { // siblings[i] is is already visited
            d.darts.Base[begin+int(siblings[i].Char)] = siblings[i].node.index
	}else{ // siblings[i] is a new node
	    newSiblings := d.fetchDAWG(siblings[i].node)
	    if len(newSiblings) == 0 {
		var value Value
		value.Freq = 1
		d.darts.Base[begin+int(siblings[i].Char)] = -len(d.darts.ValuePool) - 1
		d.darts.ValuePool = append(d.darts.ValuePool, value)
	    } else {
		h := d.insertDAWG(newSiblings)
		siblings[i].node.index = h
		d.darts.Base[begin+int(siblings[i].Char)] = h
	    }
	}
    }

    return begin
}
