package darts

import(
    "testing"
    )


func TestExactMatchSearch(t *testing.T) {
    _, err := Import("darts.txt", "darts.lib")
    if err != nil {
	t.Errorf("Test fail")
    }
}
