package darts

import(
    "testing"
    "bufio"
    "strings"
    "strconv"
    "time"
    "os"
    )


func TestExactMatchSearch(t *testing.T) {
    _, err := Import("darts.txt", "darts.lib", false)
    if err != nil {
	t.Errorf("Test fail")
    }
}
func TestPerf(t *testing.T) {
    inFile := "d.txt"
    unifile, erri := os.Open(inFile)
    if erri != nil{
	t.Error(erri)
    }
    defer unifile.Close()
    d, _ := Load("darts.lib")

    dartsKeys := make(dartsKeySlice, 0, 130000)
    uniLineReader := bufio.NewReaderSize(unifile, 400)
    line, _, bufErr := uniLineReader.ReadLine()
    for nil == bufErr {
        rst := strings.Split(string(line), "\t")
        key := []rune /*Key_type*/(rst[0])
        value, _ := strconv.Atoi(rst[1])
        dartsKeys = append(dartsKeys, dartsKey{key, value})
        line, _, bufErr = uniLineReader.ReadLine()
    }

    keys := make([][]rune /*Key_type*/, len(dartsKeys))

    for i := 0; i < len(dartsKeys); i++ {
        keys[i] = dartsKeys[i].key
    }
    t.Logf("input dict length: %v\n", len(keys))
    tm := time.Now()
    for i := 0; i < len(keys); i++ {
        if true != d.ExactMatchSearch(keys[i], 0) {
            t.Errorf("missing key %s, %v, %d", string(keys[i]), keys[i], i)
	    return
        }
    }
    t.Log(time.Since(tm))
}
