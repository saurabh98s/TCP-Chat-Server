//Building a basic TCP concurrent channel chat server.

package main
import ("net"
          "log"
        "bufio"
      "fmt")

func main()  {
  ln,err:=net.Listen("tcp","127.0.0.1:23")
  if err != nil {
    log.Println(err.Error())
  }
  //we will need channels for incoming connections,dead connections and messages.
  // so we build 3 channels as below
  aconns:=make(map[net.Conn]int)
  conns:=make(chan net.Conn)
  dconns:=make(chan net.Conn)
  msgs:=make(chan string)
  i:=0
  go func(){
      for{
        conn,err:=ln.Accept()
        if err != nil {
          log.Println(err.Error())
        }
        conns<-conn
      }
  }()
  for {
      select{
        //Reading incoming connections
      case conn:=<-conns:
        aconns[conn]=i
        i++
        //Once we have the connections we read messages from it...
        go func(conn net.Conn, i int){
          rd:=bufio.NewReader(conn)
              for{
                    m,err:=rd.ReadString('\n')
                    if err != nil {
                      break
                    }
                    msgs<-fmt.Sprintf("CLient %v:%v",i,m)
              }
          //done reading from file!!
          dconns<-conn
          }(conn,i)
        case msg:=<-msgs:
          // we have to broadcast it to all connections.
          for conn:=range aconns{
            conn.Write([]byte(msg))
          }
           case dconn:=<-dconns:
                log.Printf("CLient %v has gone!\n",aconns[dconn])
                delete(aconns,dconn)

      }

  }

}
