package runner

import (
	"dumont/parser"
	"fmt"
	"os/exec"
	"sync"
)

func SendCommand(ch chan<- CommandExecution, cmdExecution CommandExecution) {
	ch <- cmdExecution
}

func RunCommand(ch <-chan CommandExecution, tag int) {
	fmt.Println("Received On Consumer #", tag)
	for {
		cmdExecution := <-ch
		cmd := exec.Command("mariadb-binlog", cmdExecution.ExecuteArgs...)
		_ = cmd.Wait()
		out, _ := cmd.Output()
		transactions := parser.ParseTransactions(out)

		var wg sync.WaitGroup
		for i, transaction := range transactions {
			wg.Add(1)
			go worker(i, transaction, cmdExecution, &wg)
		}

		wg.Wait()
	}
}

func StartConsumers(totalConsumers int, ch <-chan CommandExecution) {
	for i := range totalConsumers {
		go RunCommand(ch, i)
	}
}

func worker(id int, transaction string, producer CommandExecution, wg *sync.WaitGroup) {
	defer wg.Done()
	t := parser.ParseTransactionQuery(transaction)
	producer.Producer.Publish(t)
	fmt.Println("Transaction parsed:", id, t)
}
