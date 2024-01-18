package RANode

import (
	"fmt"
	"testing"
)

func TestNewNode(t *testing.T) {
	ports := []string{"1000", "1001", "1002", "1003", "1004"}
	RAlist := NewNodewithlocal(ports)
	for i := 0; i < len(RAlist); i++ {
		fmt.Println("第", i, "个node中的信息")
		fmt.Println(RAlist[i].IP, "  ", RAlist[i].Port, "  ", RAlist[i].Info.threshold)
		for j := 0; j < len(RAlist[i].OtherNode); j++ {
			fmt.Println(RAlist[i].OtherNode[j].IP, "   ", RAlist[i].OtherNode[j].Port, "  ", RAlist[i].OtherNode[j].Info.threshold)
		}
	}
}
