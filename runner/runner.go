package runner

import (
	"dumont/database"
	"dumont/producers"
	"fmt"
)

type RunnerConfig struct {
	DbConnection *database.Database
	DateFilter   string
	BinlogPath   string
	Producer     *producers.Producer
}

type CommandExecution struct {
	Producer    *producers.Producer
	ExecuteArgs []string
}

func (r RunnerConfig) Execute(ch chan CommandExecution) []string {
	files, _ := r.DbConnection.GetBinLogFiles()

	transactions := []string{}
	for _, fileName := range files {
		cmd := CommandExecution{ExecuteArgs: r.GetArgs(fileName), Producer: r.Producer}
		SendCommand(ch, cmd)
	}
	return transactions
}

func (r RunnerConfig) GetArgs(fileName string) []string {
	args := []string{"--base64-output=decode-rows"}
	if r.DateFilter != "" {
		args = append(args, fmt.Sprintf("--start-datetime=\"%s\"", r.DateFilter))
	}

	args = append(args, "-vv")
	args = append(args, fmt.Sprintf("%s%s", r.BinlogPath, fileName))
	return args
}
