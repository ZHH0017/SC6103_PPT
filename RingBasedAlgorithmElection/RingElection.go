package RingBasedAlgorithmElection

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

const PARTICIPANT = 1
const NONPARTICIPANT = 2
const COORDINATOR = 3

type LocalNode struct {
	mu           sync.Mutex
	LocalPort    string
	State        int
	Version      int
	HaveElection bool
	NextNode     *LocalNode
}

func NewNode(port string) *LocalNode {
	n := &LocalNode{
		LocalPort:    port,
		HaveElection: false,
		State:        NONPARTICIPANT,
		NextNode:     nil,
	}
	return n
}

func (n *LocalNode) Clear() {
	n.HaveElection = false
	n.State = NONPARTICIPANT
}

func CreateRingByPorts(Ports []string) *LocalNode {
	lenP := len(Ports)
	if lenP == 1 {
		return NewNode(Ports[0])
	}
	start := NewNode(Ports[0])
	var end *LocalNode
	i := 1
	for i < lenP {
		n := NewNode(Ports[i])
		AppendNode := start
		for AppendNode.NextNode != nil {
			AppendNode = AppendNode.NextNode
		}
		AppendNode.AppendNode(n)
		if i == lenP-1 {
			end = n
		}
		i++
	}
	end.AppendNode(start)
	return start
}

func (n *LocalNode) AppendNode(next *LocalNode) {
	n.NextNode = next
}

func (n *LocalNode) Boot(version int) {
	n.State = NONPARTICIPANT
	n.Version = version
	fmt.Println(n.LocalPort, "version is", n.Version)
	methodname := fmt.Sprintf("LocalNode%s", n.LocalPort)
	err := rpc.RegisterName(methodname, n)
	if err != nil {
		log.Fatal(err)
	}
	address := "localhost:" + n.LocalPort
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
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

func (n *LocalNode) ForwardMsg(version int) {
	address := "localhost" + ":" + n.NextNode.LocalPort
	conn, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	methodname := fmt.Sprintf("LocalNode%s.Judge", n.NextNode.LocalPort)
	arg := &Arg{
		Version:  version,
		Election: false,
	}
	var reply Reply
	err = conn.Call(methodname, arg, &reply)
	if err != nil {
		log.Fatal(err)
	}
}

func (n *LocalNode) ForwardElection(version int, ElectionFlag bool) {
	address := "localhost" + ":" + n.NextNode.LocalPort
	conn, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	methodname := fmt.Sprintf("LocalNode%s.Judge", n.NextNode.LocalPort)
	arg := &Arg{
		Version:  version,
		Election: ElectionFlag,
	}
	var reply Reply
	err = conn.Call(methodname, arg, &reply)
	if err != nil {
		log.Fatal(err)
	}
}

func (n *LocalNode) StartElection() {
	log.Println(n.LocalPort, "start election")
	address := "localhost" + ":" + n.NextNode.LocalPort
	conn, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	methodname := fmt.Sprintf("LocalNode%s.Judge", n.NextNode.LocalPort)
	arg := &Arg{
		Version:  n.Version,
		Election: false,
	} //arg要加取地址符号
	var reply Reply
	err = conn.Call(methodname, arg, &reply)
	if err != nil {
		log.Fatal(err)
	}
}

func (n *LocalNode) Judge(arg Arg, reply *Reply) error {

	n.mu.Lock()
	defer n.mu.Unlock()
	if arg.Version > n.Version {
		if arg.Election == true && n.HaveElection == false {
			log.Println(n.LocalPort, "forward election")
			n.HaveElection = true
			go n.ForwardElection(arg.Version, true)
		} else if arg.Election == false {
			log.Println(n.LocalPort, "forward")
			go n.ForwardMsg(arg.Version)
		}

	} else if arg.Version < n.Version && n.State == NONPARTICIPANT {
		log.Println(n.LocalPort, "forward self")
		go n.ForwardMsg(n.Version)
	} else if arg.Version < n.Version && n.State == PARTICIPANT {
		log.Println(n.LocalPort, "dont forward")
		//不做任何操作
	}
	n.State = PARTICIPANT
	if arg.Version == n.Version && arg.Election == false && n.HaveElection == false {
		log.Println(n.LocalPort, "forward election")
		go n.ForwardElection(n.Version, true)
		n.HaveElection = true
		n.State = COORDINATOR
	}

	if arg.Version == n.Version && arg.Election == true && n.HaveElection == true {
		log.Println(n.LocalPort, "become coordinator")
	}
	reply.Replymsg = 1
	return nil
}
