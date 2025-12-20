package main

import (
	"dumont/config"
	"dumont/database"
	"dumont/producers"
	"dumont/runner"
	"fmt"
	"time"
)

/*
- Add Server to provide statistics and prometheus metrics(maybe redis or sqlite???)
- Add other transaction types DELETE, CREATE TABLE, ALTER TABLE, DROP TABLE
- Add pagination??(maria db binlog paginate results)
- Make Dockerfile and compose file work( currently it doesn't :c )
- Add logs
*/

func main() {
	//config.SetEnvExample() //Used to set env for debugging
	config := config.LoadEnv()
	fmt.Println(config)

	producer := producers.Connect(config)
	db := database.Connect(config)
	r, _ := db.GetPath()

	runnerConfig := runner.RunnerConfig{
		DbConnection: &db,
		BinlogPath:   r,
		Producer:     &producer,
	}

	ch := make(chan runner.CommandExecution)
	runner.StartConsumers(config.MaxConsumers, ch)

	for {
		run(runnerConfig, ch)
		time.Sleep(time.Duration(config.ExecuteInterval) * time.Second)
		runnerConfig.DateFilter = time.Now().Format("2006-01-02 15:04:05")
	}

}

func run(runnerConfig runner.RunnerConfig, ch chan runner.CommandExecution) {
	runnerConfig.Execute(ch)
}
