package CentralServerAlgorithm

import (
	"NTU_DSppt/CentralServerAlgorithm/Arg"
	"NTU_DSppt/CentralServerAlgorithm/Client"
	"fmt"
	"log"
	"net/rpc"
	"testing"
	"time"
)

func Test_CSA(t *testing.T) {
	ports := []string{"2001", "2006", "2002", "2008"} //p1,p6,p2,p8
	sleepTimes := []time.Duration{1, 1, 1, 2}         //time = 1,2,3,5
	for i := 0; i < len(ports); i++ {
		time.Sleep(sleepTimes[i] * time.Second)
		go func(p string) {
			BootClient(p)
		}(ports[i])
	}

	time.Sleep(30 * time.Second)
}

func BootClient(port string) {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply Arg.Reply
	arg := Client.RequestMsg(port)

	err = client.Call("Server.ExecuteRequest", arg, &reply)
	if err != nil {
		log.Fatal(err)
	}
	if !reply.GetToken {
		Client.Receivemsg2Server(port) //等待回复
	}
	//拿到token以后的操作
	log.Println(port, " Get token successfully")
	fmt.Println()
	Execute()
	log.Println(port, " Execute process successfully")
	arg = Client.ReleaseMsg(port)
	err = client.Call("Server.ExecuteRequest", arg, &reply)
	if err != nil {
		log.Fatal(err)
	}
	if reply.ReleaseToken {
		log.Println(port, " Release token successfully!")
	}
}

func Execute() {
	time.Sleep(3 * time.Second)
}
