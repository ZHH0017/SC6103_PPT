package RingBasedAlgorithm

import (
	"testing"
	"time"
)

func TestRing(t *testing.T) {
	Portlist := []string{"1234", "1235", "1236"}
	Clist := CreateClientByPorts(Portlist)
	for i := 0; i < len(Portlist); i++ {
		go func(c Client) {
			c.BootClient()
		}(Clist[i])
	}
	time.Sleep(100 * time.Second)
}
