package Ring

import (
	"fmt"
	"testing"
)

func TestCreateRingByPoint(t *testing.T) {
	s := CreateRingByPoint([]string{"1100", "1101", "1102", "1103", "1104"})
	for i := 0; i < 10; i++ {
		fmt.Println(s.LocalPort)
		s = s.NextNode
	}

}
