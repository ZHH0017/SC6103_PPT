package RANode

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

const RELEASED = 1
const WANTED = 2
const HELD = 3
const localIp = "localhost"

type RANode struct {
	IP        string
	Port      string
	wg        sync.WaitGroup
	OtherNode []*RANode
	Info      RALocalInfo
}

type RALocalInfo struct {
	mu           sync.Mutex
	TimeStamp    time.Time
	cond         *sync.Cond //定义加*是因为多个结构体要用相同的cond
	State        int
	Totalrecvmsg int
	threshold    int
}

func (R *RANode) Judge(arg Arg, reply *Reply) error {
	version, err := strconv.Atoi(R.Port)
	if err != nil {
		log.Fatal(err)
	}
	if R.Info.State == HELD || (R.Info.State == WANTED && (R.Info.TimeStamp.Before(arg.TimeStamp) || R.Info.TimeStamp.Equal(arg.TimeStamp) && version < arg.Version)) {
		//log.Println("wait")
		R.Info.mu.Lock()
		R.Info.cond.Wait() //cond.Wait会在内部自动解锁互斥锁，然后在Wait过后重新锁回去，所以该语句前得加lock，后加unlock!
		R.Info.mu.Unlock()
	}
	reply.ReplyMsg = 1
	return nil
}

func (R *RANode) EnterState() {
	R.Info.State = WANTED
	R.Info.TimeStamp = time.Now()
	log.Println(R.Port, "timestamp is ", R.Info.TimeStamp)
	flag := R.RequestBroadcast()
	if flag {
		R.Info.State = HELD
		log.Println(R.Port, "is HELD")
	}
}

func (R *RANode) ExistState() {
	R.Info.State = RELEASED
	log.Println(R.Port, "Released")
	R.Info.cond.Broadcast() //reply to all queued requests
}

func NewNodewithlocal(ports []string) []*RANode {
	RAlist := make([]*RANode, len(ports))
	for i := 0; i < len(ports); i++ {
		Node := &RANode{
			IP:        localIp,
			Port:      ports[i],
			OtherNode: make([]*RANode, 0),
		}
		RAlist[i] = Node
	}
	for i := 0; i < len(ports); i++ {
		for j := 0; j < len(ports); j++ {
			if j != i {
				RAlist[i].OtherNode = append(RAlist[i].OtherNode, RAlist[j])
			}
		}
	}
	for i := 0; i < len(ports); i++ {
		RAlist[i].NewLocalInfo()
	}
	return RAlist
}

func BootNodes(Nodes []*RANode) {
	for i := 0; i < len(Nodes); i++ {
		go func(index int) {
			//rpc.RegisterName("RANode", Nodes[index])//Registername中name是服务的别称
			Node := Nodes[index]
			//register别称不能重复
			name := fmt.Sprintf("RANode%s", Node.Port)//Register相当于使用registername,name参数为结构体名字
			err := rpc.RegisterName(name, Node)
			if err != nil {
				log.Fatal(err)
			}
			address := Node.IP + ":" + Node.Port
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
		}(i)
	}
}

func (R *RALocalInfo) LocalInfoClear(threshold int) {
	R.State = RELEASED
	R.Totalrecvmsg = 0
	R.threshold = threshold
}

func (R *RANode) NewLocalInfo() {
	Len := len(R.OtherNode)
	R.Info = RALocalInfo{
		State:        RELEASED,
		Totalrecvmsg: 0,
		threshold:    Len,
		cond:         sync.NewCond(&R.Info.mu),
	}
}

func (R *RANode) RequestBroadcast() bool {
	R.wg.Add(len(R.OtherNode))
	for i := 0; i < len(R.OtherNode); i++ {
		go func(index int) {
			address := R.OtherNode[index].IP + ":" + R.OtherNode[index].Port
			//fmt.Println(R.Port, "send to ", address)
			conn, err := rpc.Dial("tcp", address)
			if err != nil {
				log.Fatal(err)
			}
			version, err := strconv.Atoi(R.Port)
			if err != nil {
				log.Fatal(err)
			}
			arg := &Arg{
				RequestMsg: 1,
				TimeStamp:  R.Info.TimeStamp,
				Version:    version,
			}
			var reply Reply
			methodname := fmt.Sprintf("RANode%s.Judge", R.OtherNode[index].Port)
			err = conn.Call(methodname, arg, &reply)
			if reply.ReplyMsg == 1 {
				R.Info.mu.Lock()
				R.Info.Totalrecvmsg++
				R.Info.mu.Unlock()
				R.wg.Done()
			}
		}(i)
	}
	R.wg.Wait()
	if R.Info.Totalrecvmsg == R.Info.threshold {
		return true
	} else {
		log.Println("Totalrecvmsg is ", R.Info.Totalrecvmsg)
		log.Println("some error happen, wg done but msg num dont")
		return false
	}
}

func (R *RANode) ExecuteProcess() {
	time.Sleep(8 * time.Second)
	log.Println(R.Port, "process done!")
}

func (R *RANode) Begin() {
	R.Info.LocalInfoClear(len(R.OtherNode))
	Rand := rand.Intn(5)
	time.Sleep(time.Duration(Rand) * time.Second)
	log.Println(R.Port, "want to enter state")
	R.EnterState()
	R.ExecuteProcess()
	R.ExistState()
}
