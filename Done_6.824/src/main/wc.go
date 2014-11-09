package main

import "os"
import "fmt"
import "mapreduce"
import "strconv"
import "container/list"
import "strings"
import "unicode"       
/*
type KeyValue struct {
  Key string
  Value string
}
*/


// our simplified version of MapReduce does not supply a
// key to the Map function, as in the paper; only a value,
// which is a part of the input file contents
func Map(value string) *list.List {

     Splitted_str := strings.FieldsFunc(value, split)
     SS_len := len(Splitted_str)
     var the_list *list.List
     the_list = list.New()
     KeyMap := map[string]int{}
     for i := 0; i < SS_len; i ++ {
     	    if IsWord(Splitted_str[i]) {
	       KeyMap[Splitted_str[i]] += 1
            }
     }

     var keys []string
     for k := range KeyMap {
     	 keys = append(keys, k)
     }

     for _, k := range keys {
          element := mapreduce.KeyValue{k,strconv.Itoa(KeyMap[k])}
     	  the_list.PushBack(element)
     }
     return the_list
}

// iterate over list and add values
func Reduce(key string, values *list.List) string {
     var KeyCount int
     KeyCount = 0
     for e := values.Front(); e != nil; e = e.Next() {
     	 num_value, _ := strconv.Atoi(e.Value.(string))
     	 KeyCount += num_value
     }
     return strconv.Itoa(KeyCount)
}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master x.txt sequential)
// 2) Master (e.g., go run wc.go master x.txt localhost:7777)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
  if len(os.Args) != 4 {
    fmt.Printf("%s: see usage comments in file\n", os.Args[0])
  } else if os.Args[1] == "master" {
    if os.Args[3] == "sequential" {
      mapreduce.RunSingle(5, 3, os.Args[2], Map, Reduce)
    } else {
      mr := mapreduce.MakeMapReduce(5, 3, os.Args[2], os.Args[3])    
      // Wait until MR is done
      <- mr.DoneChannel
    }
  } else {
    mapreduce.RunWorker(os.Args[2], os.Args[3], Map, Reduce, 100)
  }
}

func split(s rune) bool {
     if !unicode.IsLetter(s) {
     	  return true
     }
     return false
}

func IsWord(w string) bool {
     s_len := len(w)
     for i := 0; i < s_len ; i ++ {
     	 if !unicode.IsLetter(rune(w[i])) {
	    return false
         }
     }
     return true
}