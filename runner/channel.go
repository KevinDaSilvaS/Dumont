package runner

import (
	"dumont/parser"
	"fmt"
	"os/exec"
)

func SendCommand(ch chan<- []string, executeArgs []string) {
	ch <- executeArgs
}

func RunCommand(ch <-chan []string, tag int) {
	for {
		executeArgs := <-ch
		cmd := exec.Command("mariadb-binlog", executeArgs...)
		_ = cmd.Wait()
		out, _ := cmd.Output()
		transactions := parser.ParseTransactions(out)
		fmt.Println("Received On Consumer #", tag, executeArgs, transactions)
	}
}

func StartConsumers(totalConsumers int, ch <-chan []string) {
	for i := range totalConsumers {
		go RunCommand(ch, i)
	}
}
