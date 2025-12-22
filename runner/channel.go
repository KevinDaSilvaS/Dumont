package runner

import (
	"dumont/parser"
	"log/slog"
	"os/exec"
	"sync"
)

func SendCommand(ch chan<- CommandExecution, cmdExecution CommandExecution) {
	ch <- cmdExecution
}

func RunCommand(ch <-chan CommandExecution, tag int) {
	slog.Info("Received command on consumer", slog.Int("#", tag))
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

func worker(id int, transaction string, executor CommandExecution, wg *sync.WaitGroup) {
	defer wg.Done()
	t := parser.ParseTransactionQuery(transaction, executor.RunnerConfig.DbConnection)
	executor.Producer.Publish(t)
	slog.Info("Transaction processed", slog.Int("worker_id", id), slog.Any("t", t))
}
