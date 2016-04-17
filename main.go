package main

import (
	"fmt"

	"github.com/mmbros/fmvno/fw"
	"github.com/mmbros/fmvno/util"
)

func main() {
	fw.InitConfigEngineMobile()

	accounts := fw.NewAccounts(100, 160, 10)
	fmt.Printf("%v\n", accounts)

	for j := 0; j < 20; j++ {
		fmt.Printf("acc[%d] = %v\n", j, accounts.ListDaSpedire[j])
	}

	q := util.NewCircularFifoQueue(10)
	q.Push(1)
	x := q.Pop().(int)

	fmt.Printf("Pop = %v\n", x)
}
