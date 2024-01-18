package RingBasedAlgorithm

import (
	"log"
	"net"
	"time"
)

type Client struct {
	HaveToken      bool
	OpenPort       string
	NextClientIp   string
	NextClientPort string
}

func NewClient(havetoken bool, openport string, nextip string, nextport string) *Client {
	return &Client{
		HaveToken:      havetoken,
		OpenPort:       openport,
		NextClientIp:   nextip,
		NextClientPort: nextport,
	}
}
func (C *Client) BootClient() {
	nextip := "localhost:" + C.NextClientPort
	localip := "localhost:" + C.OpenPort
	listener, err := net.Listen("tcp", localip)
	if err != nil {
		log.Fatal(err)
	}
	var lconn net.Conn
	var conn net.Conn
	go func() {
		lconn, err = listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		conn, err = net.Dial("tcp", nextip)
		if err != nil {
			log.Fatal(err)
		}
	}()
	//go func是共享变量的,如果使用go func(i int)开启多个进程，那么这个i相当于副本传入参数
	time.Sleep(1 * time.Second) //sleep等待所有客户端都互相连接，相当于初始化的时间

	for true {
		if C.HaveToken == false {
			C.WaitToken(lconn)
		} else {
			C.ExecuteProcess()
			C.SendToken2Next(conn)
		}
	}
}

func (C *Client) SendToken2Next(conn net.Conn) {
	_, err := conn.Write([]byte("Token"))
	//log.Println("send token to: ", C.NextClientPort)
	if err != nil {
		log.Fatal(err)
	}
	C.HaveToken = false
}

func (C *Client) WaitToken(conn net.Conn) {
	//conn, err := listener.Accept()
	//if err != nil {
	//	log.Fatal(err)
	//}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	ip := "localhost:" + C.OpenPort
	log.Println("ip: ", ip, " have ", string(data[:n]))
	C.HaveToken = true
}

func (C *Client) ExecuteProcess() {
	ip := "localhost:" + C.OpenPort
	log.Println("ip: ", ip, " is using token")
	time.Sleep(3 * time.Second)
}

func CreateClientByPorts(ports []string) []Client {
	Clist := make([]Client, len(ports))
	for i := 0; i < len(ports); i++ {
		if i == 0 {
			Clist[i] = *NewClient(true, ports[i], "localhost", ports[i+1])
		} else if i == len(ports)-1 {
			Clist[i] = *NewClient(false, ports[i], "localhost", ports[0])
		} else {
			Clist[i] = *NewClient(false, ports[i], "localhost", ports[i+1])
		}

	}
	return Clist
}
