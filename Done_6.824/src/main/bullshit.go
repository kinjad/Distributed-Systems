package main

import "fmt"
import "time"

var c chan int

func ready(w string, sec int) {
     time.Sleep(time.Duration(sec) * time.Second)
     fmt.Println(w, " is ready!")
     c <- 1
}

func main() {
     c = make(chan int)
     keys := []string{"Tea", "Coffe", "Diet Coke", "Pepsi", "Coco"}
     for j := 0; j < 5; j++ {
          go ready(keys[j], 5 - j)
     }
     fmt.Println("I am waiting, but not too long!")
     i := 0
     L: for {
		select {
			case <- c: 
			     	    i++
				    if i > 4 {
						break L
				    }
		       }
            }
}