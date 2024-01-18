package RingBasedAlgorithmElection

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRE(t *testing.T) {
	ports := []string{"1000", "1001", "1002", "1003", "1004"}
	versions := make([]int, 0)
	Map := make(map[int]int, 0)
	for i := 0; i < len(ports); i++ {
		rand.Seed(int64(i))
		Rand := rand.Intn(50)
		for true {
			if Map[Rand] == 0 {
				Map[Rand] = 1
				break
			}
		}
		versions = append(versions, Rand)
	}
	Nodes := CreateRingByPorts(ports)
	Node := Nodes
	Nodelist := make([]*LocalNode, 0)
	for i := 0; i < len(ports); i++ {
		Nodelist = append(Nodelist, Node)
		Node = Node.NextNode
	}
	fmt.Println(len(Nodelist))
	for i := 0; i < len(ports); i++ {
		go Nodelist[i].Boot(versions[i])
	}

	go Nodelist[0].StartElection()
	go Nodelist[3].StartElection()
	time.Sleep(100 * time.Second)
}
