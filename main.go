package main

import (
	"dumont/database"
	"dumont/runner"
	"fmt"
	"time"
)

func main() {
	db := database.Connect()
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
