package main

import "fmt"
import "container/list"
import "strings"
import "unicode"
import "strconv"

type WordPair struct {
     Word string
     Count int
}

func Crap(value string) *list.List {
     x := list.New()
     sth := WordPair{"SOB", 19}
     x.PushBack(sth)
     return x
}

func main(){
     var y *list.List
     y = Crap("sb")
     e := y.Front()
     w := e.Value.(WordPair)
     fmt.Printf("%v\n", w.Word)
     cot, _ := strconv.Atoi("123")
     fmt.Printf("%v\n", cot)
     
     St := "For the Lich King! You stupid ass! How do you know that? Let's ride! For the Lich King! 00::11 Entry 1"
     Af_st := strings.FieldsFunc(St, split)
     
     SbMap := map[string]int{}
     
     for i := 0; i < len(Af_st); i ++ {	
     	 if IsWord(Af_st[i]) {
	    SbMap[Af_st[i]] += 1
	 }	 
     }
     fmt.Printf("Crap\n")	
//     fmt.Printf("%v\n", Af_st[0])
//     fmt.Printf("%v\n", y.Len())
     var keys []string
     for k := range SbMap {
     	 keys = append(keys, k)
     }
     for _, k := range keys {
     	 fmt.Printf("%v ", k)
       	 fmt.Printf("%v\n", SbMap[k])
     }
}

func split(s rune) bool {
     if s == ' ' {
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