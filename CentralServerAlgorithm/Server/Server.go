package Server

import (
	"NTU_DSppt/CentralServerAlgorithm/Arg"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Server struct {
	IpQueue   Arg.Queue
	PortQueue Arg.Queue
	mu        sync.Mutex
	token     bool
}

func (s *Server) ExecuteRequest(arg Arg.Arg, reply *Arg.Reply) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if arg.RequireToken && s.token && s.PortQueue.IsEmpty() {
		log.Println("Port:", arg.WaitPort, "have token")
		s.token = false
		reply.GetToken = true
		return nil
	}

	if arg.RealeaseToken && !s.PortQueue.IsEmpty() {
		log.Println("Port:", arg.WaitPort, "realease token")
		s.token = true
		reply.ReleaseToken = true
		port := s.PortQueue.Pop()
		s.IpQueue.Pop()
		log.Println("Port:", port, "have token")
		fmt.Println()
		s.token = false
		reply.GetToken = false
		Sendmsg2Client("ok", port)
		return nil
	} else if arg.RealeaseToken {
		log.Println("Port:", arg.WaitPort, "realease token")
		s.token = true
		reply.ReleaseToken = true
	}

	if !s.token {
		log.Println("Port:", arg.WaitPort, "wait for token")
		s.IpQueue.Push(arg.WaitIp)
		s.PortQueue.Push(arg.WaitPort)
		reply.GetToken = false
		return nil
	}
	return nil
}

func Sendmsg2Client(msg string, port string) {
	//conn, err := net.Dial("tcp", "remote_server_ip:remote_server_port")
	ip := "localhost:" + port
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal("Dial fail")
	}
	defer conn.Close()
	//fmt.Println("ok")
	data := []byte(msg)
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal("send msg error:", err)
		return
	}
	return
}

func BootServer() {
	s := new(Server)
	s.token = true
	//ch := make(chan int)
	//rpc.RegisterName("Server", new(Server))
	rpc.Register(s)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listen TCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accepet error:", err)
		}

		go func() {
			//rpc.ServeConn(conn) 是一个阻塞调用，它会等待客户端的 RPC 请求，直到完成请求的处理。
			//如果在此调用之后执行了阻塞的代码，整个程序可能会停滞不前，无法响应其他请求或事件
			rpc.ServeConn(conn) //为什么不用go就会卡住？
		}()
		time.Sleep(50 * time.Millisecond)
	}

}
