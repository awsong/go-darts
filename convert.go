// +build ignore

package main

import (
    "fmt"
    "darts"
    "os"
    "bufio"
    "sort"
    "strconv"
    "strings"
    "time"
    "encoding/gob"
    )
type runeKey struct{
    key []rune
    value int
}
type runeKeySlice []runeKey

func (r runeKeySlice) Len() int{
    return len(r)
}
func (r runeKeySlice) Less(i, j int) bool{
    var l int
    if len(r[i].key) < len(r[j].key) {
	l = len(r[i].key)
    }else{
	l = len(r[j].key)
    }

    for m := 0; m < l; m++{
	if r[i].key[m] < r[j].key[m] {
	    return true
	}else if r[i].key[m] == r[j].key[m] {
	    continue
	}else {
	    return false
	}
    }
    if len(r[i].key) < len(r[j].key){
	return true
    }
    return false
}
func (r runeKeySlice) Swap(i, j int){
    r[i], r[j] = r[j], r[i]
}

func main(){
    unifile, _ := os.Open("darts.txt");
    defer unifile.Close();
    ofile, _ := os.Create("darts.lib");
    defer ofile.Close();

    runeKeys := make(runeKeySlice, 0, 130000)
    uniLineReader := bufio.NewReaderSize(unifile, 400);
    line, _, bufErr := uniLineReader.ReadLine();
    for nil == bufErr {
	rst := strings.Split(string(line), "\t")
	key := []rune(rst[0])
	value,_ := strconv.Atoi(rst[1])
	runeKeys = append(runeKeys, runeKey{key, value})
	line, _, bufErr = uniLineReader.ReadLine();
    }
    sort.Sort(runeKeys)

    keys := make([][]rune, len(runeKeys))
    values := make([]int, len(runeKeys))

    for i := 0; i < len(runeKeys); i++ {
	keys[i] = runeKeys[i].key
	values[i] = runeKeys[i].value
    }
    fmt.Printf("input dict length: %v %v\n", len(keys), len(values));
    round := len(keys)
    d := darts.Build(keys[:round], values[:round])
    d.UpdateThesaurus(keys[:round])
    fmt.Printf("build out length %v\n", len(d.Base))
    t := time.Now()
    for i := 0; i < round; i++ {
	if true != d.ExactMatchSearch(keys[i],0){
	    fmt.Println("wrong", string(keys[i]), i)
	}
    }
    fmt.Println(time.Since(t))
    enc := gob.NewEncoder(ofile);
    enc.Encode(d);

    r := d.CommonPrefixSearch([]rune("考察队员"), 0)
    m := r[len(r) -1]
    fmt.Println(m)
}
