package runner

import (
	"dumont/config"
	"dumont/database"
	"dumont/producers"
	"fmt"
)

type RunnerConfig struct {
	DbConnection         *database.Database
	DateFilter           string
	BinlogPath           string
	Producer             *producers.Producer
	ReadFromRemoteConfig *config.Config
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
	if r.ReadFromRemoteConfig != nil && r.ReadFromRemoteConfig.ReadFromRemote {
		args = append(args, "--read-from-remote-server")
		args = append(args, fmt.Sprintf("--host=%s", r.ReadFromRemoteConfig.Host))
		args = append(args, fmt.Sprintf("--user=%s", r.ReadFromRemoteConfig.User))
		args = append(args, fmt.Sprintf("--password=%s", r.ReadFromRemoteConfig.Passwd))
		args = append(args, fmt.Sprintf("--port=%s", r.ReadFromRemoteConfig.Port))
	}

	if r.DateFilter != "" {
		args = append(args, fmt.Sprintf("--start-datetime=\"%s\"", r.DateFilter))
	}

	args = append(args, "-vv")

	if r.ReadFromRemoteConfig.ReadFromRemote {
		args = append(args, fileName)
	} else {
		args = append(args, fmt.Sprintf("%s%s", r.BinlogPath, fileName))
	}

	return args
}
