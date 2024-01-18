package Client

import (
	"NTU_DSppt/CentralServerAlgorithm/Arg"
	"fmt"
	"log"
	"net"
)

// 1代表request,2代表release
type RPC struct {
	Request int
}

func Receivemsg2Server(port string) bool {
	//listener, err := net.Listen("tcp", "local_server_ip:local_server_port")
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Accepet error:", err)
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if string(data[:n]) == "ok" {
		return true
	}
	fmt.Println(string(data[:n]))
	fmt.Println("error")
	return false

}

func RequestMsg(port string) Arg.Arg {
	arg := Arg.Arg{
		WaitPort:      port,
		RequireToken:  true,
		RealeaseToken: false,
	}
	return arg
}

func ReleaseMsg(port string) Arg.Arg {
	arg := Arg.Arg{
		WaitPort:      port,
		RequireToken:  false,
		RealeaseToken: true,
	}
	return arg
}
