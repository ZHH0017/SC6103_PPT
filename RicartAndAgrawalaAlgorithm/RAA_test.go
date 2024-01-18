package RicartAndAgrawalaAlgorithm

import (
	RANode "NTU_DSppt/RicartAndAgrawalaAlgorithm/RAClient"
	"testing"
	"time"
)

func Test_RAA(t *testing.T) {
	ports := []string{"1000", "1001", "1002", "1003"}
	RAList := RANode.NewNodewithlocal(ports)
	go RANode.BootNodes(RAList)

	for i := 0; i < len(RAList); i++ {
		go RAList[i].Begin()
	}
	time.Sleep(200 * time.Second)
}
