package runner

import (
	"dumont/parser"
	"fmt"
	"os/exec"
	"sync"
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

		var wg sync.WaitGroup
		for i, transaction := range transactions {
			wg.Add(1)
			go worker(i, transaction, &wg)
		}

		wg.Wait()
		fmt.Println("Received On Consumer #", tag, executeArgs, transactions)
	}
}

func StartConsumers(totalConsumers int, ch <-chan []string) {
	for i := range totalConsumers {
		go RunCommand(ch, i)
	}
}

func worker(id int, transaction string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Transaction parsed:", id, parser.ParseTransactionQuery(transaction))
}
