package runner

import (
	"dumont/database"
	"dumont/parser"
	"fmt"
	"os/exec"
)

type RunnerConfig struct {
	DbConnection *database.Database
	DateFilter   string
	BinlogPath   string
}

func (r RunnerConfig) Execute() []string {
	files, _ := r.DbConnection.GetBinLogFiles()

	transactions := []string{}
	for _, fileName := range files {
		cmd := exec.Command("mariadb-binlog", r.GetArgs(fileName)...)
		_ = cmd.Wait()
		out, _ := cmd.Output()

		transactions = append(transactions, parser.ParseTransactions(out)...)
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
