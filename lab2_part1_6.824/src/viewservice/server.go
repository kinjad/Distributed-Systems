
package viewservice

import "net"
import "net/rpc"
import "log"
import "time"
import "sync"
import "fmt"
import "os"

type ViewServer struct {
  mu sync.Mutex
  l net.Listener
  dead bool
  me string

  // Your declarations here.
  Cur_View View
  Prim_Ack bool
  deadflag_prim bool
  deadflag_back bool
  Rec_Pin map[string]time.Time
  deadcount_prim int
  deadcount_back int
}

//
// server Ping RPC handler.
//
func (vs *ViewServer) Ping(args *PingArgs, reply *PingReply) error {

  // Your code here.

  if args.Viewnum == vs.Cur_View.Viewnum && args.Me == vs.Cur_View.Primary { //Meaning this is ack_cur from the primary
     vs.Prim_Ack = true
  }
  
  //Heard from the primary				
  if args.Me == vs.Cur_View.Primary { 
    vs.deadflag_prim = false
    value1 := vs.Cur_View.Viewnum - args.Viewnum
    if value1 != 1 && value1 != 0 { //The primary crashes
       if vs.Prim_Ack == true {
       	  if vs.Cur_View.Backup != "" {
       	     vs.Cur_View.Primary = vs.Cur_View.Backup
	     delete(vs.Rec_Pin, vs.Cur_View.Primary)
	     vs.Cur_View.Backup = ""
	     vs.Cur_View.Viewnum++	  
          }
       }
    }
  }
  //Heard from the Backup
  if args.Me == vs.Cur_View.Backup { 
    vs.deadflag_back = false
    value2 := vs.Cur_View.Viewnum - args.Viewnum
    if value2 != 1 && value2 != 0 { //The backup crashes
       if vs.Prim_Ack == true {
	  vs.Cur_View.Backup = ""
	  vs.Cur_View.Viewnum++	  
       }
    }
  }


  old_len := len(vs.Rec_Pin)
  //Update the dic
  if vs.Prim_Ack == true {
     vs.Rec_Pin[args.Me] = time.Now()  
  }
  //If there is no primary yet
  if vs.Cur_View.Primary == "" {
     if vs.Prim_Ack == true {
          vs.Cur_View.Primary = args.Me
     }
  }

  //If there is a primary yet no backup and a new guys comes
  if vs.Cur_View.Primary != "" && vs.Cur_View.Backup == "" && vs.Cur_View.Primary != args.Me { //Newcomer to become backup
     if vs.Prim_Ack == true {
        vs.Cur_View.Backup = args.Me
	/*
	fmt.Printf("New Backup\n")
        fmt.Printf("%v\n", vs.Cur_View.Backup)
        */
     }
  }

  //Other cases

  if len(vs.Rec_Pin) > old_len {
     if vs.Prim_Ack == true {
        vs.Cur_View.Viewnum++
     }
  }    
  reply.View = vs.Cur_View
  if args.Me == vs.Cur_View.Primary && args.Viewnum != vs.Cur_View.Viewnum {//Meaning this is from primary
     	vs.Prim_Ack = false
  }
  return nil
}

// 
// server Get() RPC handler.
//
func (vs *ViewServer) Get(args *GetArgs, reply *GetReply) error {

  // Your code here.
  reply.View = vs.Cur_View
  return nil
}


//
// tick() is called once per PingInterval; it should notice
// if servers have died or recovered, and change the view
// accordingly.
//
func (vs *ViewServer) tick() {

  // Your code here.
  //Monitoring Primary
  if vs.deadcount_prim == DeadPings {
     if vs.Prim_Ack == true {
        delete(vs.Rec_Pin, vs.Cur_View.Primary)
        if vs.Cur_View.Backup != "" {
     	   vs.Cur_View.Primary = vs.Cur_View.Backup
	   vs.Cur_View.Backup = ""
	   vs.Cur_View.Viewnum++
	   /*
	   fmt.Printf("The Backup takes over\n")
	   fmt.Printf("%v\n", vs.Cur_View.Primary)
	   fmt.Printf("%v\n", vs.Cur_View.Backup)
	   */
        }
        vs.deadcount_prim = 0
     }
  }
  if vs.deadflag_prim == true { //Nothing heard from primary in one turn
     vs.deadcount_prim += 1
  } else {
	     vs.deadcount_prim = 0
	     vs.deadflag_prim = true
  }

  //Monitoring Backup
  if vs.deadcount_back == DeadPings {
     if vs.Prim_Ack == true {
        delete(vs.Rec_Pin, vs.Cur_View.Backup)
        if vs.Cur_View.Backup != "" {
	   vs.Cur_View.Backup = ""
	   vs.Cur_View.Viewnum++
        }
        vs.deadcount_back = 0
     }
  }
  if vs.deadflag_back == true { //Nothing heard from backup in one turn
     vs.deadcount_back += 1
  } else {
	     vs.deadcount_back = 0
	     vs.deadflag_back = true
  }


}

//
// tell the server to shut itself down.
// for testing.
// please don't change this function.
//
func (vs *ViewServer) Kill() {
  vs.dead = true
  vs.l.Close()
}

func StartServer(me string) *ViewServer {
  vs := new(ViewServer)
  vs.me = me
  // Your vs.* initializations here.
  vs.Cur_View = *new(View)  
  vs.Cur_View.Viewnum = 0
  vs.Cur_View.Primary = ""
  vs.Cur_View.Backup = ""
  vs.Prim_Ack = true
  vs.deadflag_prim = true //Meaning nothing heard from the primary
  vs.deadcount_prim = 0
  vs.deadflag_back = true //Meaning nothing heard from the backup
  vs.deadcount_back = 0
  vs.Rec_Pin = make(map[string]time.Time)



  // tell net/rpc about our RPC server and handlers.
  rpcs := rpc.NewServer()
  rpcs.Register(vs)

  // prepare to receive connections from clients.
  // change "unix" to "tcp" to use over a network.
  os.Remove(vs.me) // only needed for "unix"
  l, e := net.Listen("unix", vs.me);
  if e != nil {
    log.Fatal("listen error: ", e);
  }
  vs.l = l

  // please don't change any of the following code,
  // or do anything to subvert it.

  // create a thread to accept RPC connections from clients.
  go func() {
    for vs.dead == false {
      conn, err := vs.l.Accept()
      if err == nil && vs.dead == false {
        go rpcs.ServeConn(conn)
      } else if err == nil {
        conn.Close()
      }
      if err != nil && vs.dead == false {
        fmt.Printf("ViewServer(%v) accept: %v\n", me, err.Error())
        vs.Kill()
      }
    }
  }()

  // create a thread to call tick() periodically.
  go func() {
    for vs.dead == false {
      vs.tick()
      time.Sleep(PingInterval)
    }
  }()

  return vs
}
