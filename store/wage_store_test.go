package store

import (
	"fmt"
	"testing"
)

func TestWageStorePull(t *testing.T) {
	s := WageStore{}
	wages, _ := s.Pull("http://www.wdol.gov/wdol/scafiles/std/05-2003.txt")
	fmt.Println(wages)
}
