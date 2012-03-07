package darts

import(
    "testing"
    )

type buildTest struct{
    inKey [][]key_type
    inValue []int
    out Darts
}

var buildTests = []buildTest{
    buildTest{
	[][]key_type{[]key_type("不学无术"),[]key_type("真功")},
	[]int{1,1},
	Darts{unit{3,5}, unit{4,6}}},
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
    in []key_type
    out bool
}

var exactMatchSearchTests = []exactMatchSearchTest{
    exactMatchSearchTest{
	[]key_type("不学无术"),
	true},
    exactMatchSearchTest{
	[]key_type("真功"),
	true},
    exactMatchSearchTest{
	[]key_type("历史"),
	false},
    exactMatchSearchTest{
	[]key_type("Adidas"),
	false},
}

func TestExactMatchSearch(t *testing.T) {
    d := Build([][]key_type{[]key_type("不学无术"),[]key_type("真功")}, []int{1,1})
    for _, dt := range exactMatchSearchTests {
	v := d.ExactMatchSearch(dt.in, 0)
	if v != dt.out {
	    t.Errorf("Key = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}
