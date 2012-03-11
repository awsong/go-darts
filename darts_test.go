package darts

import(
    "testing"
    )

type buildTest struct{
    inKey [][]byte
    inValue []int
    out Darts
}

var buildTests = []buildTest{
    buildTest{
	[][]byte{[]byte("不学无术"),[]byte("真功")},
	[]int{1,1},
	Darts{Unit{3,5}, Unit{4,6}}},
}

func testBuild(t *testing.T) {
    for _, dt := range buildTests {
	v := Build(dt.inKey, dt.inValue)
	if v != nil {
	    t.Errorf("Len: %v\n", len(v))
	    t.Errorf("Key = %v, Real = %v, want %v.", dt.inKey, v, dt.out)
	}
    }
}
type exactMatchSearchTest struct{
    in []byte
    out bool
}

var exactMatchSearchTests = []exactMatchSearchTest{
    exactMatchSearchTest{
	[]byte("不学无术"),
	true},
    exactMatchSearchTest{
	[]byte("真功"),
	true},
    exactMatchSearchTest{
	[]byte("历史"),
	false},
    exactMatchSearchTest{
	[]byte("Adidas"),
	true},
}

func TestExactMatchSearch(t *testing.T) {
    d := Build([][]byte{[]byte("Adidas"),[]byte("不学无术"),[]byte("真功")}, []int{1,1,1})
    for _, dt := range exactMatchSearchTests {
	v := d.ExactMatchSearch(dt.in, 0)
	if v != dt.out {
	    t.Errorf("Key = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}
