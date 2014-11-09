package mapreduce
import "container/list"
import "fmt"


type WorkerInfo struct {
  address string
  // You can add definitions here.
}


// Clean up all workers by sending a Shutdown RPC to each one of them Collect
// the number of jobs each work has performed.
func (mr *MapReduce) KillWorkers() *list.List {
  l := list.New()
  for _, w := range mr.Workers {
    DPrintf("DoWork: shutdown %s\n", w.address)
    args := &ShutdownArgs{}
    var reply ShutdownReply;
    ok := call(w.address, "Worker.Shutdown", args, &reply)
    if ok == false {
      fmt.Printf("DoWork: RPC %s shutdown error\n", w.address)
    } else {
      l.PushBack(reply.Njobs)
    }
  }
  return l
}

func (mr *MapReduce) RunMaster() *list.List {
  // Your code here
  //Assgining the work	 	 
  //Assigning the Map work

  for i := 0; i < mr.nMap; i++ {
      //Check if there is worker available
      cur_wk := <- mr.registerChannel
      //Distribute the job
      go func(cur_wk string, i int) {
      	 for {
              args := &DoJobArgs{mr.file, "Map", i, mr.nReduce}
	      var reply DoJobReply
	      ok := call(cur_wk, "Worker.DoJob", args, &reply)
	      if ok == true {
		       //Meaning the job has been sucessfully done
		       mr.registerChannel <- cur_wk
		       mr.WorkerChannel <- 1
		       break
	      } else {
		       //Meaning the job somehow failed
	      	       fmt.Printf("Distributing Job fails, and will go")
		       cur_wk = <- mr.registerChannel
		       continue
	      }
         }
      }(cur_wk, i)
  }     
  //Wait for all the Map job to be done
  j1 := 0
  L1: for {
	        select {
			case <- mr.WorkerChannel: 
			     	    j1++
				    if j1 > mr.nMap - 1 {
						break L1
				    }
		       }
         }	
  //Assigning the Reduce Job
  for i := 0; i < mr.nReduce; i++ {
      //Check if there is worker available
      cur_wk := <- mr.registerChannel
      //Distribute the job
      go func(cur_wk string, i int) {
      	 for {
              args := &DoJobArgs{mr.file, "Reduce", i, mr.nMap}
	      var reply DoJobReply
	      ok := call(cur_wk, "Worker.DoJob", args, &reply)
	      if ok == true {
		       //Meaning the job has been sucessfully done
		       mr.registerChannel <- cur_wk
		       mr.WorkerChannel <- 1
		       break
	      } else {
		       //Meaning the job somehow failed
		       fmt.Printf("Distributing Job fails, and will go")
		       cur_wk = <- mr.registerChannel
		       continue
	      }
 	 }
      }(cur_wk, i)
  }     
  //Wait for all the Reduce job to be done
  j2 := 0
  L2: for {
	        select {
			case <- mr.WorkerChannel: 
			     	    j2++
				    if j2 > mr.nReduce - 1 {
						break L2
				    }
		       }
         }	 

  return mr.KillWorkers()
}
