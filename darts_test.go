package darts

import(
    "testing"
    )


func TestExactMatchSearch(t *testing.T) {
    r, _ := Import("darts.txt", "darts.lib")
    if r != true {
	t.Errorf("Test fail")
    }
}
