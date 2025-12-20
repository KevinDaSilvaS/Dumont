package main

import (
	"dumont/config"
	"dumont/database"
	"dumont/runner"
	"fmt"
	"time"
)

/*
- Add producer
- Add transaction processing
- Add README.md
- Add Server to provide statistics and prometheus metrics(maybe redis or sqlite???)
- Add other transaction types DELETE, CREATE TABLE, ALTER TABLE, DROP TABLE
- Add pagination??
*/

func main() {
	config.SetEnvExample()
	config := config.LoadEnv()
	fmt.Println(config)

	db := database.Connect(config)
	r, _ := db.GetPath()

	runnerConfig := runner.RunnerConfig{
		DbConnection: &db,
		BinlogPath:   r,
	}

	ch := make(chan []string)
	runner.StartConsumers(config.MaxConsumers, ch)

	for {
		run(runnerConfig, ch)
		time.Sleep(time.Duration(config.ExecuteInterval) * time.Second)
		runnerConfig.DateFilter = time.Now().Format("2006-01-02 15:04:05")
	}

}

func run(runnerConfig runner.RunnerConfig, ch chan []string) {
	fmt.Println(len(runnerConfig.Execute(ch)))
}
