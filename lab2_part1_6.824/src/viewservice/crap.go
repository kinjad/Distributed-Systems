package viewservice

import "fmt"
import "time"


func main() {
     dick := make(map[string]time.Time)
     dick["sf"] = time.Now()
     fmt.Printf("%v\n", dick["sf"])
}