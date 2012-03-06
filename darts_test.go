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

func TestBuild(t *testing.T) {
    for _, dt := range buildTests {
	v := Build(dt.inKey, dt.inValue)
	if v != nil {
		t.Errorf("Key = %v, Real = %v, want %v.", dt.inKey, v, dt.out)
	}
    }
}
