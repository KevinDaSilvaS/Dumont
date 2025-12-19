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

	for {
		run(runnerConfig)
		time.Sleep(2 * time.Second)
		runnerConfig.DateFilter = time.Now().Format("2006-01-02 15:04:05")
	}

}

func run(runnerConfig runner.RunnerConfig) {
	fmt.Println(len(runnerConfig.Execute()))
}
